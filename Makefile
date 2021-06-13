## permanent variables
PROJECT			?= github.com/gpenaud/needys-api-need
RELEASE			?= $(shell git describe --tags --abbrev=0)
COMMIT			?= $(shell git rev-parse --short HEAD)
BUILD_TIME  ?= $(shell date -u '+%Y-%m-%d_%H:%M:%S')

## docker environment options
DOCKER_BUILD_ARGS ?= --build-arg PROJECT="${PROJECT}" --build-arg RELEASE="${RELEASE}" --build-arg COMMIT="${COMMIT}" --build-arg BUILD_TIME="${BUILD_TIME}"

## Colors
COLOR_RESET       = $(shell tput sgr0)
COLOR_ERROR       = $(shell tput setaf 1)
COLOR_COMMENT     = $(shell tput setaf 3)
COLOR_TITLE_BLOCK = $(shell tput setab 4)

## display this help text
help:
	@printf "\n"
	@printf "${COLOR_TITLE_BLOCK}${PROJECT} Makefile${COLOR_RESET}\n"
	@printf "\n"
	@printf "${COLOR_COMMENT}Usage:${COLOR_RESET}\n"
	@printf " make build\n\n"
	@printf "${COLOR_COMMENT}Available targets:${COLOR_RESET}\n"
	@awk '/^[a-zA-Z\-_0-9@]+:/ { \
				helpLine = match(lastLine, /^## (.*)/); \
				helpCommand = substr($$1, 0, index($$1, ":")); \
				helpMessage = substr(lastLine, RSTART + 3, RLENGTH); \
				printf " ${COLOR_INFO}%-15s${COLOR_RESET} %s\n", helpCommand, helpMessage; \
		} \
		{ lastLine = $$0 }' $(MAKEFILE_LIST)
	@printf "\n"

## stack - start the entire stack in background, then follow logs
start:
	docker-compose up --build --detach
	docker-compose logs --follow needys-api-need

## stack - stop the entire stack
stop:
	docker-compose down

## stack - only start the api "needys-api-need"
api-only:
	docker-compose up needys-api-need

## stack - only start the sidecars backends (rabbitmq here)
sidecars-only:
	docker-compose up mariadb rabbitmq

## docker - build the needys-api-need image
.PHONY: build
build:
	docker build ${DOCKER_BUILD_ARGS} --file Dockerfile --tag needys-api-need:latest .

## docker - enter into the needys-api-need container
enter:
	docker-compose exec needys-api-need /bin/sh

## test - display all "need" table entries
test-list:
	curl -X GET http://localhost:8010

## test - remove, then insert a need entry named "testing-need" in need table
test-all:
	curl -X DELETE http://localhost:8010?name=testing-need
	curl -d "name=testing-need&priority=high" -X POST http://localhost:8010

## test - remove all entries filtered by name "testing-need" from need table
test-delete:
	curl -X DELETE http://localhost:8010?name=testing-need

## test - insert a need entry named "testing-need" in need table
test-insert:
	curl -d "name=testing-need&priority=high" -X POST http://localhost:8010
