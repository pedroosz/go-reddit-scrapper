FROM golang:latest AS builder

WORKDIR /build

COPY . .

RUN go build -o app

ENV forum=
ENV interval=

ENTRYPOINT [ "./app" ]