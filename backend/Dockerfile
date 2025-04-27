FROM golang:1.24-alpine

RUN mkdir /app
WORKDIR /app

COPY go.* *.go ./

RUN go mod tidy
RUN go mod verify
RUN go mod download

COPY . .

