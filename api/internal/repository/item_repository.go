package repository

import "my-go-app/internal/models"

// ItemRepository defines the database interface for inventory items.
type ItemRepository interface {
	GetAll() ([]models.InventoryItem, error)
	GetByID(id string) (models.InventoryItem, error)
	Create(item models.InventoryItem) (models.InventoryItem, error)
	Update(item models.InventoryItem) (models.InventoryItem, error)
	Delete(id string) error
}
