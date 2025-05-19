FROM golang:1.24-alpine

RUN mkdir /app
WORKDIR /app

COPY go.* *.go ./

RUN go mod tidy
RUN go mod verify
RUN go mod download

COPY . .

RUN go build -o /app/main .
EXPOSE 8080
CMD [ "/app/main" ]
