package main

import (
	"context"
	"fmt"

	"github.com/dapr/go-sdk/service/common"
	"github.com/dapr/go-sdk/service/grpc"
)

var pubsubName string = "pubsub"
var pubsubTopic string = "events"

func main() {
	fmt.Println("starting service-c app (consumes from events topic)")

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

	err = s.Start()
	if err != nil {
		panic(err)
	}
}
