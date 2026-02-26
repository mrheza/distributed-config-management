FROM golang:1.22-alpine

WORKDIR /app

COPY shared ./shared
COPY worker ./worker

WORKDIR /app/worker
RUN go mod tidy
RUN go build -o worker ./cmd/main.go

EXPOSE 8082

CMD ["./worker"]
