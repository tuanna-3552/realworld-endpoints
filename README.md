# RealWorld Backend API (Go + Echo + GORM + PostgreSQL)

An exemplary backend application implementing the [RealWorld](https://github.com/gothinkster/realworld) spec built with **Go**, **Echo Framework**, **GORM**, and **PostgreSQL**, structured using **Clean Architecture** principles.

---

## 🚀 Tech Stack

- **Language:** Go (1.26+)
- **Web Framework:** [Echo v4](https://echo.labstack.com/)
- **ORM:** [GORM](https://gorm.io/)
- **Database:** PostgreSQL 16 (via Docker Compose)
- **Authentication:** JWT (`github.com/golang-jwt/jwt/v5`) & bcrypt (`golang.org/x/crypto/bcrypt`)
- **Environment Management:** [godotenv](https://github.com/joho/godotenv)

---

## 📁 Project Architecture

The project follows a standard Go Clean/Layered Architecture:

```
realworld-endpoints/
├── cmd/
│   └── api/
│       └── main.go              # Entry point for the Echo HTTP server
├── internal/
│   ├── auth/                    # JWT token generation/parsing & bcrypt password hashing
│   ├── config/                  # Environment variable configuration & DSN helper
│   ├── db/                      # GORM PostgreSQL initialization & auto-migrations
│   ├── dto/                     # Data Transfer Objects (Request/Response schemas)
│   ├── handlers/                # HTTP Handlers (Controllers)
│   ├── middleware/              # Echo Middlewares (Logger, Recover, JWT Auth)
│   ├── models/                  # GORM models (User, Article, Comment, Tag)
│   ├── repository/              # Data access layer (Interfaces & GORM implementations)
│   └── routes/                  # Echo Router setup & endpoint registrations
├── .env                         # Local environment configuration
├── .env.example                 # Example environment template
├── docker-compose.yml           # PostgreSQL Docker Compose service
├── go.mod                       # Go module dependencies
└── go.sum                       # Go checksums
```

---

## 🛠 Getting Started

### Prerequisites

- [Go](https://go.dev/dl/) `1.26` or higher installed
- [Docker](https://www.docker.com/) and Docker Compose installed

### 1. Environment Setup

Copy `.env.example` to `.env`:

```bash
cp .env.example .env
```

Ensure default environment configuration matches your local port & database setup:

```env
PORT=8080
DB_HOST=127.0.0.1
DB_PORT=54333
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=realworld
DB_SSLMODE=disable
JWT_SECRET=secret-jwt-key-change-in-production
```

### 2. Start PostgreSQL via Docker

Run PostgreSQL using Docker Compose:

```bash
docker compose up -d
```

### 3. Run the Backend API

Start the HTTP server:

```bash
go run ./cmd/api
```

The server will initialize the PostgreSQL connection, execute GORM AutoMigrate for all models, and listen on `http://localhost:8080`.

---

## 🔌 Implemented API Endpoints

| Method | Endpoint | Auth Required | Description |
| :--- | :--- | :--- | :--- |
| `GET` | `/health` | No | Server Health Status |
| `POST` | `/api/users` | No | Register a new user |
| `POST` | `/api/users/login` | No | User login (returns JWT Token) |
| `GET` | `/api/user` | **Yes (JWT)** | Get current authenticated user profile |
| `GET` | `/api/users` | No | List all users |
| `GET` | `/api/articles` | No | List articles |
| `GET` | `/api/articles/:slug` | No | Get article details by slug |
| `GET` | `/api/articles/:slug/comments` | No | Get comments for an article |
| `POST` | `/api/articles/:slug/comments` | No | Add comment to an article |
| `GET` | `/api/profiles/:username` | No | Get user profile by username |
| `GET` | `/api/tags` | No | List all tags |

---

## 🔑 Step 4: JWT Authentication Guide

### 1. Register New User (`POST /api/users`)
```bash
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{"user":{"username":"johndoe","email":"john@example.com","password":"password123"}}'
```

### 2. User Login (`POST /api/users/login`)
```bash
curl -X POST http://localhost:8080/api/users/login \
  -H "Content-Type: application/json" \
  -d '{"user":{"email":"john@example.com","password":"password123"}}'
```

### 3. Get Current User (`GET /api/user`)
```bash
curl -X GET http://localhost:8080/api/user \
  -H "Authorization: Token <YOUR_JWT_TOKEN>"
```

---

## 🧪 Testing & Verification

Run static code analysis and build checks:

```bash
# Code static analysis
go vet ./...

# Build executable
go build -o api.exe ./cmd/api
```

---

## 🗺 Roadmap

- [x] **Step 1:** Environment Setup & Echo Framework Basics (Go, Echo, GORM & PostgreSQL setup).
- [x] **Step 2:** Basic CRUD with Handlers & Data Binding (`GET /api/articles`, `GET /api/articles/:slug`, `GET /api/tags`).
- [x] **Step 3:** Model Relationships & Data Transformation (`GET/POST /api/articles/:slug/comments`, `GET /api/profiles/:username`, DTOs, GORM Preload).
- [x] **Step 4:** User Authentication via JWT (`POST /api/users`, `POST /api/users/login`, `GET /api/user`, JWT Middleware, bcrypt password hashing).
- [ ] **Step 5:** Middleware, Filtering & Pagination (Custom permission middleware, query parameters filtering, limit/offset pagination).
- [ ] **Step 6:** Social Features - Like, Follow & Comment (Favorite/Unfavorite articles, follow users, delete comments).
- [ ] **Step 7:** Testing & Optimization (Unit testing with `httptest`, Redis caching for high-frequency endpoints, `.env` management).
- [ ] **Step 8:** Summary & Final Build (Swagger/OpenAPI integration, GORM query optimization avoiding N+1 problem).
