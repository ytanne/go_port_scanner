service_name = "ps"

.PHONY: up
up: 
	docker-compose up -d --build

.PHONY: down
down:
	docker-compose down

.PHONY: ps
ps:
	docker-compose ps

.PHONY: logs
logs:
	docker-compose logs --tail 100 -f $(service_name)

.PHONY: restart
restart: down up