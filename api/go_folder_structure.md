# Go API Best Practices & Folder Structure

Just like Angular, building a scalable backend in Go requires a solid foundation. While Go is technically unopinionated, the community has settled on a few standard patterns that work incredibly well for REST APIs.

## 1. The Base Framework

For modern Go development, you have three primary choices for your base framework:

1. **Gin (`github.com/gin-gonic/gin`) - 🏆 Recommended**
   - **Why:** It is the most popular Go web framework. It is incredibly fast, has a massive ecosystem of middleware (CORS, logging, auth), and is very easy to learn.
2. **Echo (`github.com/labstack/echo`)**
   - **Why:** Very similar to Gin, but some developers prefer its cleaner API and excellent documentation.
3. **Standard Library (`net/http`)**
   - **Why:** Since Go 1.22, the standard library now supports URL path parameters (e.g., `GET /api/users/{id}`). If you want zero dependencies, this is now a highly viable option!

*(For your project, **Gin** is highly recommended as it makes handling JSON and CORS extremely simple).*

---

## 2. Recommended Folder Structure

The best practice for modern Go applications is a mix of the Standard Layout (`cmd/` and `internal/`) combined with **Domain-Driven Design**.

Here is the ideal structure for your `api/` folder:

```text
api/
├── cmd/
│   └── api/
│       └── main.go             # 1. The Entry Point. Only wires things together and starts the server.
│
├── internal/                   # 2. Private Code. Go prevents other projects from importing this.
│   ├── config/                 # Loads environment variables (e.g., DB URLs, Ports)
│   │
│   ├── core/                   # 3. Domain Models (The heart of your app)
│   │   └── models.go           # e.g., type User struct { ... }
│   │
│   ├── handlers/               # 4. HTTP Layer (Controllers)
│   │   └── user_handler.go     # Reads HTTP requests, calls services, returns JSON
│   │
│   ├── services/               # 5. Business Logic
│   │   └── user_service.go     # e.g., "Check if user exists, then hash password"
│   │
│   └── repositories/           # 6. Database Layer
│       └── user_repo.go        # SQL/GORM queries (e.g., SELECT * FROM users)
│
├── pkg/                        # 7. (Optional) Public utilities (e.g., custom loggers)
├── go.mod                      # Go dependencies
└── go.sum
```

### Why this structure? (The Flow of Data)

This structure follows **Clean Architecture**. When a request hits your server, the data flows like this:

1. **`main.go`** starts the Gin server and sets up the routes.
2. The route points to a **`Handler`** (e.g., `GetUser()`).
3. The **`Handler`** extracts the JSON/Params from the request and passes it to a **`Service`**.
4. The **`Service`** performs the actual business logic (validation, calculations) and asks the **`Repository`** for data.
5. The **`Repository`** talks to the Database and returns a **`Model`**.
6. The data flows back up and the **`Handler`** sends it to Angular as JSON.

### Moving Forward

Currently, your Go app just has a single `main.go` file. Would you like me to install **Gin** and refactor your current `api/` folder into this best-practice structure?
