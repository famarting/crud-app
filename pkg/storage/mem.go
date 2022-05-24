package storage

import (
	"github.com/famarting/crud-app/pkg/todos"
	"github.com/google/uuid"
)

type InMemoryStorage struct {
	all []*todos.Todo
}

var impl TodosStorage = &InMemoryStorage{}

func (s *InMemoryStorage) Create(todo *todos.Todo) error {
	todo.Id = uuid.New().String()
	s.all = append(s.all, todo)
	return nil
}

func (s *InMemoryStorage) Update(todo *todos.Todo) error {
	return nil
}

func (s *InMemoryStorage) Delete(todo *todos.Todo) error {
	return nil
}

func (s *InMemoryStorage) ListAll() ([]*todos.Todo, error) {
	return s.all, nil
}

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		all: []*todos.Todo{},
	}
}
