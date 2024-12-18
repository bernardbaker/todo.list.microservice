package domain

type TodoRepository interface {
	AddTodo(item *TodoItem) error
	ListTodos() ([]*TodoItem, error)
	MarkCompleted(id string) error
	DeleteTodo(id string) error
}

type TodoService interface {
	AddTodo(item *TodoItem) error
	ListTodos() ([]*TodoItem, error)
	MarkCompleted(id string) error
	DeleteTodo(id string) error
}
