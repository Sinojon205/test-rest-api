FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./
COPY main.go server.go ./

RUN go mod download

COPY pkg ./pkg
COPY config ./config

RUN go build -o /app/test-rest-api ./cmd/

EXPOSE 8080

CMD ["/app/test-rest-api"]
