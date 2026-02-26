DC_AGENT_WORKER := docker compose -p agent-worker --env-file agent/.env --env-file worker/.env -f docker-compose.agent-worker.yml
DC_CONTROLLER := docker compose -p controller --env-file controller/.env -f controller/docker-compose.yml

.PHONY: docker docker-down docker-controller docker-controller-down docker-all docker-all-down

docker:
	$(DC_AGENT_WORKER) up --build -d

docker-down:
	$(DC_AGENT_WORKER) down

docker-controller:
	$(DC_CONTROLLER) up --build -d

docker-controller-down:
	$(DC_CONTROLLER) down

docker-all:
	$(DC_CONTROLLER) up --build -d
	$(DC_AGENT_WORKER) up --build -d

docker-all-down:
	$(DC_AGENT_WORKER) down
	$(DC_CONTROLLER) down
