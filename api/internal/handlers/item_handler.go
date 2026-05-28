package handlers

import (
	"fmt"
	"my-go-app/internal/models"
	"my-go-app/internal/usecase"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// ItemHandler manages HTTP requests for inventory items.
type ItemHandler struct {
	useCase usecase.ItemUseCase
}

// NewItemHandler creates a new ItemHandler with the injected UseCase dependency.
func NewItemHandler(uc usecase.ItemUseCase) *ItemHandler {
	return &ItemHandler{
		useCase: uc,
	}
}

// GetAll handles GET /api/items
func (h *ItemHandler) GetAll(c *gin.Context) {
	items, err := h.useCase.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch items: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}

// GetByID handles GET /api/items/:id
func (h *ItemHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	item, err := h.useCase.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

// CreateInput represents the JSON body structure for creating a new item.
type CreateInput struct {
	Name         string  `json:"name" binding:"required"`
	Category     string  `json:"category" binding:"required"`
	CurrentStock float64 `json:"currentStock" binding:"required"`
	BaseQuantity float64 `json:"baseQuantity" binding:"required"`
	Unit         string  `json:"unit" binding:"required"`
}

// Create handles POST /api/items
func (h *ItemHandler) Create(c *gin.Context) {
	var input CreateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Map input to domain models with dynamic unique ID
	uniqueID := fmt.Sprintf("item_%d", time.Now().UnixNano())
	newItem := models.InventoryItem{
		ID:           uniqueID,
		Name:         input.Name,
		Category:     input.Category,
		CurrentStock: input.CurrentStock,
		BaseQuantity: input.BaseQuantity,
		Unit:         input.Unit,
		LastRestockDate: time.Now(),
	}

	savedItem, err := h.useCase.Create(newItem)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create item: " + err.Error()})
		return
	}
	c.JSON(http.StatusCreated, savedItem)
}

// ConsumeInput represents the request payload for cutting stock.
type ConsumeInput struct {
	Amount float64 `json:"amount"`
}

// Consume handles POST /api/items/:id/consume
func (h *ItemHandler) Consume(c *gin.Context) {
	id := c.Param("id")
	
	var input ConsumeInput
	// It is fine if JSON is empty, we will default the consume amount
	_ = c.ShouldBindJSON(&input)

	// If amount is not specified, let's default to a sensible value:
	// 1.0 unit (for items like bottles/pcs) or 10% of stock.
	amount := input.Amount
	if amount <= 0 {
		amount = 1.0 // Default to consuming 1 unit
	}

	updatedItem, err := h.useCase.ConsumeStock(id, amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to consume item stock: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, updatedItem)
}

// Restock handles POST /api/items/:id/restock
func (h *ItemHandler) Restock(c *gin.Context) {
	id := c.Param("id")
	updatedItem, err := h.useCase.RestockItem(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to restock item: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, updatedItem)
}
