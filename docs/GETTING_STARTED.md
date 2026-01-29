# Getting Started Guide

## Quick Start

### 1. Prerequisites
- Go 1.25.6 or higher
- PostgreSQL 15 or higher
- (Optional) Docker & Docker Compose

### 2. Clone/Setup the Project
```bash
cd gin-demo
```

### 3. Install Dependencies
```bash
go mod download
```

### 4. Configure Environment
Copy the example environment file:
```bash
cp .env.example .env
```

Edit `.env` with your database credentials:
```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=gin_db
```

### 5. Setup Database

**Option A: Manual Setup**
```bash
# Create database
createdb gin_db

# The application will auto-migrate tables on startup
```

**Option B: Docker Compose**
```bash
docker-compose up -d postgres
```

### 6. Run the Application

**Option A: Direct Execution**
```bash
go run ./cmd/api
```

**Option B: PowerShell Script (Windows)**
```powershell
.\scripts\run_local.ps1
```

**Option C: Docker**
```bash
docker-compose up
```

### 7. Verify the Application
```bash
curl http://localhost:8080/health
```

Expected response:
```json
{
  "status": "OK"
}
```

## Testing the API

### Create a User
```bash
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "password": "password123"
  }'
```

### Get All Users
```bash
curl http://localhost:8080/api/users
```

### Get User by ID
```bash
curl http://localhost:8080/api/users/1
```

### Update User
```bash
curl -X PUT http://localhost:8080/api/users/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Jane Doe",
    "email": "jane@example.com"
  }'
```

### Delete User
```bash
curl -X DELETE http://localhost:8080/api/users/1
```

## Project Structure Overview

```
gin-demo/
├── cmd/api/              # Application entry point
├── internal/
│   ├── app/              # Server, router, middleware setup
│   ├── config/           # Configuration management
│   ├── domain/           # Data models and DTOs
│   ├── handler/          # HTTP handlers/controllers
│   ├── service/          # Business logic
│   ├── repository/       # Data access layer
│   └── database/         # Database setup and migrations
├── pkg/
│   ├── logger/           # Logging utilities
│   ├── response/         # Response formatting
│   └── utils/            # Helper functions
├── migrations/           # SQL migration files
├── configs/              # Configuration files
├── scripts/              # Helper scripts
└── docs/                 # Documentation
```

## Development Workflow

### Adding a New Feature

1. **Define your domain model** in `internal/domain/`
2. **Create repository methods** in `internal/repository/`
3. **Implement business logic** in `internal/service/`
4. **Create HTTP handlers** in `internal/handler/`
5. **Register routes** in `internal/app/server.go`
6. **Create database migrations** in `migrations/`

### Example: Adding a Product Entity

1. Create `internal/domain/product.go`
2. Create `internal/repository/product_repo.go`
3. Create `internal/service/product_service.go`
4. Create `internal/handler/product_handler.go`
5. Add routes in `internal/app/server.go`
6. Create migration files in `migrations/`

## Troubleshooting

### Database Connection Issues
- Verify PostgreSQL is running
- Check database credentials in `.env`
- Ensure database exists: `createdb gin_db`

### Port Already in Use
- Change `SERVER_PORT` in `.env`
- Or kill the process using port 8080

### Module Not Found Error
```bash
go mod download
go mod tidy
```

## Additional Resources

- [Gin Framework Docs](https://github.com/gin-gonic/gin)
- [PostgreSQL Docs](https://www.postgresql.org/docs/)
- [Go Documentation](https://golang.org/doc/)

## Support

For issues or questions, please check the main README.md in the docs directory.
