package server

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	s "github.com/famartinrh/crud-app/pkg/storage"
	"github.com/famartinrh/crud-app/pkg/todos"
	"github.com/gin-gonic/gin"
)

type Server struct {
	Port    int
	Storage s.TodosStorage
}

func (server *Server) Start() {

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
