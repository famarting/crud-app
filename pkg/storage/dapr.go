package storage

import (
	"context"
	"encoding/json"
	"fmt"

	dapr "github.com/dapr/go-sdk/client"
	"github.com/famarting/crud-app/pkg/todos"
	"github.com/google/uuid"
)

type DaprStorage struct {
	client   dapr.Client
	maxItems int
}

var statestoreName string = "statestore"
var pubsubName string = "pubsub"
var pubsubTopic string = "todos"
var indexKey string = "index"

var daprImpl TodosStorage = &DaprStorage{}

func (s *DaprStorage) Create(todo *todos.Todo) error {
	todo.Id = uuid.New().String()
	bytes, err := json.Marshal(todo)
	if err != nil {
		return err
	}

	allTodos, err := s.getIndexState()
	if err != nil {
		return err
	}

	if len(allTodos) >= s.maxItems {
		allTodos = allTodos[1:]
	}

	allTodos = append(allTodos, todo.Id)
	atb, err := json.Marshal(allTodos)
	if err != nil {
		return err
	}

	err = s.newDaprClient().SaveBulkState(context.Background(), statestoreName,
		&dapr.SetStateItem{Key: todo.Id, Value: bytes},
		&dapr.SetStateItem{Key: indexKey, Value: atb})
	if err != nil {
		return err
	}

	err = s.publishEvent(todo)
	if err != nil {
		return err
	}

	return nil
}

func (s *DaprStorage) getIndexState() ([]string, error) {
	index, err := s.newDaprClient().GetState(context.TODO(), statestoreName, indexKey)
	if err != nil {
		return nil, fmt.Errorf("get index state from dapr: %w", err)
	}
	if index == nil || index.Value == nil {
		return make([]string, 0), nil
	}
	var a []string = make([]string, 0)
	err = json.Unmarshal(index.Value, &a)
	if err != nil {
		return nil, fmt.Errorf("index json deserialization: %w", err)
	}
	return a, nil
}

func (s *DaprStorage) Update(todo *todos.Todo) error {
	bytes, err := json.Marshal(todo)
	if err != nil {
		return err
	}

	err = s.newDaprClient().SaveState(context.Background(), statestoreName, todo.Id, bytes)
	if err != nil {
		return err
	}

	err = s.publishEvent(todo)
	if err != nil {
		return err
	}

	return nil
}

func (s *DaprStorage) Delete(todo *todos.Todo) error {
	allTodos, err := s.getIndexState()
	if err != nil {
		return err
	}

	var newtodos []string = make([]string, 0)
	for _, t := range allTodos {
		if t != todo.Id {
			newtodos = append(newtodos, t)
		}
	}

	atb, err := json.Marshal(newtodos)
	if err != nil {
		return err
	}

	fmt.Printf("deleting todo: %v\n", todo.Id)
	err = s.newDaprClient().DeleteState(context.Background(), statestoreName, todo.Id)
	if err != nil {
		return err
	}
	err = s.newDaprClient().SaveState(context.Background(), statestoreName, indexKey, atb)
	if err != nil {
		return err
	}

	err = s.publishEvent(todo)
	if err != nil {
		return err
	}

	return nil
}

func (s *DaprStorage) publishEvent(t *todos.Todo) error {
	return s.newDaprClient().PublishEvent(context.Background(), pubsubName, pubsubTopic, t, dapr.PublishEventWithMetadata(map[string]string{"rawPayload": "true"}))
}

func (s *DaprStorage) ListAll() ([]*todos.Todo, error) {

	keys, err := s.getIndexState()
	if err != nil {
		return nil, err
	}

	if len(keys) == 0 {
		return make([]*todos.Todo, 0), nil
	}

	items, err := s.newDaprClient().GetBulkState(context.Background(), statestoreName, keys, nil, 1)
	if err != nil {
		return nil, err
	}
	var all []*todos.Todo = []*todos.Todo{}

	for _, item := range items {
		var t todos.Todo = todos.Todo{}
		json.Unmarshal(item.Value, &t)
		all = append(all, &t)
	}
	return all, nil
}

func (s *DaprStorage) newDaprClient() dapr.Client {
	if s.client == nil {
		client, err := dapr.NewClient()
		if err != nil {
			panic(err)
		}
		s.client = client
	}
	return s.client
}

func NewDaprStorage(maxItems int) *DaprStorage {
	return &DaprStorage{maxItems: maxItems}
}
