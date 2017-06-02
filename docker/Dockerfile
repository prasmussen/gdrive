FROM golang:1.8.3-alpine
MAINTAINER Hernan Garcia <hernandanielg@gmail.com>

RUN apk add --update git && \
    rm -rf /var/cache/apk/*

RUN go get github.com/prasmussen/gdrive

ENTRYPOINT gdrive
