package server

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	s "github.com/famarting/crud-app/pkg/storage"
	"github.com/famarting/crud-app/pkg/todos"
	"github.com/gin-gonic/gin"
)

type Server struct {
	Port                   int
	Storage                s.TodosStorage
	CleanupIntervalSeconds int
}

func (server *Server) Start() {

	go generateLoad(context.TODO(), server)
	go cleanup(context.TODO(), server)

	engine := gin.Default()

	group := engine.Group("/api/v1")

	todosGroup := group.Group("/todos")

	todosGroup.GET("", func(c *gin.Context) {
		todos, err := server.Storage.ListAll()
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, todos)
	})

	todosGroup.POST("", func(c *gin.Context) {
		var json todos.Todo
		if err := c.ShouldBindJSON(&json); err != nil {
			fmt.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		json.CreatedAt = time.Now()
		json.UpdatedAt = time.Now()
		err := server.Storage.Create(&json)
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, json)
	})

	todosGroup.PUT("", func(c *gin.Context) {
		var json todos.Todo
		if err := c.ShouldBindJSON(&json); err != nil {
			fmt.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		json.UpdatedAt = time.Now()
		err := server.Storage.Update(&json)
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, json)
	})

	todosGroup.DELETE("", func(c *gin.Context) {
		var json todos.Todo
		if err := c.ShouldBindJSON(&json); err != nil {
			fmt.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		json.UpdatedAt = time.Now()
		json.Deleted = "true"
		err := server.Storage.Delete(&json)
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, json)
	})

	engine.Run("0.0.0.0:" + strconv.Itoa(server.Port))

}

func generateLoad(ctx context.Context, server *Server) {
	ticker := time.NewTicker(30 * time.Second)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			fmt.Println("generating load")
			to := &todos.Todo{
				Text: "foo",
			}
			to.CreatedAt = time.Now()
			to.UpdatedAt = time.Now()
			err := server.Storage.Create(to)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

func cleanup(ctx context.Context, server *Server) {
	ticker := time.NewTicker(time.Duration(server.CleanupIntervalSeconds) * time.Second)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			fmt.Println("cleaning data")

			all, err := server.Storage.ListAll()
			if err != nil {
				fmt.Println(err)
				continue
			}
			for _, item := range all {
				err = server.Storage.Delete(item)
				if err != nil {
					fmt.Println(err)
				}
			}
		}
	}
}
