# คู่มือโครงสร้างและแนวทางปฏิบัติที่ดีที่สุดสำหรับ Go + Angular (Best Practices Guide)

เอกสารนี้รวบรวมแนวทางปฏิบัติที่ดีที่สุด (Best Practices) สำหรับการจัดโครงสร้างโปรเจกต์ทั้งฝั่ง Frontend (Angular) และ Backend (Go) เพื่อใช้สำหรับการเรียนรู้และสามารถนำไปอัปโหลดลงใน NotebookLM ได้

---

## ส่วนที่ 1: โครงสร้างโปรเจกต์ Angular (Frontend)

ใน Angular การจัดโครงสร้างโค้ดให้เป็นระเบียบมีความสำคัญมากเมื่อโปรเจกต์มีขนาดใหญ่ขึ้น รูปแบบที่ได้รับการยอมรับมากที่สุดคือสถาปัตยกรรมแบบ **Core / Shared / Features** 

### โครงสร้างโฟลเดอร์ที่แนะนำ (`src/app/`)

```text
src/
└── app/
    ├── core/                   # 1. Core Module (โหลดครั้งเดียว)
    │   ├── auth/               # จัดการการเข้าสู่ระบบ
    │   ├── http/               # Services สำหรับเรียก API (เช่น api.service.ts)
    │   └── layout/             # โครงสร้างหลักของเว็บ (Header, Footer, Sidebar)
    │
    ├── shared/                 # 2. Shared Module (ใช้ซ้ำได้หลายที่)
    │   ├── components/         # UI Components โง่ๆ (Dumb Components) เช่น ปุ่ม, การ์ด
    │   ├── pipes/              # ตัวแปลงข้อมูล
    │   └── utils/              # ฟังก์ชันช่วยเหลือต่างๆ
    │
    ├── features/               # 3. Features (หน้าเพจต่างๆ)
    │   ├── home/               # หน้า Home
    │   └── dashboard/          # หน้า Dashboard (Lazy Loaded)
    │
    ├── app.ts                  # Root Component (ตัวจัดการ Layout หลัก)
    ├── app.config.ts           # ไฟล์ตั้งค่า (Providers)
    └── app.routes.ts           # ไฟล์จัดการ Routing หลัก
```

### สรุปคอนเซปต์สำคัญใน Angular
1. **Core:** เก็บสิ่งที่มีเพียง "ชิ้นเดียว" ในแอป (Singletons) เช่น Layout หลัก (Header/Footer) หรือ Service ที่ติดต่อกับ Backend
2. **Shared:** เก็บ UI Components ที่ไม่มี Logic ผูกติดกับข้อมูลใดๆ (รับค่าผ่าน `@Input` และส่งค่ากลับด้วย `@Output`)
3. **Features (Pages):** เทียบได้กับโฟลเดอร์ `pages/` ใน React ที่นี่จะเก็บ Component ที่มีความฉลาด (Smart Components) ที่ดึงข้อมูลจาก Services และแสดงผล
4. **Root Component (`app.ts`):** ทำหน้าที่เป็นเปลือกหุ้ม (Shell) โดยการเรียกใช้ `<app-header>`, `<router-outlet>`, และ `<app-footer>`

---

## ส่วนที่ 2: โครงสร้างโปรเจกต์ Go API (Backend)

สำหรับ Go ภาษาไม่ได้บังคับโครงสร้างที่ตายตัว แต่ในระดับสากลนิยมใช้ **Clean Architecture** หรือ Domain-Driven Design เพื่อแยกส่วนการทำงานออกจากกันให้ชัดเจน

### โครงสร้างโฟลเดอร์ที่แนะนำ

```text
api/
├── cmd/
│   └── server/
│       └── main.go             # 1. จุดเริ่มต้นของโปรแกรม (Entry Point)
│
├── internal/                   # 2. โค้ดส่วนตัวของแอป (Private Code)
│   ├── handlers/               # HTTP Layer (รับ Request และตอบกลับเป็น JSON)
│   ├── models/                 # โครงสร้างข้อมูล (Structs)
│   ├── services/               # โลจิกทางธุรกิจ (Business Logic)
│   └── repositories/           # การจัดการฐานข้อมูล (Database Layer)
│
├── go.mod                      # ไฟล์จัดการ Dependencies
└── go.sum
```

### ทำไมถึงต้องใช้โครงสร้างนี้?
เมื่อมีคำขอ (Request) เข้ามา การไหลของข้อมูลจะเป็นดังนี้:
1. **`cmd/server/main.go`**: รับหน้าที่ตั้งค่า Framework, เปิดพอร์ต และกำหนด Routing
2. **`handlers/`**: รับ Request ตรวจสอบความถูกต้อง แล้วส่งต่อให้ Service
3. **`services/`**: คิดคำนวณและประมวลผลโลจิกหลัก (เช่น ตรวจสอบรหัสผ่าน)
4. **`repositories/`**: คุยกับ Database (SQL/NoSQL) เพื่อดึงหรือบันทึกข้อมูล
5. ข้อมูลไหลกลับไปที่ Handler เพื่อแปลงเป็น JSON ส่งกลับให้ Angular

### Framework ที่แนะนำ: Gin (`gin-gonic/gin`)
แม้ Go จะมี Standard Library ที่ดี แต่ **Gin** คือ Framework ที่ได้รับความนิยมสูงสุด เนื่องจาก:
- ทำงานได้เร็วมาก (High Performance)
- จัดการ JSON และ Path Parameters ได้ง่าย
- มี Middleware ครบครัน (เช่น CORS, Authentication, Logging)

---

## ส่วนที่ 3: การเชื่อมต่อระหว่าง Angular และ Go (Integration)

เมื่อนำทั้งสองส่วนมาต่อกัน หัวใจสำคัญคือการคุยกันผ่าน HTTP Protocols และการจัดการ CORS (Cross-Origin Resource Sharing)

### ฝั่ง Backend (Go)
ต้องตั้งค่า CORS อนุญาตให้ Angular (ซึ่งปกติรันที่ `localhost:4200`) สามารถดึงข้อมูลได้:
```go
// ใน cmd/server/main.go
config := cors.DefaultConfig()
config.AllowOrigins = []string{"http://localhost:4200"}
router.Use(cors.New(config))
```

### ฝั่ง Frontend (Angular)
ใช้ `HttpClient` ซึ่งตั้งค่าไว้ใน `app.config.ts` และถูกเรียกใช้ผ่าน Service ในโฟลเดอร์ `core/http/`:
```typescript
// ใน core/http/api.service.ts
this.http.get<{ text: string }>('http://localhost:8080/api/hello')
```

---
*หมายเหตุ: เอกสารนี้ถูกสร้างมาให้มีโครงสร้างที่ชัดเจน เหมาะสำหรับการนำไปประมวลผลต่อในเครื่องมือ AI อย่าง NotebookLM เพื่อใช้สร้างแบบทดสอบหรือถาม-ตอบสำหรับทบทวนความรู้*
