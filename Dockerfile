FROM golang:latest

ENV GOPATH /go

RUN mkdir -p "$GOPATH/src/api-response-time"
ADD . "$GOPATH/src/api-response-time"

WORKDIR $GOPATH/src/api-response-time/app
RUN chmod +x script/*
RUN ./script/build

CMD ./app