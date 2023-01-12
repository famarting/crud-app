package main

import (
	"context"
	"fmt"
	"time"

	dapr "github.com/dapr/go-sdk/client"
)

var statestoreName string = "statestore"
var client dapr.Client

func main() {
	fmt.Println("starting datagen app")
	go generateReads(context.TODO())

	// leaving this one in the main goroutine to keep the process alive :)
	generateWrites(context.TODO())
}

func generateReads(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Second)

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
