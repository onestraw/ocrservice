package main

import (
	"flag"
	"io/ioutil"
	"os"
	"strings"

	"github.com/otiai10/gosseract"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"

	"github.com/onestraw/ocrservice/rpc"
)

const version = "0.2.0"

var (
	flagRabbitMQ  = flag.String("rabbitmq", "amqp://guest:guest@localhost:5672/", "RabbitMQ Address")
	flagQueueName = flag.String("queue_name", "ocrimage", "Queue name for OCR image")
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func worker(data []byte) ([]byte, error) {
	var req rpc.OCRImageRequest
	err := req.Decode(data)
	if err != nil {
		log.Errorf("Decode OCRImageRequest error: %v", err)
		return nil, err
	}

	tempfile, err := ioutil.TempFile("", "ocrserver"+"-")
	if err != nil {
		return nil, err
	}
	defer func() {
		tempfile.Close()
		os.Remove(tempfile.Name())
	}()

	if _, err := tempfile.Write(req.Image); err != nil {
		return nil, err
	}

	client := gosseract.NewClient()
	defer client.Close()

	client.SetImage(tempfile.Name())
	client.Languages = []string{"eng"}
	if langs := req.Lang; langs != "" {
		client.Languages = strings.Split(langs, ",")
	}
	if whitelist := req.Whitelist; whitelist != "" {
		client.SetWhitelist(whitelist)
	}

	text, err := client.Text()
	if err != nil {
		return nil, err
	}

	resp := rpc.OCRImageResponse{
		Version: version,
		Text:    text,
	}
	return resp.Encode()
}

func main() {
	flag.Parse()

	conn, err := amqp.Dial(*flagRabbitMQ)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		*flagQueueName, // name
		false,          // durable
		false,          // delete when unused
		false,          // exclusive
		false,          // no-wait
		nil,            // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	failOnError(err, "Failed to set QoS")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Infof("Get a task...")
			resp, err := worker(d.Body)
			if err != nil {
				log.Errorf("Worker errror: %v", err)
			}

			err = ch.Publish(
				"",        // exchange
				d.ReplyTo, // routing key
				false,     // mandatory
				false,     // immediate
				amqp.Publishing{
					ContentType:   "text/plain",
					CorrelationId: d.CorrelationId,
					Body:          resp,
				})
			failOnError(err, "Failed to publish a message")

			d.Ack(false)
		}
	}()

	log.Printf(" [*] Awaiting RPC requests")
	<-forever
}
