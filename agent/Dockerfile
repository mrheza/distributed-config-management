FROM golang:1.22-alpine

WORKDIR /app

COPY . .

RUN go mod tidy
RUN go build -o agent ./cmd/main.go

EXPOSE 8081

CMD ["./agent"]
