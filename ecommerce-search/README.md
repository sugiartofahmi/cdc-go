# go-service

A Go REST API boilerplate/starter: layered architecture (controller → service → repository), JWT auth, role-based access control, database migrations, and a JSON-driven seeder, ready to use as the foundation for a new project.

## Features

- **JWT authentication**: register & login, passwords are bcrypt-hashed automatically.
- **Role-based access control**: authorization middleware + role checks (`Admin`, `Super Admin`, `Developer`) ready to use.
- **User & Role management**: full CRUD with pagination.
- **Layered architecture**: a consistent controller/service/repository pattern (Query + Store split) to copy for new modules.
- **Database migrations**: up/down SQL migrations via `golang-migrate`, or auto-migrate from entity structs for fast development.
- **Seeder**: JSON-driven initial data (roles & users), easy to extend.

## Tech Stack

- [Go](https://go.dev/) 1.25
- [Gin](https://github.com/gin-gonic/gin): HTTP web framework
- [GORM](https://gorm.io/) + PostgreSQL: ORM & database
- [golang-migrate](https://github.com/golang-migrate/migrate): database migrations
- [golang-jwt](https://github.com/golang-jwt/jwt): JWT auth
- [go-redis](https://github.com/redis/go-redis): caching
- [go-playground/validator](https://github.com/go-playground/validator): request validation

## Prerequisites

- Go 1.25+
- PostgreSQL
- Redis

## Getting Started

1. **Copy the environment file and fill in the config:**

   ```bash
   cp .env-example .env
   ```

   | Group | Variable | Description |
   |---|---|---|
   | Application | `APP_NAME`, `APP_ENV`, `APP_PORT`, `APP_GIN_MODE` | Basic app config (default port `8080`) |
   | Database | `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`, `DB_TIMEZONE` | PostgreSQL connection |
   | Database Pool | `DB_MAX_IDLE_CONNS`, `DB_MAX_OPEN_CONNS`, `DB_MAX_IDLE_CONNS_IN_MINUTES`, `DB_MAX_LIFETIME_CONNS_IN_MINUTES` | Connection pool settings |
   | JWT | `JWT_SECRET`, `JWT_EXPIRES_IN` | Signing key and token lifetime |
   | Redis | `REDIS_HOST`, `REDIS_PORT`, `REDIS_PASS` | Redis connection |

2. **Install dependencies:**

   ```bash
   go mod download
   ```

3. **Run migrations:**

   ```bash
   go run . --migration=true --exec=up
   ```

   Other available variants:
   - `--exec=down`: roll back one migration
   - `--exec=fresh`: roll back everything, then migrate from scratch
   - `--exec=create --fileName=<name>`: create a new migration file

   For fast development, you can sync the schema directly from entity structs instead of migration files:

   ```bash
   go run . --automigrate=true
   ```

4. **Seed initial data** (roles + one example account per role):

   ```bash
   go run . --seed=true                                 # run all seeders
   go run . --seed=true --target=RoleSeeder,UserSeeder   # run specific seeders
   ```

   Example accounts seeded from `seeder/files/users.json`, use these for your first login:

   | Email | Role |
   |---|---|
   | thomas.anderson@mail.com | Super Admin |
   | sarah.chen@mail.com | Admin |
   | marcus.rivera@mail.com | Developer |

   See `seeder/files/users.json` for each account's password. **Replace or remove these example accounts before using this in an environment anyone else can reach.**

5. **Run the server:**

   ```bash
   go run .
   ```

   The server listens on the port set by `APP_PORT` (default `8080`).

   Alternatively, build a binary first and run that instead:

   ```bash
   go build -o go-service .
   ./go-service
   ```

## API Endpoints

| Method | Path | Auth | Description |
|---|---|---|---|
| GET | `/health` | Public | Health check |
| POST | `/api/v1/auth/register` | Public | Register a new account (auto-assigned the `Developer` role) |
| POST | `/api/v1/auth/login` | Public | Log in, receive a JWT token |
| GET | `/api/v1/roles` | Admin/Super Admin | List roles (paginated) |
| GET | `/api/v1/roles/:id` | Admin/Super Admin | Get a single role |
| POST | `/api/v1/roles` | Admin/Super Admin | Create a role |
| PUT | `/api/v1/roles/:id` | Admin/Super Admin | Update a role |
| DELETE | `/api/v1/roles/:id` | Admin/Super Admin | Delete a role |
| GET | `/api/v1/users` | Admin/Super Admin | List users (paginated) |
| GET | `/api/v1/users/:id` | Admin/Super Admin | Get a single user |
| POST | `/api/v1/users` | Admin/Super Admin | Create a user |
| PUT | `/api/v1/users/:id` | Admin/Super Admin | Update a user |
| DELETE | `/api/v1/users/:id` | Admin/Super Admin | Delete a user |

Endpoints that require auth expect an `Authorization: Bearer <token>` header (the token comes from the login response).

### Login example

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"thomas.anderson@mail.com","password":"<see seeder/files/users.json>"}'
```

## Testing

```bash
go test ./...
```

## Project Structure

```
domain/           Business logic per module (auth, role, user): constants, dtos, interfaces, repositories, services
entities/         Entity structs mapped to database tables (GORM)
infrastructure/   Technical integrations: config, database, redis, middleware, exceptions, singleton, utils, validators
migration/        SQL migration files and runner (golang-migrate)
presentation/     HTTP layer: controllers and routing (Gin)
seeder/           JSON-driven database seeders
```

## Extending This Boilerplate

Add a new module by following the pattern in `domain/auth`, `domain/role`, or `domain/user`. In short: create `domain/<module>/{constants,dtos,interfaces,repositories,services}`, split repositories into Query (read) and Store (write), then manually register the new repository/service/controller in `main.go`.
