package tests

import (
	"context"
	"testing"

	"github.com/bernardbaker/todo.list.microservice/adapters/grpc"
	"github.com/bernardbaker/todo.list.microservice/application"
	"github.com/bernardbaker/todo.list.microservice/domain"
	"github.com/bernardbaker/todo.list.microservice/proto"

	"github.com/stretchr/testify/assert"
)

type mockTodoRepository struct {
	todos map[string]*domain.TodoItem
}

func (m *mockTodoRepository) AddTodo(item *domain.TodoItem) error {
	m.todos[item.ID] = item
	return nil
}

func (m *mockTodoRepository) ListTodos() ([]*domain.TodoItem, error) {
	var items []*domain.TodoItem
	for _, v := range m.todos {
		items = append(items, v)
	}
	return items, nil
}

func (m *mockTodoRepository) MarkCompleted(id string) error {
	if item, exists := m.todos[id]; exists {
		item.Completed = true
		return nil
	}
	return domain.ErrTodoNotFound
}

func (m *mockTodoRepository) DeleteTodo(id string) error {
	if _, exists := m.todos[id]; exists {
		delete(m.todos, id)
		return nil
	}
	return domain.ErrTodoNotFound
}

func newMockRepository() *mockTodoRepository {
	return &mockTodoRepository{todos: make(map[string]*domain.TodoItem)}
}

func TestAddTodo(t *testing.T) {
	repo := newMockRepository()
	service := application.NewTodoService(repo)
	handler := grpc.NewGRPCTodoHandler(service)

	req := &proto.TodoItem{Id: "1", Title: "Test Todo", Completed: false}
	res, err := handler.AddTodo(context.Background(), req)

	assert.NoError(t, err)
	assert.Equal(t, "Todo added successfully", res.Message)
	assert.Equal(t, "Test Todo", repo.todos["1"].Title)
}

func TestListTodos(t *testing.T) {
	repo := newMockRepository()
	repo.AddTodo(&domain.TodoItem{ID: "1", Title: "First Todo", Completed: false})
	repo.AddTodo(&domain.TodoItem{ID: "2", Title: "Second Todo", Completed: true})

	service := application.NewTodoService(repo)
	handler := grpc.NewGRPCTodoHandler(service)

	req := &proto.Empty{}
	res, err := handler.ListTodos(context.Background(), req)

	assert.NoError(t, err)
	assert.Len(t, res.Todos, 2)
	assert.Equal(t, "First Todo", res.Todos[0].Title)
	assert.Equal(t, "Second Todo", res.Todos[1].Title)
}

func TestMarkCompleted(t *testing.T) {
	repo := newMockRepository()
	repo.AddTodo(&domain.TodoItem{ID: "1", Title: "Incomplete Todo", Completed: false})

	service := application.NewTodoService(repo)
	handler := grpc.NewGRPCTodoHandler(service)

	req := &proto.TodoRequest{Id: "1"}
	res, err := handler.MarkCompleted(context.Background(), req)

	assert.NoError(t, err)
	assert.Equal(t, "Todo marked as completed", res.Message)
	assert.True(t, repo.todos["1"].Completed)
}

func TestDeleteTodo_Success(t *testing.T) {
	repo := newMockRepository()
	repo.AddTodo(&domain.TodoItem{ID: "1", Title: "Todo to delete", Completed: false})

	service := application.NewTodoService(repo)
	handler := grpc.NewGRPCTodoHandler(service)

	req := &proto.TodoRequest{Id: "1"}
	res, err := handler.DeleteTodo(context.Background(), req)

	assert.NoError(t, err)
	assert.Equal(t, "Todo deleted successfully", res.Message)
	assert.NotContains(t, repo.todos, "1")
}

func TestDeleteTodo_NotFound(t *testing.T) {
	repo := newMockRepository()

	service := application.NewTodoService(repo)
	handler := grpc.NewGRPCTodoHandler(service)

	req := &proto.TodoRequest{Id: "non-existent-id"}
	_, err := handler.DeleteTodo(context.Background(), req)

	assert.Error(t, err)
	assert.Equal(t, domain.ErrTodoNotFound.Error(), err.Error())
}
