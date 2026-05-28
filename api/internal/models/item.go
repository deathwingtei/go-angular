package models

import "time"

// InventoryItem represents a household item in the inventory tracker.
type InventoryItem struct {
	ID                string    `json:"id"`
	Name              string    `json:"name"`
	Category          string    `json:"category"`           // e.g., "Kitchen", "Bathroom", "Snacks"
	CurrentStock      float64   `json:"currentStock"`       // Remaining stock
	BaseQuantity      float64   `json:"baseQuantity"`       // Total amount when restocked (e.g., 500)
	Unit              string    `json:"unit"`               // e.g., "ml", "pcs", "bags", "kg"
	LastRestockDate   time.Time `json:"lastRestockDate"`    // Date when it was last fully restocked
	ReplenishmentDate time.Time `json:"replenishmentDate"` // Calculated date when stock is predicted to run out
}
