.PHONY: proto build run test clean migrate-up migrate-down docker-up docker-down

# Generate protobuf code
proto:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/*.proto

# Build the server
build:
	go build -o bin/server cmd/server/main.go

# Run the server
run:
	go run cmd/server/main.go

# Kill process on port 50051
kill-port:
	@lsof -ti:50051 | xargs kill -9 2>/dev/null || echo "No process on port 50051"

# Run tests
test:
	go test ./... -v

# Clean build artifacts
clean:
	rm -rf bin/

# Run database migrations up
migrate-up:
	go run cmd/migrate/main.go -direction=up

# Run database migrations down
migrate-down:
	go run cmd/migrate/main.go -direction=down

# Start Docker services
docker-up:
	docker-compose up -d

# Stop Docker services
docker-down:
	docker-compose down

# Install dependencies
deps:
	go mod download
	go mod tidy

# Format code
fmt:
	go fmt ./...

# Run the full development setup
dev: docker-up migrate-up run
