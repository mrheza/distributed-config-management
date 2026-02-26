# Worker Service

Worker receives config from agent and executes `/hit` based on latest config.

## Endpoints

- `POST /config` (agent auth required)
- `GET /hit`
- `GET /state`
- `GET /swagger/*any`

## Authentication

Header: `X-API-Key`
- Required for `POST /config`
- Key is configured by `AGENT_API_KEY`

## Environment Variables

| Variable | Default | Description |
|---|---|---|
| `REQUEST_TIMEOUT_SECONDS` | `10` | Timeout when worker fetches configured URL |
| `AGENT_API_KEY` | `worker-secret` | API key for `POST /config` |
| `GIN_MODE` | `release` | Gin mode |
| `PORT` | `8082` | HTTP port |

## Local Run

```bash
cd worker
make run
```

## Test

```bash
cd worker
make test
make coverage
```

## Swagger

```bash
cd worker
make swagger
```

Swagger route: `http://localhost:8082/swagger/index.html`

## Docker

Worker uses shared compose with agent:

```bash
cd worker
make docker
make docker-down
```

Compose file: `../docker-compose.agent-worker.yml`
