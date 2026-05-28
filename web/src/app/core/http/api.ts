import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { InventoryItem } from '../models/inventory-item';

@Injectable({
  providedIn: 'root',
})
export class Api {
  private baseUrl = 'http://localhost:8080/api';

  constructor(private http: HttpClient) {}

  getHelloMessage(): Observable<{ text: string }> {
    return this.http.get<{ text: string }>(`${this.baseUrl}/hello`);
  }

  // Get all inventory items
  getItems(): Observable<InventoryItem[]> {
    return this.http.get<InventoryItem[]>(`${this.baseUrl}/items`);
  }

  // Add a new inventory item
  createItem(item: Omit<InventoryItem, 'id' | 'lastRestockDate' | 'replenishmentDate'>): Observable<InventoryItem> {
    return this.http.post<InventoryItem>(`${this.baseUrl}/items`, item);
  }

  // Consume stock (Quick action)
  consumeStock(id: string, amount: number): Observable<InventoryItem> {
    return this.http.post<InventoryItem>(`${this.baseUrl}/items/${id}/consume`, { amount });
  }

  // Restock an item
  restockItem(id: string): Observable<InventoryItem> {
    return this.http.post<InventoryItem>(`${this.baseUrl}/items/${id}/restock`, {});
  }
}
