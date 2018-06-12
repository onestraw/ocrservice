package main

import (
	"flag"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"github.com/onestraw/ocrservice/rpc"
)

const (
	MAX_FILE_SIZE = 4 * 1024 * 1024
)

var (
	flagAddr      = flag.String("addr", ":10002", "listening address")
	flagRabbitMQ  = flag.String("rabbitmq", "amqp://guest:guest@localhost:5672/", "RabbitMQ Address")
	flagQueueName = flag.String("queue_name", "ocrimage", "Queue name for OCR image")
)

func main() {
	flag.Parse()
	gin.DisableConsoleColor()

	router := gin.Default()
	router.POST("/ocrimage", OCRImage)
	router.Run(*flagAddr)
}

func OCRImage(ctx *gin.Context) {
	image, err := ctx.FormFile("file")
	if err != nil {
		log.Errorf("Fetch image error: %v", err)
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
	log.Infof("Image: %s, size: %d", image.Filename, image.Size)
	if image.Size > MAX_FILE_SIZE {
		msg := fmt.Sprintf("Image size (%d) bytes exceed %d bytes", image.Size, MAX_FILE_SIZE)
		log.Errorf(msg)
		ctx.JSON(http.StatusBadRequest, msg)
		return
	}

	buf := make([]byte, image.Size)
	fd, err := image.Open()
	if err != nil {
		log.Errorf("Open image error: %v", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	defer fd.Close()
	count, err := fd.Read(buf)
	if err != nil {
		log.Errorf("Read image error: %v", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	log.Infof("Read %d bytes from %s", count, image.Filename)

	req := rpc.OCRImageRequest{
		Lang:      ctx.PostForm("languages"),
		Whitelist: ctx.PostForm("whitelist"),
		Image:     buf,
	}
	log.Infof("RPC: send task to %s/%s", *flagRabbitMQ, *flagQueueName)
	resp, err := OCR_RPC(*flagRabbitMQ, *flagQueueName, req)
	if err != nil {
		log.Errorf("OCR image rpc error: %v", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"result":  strings.Trim(resp.Text, ctx.PostForm("trim")),
		"version": resp.Version,
	})
}
