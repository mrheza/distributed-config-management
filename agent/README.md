# Agent Service

Agent bridges controller and worker:
- registers to controller
- polls latest config with ETag
- forwards config updates to worker
- persists local state for resilience

## Endpoints

- `GET /state`
- `GET /swagger/*any`

## Environment Variables

| Variable | Default | Description |
|---|---|---|
| `CONTROLLER_BASE_URL` | `http://localhost:8080` | Controller base URL |
| `CONTROLLER_API_KEY` | `agent-secret` | API key for controller agent endpoints |
| `WORKER_BASE_URL` | `http://localhost:8082` | Worker base URL |
| `WORKER_API_KEY` | `worker-secret` | API key sent to worker `POST /config` |
| `POLL_URL` | `/config` | Poll path on controller |
| `POLL_INTERVAL_SECONDS` | `30` | Initial poll interval |
| `STATE_PATH` | `data/agent_state.json` | Local persisted state file |
| `MAX_BACKOFF_SECONDS` | `60` | Max backoff limit |
| `BACKOFF_JITTER_PERCENT` | `20` | Backoff jitter percent |
| `REQUEST_TIMEOUT_SECONDS` | `10` | Outbound HTTP timeout |
| `GIN_MODE` | `release` | Gin mode |
| `PORT` | `8081` | HTTP port |

## Local Run

```bash
cd agent
make run
```

## Test

```bash
cd agent
make test
make coverage
```

## Swagger

```bash
cd agent
make swagger
```

Swagger route: `http://localhost:8081/swagger/index.html`

## Docker

Agent uses shared compose with worker:

```bash
cd agent
make docker
make docker-down
```

Compose file: `../docker-compose.agent-worker.yml`
