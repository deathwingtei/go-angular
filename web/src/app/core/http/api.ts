import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';

@Injectable({
  providedIn: 'root',
})
export class Api {
  constructor(private http: HttpClient) {}

  getHelloMessage(): Observable<{ text: string }> {
    return this.http.get<{ text: string }>('http://localhost:8080/api/hello');
  }
}
