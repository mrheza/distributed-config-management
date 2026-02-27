# Worker Service

## Overview
Worker receives config from agent and executes `/hit` using the latest applied URL.

Public URL: `https://worker-4ldb.onrender.com`

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
| Variable | Required | Description |
|---|---|---|
| `REQUEST_TIMEOUT_SECONDS` | Yes | Timeout when worker calls configured URL |
| `AGENT_API_KEY` | Yes | API key for `POST /config` |
| `GIN_MODE` | Yes | Gin mode (`debug`/`release`) |
| `PORT` | Yes | HTTP port |

## Local Development
### Run
```bash
cd worker
make run
```

### Generate Mocks
```bash
cd worker
make mocks
```

### Test
```bash
cd worker
make test
make coverage
```

### Generate Swagger
```bash
cd worker
make swagger
```

Swagger:
- Local: `http://localhost:8082/swagger/index.html`
- Public: `https://worker-4ldb.onrender.com/swagger/index.html`

## Docker
Worker uses shared compose with agent:

```bash
cd worker
make docker
make docker-down
```

Compose file: `../docker-compose.agent-worker.yml`

## Notes
- Config is stored in memory (reapplied by agent after startup if available).
- Keep key aligned: `AGENT_API_KEY == agent.WORKER_API_KEY`.
