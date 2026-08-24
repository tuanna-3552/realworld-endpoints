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
│   ├── middleware/              # Echo Middlewares (Logger, Recover, JWT Auth, Optional Auth)
│   ├── models/                  # GORM models (User, Article, Comment, Tag, UserFollow, ArticleFavorite)
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
| `GET` | `/api/articles` | Optional | List articles (supports `tag`, `author`, `favorited`, `limit`, `offset`) |
| `GET` | `/api/articles/feed` | **Yes (JWT)** | Get articles feed from followed users |
| `GET` | `/api/articles/:slug` | Optional | Get article details by slug |
| `POST` | `/api/articles/:slug/favorite` | **Yes (JWT)** | Favorite an article |
| `DELETE` | `/api/articles/:slug/favorite` | **Yes (JWT)** | Unfavorite an article |
| `GET` | `/api/articles/:slug/comments` | Optional | Get comments for an article |
| `POST` | `/api/articles/:slug/comments` | **Yes (JWT)** | Add comment to an article |
| `DELETE` | `/api/articles/:slug/comments/:id` | **Yes (JWT)** | Delete personal comment (ownership enforced) |
| `GET` | `/api/profiles/:username` | Optional | Get user profile by username |
| `POST` | `/api/profiles/:username/follow` | **Yes (JWT)** | Follow a user |
| `DELETE` | `/api/profiles/:username/follow` | **Yes (JWT)** | Unfollow a user |
| `GET` | `/api/tags` | No | List all tags |

---

## 🔑 Authentication & Features Usage Guide

### 1. User Register & Login
```bash
# Register
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{"user":{"username":"johndoe","email":"john@example.com","password":"password123"}}'

# Login
curl -X POST http://localhost:8080/api/users/login \
  -H "Content-Type: application/json" \
  -d '{"user":{"email":"john@example.com","password":"password123"}}'
```

### 2. Filtering & Pagination (`GET /api/articles`)
```bash
# Filter by tag and author with limit/offset
curl "http://localhost:8080/api/articles?tag=dragons&author=johndoe&limit=10&offset=0"
```

### 3. User Feed (`GET /api/articles/feed`)
```bash
curl -X GET "http://localhost:8080/api/articles/feed?limit=10&offset=0" \
  -H "Authorization: Token <YOUR_JWT_TOKEN>"
```

### 4. Follow / Unfollow User
```bash
# Follow User
curl -X POST http://localhost:8080/api/profiles/johndoe/follow \
  -H "Authorization: Token <YOUR_JWT_TOKEN>"

# Unfollow User
curl -X DELETE http://localhost:8080/api/profiles/johndoe/follow \
  -H "Authorization: Token <YOUR_JWT_TOKEN>"
```

### 5. Favorite / Unfavorite Article
```bash
# Favorite Article
curl -X POST http://localhost:8080/api/articles/how-to-train-your-dragon/favorite \
  -H "Authorization: Token <YOUR_JWT_TOKEN>"

# Unfavorite Article
curl -X DELETE http://localhost:8080/api/articles/how-to-train-your-dragon/favorite \
  -H "Authorization: Token <YOUR_JWT_TOKEN>"
```

### 6. Delete Personal Comment
```bash
curl -X DELETE http://localhost:8080/api/articles/how-to-train-your-dragon/comments/1 \
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
- [x] **Step 5:** Middleware, Filtering & Pagination (Optional Auth Middleware, query parameters filtering by `tag`/`author`/`favorited`, `limit`/`offset` pagination, `GET /api/articles/feed`).
- [x] **Step 6:** Social Features - Like, Follow & Comment (Favorite/Unfavorite articles, follow/unfollow users, delete comments with ownership check).
- [x] **Step 7:** Testing & Optimization (Unit testing with `httptest`, Redis caching for high-frequency endpoints, `.env` management).
- [ ] **Step 8:** Summary & Final Build (Swagger/OpenAPI integration, GORM query optimization avoiding N+1 problem).
