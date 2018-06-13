FROM golang:1.9

MAINTAINER onestraw <hexiaowei91@gmail.com>

RUN apt-get -qq update
RUN apt-get install -y libleptonica-dev libtesseract-dev tesseract-ocr

# Load languages
RUN apt-get install -y \
  tesseract-ocr-eng \
  tesseract-ocr-chi-sim \
  tesseract-ocr-chi-tra

ADD . $GOPATH/src/github.com/onestraw/ocrservice
WORKDIR $GOPATH/src/github.com/onestraw/ocrservice
#RUN go get ./...
RUN go get github.com/otiai10/gosseract
RUN go get github.com/gin-gonic/gin
RUN go get github.com/streadway/amqp
RUN go get github.com/sirupsen/logrus

RUN go install github.com/onestraw/ocrservice/worker
RUN go install github.com/onestraw/ocrservice/backend
RUN go install github.com/onestraw/ocrservice/frontend

RUN ln -s $GOPATH/src/github.com/onestraw/ocrservice/frontend/app/ $GOPATH/src/github.com/onestraw/ocrservice/app
