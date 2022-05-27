FROM golang:1.18.2-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
RUN go install github.com/cosmtrek/air@latest

COPY . .

EXPOSE 8081

RUN go build -o relation-service

ENTRYPOINT ./relation-service
