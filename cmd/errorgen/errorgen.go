package main

import (
	"context"
	"fmt"
	"time"

	dapr "github.com/dapr/go-sdk/client"
)

var client dapr.Client

func main() {
	fmt.Println("starting error generator app")
	go invokeBinding(context.TODO())
	go invokeMethodErrors(context.TODO())
	go invokeMethod(context.TODO())

	// leaving this one in the main goroutine to keep the process alive :)
	publishEvent(context.TODO())
}

func invokeBinding(ctx context.Context) {
	ticker := time.NewTicker(20 * time.Second)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			fmt.Println("generating invoke binding errors")

			_, err := getDaprClient().InvokeBinding(ctx, &dapr.InvokeBindingRequest{
				Name:      "dummy-binding",
				Operation: "dummy-operation",
				Data:      []byte("ooo"),
			})
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

func invokeMethodErrors(ctx context.Context) {
	ticker := time.NewTicker(10 * time.Second)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			fmt.Println("generating invoke method errors")

			_, err := getDaprClient().InvokeMethod(ctx, "dummy-app", "do-something", "get")
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

// TODO add a function to publish events to an unauthorized pubsub
func publishEvent(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Second)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			fmt.Println("generating publish event errors")

			err := getDaprClient().PublishEvent(ctx, "dummy-pubsub", "dummy-topic", []byte("data"))
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("publish event succeeded")
			}
		}
	}
}

func invokeMethod(ctx context.Context) {
	ticker := time.NewTicker(15 * time.Second)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			fmt.Println("invoking crud-app")

			data, err := getDaprClient().InvokeMethod(ctx, "crud-app", "/api/v1/todos", "GET")
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("invoke crud-app succeeded")
			fmt.Println(string(data))
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
