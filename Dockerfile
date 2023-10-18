ARG GO_VERSION=1.21

FROM golang:${GO_VERSION}-alpine AS builder

RUN apk update && apk add alpine-sdk git && rm -rf /var/cache/apk/*

RUN mkdir -p /usr/local/bin/app
WORKDIR /usr/local/bin/app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build -o ./app ./main.go

FROM alpine:latest

ENV GIN_MODE=release
ENV USER=user
ENV UID=1000


RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

RUN mkdir -p /usr/local/bin/app

WORKDIR /usr/local/bin/app

COPY --from=builder /usr/local/bin/app/app .

EXPOSE 8080

ENTRYPOINT ["./app"]