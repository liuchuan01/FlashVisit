package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize task manager
	taskManager := NewTaskManager()

	// Create gin router
	r := gin.Default()

	// Enable CORS
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// API routes
	api := r.Group("/api")
	{
		api.POST("/tasks", taskManager.CreateTask)
		api.GET("/tasks", taskManager.GetTasks)
		api.GET("/tasks/:id/urls", taskManager.GetTaskDetails)
		api.POST("/tasks/:id/stop", taskManager.StopTask)
	}

	// Start server
	log.Fatal(r.Run(":8080"))
}
