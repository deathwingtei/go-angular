import { Component, OnInit } from '@angular/core';
import { Api } from '../../core/http/api';

@Component({
  selector: 'app-home',
  imports: [],
  templateUrl: './home.html',
  styleUrl: './home.css',
})
export class Home implements OnInit {
  message = 'กำลังโหลดข้อมูลจาก Go Backend...';

  constructor(private api: Api) {}

  ngOnInit() {
    this.api.getHelloMessage().subscribe({
      next: (data) => (this.message = data.text),
      error: (err) => (this.message = 'ไม่สามารถเชื่อมต่อ Go Backend ได้'),
    });
  }
}
