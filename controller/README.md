# Controller Service

Controller is the source of truth for configuration and agent registration.

## Endpoints

- `POST /register` (agent auth required)
- `GET /config` (agent auth required, ETag support)
- `POST /config` (admin auth required)
- `GET /swagger/*any`

## Authentication

Header: `X-API-Key`
- Agent routes use `AGENT_API_KEY`
- Admin route uses `ADMIN_API_KEY`

## Environment Variables

| Variable | Default | Description |
|---|---|---|
| `ADMIN_API_KEY` | `admin-secret` | API key for `POST /config` |
| `AGENT_API_KEY` | `agent-secret` | API key for `POST /register`, `GET /config` |
| `POLL_URL` | `/config` | Poll path returned to agents |
| `GIN_MODE` | `release` | Gin mode |
| `DATABASE_URL` | `postgresql://postgres:postgres@localhost:5432/postgres?sslmode=disable` | PostgreSQL connection string |
| `PORT` | `8080` | HTTP port |

## Local Run

```bash
cd controller
make run
```

## Test

```bash
cd controller
make test
make coverage
```

## Swagger

```bash
cd controller
make swagger
```

Swagger route: `http://localhost:8080/swagger/index.html`

## Docker

```bash
cd controller
make docker
make docker-down
```

Compose file: `controller/docker-compose.yml`

Note: controller now uses PostgreSQL via `DATABASE_URL` (not SQLite file path).
