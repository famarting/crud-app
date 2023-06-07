package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/famarting/crud-app/pkg/server"
	"github.com/famarting/crud-app/pkg/storage"
	"github.com/peterbourgon/ff/v3/ffcli"
)

func main() {

	serveFlagSet := flag.NewFlagSet("app serve", flag.ExitOnError)
	var serverPort *int = serveFlagSet.Int("port", 8080, "port for the server to listen to")
	var connStr *string = serveFlagSet.String("connStr", "mongodb://localhost:27017", "connection string for storage")
	var maxItems *int = serveFlagSet.Int("maxItems", 20, "maximum numbers of items to store")
	var cleanupIntervalSeconds *int = serveFlagSet.Int("cleanupIntervalSeconds", 600, "seconds to wait between cleanup runs")

	serve := &ffcli.Command{
		Name:     "serve",
		LongHelp: "Runs the http server for this app",
		FlagSet:  serveFlagSet,
		Exec: func(ctx context.Context, args []string) error {

			var s storage.TodosStorage

			if *connStr == "" || *connStr == "mem" {
				fmt.Println("Using inmemory storage")
				s = storage.NewInMemoryStorage(*maxItems)
			} else if *connStr == "dapr" {
				s = storage.NewDaprStorage(*maxItems)
			} else {
				s = storage.NewMongoStorage(*connStr, *maxItems)
			}

			server := server.Server{
				Port:                   *serverPort,
				Storage:                s,
				CleanupIntervalSeconds: *cleanupIntervalSeconds,
			}
			server.Start()
			return nil
		},
	}

	root := &ffcli.Command{
		Name:        "app",
		LongHelp:    "management cli for basic crud-app",
		Subcommands: []*ffcli.Command{serve},
	}

	root.ParseAndRun(context.Background(), os.Args[1:])

}
