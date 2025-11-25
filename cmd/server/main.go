package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	userApp "github.com/nguyenphuoc/super-salary-sacrifice/internal/features/user-management/application"
	userEntrypoints "github.com/nguyenphuoc/super-salary-sacrifice/internal/features/user-management/entrypoints"
	userInfra "github.com/nguyenphuoc/super-salary-sacrifice/internal/features/user-management/infrastructure"
	"github.com/nguyenphuoc/super-salary-sacrifice/internal/shared/infrastructure"
	userProto "github.com/nguyenphuoc/super-salary-sacrifice/proto"
)

func main() {
	log.Println("🚀 Starting Super & Salary Sacrifice CRUD Server...")

	// Load configuration from environment
	config, err := infrastructure.LoadConfig()
	if err != nil {
		log.Fatalf("❌ Failed to load configuration: %v", err)
	}

	log.Printf("📝 Environment: %s", config.App.Environment)
	log.Printf("📝 Log Level: %s", config.App.LogLevel)

	// Initialize database connection (Infrastructure Adapter)
	db, err := infrastructure.NewDatabase(config.Database)
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}

	// Initialize repositories (Infrastructure Layer - Adapters)
	userRepo := userInfra.NewGormUserRepository(db)

	// Initialize application services (Application Layer - Use Cases)
	userService := userApp.NewUserService(userRepo)

	// Create gRPC server
	grpcServer := grpc.NewServer()

	// Register gRPC services (Entrypoints - Inbound Adapters)
	userProto.RegisterUserServiceServer(grpcServer, userEntrypoints.NewUserServiceServer(userService))

	// Enable gRPC reflection for development (allows grpcurl to work)
	reflection.Register(grpcServer)

	// Start TCP listener
	listener, err := net.Listen("tcp", config.Server.Address())
	if err != nil {
		log.Fatalf("❌ Failed to listen on %s: %v", config.Server.Address(), err)
	}

	// Handle graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		sig := <-sigChan
		log.Printf("⚠️  Received signal: %v. Shutting down gracefully...", sig)
		grpcServer.GracefulStop()
		log.Println("✅ Server stopped")
	}()

	// Start serving
	log.Printf("✅ gRPC server listening on %s", config.Server.Address())
	log.Println("💡 Use grpcurl to test: grpcurl -plaintext localhost:50051 list")
	
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("❌ Failed to serve: %v", err)
	}
}
