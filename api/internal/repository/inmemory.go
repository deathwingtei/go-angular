package repository

import (
	"errors"
	"my-go-app/internal/models"
	"sync"
	"time"
)

// InMemoryItemRepository is a thread-safe, in-memory implementation of ItemRepository.
type InMemoryItemRepository struct {
	mu    sync.RWMutex
	items map[string]models.InventoryItem
}

// NewInMemoryItemRepository creates a new in-memory repository pre-populated with seed data.
func NewInMemoryItemRepository() *InMemoryItemRepository {
	repo := &InMemoryItemRepository{
		items: make(map[string]models.InventoryItem),
	}
	repo.seed()
	return repo
}

// seed populates initial household inventory items with realistic last restock times.
func (r *InMemoryItemRepository) seed() {
	now := time.Now()

	// Item 1: Fresh Milk (ห้องครัว / ตู้เย็น)
	item1 := models.InventoryItem{
		ID:              "item_1",
		Name:            "นมสดพาสเจอร์ไรส์ Meiji (จืด)",
		Category:        "ตู้เย็น / เครื่องดื่ม",
		CurrentStock:    2.0,
		BaseQuantity:    4.0,
		Unit:            "ขวด",
		LastRestockDate: now.AddDate(0, 0, -4), // 4 days ago
	}

	// Item 2: Fabric Softener (ห้องน้ำ / ซักรีด)
	item2 := models.InventoryItem{
		ID:              "item_2",
		Name:            "น้ำยาปรับผ้านุ่ม Downy (สีฟ้า)",
		Category:        "ห้องน้ำ / ซักรีด",
		CurrentStock:    120.0,
		BaseQuantity:    500.0,
		Unit:            "ml",
		LastRestockDate: now.AddDate(0, 0, -10), // 10 days ago
	}

	// Item 3: Jasmine Rice (ห้องครัว / เสบียงแห้ง)
	item3 := models.InventoryItem{
		ID:              "item_3",
		Name:            "ข้าวหอมมะลิ 100% ตราฉัตร",
		Category:        "เสบียง / ของแห้ง",
		CurrentStock:    4.2,
		BaseQuantity:    5.0,
		Unit:            "kg",
		LastRestockDate: now.AddDate(0, 0, -15), // 15 days ago
	}

	// Item 4: Liquid Body Wash (ห้องน้ำ)
	item4 := models.InventoryItem{
		ID:              "item_4",
		Name:            "ครีมอาบน้ำ Shokubutsu Monogatari",
		Category:        "ห้องน้ำ / สบู่",
		CurrentStock:    320.0,
		BaseQuantity:    400.0,
		Unit:            "ml",
		LastRestockDate: now.AddDate(0, 0, -3), // 3 days ago
	}

	r.items[item1.ID] = item1
	r.items[item2.ID] = item2
	r.items[item3.ID] = item3
	r.items[item4.ID] = item4
}

// GetAll returns all inventory items.
func (r *InMemoryItemRepository) GetAll() ([]models.InventoryItem, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	items := make([]models.InventoryItem, 0, len(r.items))
	for _, item := range r.items {
		items = append(items, item)
	}
	return items, nil
}

// GetByID retrieves a single item by its ID.
func (r *InMemoryItemRepository) GetByID(id string) (models.InventoryItem, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	item, exists := r.items[id]
	if !exists {
		return models.InventoryItem{}, errors.New("item not found")
	}
	return item, nil
}

// Create inserts a new item into the repository.
func (r *InMemoryItemRepository) Create(item models.InventoryItem) (models.InventoryItem, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.items[item.ID] = item
	return item, nil
}

// Update modifies an existing item in the repository.
func (r *InMemoryItemRepository) Update(item models.InventoryItem) (models.InventoryItem, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.items[item.ID]; !exists {
		return models.InventoryItem{}, errors.New("item not found for update")
	}
	r.items[item.ID] = item
	return item, nil
}

// Delete removes an item from the repository.
func (r *InMemoryItemRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.items[id]; !exists {
		return errors.New("item not found for deletion")
	}
	delete(r.items, id)
	return nil
}
