DOCKER_OPTIONS := --env-file=goapi/.env
API_CONTAINER := api

ifeq (exec, $(firstword $(MAKECMDGOALS)))
  EXEC_ARGS := $(wordlist 2, $(words $(MAKECMDGOALS)), $(MAKECMDGOALS))
  $(eval $(EXEC_ARGS):;@true)
endif

ifeq (logs, $(firstword $(MAKECMDGOALS)))
  LOGS_ARGS := $(wordlist 2, $(words $(MAKECMDGOALS)), $(MAKECMDGOALS))
  $(eval $(EXEC_ARGS):;@true)
endif

build:
	@docker compose $(DOCKER_OPTIONS) build --no-cache

up:
	@docker compose $(DOCKER_OPTIONS) up -d

down:
	@docker compose $(DOCKER_OPTIONS) down

start:
	@docker compose $(DOCKER_OPTIONS) start

stop:
	@docker compose $(DOCKER_OPTIONS) stop

restart:
	@docker compose $(DOCKER_OPTIONS) restart $(API_CONTAINER)

ps:
	@docker compose $(DOCKER_OPTIONS) ps

logs:
	@docker compose $(DOCKER_OPTIONS) logs $(LOGS_ARGS)

exec:
	@docker compose $(DOCKER_OPTIONS) exec $(EXEC_ARGS)
