package storage

import "github.com/famarting/crud-app/pkg/todos"

type TodosStorage interface {
	Create(todo *todos.Todo) error
	Update(todo *todos.Todo) error
	Delete(todo *todos.Todo) error
	ListAll() ([]*todos.Todo, error)
}
