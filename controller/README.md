# FULL Production Controller

Features:

- Gin
- Clean Architecture
- Swagger auto generated (swag)
- Agent register API
- ETag support
- API key auth
- CORS middleware
- Docker + docker compose

Generate swagger:

swag init -g cmd/main.go

Run docker:

docker-compose up --build
