PKGS := $(shell go run ./scripts/pkgfilter -sep=space)
DC := docker compose

mocks:
	go run scripts/clean.go
	mockery --dir=internal/service --output=internal/mocks/service --outpkg=service_mocks --all --with-expecter
	mockery --dir=internal/repository --output=internal/mocks/repository --outpkg=repository_mocks --all --with-expecter

run:
	go run ./cmd/main.go

swagger:
	go run github.com/swaggo/swag/cmd/swag@v1.16.2 init -g main.go -d cmd,internal/handler,internal/model,internal/httpresponse -o docs

test:
	go test $(PKGS) -count=1

coverage:
	go test $(PKGS) -cover -count=1

docker:
	$(DC) up --build -d

docker-down:
	$(DC) down
	
