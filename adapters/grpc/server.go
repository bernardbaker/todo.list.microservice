package grpc

import (
	"net"
	"os"

	"github.com/bernardbaker/todo.list.microservice/proto"

	"google.golang.org/grpc"
)

func StartGRPCServer(handler proto.TodoServiceServer) error {
	lis, err := net.Listen("tcp", ":"+os.Getenv("PORT"))
	if err != nil {
		return err
	}
	grpcServer := grpc.NewServer()
	proto.RegisterTodoServiceServer(grpcServer, handler)
	return grpcServer.Serve(lis)
}
