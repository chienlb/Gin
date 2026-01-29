# Project Structure Summary

## 📁 Complete Project Layout

```
gin-demo/
├── cmd/
│   └── api/
│       └── main.go                 # Application entry point
├── internal/
│   ├── app/
│   │   ├── server.go              # Server initialization and route setup
│   │   ├── router.go              # Router configuration
│   │   ├── middleware.go          # Middleware setup
│   │   └── route.go               # Route definitions
│   ├── config/
│   │   └── config.go              # Configuration management
│   ├── domain/
│   │   └── user.go                # User domain model and DTOs
│   ├── handler/
│   │   └── user_handler.go        # HTTP request handlers
│   ├── service/
│   │   └── user_service.go        # Business logic implementation
│   ├── repository/
│   │   └── user_repo.go           # Database access layer
│   └── database/
│       ├── postgres.go            # PostgreSQL initialization
│       └── migration.go           # Database migrations
├── pkg/
│   ├── logger/
│   │   └── logger.go              # Logging utilities
│   ├── response/
│   │   └── response.go            # API response helpers
│   └── utils/
│       └── utils.go               # Utility functions
├── migrations/
│   ├── 0001_init.up.sql          # Create users table
│   └── 0001_init.down.sql        # Drop users table
├── configs/
│   └── local.env                  # Local environment configuration
├── scripts/
│   └── run_local.ps1              # PowerShell run script
├── docs/
│   ├── README.md                  # Complete documentation
│   ├── GETTING_STARTED.md         # Quick start guide
│   └── PROJECT_STRUCTURE.md       # This file
├── .env.example                   # Environment variables template
├── .gitignore                     # Git ignore rules
├── Dockerfile                     # Docker image definition
├── docker-compose.yml             # Docker compose configuration
├── go.mod                         # Go module definition
└── go.sum                         # Go module checksums
```

## 📦 Dependencies

### Direct Dependencies
- `github.com/gin-gonic/gin v1.11.0` - Web framework
- `github.com/lib/pq v1.10.9` - PostgreSQL driver

## 🏗️ Architecture Layers

### 1. **Domain Layer** (`internal/domain/`)
- `User` - Main entity
- `CreateUserRequest` - Input DTO
- `UpdateUserRequest` - Update DTO
- `UserResponse` - Output DTO

### 2. **Repository Layer** (`internal/repository/`)
- `UserRepository` - Database access operations
  - `Create()` - Insert user
  - `GetByID()` - Fetch by ID
  - `GetByEmail()` - Fetch by email
  - `GetAll()` - List all users
  - `Update()` - Update user
  - `Delete()` - Delete user

### 3. **Service Layer** (`internal/service/`)
- `UserService` - Business logic
  - `CreateUser()` - User creation with validation
  - `GetUser()` - Fetch single user
  - `GetAllUsers()` - List users
  - `UpdateUser()` - Update with validation
  - `DeleteUser()` - Delete user

### 4. **Handler Layer** (`internal/handler/`)
- `UserHandler` - HTTP request handlers
  - `CreateUser()` - POST /api/users
  - `GetUser()` - GET /api/users/:id
  - `GetAllUsers()` - GET /api/users
  - `UpdateUser()` - PUT /api/users/:id
  - `DeleteUser()` - DELETE /api/users/:id

### 5. **Application Layer** (`internal/app/`)
- `Server` - Server lifecycle management
- `Router` - Route configuration
- `Middleware` - Request/response middleware

### 6. **Database Layer** (`internal/database/`)
- `Init()` - Database connection setup
- `GetDB()` - Get database connection
- `RunMigrations()` - Execute schema migrations
- `Close()` - Close database connection

### 7. **Support Layers**
- **Config** - Environment-based configuration
- **Logger** - Structured logging
- **Response** - Standardized API responses
- **Utils** - Helper functions (password hashing, email normalization)

## 🔄 Request Flow

```
HTTP Request
    ↓
Handler (user_handler.go)
    ↓
Service (user_service.go) - Business logic & validation
    ↓
Repository (user_repo.go) - Database operations
    ↓
Database (PostgreSQL)
    ↓
Response → Handler → HTTP Response
```

## 📊 Database Schema

### Users Table
```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_email ON users(email);
```

## 🚀 API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | Health check |
| POST | `/api/users` | Create user |
| GET | `/api/users` | Get all users |
| GET | `/api/users/:id` | Get user by ID |
| PUT | `/api/users/:id` | Update user |
| DELETE | `/api/users/:id` | Delete user |

## 🔧 Configuration

### Environment Variables
```env
# Server Configuration
SERVER_HOST=0.0.0.0
SERVER_PORT=8080

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=gin_db
DB_SSL_MODE=disable

# Logger Configuration
LOG_LEVEL=info
```

## 🎯 Key Features Implemented

✅ **Clean Architecture** - Clear separation of concerns
✅ **Dependency Injection** - Loose coupling between layers
✅ **Repository Pattern** - Abstraction for data access
✅ **Service Layer Pattern** - Business logic encapsulation
✅ **Error Handling** - Comprehensive error management
✅ **Logging** - Structured logging throughout
✅ **Database Migrations** - Schema version control
✅ **Configuration Management** - Environment-based config
✅ **Docker Support** - Container deployment
✅ **Input Validation** - Request validation using Gin bindings
✅ **Password Security** - SHA256 password hashing
✅ **Connection Pooling** - Efficient database connections
✅ **API Documentation** - Handler comments with Swagger format

## 🛠️ Development Best Practices

1. **Single Responsibility** - Each layer has one purpose
2. **DRY Principle** - Reusable code in pkg/
3. **Error Handling** - Descriptive error messages
4. **Configuration** - Environment-based configuration
5. **Logging** - Request/response logging
6. **Security** - Password hashing, input validation
7. **Testing** - Prepared for unit testing structure

## 📝 Running the Application

### Local Development
```bash
cp .env.example .env
go run ./cmd/api
```

### Docker
```bash
docker-compose up
```

### Scripts
```bash
.\scripts\run_local.ps1
```

## 🔍 File Purposes

| File | Purpose |
|------|---------|
| `main.go` | Application bootstrap |
| `server.go` | Server lifecycle and routes |
| `config.go` | Configuration loading |
| `user.go` | Domain models |
| `user_handler.go` | HTTP handlers |
| `user_service.go` | Business logic |
| `user_repo.go` | Database queries |
| `postgres.go` | DB connection |
| `migration.go` | Schema creation |
| `logger.go` | Logging utilities |
| `response.go` | Response formatting |
| `utils.go` | Helper functions |

## 🚦 Next Steps

1. **Run the application** - `go run ./cmd/api`
2. **Test endpoints** - Use provided curl examples
3. **Add new features** - Follow the same pattern for new entities
4. **Deploy** - Use Docker compose or Kubernetes
5. **Monitor** - Add additional logging and metrics

---

For detailed instructions, see [GETTING_STARTED.md](GETTING_STARTED.md)
For complete documentation, see [README.md](README.md)
