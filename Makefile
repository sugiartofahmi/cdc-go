.PHONY: up down clean register unregister status topics os-count logs

COMPOSE := docker compose
CONNECTOR := ecommerce-postgres-source

up:
	$(COMPOSE) up -d

down:
	$(COMPOSE) down

clean:
	$(COMPOSE) down -v

register:
	curl -sS -X POST -H "Content-Type: application/json" \
		--data @debezium/products-connector.json \
		http://localhost:8083/connectors && echo ""

unregister:
	curl -sS -X DELETE http://localhost:8083/connectors/$(CONNECTOR) && echo "deleted"

status:
	@curl -sS http://localhost:8083/connectors/$(CONNECTOR)/status && echo ""

topics:
	@docker exec cdc-kafka /opt/kafka/bin/kafka-topics.sh --bootstrap-server localhost:9092 --list 2>/dev/null

os-count:
	@curl -sS http://localhost:9200/products-index/_count 2>/dev/null || echo "index not found"

logs:
	$(COMPOSE) logs -f --tail=100
