# Agent Service

## Overview
Agent bridges controller and worker:
- registers to controller
- polls config using ETag
- forwards new config to worker
- persists local runtime state for resilience

Public URL: `https://agent-awcy.onrender.com`

## Endpoints
- `GET /state`
- `GET /swagger/*any`

## Environment Variables
| Variable | Required | Description |
|---|---|---|
| `CONTROLLER_BASE_URL` | Yes | Controller base URL |
| `CONTROLLER_API_KEY` | Yes | API key for controller agent endpoints |
| `WORKER_BASE_URL` | Yes | Worker base URL |
| `WORKER_API_KEY` | Yes | API key sent to worker `POST /config` |
| `POLL_URL` | Yes | Poll path on controller |
| `POLL_INTERVAL_SECONDS` | Yes | Initial poll interval |
| `STATE_PATH` | Yes | Local state file path |
| `MAX_BACKOFF_SECONDS` | Yes | Max exponential backoff |
| `BACKOFF_JITTER_PERCENT` | Yes | Jitter percent for backoff |
| `REQUEST_TIMEOUT_SECONDS` | Yes | Outbound HTTP timeout |
| `GIN_MODE` | Yes | Gin mode (`debug`/`release`) |
| `PORT` | Yes | HTTP port |

## Local Development
### Run
```bash
cd agent
make run
```

### Generate Mocks
```bash
cd agent
make mocks
```

### Test
```bash
cd agent
make test
make coverage
```

### Generate Swagger
```bash
cd agent
make swagger
```

Swagger:
- Local: `http://localhost:8081/swagger/index.html`
- Public: `https://agent-awcy.onrender.com/swagger/index.html`

## Docker
Agent uses shared compose with worker:

```bash
cd agent
make docker
make docker-down
```

Compose file: `../docker-compose.agent-worker.yml`

## Notes
- Keep keys aligned:
  - `CONTROLLER_API_KEY == controller.AGENT_API_KEY`
  - `WORKER_API_KEY == worker.AGENT_API_KEY`
