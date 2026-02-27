# Distributed Config Management

Distributed configuration system in a single monorepo with 3 services:
- `controller`: source of truth for config + agent registration
- `agent`: registers to controller, polls config (ETag), forwards updates to worker
- `worker`: executes `/hit` using the latest config from agent

## Service Docs
- [controller/README.md](controller/README.md)
- [agent/README.md](agent/README.md)
- [worker/README.md](worker/README.md)

## Repository Structure
- `controller/`
- `agent/`
- `worker/`
- `shared/`

## End-to-End Flow
1. Agent calls `POST /register` to controller.
2. Controller returns `agent_id`, `poll_url`, and `poll_interval_seconds`.
3. Agent polls `GET /config` with `If-None-Match`.
4. If config changes, agent pushes config to worker via `POST /config`.
5. User calls worker `GET /hit`; worker requests configured URL and returns raw body.

## Prerequisites
- Go `1.22.x`
- GNU Make
- Docker + Docker Compose (optional)

Install `make` (Windows/macOS):

macOS:
```bash
brew install make
```

Windows (Chocolatey):
```powershell
choco install make
```

## Quick Start (Local)
### 1) Prepare env files
Create:
- `controller/.env` from `controller/.env.example`
- `agent/.env` from `agent/.env.example`
- `worker/.env` from `worker/.env.example`

Required key alignment:
- `controller.AGENT_API_KEY == agent.CONTROLLER_API_KEY`
- `worker.AGENT_API_KEY == agent.WORKER_API_KEY`

### 2) Build binaries
```bash
cd controller && go build ./cmd/main.go
cd ../agent && go build ./cmd/main.go
cd ../worker && go build ./cmd/main.go
```

### 3) Run services (3 terminals)
```bash
cd controller && make run
cd agent && make run
cd worker && make run
```

### 4) Run tests
```bash
cd controller && make test
cd ../agent && make test
cd ../worker && make test
```

## Make Commands
### Root level
```bash
make docker                # up agent + worker
make docker-down           # down agent + worker
make docker-controller     # up controller
make docker-controller-down
make docker-all            # up all services
make docker-all-down       # down all services
```

### Per-service utility commands
```bash
cd controller && make mocks && make swagger && make test && make coverage
cd agent && make mocks && make swagger && make test && make coverage
cd worker && make mocks && make swagger && make test && make coverage
```

## Docker Layout
- Controller compose: `controller/docker-compose.yml`
- Agent + Worker compose: `docker-compose.agent-worker.yml`

## CI/CD (GitHub Actions -> Render)
Tag-based deploy per service:
- `controller/v*.*.*` -> controller deploy
- `agent/v*.*.*` -> agent deploy
- `worker/v*.*.*` -> worker deploy

Workflow files:
- `.github/workflows/deploy-controller.yml`
- `.github/workflows/deploy-agent.yml`
- `.github/workflows/deploy-worker.yml`

Required GitHub secrets:
- `RENDER_DEPLOY_HOOK_CONTROLLER`
- `RENDER_DEPLOY_HOOK_AGENT`
- `RENDER_DEPLOY_HOOK_WORKER`

Trigger deploy:
```bash
git tag controller/v1.0.0
git push origin controller/v1.0.0

git tag agent/v1.0.0
git push origin agent/v1.0.0

git tag worker/v1.0.0
git push origin worker/v1.0.0
```

CI/CD notes:
- Disable Render auto-deploy-on-commit for strict tag-only deployment.
- Use repository root as Render build context for `agent` and `worker` (both depend on `shared/`).

## Public Deployment
- Controller: `https://controller-8hwn.onrender.com`
- Agent: `https://agent-awcy.onrender.com`
- Worker: `https://worker-4ldb.onrender.com`

Swagger:
- Controller: `https://controller-8hwn.onrender.com/swagger/index.html`
- Agent: `https://agent-awcy.onrender.com/swagger/index.html`
- Worker: `https://worker-4ldb.onrender.com/swagger/index.html`

## Verification Checklist (Public URLs)
- `POST https://controller-8hwn.onrender.com/register`
- `GET https://controller-8hwn.onrender.com/config`
- `POST https://worker-4ldb.onrender.com/config`
- `GET https://worker-4ldb.onrender.com/hit`
- `GET https://agent-awcy.onrender.com/state`
- `GET https://worker-4ldb.onrender.com/state`

## Important Notes
- Controller persistence is PostgreSQL (`DATABASE_URL`).
- Agent stores local runtime state in file; worker stores active config in memory.
- On Render free tier, instances may sleep/restart; runtime state can reset, but controller data remains in PostgreSQL.
