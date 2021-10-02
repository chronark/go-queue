build: export COMPOSE_DOCKER_CLI_BUILD=1
build:
	docker-compose build


up: down build
	docker-compose up

down:
	docker-compose down