# Go HTTP Server

A learning project for building HTTP servers and REST APIs in Go using the standard library. The project explores HTTP routing, request handling, middleware, authentication, PostgreSQL integration, and API documentation.

## Requirements

- Go
- PostgreSQL
- [Goose](https://github.com/pressly/goose) for database migrations
- [SQLC](https://sqlc.dev/) for generating type-safe database code
- [Swaggo](https://github.com/swaggo/swag) for Swagger documentation

## Getting Started

### 1. Start PostgreSQL

From the project root:

```bash
./startPostgres.sh
```

### 2. Connect to PostgreSQL (optional)

```bash
psql -h localhost -p 5432 -U postgres
```

### 3. Run database migrations

```bash
export GOOSE_MIGRATION_DIR="$(pwd)/sql/schema"
goose postgres "postgres://postgres:postgres@localhost:5432/chirpy?sslmode=disable" up
```

### 4. Start the HTTP server

```bash
go run cmd/main.go
```

## API Documentation

The API documentation is generated using Swaggo annotations in the Go source code.

Regenerate the documentation after adding or changing handlers:

```bash
swag init -g internal/api/swagger.go --parseInternal
```

The generated files are written to `docs/`.

Open Swagger UI in your browser:

```text
http://localhost:<port>/swagger/index.html
```

Replace `<port>` with the port configured by the application.

## Available Endpoints

| Method | Endpoint              | Purpose                                     |
| ------ | --------------------- | ------------------------------------------- |
| GET    | `/api/healthz`        | Health check                                |
| POST   | `/api/users`          | Create a user                               |
| PUT    | `/api/users`          | Update a user                               |
| POST   | `/api/login`          | Authenticate a user                         |
| POST   | `/api/refresh`        | Refresh an access token                     |
| POST   | `/api/revoke`         | Revoke a refresh token                      |
| GET    | `/api/chirps`         | List records, optionally filtered by author |
| POST   | `/api/chirps`         | Create a record                             |
| GET    | `/api/chirps/{id}`    | Retrieve a record by ID                     |
| DELETE | `/api/chirps/{id}`    | Delete a record                             |
| POST   | `/api/polka/webhooks` | Handle an incoming webhook                  |
| GET    | `/admin/metrics`      | View HTTP server metrics                    |
| POST   | `/admin/reset`        | Reset development data and metrics          |

## Topics Explored

- HTTP servers and routing with Go's `net/http`
- HTTP methods, status codes, headers, and JSON
- Middleware and request logging
- Authentication with JWTs and refresh tokens
- Request validation and error handling
- PostgreSQL integration and SQLC-generated queries
- Database migrations with Goose
- API documentation with OpenAPI and Swagger UI

## Development Workflow

When adding a new endpoint:

1. Implement the handler and register its route.
2. Document the handler using Swaggo annotations.
3. Test successful requests and error cases.
4. Regenerate the Swagger documentation.
5. Verify the endpoint in Swagger UI.
