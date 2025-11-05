FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY internal ./internal
COPY cmd ./cmd
COPY config ./config
COPY ui ./ui

RUN go build -o /app/tech-e-market ./cmd/

EXPOSE 8080

CMD ["/app/tech-e-market"]