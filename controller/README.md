# Controller Service

## Overview
Controller is the source of truth for global configuration and agent registration.

Public URL: `https://controller-8hwn.onrender.com`

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
| Variable | Required | Description |
|---|---|---|
| `ADMIN_API_KEY` | Yes | API key for `POST /config` |
| `AGENT_API_KEY` | Yes | API key for `POST /register` and `GET /config` |
| `POLL_URL` | Yes | Poll path returned to agents |
| `GIN_MODE` | Yes | Gin mode (`debug`/`release`) |
| `DATABASE_URL` | Yes | PostgreSQL connection string |
| `PORT` | Yes | HTTP port |

## Local Development
### Run
```bash
cd controller
make run
```

### Generate Mocks
```bash
cd controller
make mocks
```

### Test
```bash
cd controller
make test
make coverage
```

### Generate Swagger
```bash
cd controller
make swagger
```

Swagger:
- Local: `http://localhost:8080/swagger/index.html`
- Public: `https://controller-8hwn.onrender.com/swagger/index.html`

## Docker
```bash
cd controller
make docker
make docker-down
```

Compose file: `controller/docker-compose.yml`

## Notes
- Persistence uses PostgreSQL via `DATABASE_URL`.
- Ensure `AGENT_API_KEY` is aligned with `agent` service (`CONTROLLER_API_KEY`).
