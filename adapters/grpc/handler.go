package grpc

import (
	"context"

	"github.com/bernardbaker/todo.list.microservice/domain"
	"github.com/bernardbaker/todo.list.microservice/proto"
)

type GRPCTodoHandler struct {
	service domain.TodoService
	proto.UnimplementedTodoServiceServer
}

func NewGRPCTodoHandler(service domain.TodoService) *GRPCTodoHandler {
	return &GRPCTodoHandler{service: service}
}

func (h *GRPCTodoHandler) AddTodo(ctx context.Context, req *proto.TodoItem) (*proto.TodoResponse, error) {
	item := &domain.TodoItem{ID: req.Id, Title: req.Title, Completed: req.Completed}
	err := h.service.AddTodo(item)
	if err != nil {
		return nil, err
	}
	return &proto.TodoResponse{Message: 200}, nil
}

func (h *GRPCTodoHandler) ListTodos(ctx context.Context, req *proto.Empty) (*proto.TodoList, error) {
	items, _ := h.service.ListTodos()
	var protoItems []*proto.TodoItem
	for _, i := range items {
		protoItems = append(protoItems, &proto.TodoItem{
			Id:        i.ID,
			Title:     i.Title,
			Completed: i.Completed,
		})
	}
	return &proto.TodoList{Todos: protoItems}, nil
}

// Implement MarkCompleted and DeleteTodo similarly.
func (h *GRPCTodoHandler) MarkCompleted(ctx context.Context, req *proto.TodoRequest) (*proto.TodoResponse, error) {
	err := h.service.MarkCompleted(req.Id)
	if err != nil {
		return nil, err
	}
	return &proto.TodoResponse{Message: 200}, nil
}

func (h *GRPCTodoHandler) DeleteTodo(ctx context.Context, req *proto.TodoDeleteRequest) (*proto.TodoResponse, error) {
	err := h.service.DeleteTodo(req.Id)
	if err != nil {
		return nil, err
	}
	return &proto.TodoResponse{Message: 200}, nil
}
