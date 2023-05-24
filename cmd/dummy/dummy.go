package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("starting dummy app (accepts requests at the /events path)")

	engine := gin.Default()

	group := engine.Group("/events")

	group.GET("", func(c *gin.Context) {
		c.JSON(200, map[string]interface{}{
			"result": true,
		})
	})

	group.POST("", func(c *gin.Context) {
		c.JSON(200, map[string]interface{}{
			"result": true,
		})
	})

	engine.Run("0.0.0.0:8080")
}
