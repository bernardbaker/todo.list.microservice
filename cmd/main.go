package main

import (
	"log"
	"os"

	"github.com/bernardbaker/todo.list.microservice/adapters/grpc"
	"github.com/bernardbaker/todo.list.microservice/adapters/memory"
	"github.com/bernardbaker/todo.list.microservice/application"
)

func main() {
	// Initialize Adapters
	repo := memory.NewInMemoryTodoRepository()
	service := application.NewTodoService(repo)
	handler := grpc.NewGRPCTodoHandler(service)

	// Start gRPC Server
	log.Println("Starting gRPC server on :" + os.Getenv("PORT"))
	if err := grpc.StartGRPCServer(&handler); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
