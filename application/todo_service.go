package application

import (
	"github.com/bernardbaker/todo.list.microservice/domain"
)

// Application Service: Coordinates operations between domain and ports
type TodoService struct {
	repo domain.TodoRepository // Port to storage
}

// NewTodoService creates a new instance of TodoService
func NewTodoService(repo domain.TodoRepository) *TodoService {
	return &TodoService{repo: repo}
}

func (s *TodoService) AddTodo(item *domain.TodoItem) error {
	return s.repo.AddTodo(item)
}

func (s *TodoService) ListTodos() ([]*domain.TodoItem, error) {
	return s.repo.ListTodos()
}

func (s *TodoService) MarkCompleted(id string) error {
	return s.repo.MarkCompleted(id)
}

func (s *TodoService) DeleteTodo(id string) error {
	return s.repo.DeleteTodo(id)
}
