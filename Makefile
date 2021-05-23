ONESHELL:

## HELP
PROJECT           = needys-api-need
## Colors
COLOR_RESET       = $(shell tput sgr0)
COLOR_ERROR       = $(shell tput setaf 1)
COLOR_COMMENT     = $(shell tput setaf 3)
COLOR_TITLE_BLOCK = $(shell tput setab 4)

## Display this help text
help:
	@printf "\n"
	@printf "${COLOR_TITLE_BLOCK}${PROJECT} Makefile${COLOR_RESET}\n"
	@printf "\n"
	@printf "${COLOR_COMMENT}Usage:${COLOR_RESET}\n"
	@printf " make build\n\n"
	@printf "${COLOR_COMMENT}Available targets:${COLOR_RESET}\n"
	@awk '/^[a-zA-Z\-\_0-9\@]+:/ { \
				helpLine = match(lastLine, /^## (.*)/); \
				helpCommand = substr($$1, 0, index($$1, ":")); \
				helpMessage = substr(lastLine, RSTART + 3, RLENGTH); \
				printf " ${COLOR_INFO}%-20s${COLOR_RESET} %s\n", helpCommand, helpMessage; \
		} \
		{ lastLine = $$0 }' $(MAKEFILE_LIST)
	@printf "\n"

start:
	docker-compose up --build --detach
	docker-compose logs --follow

stop:
	docker-compose down

start-services-only:
	docker-compose up mariadb rabbitmq

test-list:
	curl -X GET http://localhost:8010

test-all:
	curl -X DELETE http://localhost:8010?name=testing-need
	curl -d "name=testing-need&priority=high" -X POST http://localhost:8010

test-delete:
	curl -X DELETE http://localhost:8010?name=testing-need

test-insert:
	curl -d "name=testing-need&priority=high" -X POST http://localhost:8010
