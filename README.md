# Feature Flag API Service

A scalable feature flag management system with REST and gRPC APIs. This project is designed for extensibility, observability, and production-readiness, supporting PostgreSQL and etcd as backends.

---

## Architecture Overview


```plaintext
                    +----------------+
                    |     CLI Tool   |
                    +-------+--------+
                            |
                +-----------+-----------+
                |   REST/gRPC/WebSocket|
                +-----------+-----------+
                            |
         +------------------+-------------------+
         |       API Gateway & Auth Service     |
         +------------------+-------------------+
                            |
                 +----------+-----------+
                 |    Backend Services  |
                 | (Flag, Config, Audit|
                 |  Rule Engine, RBAC) |
                 +----------+-----------+
                            |
                 +----------+-----------+
                 |  Distributed Storage |
                 | (etcd/Consul/Redis/ |
                 |   PostgreSQL)       |
                 +---------------------+
```

### Core Components

- **REST API Server**  
  Exposes OpenAPI-compliant HTTP endpoints for managing feature flags, configurations, audits, and RBAC controls.

- **gRPC API Server**  
  Provides a performant and type-safe interface for client integration.

- **Service Layer**  
  - `flag`: Business logic for feature flag CRUD and evaluation.
  - `config`: Manages dynamic configuration.
  - `audit`: Stores and queries audit logs for change tracking.
  - `rbac`: Role-based access control for secure administration.

- **Persistence Layer**  
  - **etcd**: Stores flag state for low-latency access and high availability.
  - **PostgreSQL**: Stores configuration, audit logs, and RBAC data.

- **Shared Logger**  
  Centralized logging utility for all components.

---

## Project Structure

```plaintext
/
├── cmd/
│   └── api/                  # Entrypoint for the API service
├── internal/
│   ├── flag/                 # Feature flag logic and interface
│   ├── config/               # Dynamic configuration logic
│   ├── audit/                # Auditing logic
│   ├── rbac/                 # RBAC logic
│   ├── logger/               # Shared logger package
│   └── server/               # REST and gRPC server wiring
├── api/grpc/v1/              # Protobuf (gRPC) definitions
├── Dockerfile                # API service Dockerfile
├── docker-compose.yaml       # Development/test stack
├── compose/                  # Optional: production and override compose files
├── deploy/                   # Infrastructure-as-code, k8s manifests, etc.
├── README.md                 # This file
```

---

## Local Development

### Prerequisites

- [Go 1.24+](https://golang.org)
- [Docker](https://www.docker.com/)
- [docker-compose](https://docs.docker.com/compose/)
- [protoc](https://grpc.io/docs/protoc-installation/) and `protoc-gen-go-grpc`

### Running with Docker Compose

```sh
docker-compose up --build
```

- Feature Flag API: http://localhost:8080
- etcd: localhost:2379
- Postgres: localhost:5432

Set environment variables in a `.env` file or in your shell:

```env
DB_URL=postgres://featureuser:featurepass@postgres:5432/featureflags?sslmode=disable
JWT_SECRET=your_jwt_secret_here
```

---

## API Entrypoint

The main entrypoint is in [`cmd/api/main.go`](cmd/api/main.go):

- Loads configuration from environment
- Initializes services (`flag`, `config`, `audit`, `rbac`)
- Starts both REST and gRPC servers
- Handles graceful shutdown on SIGINT/SIGTERM

---

## Server Package Example

The [`internal/server`](internal/server) package wires up both REST (using [gorilla/mux](https://github.com/gorilla/mux)) and gRPC (using [pgx](https://github.com/jackc/pgx) and generated protobuf code):

- `StartREST`: Sets up HTTP routes and handlers
- `RegisterGRPC`: Registers gRPC services with the gRPC server

---

## Persistence

- **Feature Flags:** Stored in etcd for distributed consistency and fast reads.
- **Configs, Audits, RBAC:** Stored in PostgreSQL using the [pgx](https://github.com/jackc/pgx) driver for performance and reliability.

---

## Logging

A shared logger is provided in [`internal/logger`](internal/logger/logger.go). Use `logger.Get()` for consistent, thread-safe logging across the codebase.

---

## Extensibility

- Add new services by implementing the appropriate interface and wiring them into the server package.
- Swap out the persistence layer by providing alternative implementations of the service interfaces.

---

## Security

- JWT authentication and RBAC for admin endpoints.
- Sensitive configuration via environment variables.

