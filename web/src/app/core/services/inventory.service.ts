import { Injectable, signal, computed } from '@angular/core';
import { Api } from '../http/api';
import { InventoryItem } from '../models/inventory-item';
import { finalize } from 'rxjs/operators';

@Injectable({
  providedIn: 'root',
})
export class InventoryService {
  // 1. Core State Signals
  private _items = signal<InventoryItem[]>([]);
  private _loading = signal<boolean>(false);
  private _error = signal<string | null>(null);

  // 2. Public Read-Only State Expositions
  items = this._items.asReadonly();
  loading = this._loading.asReadonly();
  error = this._error.asReadonly();

  // 3. Computed Signals (Reactive Selectors)
  totalItems = computed(() => this._items().length);

  // Filter items running low on stock (<= 20% of base capacity)
  lowStockItems = computed(() =>
    this._items().filter((item) => item.currentStock / item.baseQuantity <= 0.20)
  );

  lowStockCount = computed(() => this.lowStockItems().length);

  // Filter critical stock items (<= 5% of base capacity)
  criticalCount = computed(() =>
    this._items().filter((item) => item.currentStock / item.baseQuantity <= 0.05).length
  );

  // Group items by category for visual organization
  categories = computed(() => {
    const list = this._items();
    const map = new Map<string, InventoryItem[]>();
    list.forEach((item) => {
      const group = map.get(item.category) || [];
      group.push(item);
      map.set(item.category, group);
    });
    return Array.from(map.entries()).map(([name, items]) => ({ name, items }));
  });

  constructor(private api: Api) {}

  // 4. Actions (State Mutation Methods)
  
  // Load all items from Go backend
  loadItems() {
    this._loading.set(true);
    this._error.set(null);
    this.api.getItems()
      .pipe(finalize(() => this._loading.set(false)))
      .subscribe({
        next: (data) => {
          // Sort items by predicted replenishment date (soonest depletion first)
          const sorted = [...data].sort((a, b) => 
            new Date(a.replenishmentDate).getTime() - new Date(b.replenishmentDate).getTime()
          );
          this._items.set(sorted);
        },
        error: (err) => {
          console.error('Error loading inventory items', err);
          this._error.set('ไม่สามารถเชื่อมต่อคลังสินค้าได้ โปรดลองอีกครั้ง');
        }
      });
  }

  // Quick Action: Consume stock
  consume(id: string, amount: number) {
    // Optional: Optimistic update can be done here, but standard API-driven update works great:
    this.api.consumeStock(id, amount).subscribe({
      next: (updatedItem) => {
        this.updateItemInState(updatedItem);
      },
      error: (err) => {
        console.error(`Error consuming stock for item ${id}`, err);
        this._error.set('ไม่สามารถตัดสต๊อกสินค้าได้');
      }
    });
  }

  // Action: Restock item
  restock(id: string) {
    this.api.restockItem(id).subscribe({
      next: (updatedItem) => {
        this.updateItemInState(updatedItem);
      },
      error: (err) => {
        console.error(`Error restocking item ${id}`, err);
        this._error.set('ไม่สามารถเติมสต๊อกสินค้าได้');
      }
    });
  }

  // Action: Add new item
  addItem(item: Omit<InventoryItem, 'id' | 'lastRestockDate' | 'replenishmentDate'>) {
    this._loading.set(true);
    this.api.createItem(item)
      .pipe(finalize(() => this._loading.set(false)))
      .subscribe({
        next: (newItem) => {
          // Add to current state list and sort
          const updated = [...this._items(), newItem].sort((a, b) => 
            new Date(a.replenishmentDate).getTime() - new Date(b.replenishmentDate).getTime()
          );
          this._items.set(updated);
        },
        error: (err) => {
          console.error('Error adding item', err);
          this._error.set('ไม่สามารถเพิ่มสินค้าชิ้นใหม่ได้');
        }
      });
  }

  // Helper to replace an updated item in our array and maintain state sorting
  private updateItemInState(updatedItem: InventoryItem) {
    const currentList = this._items();
    const updatedList = currentList.map((item) => 
      item.id === updatedItem.id ? updatedItem : item
    );
    
    // Re-sort so that items running out first remain at the top
    const sorted = updatedList.sort((a, b) => 
      new Date(a.replenishmentDate).getTime() - new Date(b.replenishmentDate).getTime()
    );
    
    this._items.set(sorted);
  }
}
