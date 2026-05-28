import { Component, OnInit, signal, computed, ViewChild, ElementRef } from '@angular/core';
import { CommonModule } from '@angular/common';
import { InventoryService } from '../../core/services/inventory.service';
import { InventoryItem } from '../../core/models/inventory-item';

@Component({
  selector: 'app-dashboard',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './dashboard.html',
  styleUrl: './dashboard.css',
})
export class Dashboard implements OnInit {
  // 1. Local UI State Signals
  searchQuery = signal<string>('');
  selectedCategory = signal<string>('');
  isAddModalOpen = signal<boolean>(false);

  // Reference form fields for adding items
  @ViewChild('nameInput') nameInput!: ElementRef<HTMLInputElement>;
  @ViewChild('categorySelect') categorySelect!: ElementRef<HTMLSelectElement>;
  @ViewChild('currentInput') currentInput!: ElementRef<HTMLInputElement>;
  @ViewChild('baseInput') baseInput!: ElementRef<HTMLInputElement>;
  @ViewChild('unitSelect') unitSelect!: ElementRef<HTMLSelectElement>;

  // 2. Filtered Items Selector (Computed Signal)
  filteredItems = computed(() => {
    let items = this.inventoryService.items();
    const query = this.searchQuery().toLowerCase().trim();
    const category = this.selectedCategory();

    // Apply Search Query filter
    if (query) {
      items = items.filter((item) => item.name.toLowerCase().includes(query));
    }

    // Apply Category filter
    if (category) {
      items = items.filter((item) => item.category === category);
    }

    return items;
  });

  constructor(public inventoryService: InventoryService) {}

  ngOnInit() {
    this.inventoryService.loadItems();
  }

  // 3. Search & Filter Handlers
  updateSearch(event: Event) {
    const input = event.target as HTMLInputElement;
    this.searchQuery.set(input.value);
  }

  selectCategory(category: string) {
    this.selectedCategory.set(category);
  }

  trackByItemId(index: number, item: InventoryItem): string {
    return item.id;
  }

  // 4. Stock Math & Styling Helpers
  
  // Calculate stock percent (0-100)
  getStockPercent(item: InventoryItem): number {
    if (item.baseQuantity <= 0) return 0;
    const pct = (item.currentStock / item.baseQuantity) * 100;
    return Math.min(100, Math.max(0, Math.round(pct)));
  }

  // SVG Progress Ring Circle math: Circumference = 2 * PI * r = 2 * 3.14159 * 30 = 188.4
  getStrokeDashoffset(item: InventoryItem): number {
    const pct = this.getStockPercent(item);
    const circumference = 188.4;
    return circumference - (pct / 100) * circumference;
  }

  // Categorize Stock Level for visual tags
  getStockLevelClass(item: InventoryItem): 'safe' | 'warning' | 'critical' {
    const ratio = item.currentStock / item.baseQuantity;
    if (ratio <= 0.05) return 'critical';
    if (ratio <= 0.20) return 'warning';
    return 'safe';
  }

  // Text color helper
  getStockTextClass(item: InventoryItem): string {
    const level = this.getStockLevelClass(item);
    if (level === 'critical') return 'critical-color';
    if (level === 'warning') return 'warning-color';
    return 'safe-color';
  }

  // Thai stock status description label
  getStockStatusLabel(item: InventoryItem): string {
    const level = this.getStockLevelClass(item);
    if (level === 'critical') return 'วิกฤต! ใกล้หมดแล้ว';
    if (level === 'warning') return 'เตือน: ของใกล้หมด';
    return 'ปกติ: มีเพียงพอ';
  }

  // Return a readable days remaining prediction label
  getDaysRemainingLabel(item: InventoryItem): string {
    const now = new Date();
    const depletion = new Date(item.replenishmentDate);
    
    // Reset hours to compare dates cleanly
    now.setHours(0, 0, 0, 0);
    depletion.setHours(0, 0, 0, 0);

    const diffTime = depletion.getTime() - now.getTime();
    const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24));

    if (diffDays <= 0) {
      return 'ของหมดวันนี้!';
    }
    if (diffDays === 1) {
      return 'จะหมดในวันพรุ่งนี้';
    }
    if (diffDays <= 5) {
      return `จะหมดในอีก ${diffDays} วัน (ซื้อด่วน)`;
    }
    return `มีพอใช้อีก ${diffDays} วัน`;
  }

  // Return a default decrement amount based on unit
  getConsumeAmount(item: InventoryItem): number {
    if (item.unit === 'ml') return 50; // ml is high volume
    if (item.unit === 'kg') return 0.5; // kg is decimal
    return 1; // standard bottles, items, box etc.
  }

  // 5. Actions Handlers
  quickConsume(item: InventoryItem) {
    const amount = this.getConsumeAmount(item);
    this.inventoryService.consume(item.id, amount);
  }

  restockItem(id: string) {
    this.inventoryService.restock(id);
  }

  // 6. Date Formatting
  formatDate(dateStr: string): string {
    if (!dateStr) return '-';
    const date = new Date(dateStr);
    
    // Thai months mapping
    const thaiMonths = [
      'ม.ค.', 'ก.พ.', 'มี.ค.', 'เม.ย.', 'พ.ค.', 'มิ.ย.',
      'ก.ค.', 'ส.ค.', 'ก.ย.', 'ต.ค.', 'พ.ย.', 'ธ.ค.'
    ];

    const day = date.getDate();
    const month = thaiMonths[date.getMonth()];
    const year = date.getFullYear() + 543; // Buddhist Era
    
    return `${day} ${month} ${year}`;
  }

  // 7. Modal management
  openAddModal() {
    this.isAddModalOpen.set(true);
  }

  closeAddModal() {
    this.isAddModalOpen.set(false);
  }

  submitNewItem(event: Event) {
    event.preventDefault();

    const name = this.nameInput.nativeElement.value.trim();
    const category = this.categorySelect.nativeElement.value;
    const currentStock = parseFloat(this.currentInput.nativeElement.value);
    const baseQuantity = parseFloat(this.baseInput.nativeElement.value);
    const unit = this.unitSelect.nativeElement.value;

    if (!name || isNaN(currentStock) || isNaN(baseQuantity)) {
      alert('โปรดกรอกข้อมูลให้ครบถ้วนและถูกต้อง');
      return;
    }

    // Call service to trigger POST and refresh list
    this.inventoryService.addItem({
      name,
      category,
      currentStock,
      baseQuantity,
      unit,
    });

    // Close and reset
    this.closeAddModal();
  }
}
