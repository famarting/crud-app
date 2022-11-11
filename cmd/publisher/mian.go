package main

import (
	"context"
	"fmt"
	"time"

	dapr "github.com/dapr/go-sdk/client"
	"github.com/dapr/go-sdk/service/grpc"
)

var pubsubName string = "pubsub"
var pubsubTopic string = "events"

var client dapr.Client

func main() {
	fmt.Println("starting event publisher app")

	s, err := grpc.NewService(":8080")
	if err != nil {
		panic(err)
	}

	go publishEvent(context.TODO())

	err = s.Start()
	if err != nil {
		panic(err)
	}
}

func publishEvent(ctx context.Context) {
	ticker := time.NewTicker(10 * time.Second)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			fmt.Println("generating event")

			err := getDaprClient().PublishEvent(ctx, pubsubName, pubsubTopic, "data")
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
