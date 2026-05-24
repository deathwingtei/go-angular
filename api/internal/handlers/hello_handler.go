package handlers

import (
	"my-go-app/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetHello(c *gin.Context) {
	// Create the message model
	response := models.Message{
		Text: "สวัสดีจาก Go Backend! 🐹",
	}

	// Return as JSON
	c.JSON(http.StatusOK, response)
}
