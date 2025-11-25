# Super & Salary Sacrifice CRUD System

A production-grade CRUD system for managing Australian Superannuation funds and Salary Sacrifice arrangements, built with **Hexagonal Architecture** (Ports & Adapters).

## 🏗️ Architecture

This project follows **Hexagonal Architecture** with **Vertical Slice** organization:

```
internal/
├── features/                    # Vertical slices by feature
│   ├── user-management/
│   │   ├── domain/             # 🔵 Core business logic (entities, interfaces)
│   │   ├── application/        # 🟢 Use cases (services, commands, queries)
│   │   ├── infrastructure/     # 🟡 Adapters (GORM repositories)
│   │   └── entrypoints/        # 🟠 Inbound adapters (gRPC handlers)
│   ├── super-management/
│   └── sacrifice-management/
└── shared/
    ├── domain/                 # Shared value objects
    └── infrastructure/         # Shared adapters (DB, Config)
```

### Hexagonal Architecture Layers

- **Domain** (Core): Pure business logic, no dependencies
- **Application**: Use cases that orchestrate domain logic
- **Infrastructure**: Adapters for external systems (DB, APIs)
- **Entrypoints**: Inbound adapters (gRPC, REST, CLI)

## 🚀 Quick Start

### Prerequisites

- Go 1.21+
- Docker & Docker Compose
- `protoc` (Protocol Buffers compiler)

### Setup

```bash
# 1. Clone and install dependencies
go mod download

# 2. Copy environment file
cp .env.example .env

# 3. Start PostgreSQL
make docker-up

# 4. Run migrations
make migrate-up

# 5. Start the server
make run
```

## 📦 Tech Stack

- **Language**: Go 1.21+
- **API**: gRPC with Protocol Buffers
- **Database**: PostgreSQL 15+ with GORM
- **Architecture**: Hexagonal (Ports & Adapters)
- **Deployment**: Docker, Kubernetes (planned)

## 🎯 MVP Scope (5 hours)

### Entities
- ✅ **User** (Employee)
- ✅ **SuperAccount** (Superannuation Fund)
- ✅ **SalarySacrifice** (Salary Sacrifice Arrangement)

### Operations (per entity)
- ✅ Create
- ✅ Get by ID
- ✅ List (with pagination)

## 📝 Development

### Generate Protobuf Code

```bash
make proto
```

### Run Migrations

```bash
# Up
make migrate-up

# Down
make migrate-down
```

### Test with grpcurl

```bash
# List services
grpcurl -plaintext localhost:50051 list

# Create user
grpcurl -plaintext -d '{"email":"test@example.com","first_name":"John","last_name":"Doe"}' \
  localhost:50051 user.UserService/CreateUser
```

## 🗂️ Project Status

See [TASK_PLAN.md](TASK_PLAN.md) for detailed progress.

## 📚 Resources

- [Hexagonal Architecture](https://alistair.cockburn.us/hexagonal-architecture/)
- [gRPC Go Quickstart](https://grpc.io/docs/languages/go/quickstart/)
- [GORM Documentation](https://gorm.io/docs/)
