package main

import (
	"log"
	"my-go-app/internal/handlers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Initialize Gin Router
	router := gin.Default()

	// 2. Configure CORS middleware for Angular
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:4200"}
	router.Use(cors.New(config))

	// 3. Register Routes
	api := router.Group("/api")
	{
		api.GET("/hello", handlers.GetHello)
	}

	// 4. Start Server
	log.Println("Backend server started at :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
