FROM golang:latest

WORKDIR /app

COPY . .
RUN go mod download

RUN apt-get update && apt-get install -y \
  vim