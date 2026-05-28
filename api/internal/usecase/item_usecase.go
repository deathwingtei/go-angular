package usecase

import (
	"math"
	"my-go-app/internal/models"
	"my-go-app/internal/repository"
	"time"
)

// ItemUseCase defines the business logic operations for inventory items.
type ItemUseCase interface {
	GetAll() ([]models.InventoryItem, error)
	GetByID(id string) (models.InventoryItem, error)
	Create(item models.InventoryItem) (models.InventoryItem, error)
	ConsumeStock(id string, amount float64) (models.InventoryItem, error)
	RestockItem(id string) (models.InventoryItem, error)
}

type itemUseCase struct {
	repo repository.ItemRepository
}

// NewItemUseCase creates a new ItemUseCase instance with the given repository.
func NewItemUseCase(repo repository.ItemRepository) ItemUseCase {
	return &itemUseCase{
		repo: repo,
	}
}

// GetAll fetches all items and calculates their predicted depletion dates.
func (u *itemUseCase) GetAll() ([]models.InventoryItem, error) {
	items, err := u.repo.GetAll()
	if err != nil {
		return nil, err
	}

	// Calculate predicted dates dynamically on retrieval
	for i := range items {
		items[i].ReplenishmentDate = u.PredictReplenishmentDate(items[i])
	}
	return items, nil
}

// GetByID fetches a single item and calculates its predicted depletion date.
func (u *itemUseCase) GetByID(id string) (models.InventoryItem, error) {
	item, err := u.repo.GetByID(id)
	if err != nil {
		return models.InventoryItem{}, err
	}

	item.ReplenishmentDate = u.PredictReplenishmentDate(item)
	return item, nil
}

// Create stores a new item and calculates its initial prediction.
func (u *itemUseCase) Create(item models.InventoryItem) (models.InventoryItem, error) {
	if item.LastRestockDate.IsZero() {
		item.LastRestockDate = time.Now()
	}
	item.ReplenishmentDate = u.PredictReplenishmentDate(item)

	return u.repo.Create(item)
}

// ConsumeStock decreases an item's current stock and recalibrates its depletion rate.
func (u *itemUseCase) ConsumeStock(id string, amount float64) (models.InventoryItem, error) {
	item, err := u.repo.GetByID(id)
	if err != nil {
		return models.InventoryItem{}, err
	}

	item.CurrentStock -= amount
	if item.CurrentStock < 0 {
		item.CurrentStock = 0
	}

	// Dynamic calculation based on new usage
	item.ReplenishmentDate = u.PredictReplenishmentDate(item)

	return u.repo.Update(item)
}

// RestockItem resets the current stock back to the base quantity and updates the restock date.
func (u *itemUseCase) RestockItem(id string) (models.InventoryItem, error) {
	item, err := u.repo.GetByID(id)
	if err != nil {
		return models.InventoryItem{}, err
	}

	item.CurrentStock = item.BaseQuantity
	item.LastRestockDate = time.Now()
	item.ReplenishmentDate = u.PredictReplenishmentDate(item)

	return u.repo.Update(item)
}

// PredictReplenishmentDate calculates the predicted date when the stock will be depleted.
// Uses real consumption rate: usedQuantity / daysPassed
func (u *itemUseCase) PredictReplenishmentDate(item models.InventoryItem) time.Time {
	now := time.Now()
	
	// Avoid division by zero by setting minimum duration as 2 hours
	durationPassed := now.Sub(item.LastRestockDate)
	daysPassed := durationPassed.Hours() / 24.0
	if daysPassed <= 0.05 { // approx 1.2 hours
		daysPassed = 0.05
	}

	usedQuantity := item.BaseQuantity - item.CurrentStock

	// If no quantity has been used or current stock is higher than base, default to 30 days out
	if usedQuantity <= 0 {
		return now.AddDate(0, 0, 30)
	}

	consumptionRatePerDay := usedQuantity / daysPassed
	if consumptionRatePerDay <= 0 {
		return now.AddDate(0, 0, 30)
	}

	daysLeft := item.CurrentStock / consumptionRatePerDay

	// Calculate and round target replenishment date
	daysLeftRounded := int(math.Round(daysLeft))
	if daysLeftRounded < 0 {
		daysLeftRounded = 0
	}
	
	return now.AddDate(0, 0, daysLeftRounded)
}
