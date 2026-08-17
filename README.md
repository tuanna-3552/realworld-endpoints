# RealWorld Backend API (Go + Echo + GORM + PostgreSQL)

An exemplary backend application implementing the [RealWorld](https://github.com/gothinkster/realworld) spec built with **Go**, **Echo Framework**, **GORM**, and **PostgreSQL**, structured using **Clean Architecture** principles.

---

## 🚀 Tech Stack

- **Language:** Go (1.26+)
- **Web Framework:** [Echo v4](https://echo.labstack.com/)
- **ORM:** [GORM](https://gorm.io/)
- **Database:** PostgreSQL 16 (via Docker Compose)
- **Environment Management:** [godotenv](https://github.com/joho/godotenv)

---

## 📁 Project Architecture

The project follows a standard Go Clean/Layered Architecture:

```
realworld-endpoints/
├── cmd/
│   └── api/
│       └── main.go              # Entry point for the Echo server
├── internal/
│   ├── config/                  # Environment variable configuration & DSN helper
│   ├── db/                      # GORM PostgreSQL initialization & auto-migrations
│   ├── models/                  # GORM models (User, Article, Comment, Tag)
│   ├── repository/              # Data access layer (Interfaces & GORM implementations)
│   ├── handlers/                # HTTP Handlers (Controllers)
│   └── routes/                  # Echo Router setup & Middleware configuration
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
DB_HOST=localhost
DB_PORT=54333
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=realworld
DB_SSLMODE=disable
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

## 🔌 Initial API Endpoints

| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `GET` | `/health` | Server Health Status |
| `GET` | `/api/users` | List Users |
| `GET` | `/api/articles` | List Articles |
| `GET` | `/api/tags` | List Tags |

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

- [x] **Step 1:** Environment setup, Echo Framework setup, GORM PostgreSQL connection, Clean Architecture directory structure, Initial Models & API routes.
- [ ] **Step 2:** User Authentication (JWT), Registration, Login & Profile Endpoints.
- [ ] **Step 3:** Articles CRUD operations, Slugs, & Pagination.
- [ ] **Step 4:** Comments API & Tagging System.
- [ ] **Step 5:** Favorite Articles & Follow/Unfollow User System.
