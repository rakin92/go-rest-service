# Simple standard golang service
FROM golang:1.17.1-alpine

RUN apk update && apk upgrade && apk add --no-cache bash git openssh curl

WORKDIR /go-service

COPY . /go-service/
RUN go mod download

ENTRYPOINT ["bash", "./scripts/run-dev.sh"]