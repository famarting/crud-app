package main

import (
	"context"
	"fmt"
	"time"

	dapr "github.com/dapr/go-sdk/client"
	"github.com/dapr/go-sdk/service/common"
	"github.com/dapr/go-sdk/service/grpc"
)

var statestoreName string = "statestore"
var pubsubName string = "pubsub"
var pubsubTopic string = "events"

var client dapr.Client

func main() {
	fmt.Println("starting event consumer app (consumes from events topic and uses statestore)")

	s, err := grpc.NewService(":8080")
	if err != nil {
		panic(err)
	}
	s.AddTopicEventHandler(&common.Subscription{
		PubsubName: pubsubName,
		Topic:      pubsubTopic,
	}, func(ctx context.Context, e *common.TopicEvent) (retry bool, err error) {
		fmt.Printf("event consumed %s %s\n", e.DataContentType, e.ID)
		return false, nil
	})

	go generateReads(context.TODO())
	go generateWrites(context.TODO())

	err = s.Start()
	if err != nil {
		panic(err)
	}
}

func generateReads(ctx context.Context) {
	ticker := time.NewTicker(10 * time.Second)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			fmt.Println("generating reads")
			_, err := getDaprClient().GetState(ctx, statestoreName, "foo")
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

func generateWrites(ctx context.Context) {
	ticker := time.NewTicker(100 * time.Second)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			fmt.Println("generating writes")
			err := getDaprClient().SaveState(ctx, statestoreName, "foo", []byte("baz"))
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

func getDaprClient() dapr.Client {
	if client == nil {
		c, err := dapr.NewClient()
		if err != nil {
			panic(err)
		}
		client = c
	}
	return client
}
