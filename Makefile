service_name = "ps"

.PHONY: up
up: 
	docker-compose up -d --build

.PHONY: restart
restart:
	docker-compose restart

.PHONY: down
down:
	docker-compose down

.PHONY: ps
ps:
	docker-compose ps

.PHONY: logs
logs:
	docker-compose logs --tail 100 -f $(service_name)

.PHONY: lint
lint:
	golangci-lint run --fix --allow-parallel-runners -v

.PHONY: re
restart: down up