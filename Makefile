NO_COLOR=\033[0m
OK_COLOR=\033[32;01m
ERROR_COLOR=\033[31;01m
WARN_COLOR=\033[33;01m

.PHONY: dependencies build unit-test integration-test docker-it-up docker-down clear docker-rmi docker-rmv

dependencies:
	@GO111MODULE=on go mod download
	@GO111MODULE=on go mod verify

build: dependencies
	@echo "$(OK_COLOR)==> Building... $(NO_COLOR)"
	@docker build . -t local/checkout-service

unit-tests:
	@GO111MODULE=on go test -v -short ./...

integration-tests: docker-it-up
	@./waitForContainer.sh
	@echo "$(OK_COLOR)==> Running ITs$(NO_COLOR)"
	@GO111MODULE=on go test -v ./internal/tests/integration/...; docker-compose -f ./internal/tests/docker-compose-it.yml down

docker-it-up:
	@docker-compose -f ./internal/tests/docker-compose-it.yml up -d

docker-it-down:
	@docker-compose -f ./internal/tests/docker-compose-it.yml down

docker-up:
	@docker-compose up -d
	@echo "$(WARN_COLOR)==> Waiting for service to be ready$(NO_COLOR)"
	@chmod +x waitForContainer.sh
	@./waitForContainer.sh

docker-down:
	@docker-compose down

clear: docker-down

docker-rmi:
	@docker rmi $$(docker images -f "dangling=true" -q) -f

docker-rmv:
	@docker volume rm $$(docker volume ls -q -f dangling=true)

