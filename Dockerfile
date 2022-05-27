FROM golang:1.18.1-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
RUN go install github.com/cosmtrek/air@latest

COPY ./services/relation ./services/relation
COPY ./packages ./packages

EXPOSE 8081

WORKDIR /app/services/relation

RUN go build -o relation-service

ENTRYPOINT ./relation-service
