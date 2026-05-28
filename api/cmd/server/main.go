package main

import (
	"log"
	"my-go-app/internal/handlers"
	"my-go-app/internal/repository"
	"my-go-app/internal/usecase"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Initialize Core Clean Architecture Layers
	itemRepo := repository.NewInMemoryItemRepository()
	itemUseCase := usecase.NewItemUseCase(itemRepo)
	itemHandler := handlers.NewItemHandler(itemUseCase)

	// 2. Initialize Gin Router
	router := gin.Default()

	// 3. Configure CORS middleware for Angular
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:4200"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	router.Use(cors.New(config))

	// 4. Register Routes
	api := router.Group("/api")
	{
		api.GET("/hello", handlers.GetHello)

		// Inventory REST API Endpoints
		api.GET("/items", itemHandler.GetAll)
		api.GET("/items/:id", itemHandler.GetByID)
		api.POST("/items", itemHandler.Create)
		api.POST("/items/:id/consume", itemHandler.Consume)
		api.POST("/items/:id/restock", itemHandler.Restock)
	}

	// 5. Start Server
	log.Println("Backend server started at :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
