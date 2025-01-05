package main

import (
	"log"
	"net"
	"os"

	adapters "github.com/bernardbaker/todo.list.microservice/adapters/grpc"
	"github.com/bernardbaker/todo.list.microservice/adapters/memory"
	"github.com/bernardbaker/todo.list.microservice/application"
	"github.com/bernardbaker/todo.list.microservice/domain"
	pb "github.com/bernardbaker/todo.list.microservice/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// Get port number from environment variable
	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
		log.Printf("defaulting to port %s", port)
	}
	// Initialize Adapters
	repo := memory.NewInMemoryTodoRepository()
	service := application.NewTodoService(repo)
	handler := adapters.NewGRPCTodoHandler(service)

	// Start gRPC Server
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	server := grpc.NewServer(
		grpc.UnaryInterceptor(domain.HmacInterceptor),
	)
	pb.RegisterTodoServiceServer(server, handler)
	reflection.Register(server)

	log.Println("Starting gRPC server on port " + port)
	if err := server.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
