package domain

import "errors"

type TodoItem struct {
	ID        string
	Title     string
	Completed bool
}

var ErrTodoNotFound = errors.New("todo not found")
