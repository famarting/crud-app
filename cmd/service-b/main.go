package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/dapr/go-sdk/service/common"
	daprd "github.com/dapr/go-sdk/service/http"

	dapr "github.com/dapr/go-sdk/client"
)

var pubsubName string = "pubsub"
var pubsubTopic string = "events"

var client dapr.Client

func main() {
	fmt.Println("starting service B app (receives service invocation and publishes to events topic)")

	s := daprd.NewService(":8080")

	if err := s.AddServiceInvocationHandler("/hello", echoHandler); err != nil {
		log.Fatalf("error adding invocation handler: %v", err)
	}

	go publishEvent(context.TODO())

	if err := s.Start(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("error: %v", err)
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

func echoHandler(ctx context.Context, in *common.InvocationEvent) (*common.Content, error) {
	log.Printf("echo - ContentType:%s, Verb:%s, QueryString:%s, %+v \n", in.ContentType, in.Verb, in.QueryString, string(in.Data))

	r := rand.Intn(100)
	if r >= 45 {
		log.Println("randomly sleeping")
		// randomly sleep for more than 1 second, so the resiliency policy kicks in
		time.Sleep(1500 * time.Millisecond)
		log.Println("returning response")
	}

	// do something with the invocation here
	return &common.Content{
		Data:        in.Data,
		ContentType: in.ContentType,
		DataTypeURL: in.DataTypeURL,
	}, nil
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
