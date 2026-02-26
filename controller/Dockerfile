FROM golang:1.22-alpine

WORKDIR /app

COPY shared ./shared
COPY controller ./controller

WORKDIR /app/controller
RUN go mod tidy

RUN go build -o controller ./cmd/main.go

EXPOSE 8080

CMD ["./controller"]
