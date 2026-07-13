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
	curl -fsS -X POST -H "Content-Type: application/json" \
		--data @debezium/products-connector.json \
		http://localhost:8083/connectors && echo ""

unregister:
	curl -fsS -X DELETE http://localhost:8083/connectors/$(CONNECTOR) && echo "deleted"

status:
	@curl -fsS http://localhost:8083/connectors/$(CONNECTOR)/status

topics:
	@docker exec cdc-kafka /opt/kafka/bin/kafka-topics.sh --bootstrap-server localhost:9092 --list

os-count:
	@curl -fsS http://localhost:9200/products/_count

logs:
	$(COMPOSE) logs -f --tail=100
