package storage

import (
	"context"

	"github.com/famarting/crud-app/pkg/todos"
	"github.com/google/uuid"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func init() {
	// Setup the mgm default config
	err := mgm.SetDefaultConfig(nil, "mgm_lab", options.Client().ApplyURI("mongodb://root:12345@localhost:27017"))
	if err != nil {
		panic(err)
	}
}

type MongoStorage struct {
	coll *mgm.Collection
}

var mongoImpl TodosStorage = &MongoStorage{}

func (s *MongoStorage) Create(todo *todos.Todo) error {
	todo.Id = uuid.New().String()
	return s.coll.Create(todo)
}

func (s *MongoStorage) Update(todo *todos.Todo) error {
	return nil
}

func (s *MongoStorage) Delete(todo *todos.Todo) error {
	return nil
}

func (s *MongoStorage) ListAll() ([]*todos.Todo, error) {
	cursor, err := s.coll.Find(context.Background(), bson.D{})
	if err != nil {
		return nil, err
	}
	var all []*todos.Todo = []*todos.Todo{}
	err = cursor.All(context.Background(), &all)
	if err != nil {
		return nil, err
	}
	return all, nil
}

func NewMongoStorage(connStr string) *MongoStorage {
	err := mgm.SetDefaultConfig(nil, "mgm_lab", options.Client().ApplyURI(connStr))
	if err != nil {
		panic(err)
	}

	return &MongoStorage{
		coll: mgm.Coll(&todos.Todo{}),
	}
}
