package memory

import (
	"errors"

	"github.com/bernardbaker/todo.list.microservice/domain"
)

type InMemoryTodoRepository struct {
	todos map[string]*domain.TodoItem
}

func NewInMemoryTodoRepository() *InMemoryTodoRepository {
	return &InMemoryTodoRepository{todos: make(map[string]*domain.TodoItem)}
}

func (r *InMemoryTodoRepository) AddTodo(item *domain.TodoItem) error {
	r.todos[item.ID] = item
	return nil
}

func (r *InMemoryTodoRepository) ListTodos() ([]*domain.TodoItem, error) {
	var items []*domain.TodoItem
	for _, v := range r.todos {
		items = append(items, v)
	}
	return items, nil
}

func (r *InMemoryTodoRepository) MarkCompleted(id string) error {
	if item, exists := r.todos[id]; exists {
		item.Completed = true
		return nil
	}
	return errors.New("todo not found")
}

func (r *InMemoryTodoRepository) DeleteTodo(id string) error {
	if _, exists := r.todos[id]; exists {
		delete(r.todos, id)
		return nil
	}
	return errors.New("todo not found")
}
