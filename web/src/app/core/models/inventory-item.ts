export interface InventoryItem {
  id: string;
  name: string;
  category: string;
  currentStock: number;
  baseQuantity: number;
  unit: string;
  lastRestockDate: string; // ISO Date String
  replenishmentDate: string; // ISO Date String (Calculated predicted date)
}
