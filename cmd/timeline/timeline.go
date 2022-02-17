package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/famartinrh/crud-app/pkg/timeline"
	"github.com/famartinrh/crud-app/pkg/todos"
	"github.com/gin-gonic/gin"

	cloudevents "github.com/cloudevents/sdk-go/v2"
)

func main() {

	serveFlagSet := flag.NewFlagSet("timeline app serve", flag.ExitOnError)
	var serverPort *int = serveFlagSet.Int("port", 8080, "port for the server to listen to")

	var tl timeline.Timeline = timeline.New()

	engine := gin.Default()

	group := engine.Group("/")

	eventsGroup := group.Group("todos")

	eventsGroup.GET("", func(c *gin.Context) {
		c.JSON(200, tl.Timeline())
	})

	eventsGroup.POST("", func(c *gin.Context) {

		var bodyBytes []byte
		if c.Request.Body != nil {
			bodyBytes, _ = ioutil.ReadAll(c.Request.Body)
		}

		fmt.Println("Received event", string(bodyBytes))

		event := cloudevents.NewEvent()

		var todo todos.Todo

		err := json.Unmarshal(bodyBytes, &event)
		if err != nil {
			// is it a rawPayload or a binary cloudevent
			err = json.Unmarshal(bodyBytes, &todo)
			if err != nil {
				fmt.Println(err)
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
		} else {
			// it's an structured cloudevent

			err = json.Unmarshal(event.Data(), &todo)
			if err != nil {
				fmt.Println(err)
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
		}

		tl.Handle(todo)

		c.Writer.WriteHeader(200)
	})

	fmt.Println("Starting timeline server on port %n", *serverPort)
	engine.Run("0.0.0.0:" + strconv.Itoa(*serverPort))
}
