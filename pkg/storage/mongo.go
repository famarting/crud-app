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
	coll     *mgm.Collection
	maxItems int
}

var mongoImpl TodosStorage = &MongoStorage{}

func (s *MongoStorage) Create(todo *todos.Todo) error {
	todo.Id = uuid.New().String()
	count, err := s.coll.CountDocuments(context.Background(), bson.D{}, nil)
	if err != nil {
		return err
	}
	if count >= int64(s.maxItems) {
		// set the sort to -1 to sort descending and get the last item for deletion
		deleteOptions := &options.FindOneAndDeleteOptions{}
		deleteOptions = deleteOptions.SetSort(bson.D{{"_id", -1}})
		s.coll.FindOneAndDelete(context.Background(), bson.D{}, deleteOptions)
	}
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

func NewMongoStorage(connStr string, maxItems int) *MongoStorage {
	err := mgm.SetDefaultConfig(nil, "mgm_lab", options.Client().ApplyURI(connStr))
	if err != nil {
		panic(err)
	}

	return &MongoStorage{
		coll:     mgm.Coll(&todos.Todo{}),
		maxItems: maxItems,
	}
}
