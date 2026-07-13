# CDC-Go

Change Data Capture pipeline: PostgreSQL → Debezium → Kafka → OpenSearch.

Two Go services:
- **ecommerce** (port 8080): REST API CRUD products & categories → PostgreSQL
- **ecommerce-search** (port 8081): OpenSearch search API + Kafka consumer → OpenSearch

## Architecture

```
┌──────────────┐     ┌──────────────┐     ┌──────────────┐     ┌──────────┐     ┌──────────────┐     ┌────────────┐
│  ecommerce   │────▶│  PostgreSQL  │────▶│   Debezium   │────▶│  Kafka   │────▶│ecommerce-    │────▶│ OpenSearch │
│  (port 8080) │     │ (wal=logical)│     │  (connect)   │     │ (broker) │     │search        │     │  :9200     │
│  REST API    │     │              │     │   :8083       │     │  :9092   │     │(port 8081)   │     │products-   │
│  CRUD        │     │              │     │              │     │          │     │consumer      │     │index       │
└──────────────┘     └──────────────┘     └──────────────┘     └──────────┘     │search API    │     └────────────┘
                                                                                └──────────────┘
```

## How CDC Works

### PostgreSQL WAL (Write-Ahead Log)

PostgreSQL mencatat setiap perubahan data (INSERT/UPDATE/DELETE) ke WAL sebelum ditulis ke tabel.
Ada 3 level WAL:

| Level   | Output                        | Debezium Support |
|---------|-------------------------------|------------------|
| minimal | Crash recovery only           | No               |
| replica | Block-level changes           | No               |
| logical | Row-level changes with data   | **Yes**          |

Hanya `wal_level=logical` yang bisa dibaca Debezium — karena ngirim row-level events lengkap
seperti `{op: "u", id: "xxx", name: "iPhone 15"}`.

### How Debezium Listens to PostgreSQL

```
┌──────────────┐     (WAL stream)     ┌──────────────┐     (Kafka messages)     ┌──────────┐
│  PostgreSQL  │ ───────────────────▶ │   Debezium   │ ───────────────────────▶ │  Kafka   │
│              │                      │   Connect    │                          │  Broker  │
│ publication  │     consumer         │              │     producer             │          │
│   + slot     │                      │              │                          │          │
└──────────────┘                      └──────────────┘                          └──────────┘
```

Debezium acts as both:

- **Consumer** (from PostgreSQL): reads WAL stream via replication slot using the `pgoutput` plugin
- **Producer** (to Kafka): writes transformed events to Kafka topics

3 key components:

1. **Publication** — declares which tables to stream (`CREATE PUBLICATION FOR TABLE products, categories`)
2. **Replication slot** — bookmarks Debezium's WAL position, prevents data loss during disconnects
3. **pgoutput plugin** — decodes WAL binary into row-level JSON events

---

## Prerequisites

- Go 1.25+
- Docker & Docker Compose
- PostgreSQL (Docker or local install)

---

## Setup

### 1. PostgreSQL

#### 1.1 Enable wal_level=logical

```bash
# Find the PostgreSQL config file (Docker PostgreSQL 18 example)
docker exec postgres find / -name postgresql.conf 2>/dev/null
# → /var/lib/postgresql/18/docker/postgresql.conf

# Enable logical replication
docker exec postgres sed -i "s/#wal_level = replica/wal_level = logical/" /var/lib/postgresql/18/docker/postgresql.conf

# Restart PostgreSQL to apply changes
docker restart postgres
```

#### 1.2 Create publication

```bash
docker exec postgres psql -U postgres -d ecommerce -c "CREATE PUBLICATION dbz_publication FOR TABLE public.categories, public.products;"
```

#### 1.3 Re-register Debezium connector

PostgreSQL restart invalidates the replication slot. Re-register the connector:

```bash
# Remove existing connector
curl -sS -X DELETE http://localhost:8083/connectors/ecommerce-postgres-source

# Register connector
curl -sS -X POST -H "Content-Type: application/json" \
  --data @debezium/products-connector.json \
  http://localhost:8083/connectors

# Check status
make status   # or: curl -sS http://localhost:8083/connectors/ecommerce-postgres-source/status
```

#### 1.4 Verify

```sql
SHOW wal_level;              -- should return 'logical'
SELECT * FROM pg_publication; -- should list 'dbz_publication'
```

### 2. Infrastructure (Kafka + Debezium + OpenSearch)

```bash
# Start all services
docker compose up -d

# Register Debezium connector (if not already done in step 1.3)
curl -sS -X POST -H "Content-Type: application/json" \
  --data @debezium/products-connector.json \
  http://localhost:8083/connectors

# Verify everything is running
curl -sS http://localhost:8083/connectors/ecommerce-postgres-source/status
# → {"connector":{"state":"RUNNING"},"tasks":[{"state":"RUNNING"}],"type":"source"}

docker exec cdc-kafka /opt/kafka/bin/kafka-topics.sh --bootstrap-server localhost:9092 --list
# → ecommerce.public.products
# → ecommerce.public.categories
```

### 3. ecommerce Service (port 8080)

```bash
cd ecommerce
cp .env-example .env
go run main.go
```

### 4. ecommerce-search Service (port 8081)

```bash
cd ecommerce-search
cp .env-example .env
go run main.go
```

---

## API Examples

### Create a product

```bash
# Get an existing category ID first
curl -sS http://localhost:8080/api/v1/categories

# Create a product
curl -sS -X POST http://localhost:8080/api/v1/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "iPhone 15",
    "description": "Latest iPhone",
    "price": 15000000,
    "stock": 100,
    "category_id": "<CATEGORY_ID>"
  }'
```

### Delete a product

```bash
curl -sS -X DELETE http://localhost:8080/api/v1/products/<PRODUCT_ID>
```

Deletion uses GORM soft delete (`deleted_at` is set). Debezium sends it as an update (`op=u`),
and the ecommerce-search consumer detects `deleted_at != nil` to remove from OpenSearch.

### Search products (via OpenSearch)

```bash
# Full-text search
curl -sS "http://localhost:8081/api/v1/products?search=iphone"

# Filter by category and price
curl -sS "http://localhost:8081/api/v1/products?category=smartphones&min_price=10000000&max_price=20000000"

# Pagination
curl -sS "http://localhost:8081/api/v1/products?page=1&per_page=20&sort_by=price&order=asc"
```

### Verify OpenSearch sync

```bash
# Check document count
curl -sS http://localhost:9200/products-index/_count

# Clean all documents (if needed)
curl -X POST http://localhost:9200/products-index/_delete_by_query \
  -H "Content-Type: application/json" \
  -d '{"query":{"match_all":{}}}'
```

---

## Makefile Commands

The Makefile wraps common operations. If you don't have `make`, use the equivalent commands below.

| Command | Description | Equivalent |
|---------|-------------|------------|
| `make up` | Start Kafka, Debezium, OpenSearch | `docker compose up -d` |
| `make down` | Stop all infrastructure | `docker compose down` |
| `make clean` | Stop and remove volumes | `docker compose down -v` |
| `make register` | Register Debezium connector | `curl -sS -X POST -H "Content-Type: application/json" --data @debezium/products-connector.json http://localhost:8083/connectors` |
| `make unregister` | Remove Debezium connector | `curl -sS -X DELETE http://localhost:8083/connectors/ecommerce-postgres-source` |
| `make status` | Check connector status | `curl -sS http://localhost:8083/connectors/ecommerce-postgres-source/status` |
| `make topics` | List Kafka topics | `docker exec cdc-kafka /opt/kafka/bin/kafka-topics.sh --bootstrap-server localhost:9092 --list` |
| `make os-count` | Count OpenSearch documents | `curl -sS http://localhost:9200/products-index/_count` |
| `make logs` | Tail all infrastructure logs | `docker compose logs -f --tail=100` |

---

## Service Structure

```
cdc-go/
├── docker-compose.yml              # Kafka, Debezium, OpenSearch
├── Makefile                        # Infrastructure management commands
├── debezium/
│   └── products-connector.json     # Debezium connector configuration
├── ecommerce/                      # Write side (CRUD API → PostgreSQL)
│   ├── main.go
│   ├── entities/                   # ProductEntity, CategoryEntity
│   ├── domain/
│   │   ├── product/                # DTOs, interfaces, repositories, services
│   │   └── category/
│   ├── presentation/               # Controllers (Gin)
│   └── infrastructure/             # Postgres, Redis, middleware, utils
└── ecommerce-search/               # Read side (OpenSearch + Kafka consumer)
    ├── main.go
    ├── domain/product/
    │   ├── dtos/                   # ProductResultDto, ProductEventDto, query/request DTOs
    │   ├── interfaces/             # Query, store, event service interfaces
    │   ├── repositories/           # OpenSearch query + store repositories
    │   └── services/               # Search, event handling services
    ├── infrastructure/
    │   ├── config/                 # Env config (Redis, OpenSearch, Kafka)
    │   ├── kafka/                  # Kafka factory, consumer interface & service
    │   ├── redis/                  # Redis factory, cache interface & service
    │   ├── databases/              # OpenSearch connection
    │   ├── singleton/              # DI singletons
    │   └── utils/                  # Pagination, response helpers
    └── presentation/               # Controllers (health, product search)
```
