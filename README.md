# Distributed Config Management

Monorepo for a distributed configuration system with three services:
- `controller`: source of truth for global config and agent registration
- `agent`: polls controller, caches state, and pushes config to worker
- `worker`: executes `/hit` using latest config received from agent

## Repository Structure

- `controller/`
- `agent/`
- `worker/`
- `shared/` (shared middleware and HTTP response helpers)

## Architecture Flow

1. Agent registers to controller via `POST /register`.
2. Controller returns `agent_id`, `poll_url`, and `poll_interval_seconds`.
3. Agent polls controller `GET /config` with ETag.
4. On config changes, agent forwards config to worker via `POST /config`.
5. Worker serves `GET /hit` by calling configured URL and proxying response body.

## Architecture Notes

- Controller is the global source of truth for configuration.
- Agent persists local state (`agent_id`, ETag, last version, cached config URL) in file storage.
- Worker keeps active config in memory and applies updates from agent.
- Agent uses exponential backoff + jitter for request failures (`controller` / `worker` targets), while local file errors use fixed retry interval.
- If controller is unavailable, agent and worker continue using existing cached/applied config.
- ETag is used to avoid resending unchanged config (`304 Not Modified` path).

## Service Docs

- [controller/README.md](controller/README.md)
- [agent/README.md](agent/README.md)
- [worker/README.md](worker/README.md)

## Monorepo Commands

From repo root:

```bash
make docker                # up agent + worker
make docker-down           # down agent + worker
make docker-controller     # up controller
make docker-controller-down
make docker-all            # up all services
make docker-all-down       # down all services
```

## Compile and Setup (Local)

### Prerequisites

- Go `1.22.x`
- GNU Make
- Docker + Docker Compose (optional, for containerized run)

### 1) Prepare Environment Files

Create env files from examples:

- `controller/.env` from `controller/.env.example`
- `agent/.env` from `agent/.env.example`
- `worker/.env` from `worker/.env.example`

Ensure shared secrets are aligned:

- `controller.AGENT_API_KEY == agent.CONTROLLER_API_KEY`
- `worker.AGENT_API_KEY == agent.WORKER_API_KEY`

### 2) Compile

```bash
cd controller && go build ./cmd/main.go
cd ../agent && go build ./cmd/main.go
cd ../worker && go build ./cmd/main.go
```

### 3) Run (without Docker)

Use separate terminals:

```bash
cd controller && make run
cd agent && make run
cd worker && make run
```

### 4) Run Tests

```bash
cd controller && make test
cd ../agent && make test
cd ../worker && make test
```

## Docker Compose Layout

- Controller compose: `controller/docker-compose.yml`
- Agent + Worker compose: `docker-compose.agent-worker.yml`

## Deployment Details (Public Access)

Fill this section with your real deployed URLs before submission:

- Controller public URL: `<https://...>`
- Agent public URL: `<https://...>`
- Worker public URL: `<https://...>`

Example endpoint checks:

- `POST <controller-url>/register`
- `GET <controller-url>/config`
- `POST <worker-url>/config`
- `GET <worker-url>/hit`
- `GET <agent-url>/state`
- `GET <worker-url>/state`

## Notes

- Controller uses PostgreSQL (`DATABASE_URL`) for persistent configuration storage.
- Keep API keys synchronized:
  - `controller.AGENT_API_KEY == agent.CONTROLLER_API_KEY`
  - `worker.AGENT_API_KEY == agent.WORKER_API_KEY`
- If using Render free tier for demo/testing, SQLite persistence is not guaranteed (no persistent disk on free instances).
  Re-seed configuration (`POST /config`) after service restart/redeploy before running end-to-end checks.


