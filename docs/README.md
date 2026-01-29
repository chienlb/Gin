# Gin Demo API

A complete RESTful API built with Go and Gin framework following clean architecture principles.

## Architecture

```
gin-demo/
├── cmd/
│   └── api/              # Application entry point
├── internal/
│   ├── app/              # Application core (server, router, middleware)
│   ├── config/           # Configuration management
│   ├── domain/           # Domain models and entities
│   ├── handler/          # HTTP handlers/controllers
│   ├── service/          # Business logic layer
│   ├── repository/       # Data access layer
│   └── database/         # Database initialization and migrations
├── pkg/
│   ├── logger/           # Logging utilities
│   ├── response/         # Response formatting helpers
│   └── utils/            # Common utilities
├── migrations/           # Database migration files
├── configs/              # Configuration files
├── scripts/              # Helper scripts
└── docs/                 # Documentation
```

## Project Structure Details

### cmd/api
Entry point of the application. Initializes the server and starts the HTTP listener.

### internal/
Contains all application-specific code:
- **app**: Server configuration, routing setup, and middleware
- **config**: Environment-based configuration management
- **domain**: Data models and DTOs
- **handler**: HTTP request handlers
- **service**: Business logic implementation
- **repository**: Database access layer
- **database**: DB initialization and schema migrations

### pkg/
Reusable packages that could be shared across projects:
- **logger**: Structured logging
- **response**: Standardized response formatting
- **utils**: Helper functions (password hashing, email normalization, etc.)

## Features

- ✅ User CRUD operations
- ✅ PostgreSQL database integration
- ✅ Structured logging
- ✅ Clean architecture pattern
- ✅ Error handling
- ✅ Database migrations
- ✅ Docker support
- ✅ Environment configuration

## Prerequisites

- Go 1.25.6+
- PostgreSQL 15+
- Docker & Docker Compose (optional)

## Environment Variables

Create a `.env` file based on `.env.example`:

```bash
cp .env.example .env
```

### Configuration Options

```env
# Server
SERVER_HOST=0.0.0.0
SERVER_PORT=8080

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=gin_db
DB_SSL_MODE=disable

# Logger
LOG_LEVEL=info
```

## Setup & Running

### Local Development

1. **Install dependencies:**
   ```bash
   go mod download
   ```

2. **Set up PostgreSQL:**
   ```bash
   # Create database
   createdb gin_db
   ```

3. **Run the application:**
   ```bash
   # Using PowerShell script (Windows)
   .\scripts\run_local.ps1

   # Or directly with go
   go run ./cmd/api
   ```

4. **Server will start at:** `http://localhost:8080`

### Docker Compose

```bash
# Start services
docker-compose up

# Stop services
docker-compose down

# Remove volumes
docker-compose down -v
```

## API Endpoints

### Health Check
- `GET /health` - Application health status

### Users API
- `POST /api/users` - Create a new user
- `GET /api/users` - Get all users
- `GET /api/users/:id` - Get user by ID
- `PUT /api/users/:id` - Update user
- `DELETE /api/users/:id` - Delete user

### Example Requests

**Create User:**
```bash
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "password": "password123"
  }'
```

**Get All Users:**
```bash
curl http://localhost:8080/api/users
```

**Get User by ID:**
```bash
curl http://localhost:8080/api/users/1
```

**Update User:**
```bash
curl -X PUT http://localhost:8080/api/users/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Jane Doe",
    "email": "jane@example.com"
  }'
```

**Delete User:**
```bash
curl -X DELETE http://localhost:8080/api/users/1
```

## Database Schema

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

## Project Development

### Adding New Features

1. **Define Domain Model** (`internal/domain/`)
2. **Create Repository** (`internal/repository/`)
3. **Create Service** (`internal/service/`)
4. **Create Handler** (`internal/handler/`)
5. **Add Routes** (`internal/app/server.go`)
6. **Database Migrations** (`migrations/`)

## Best Practices Implemented

- ✅ Dependency Injection
- ✅ Repository Pattern
- ✅ Service Layer Pattern
- ✅ Clean Architecture
- ✅ Error Handling
- ✅ Logging
- ✅ Configuration Management
- ✅ Database Connection Pooling

## Technologies Used

- **Framework**: Gin Web Framework
- **Database**: PostgreSQL
- **Logging**: Standard Go Logger
- **Password Hashing**: SHA256
- **Container**: Docker & Docker Compose

## License

This project is licensed under the MIT License.
