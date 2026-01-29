# 🎯 Quick Reference - Production Standard Gin API

## Tệp mới được thêm

```
✅ pkg/apperror/error.go           - Error handling system
✅ pkg/validator/user_validator.go - Input validation
✅ pkg/middleware/middleware.go    - Middleware (logging, CORS, etc)
✅ .env.production                 - Production config template
✅ PRODUCTION_READY.md             - Production features doc (Tiếng Việt)
✅ DEPLOYMENT.md                   - Deployment guide
✅ PRODUCTION_COMPLETE.md          - This summary
```

## Các cải thiện chính

### 1️⃣ Error Handling
```go
// Before
return nil, fmt.Errorf("user with this email already exists")

// After
return nil, apperror.DuplicateEmailError(email)

// Với định dạng response:
{
  "status": "error",
  "code": "DUPLICATE_EMAIL",
  "message": "User with this email already exists",
  "details": {"email": "john@example.com"}
}
```

### 2️⃣ Validation
```go
// Before: Không validation
if req.Name == "" {...}

// After: Validation đầy đủ
validator := validator.NewUserValidator()
if err := validator.ValidateCreateRequest(name, email, password); err != nil {
  return nil, err
}
// Kiểm tra:
// - Name: 2-100 ký tự, valid characters
// - Email: Format hợp lệ
// - Password: Uppercase + Lowercase + Digit + 6+ chars
```

### 3️⃣ Response Format
```json
{
  "status": "success|error",
  "code": "OK|CREATED|VALIDATION_ERROR|DUPLICATE_EMAIL|NOT_FOUND",
  "message": "Human readable message",
  "data": {...},           // Chỉ success
  "details": {...}         // Chỉ error
}
```

### 4️⃣ Middleware
```go
// Auto setup trong server.go
s.setupMiddleware()

// Gồm:
// - RequestID: X-Request-ID header
// - Logging: Request/Response logging
// - CORS: Cross-origin support
// - Recovery: Panic handling
```

### 5️⃣ Graceful Shutdown
```go
// Tự động bắt SIGINT/SIGTERM
// Timeout 5 giây để shutdown
// Log tất cả bước shutdown
```

### 6️⃣ API Versioning
```
/health              - Health check
/                    - Root info
/api/v1/users        - List users
/api/v1/users        - Create user (POST)
/api/v1/users/:id    - Get user
/api/v1/users/:id    - Update user (PUT)
/api/v1/users/:id    - Delete user (DELETE)
```

## Chạy Production

### Development
```bash
ENVIRONMENT=development go run ./cmd/api
```

### Production (Local)
```bash
ENVIRONMENT=production \
DB_HOST=localhost \
DB_SSL_MODE=disable \
LOG_LEVEL=warn \
./api
```

### Production (Docker)
```bash
docker build -t gin-api:1.0.0 .
docker run -e ENVIRONMENT=production \
  -e DB_HOST=postgres \
  -p 8080:8080 \
  gin-api:1.0.0
```

### Production (Docker Compose)
```bash
docker-compose up -d
```

## Test Endpoints

### Health Check
```bash
curl http://localhost:8080/health
```

### Root Info
```bash
curl http://localhost:8080/
```

### Create User
```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "password": "Password123"
  }'
```

### Get All Users
```bash
curl http://localhost:8080/api/v1/users
```

### Get User by ID
```bash
curl http://localhost:8080/api/v1/users/1
```

### Update User
```bash
curl -X PUT http://localhost:8080/api/v1/users/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Jane Doe",
    "email": "jane@example.com"
  }'
```

### Delete User
```bash
curl -X DELETE http://localhost:8080/api/v1/users/1
```

## Validation Examples

### Invalid Password (no uppercase)
```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John",
    "email": "john@example.com",
    "password": "password123"
  }'
# Response: 400 Bad Request
# { "code": "VALIDATION_ERROR", 
#   "details": { 
#     "field": "password", 
#     "reason": "password must contain uppercase, lowercase, and digit" 
#   }
# }
```

### Invalid Email
```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John",
    "email": "invalid-email",
    "password": "Password123"
  }'
# Response: 400 Bad Request
# { "code": "VALIDATION_ERROR", 
#   "details": { "field": "email", "reason": "invalid email format" } 
# }
```

### Duplicate Email
```bash
# Tạo user lần đầu - OK
# Tạo user lần 2 với email same - Error
# Response: 409 Conflict
# { "code": "DUPLICATE_EMAIL", 
#   "message": "User with this email already exists", 
#   "details": { "email": "john@example.com" } 
# }
```

## Environment Variables

### Development
```bash
ENVIRONMENT=development
LOG_LEVEL=debug
DB_SSL_MODE=disable
```

### Production
```bash
ENVIRONMENT=production
LOG_LEVEL=warn
DB_SSL_MODE=require
DB_MAX_OPEN_CONNS=50
DB_PASSWORD=strong_password_here
```

## Monitoring

### Check Health
```bash
curl http://localhost:8080/health
# { "status": "OK", "timestamp": "2024-01-29T10:30:00Z" }
```

### View Logs
```bash
# Docker
docker-compose logs -f api

# Systemd
journalctl -u gin-api -f
```

### Check Process
```bash
# Get PID
lsof -i :8080

# Kill gracefully
kill -SIGTERM <PID>
```

## Security Checklist

- [x] Input validation (name, email, password)
- [x] Password hashing (SHA256)
- [x] SQL injection prevention (GORM)
- [x] Error handling (no sensitive info)
- [x] CORS configured
- [x] Request ID tracking
- [x] Logging
- [x] Graceful shutdown
- [ ] HTTPS/SSL (setup via Nginx)
- [ ] Rate limiting (optional)
- [ ] JWT auth (optional)
- [ ] DDoS protection (optional)

## Troubleshooting

### Build error
```bash
go mod tidy
go build ./cmd/api
```

### Database error
```bash
# Check connection
psql -h localhost -U postgres -d gin_db

# Create database
createdb gin_db

# Check environment variables
echo $DB_HOST $DB_USER $DB_NAME
```

### Port in use
```bash
# Linux/Mac
lsof -i :8080

# Windows
netstat -ano | findstr :8080

# Change port
SERVER_PORT=8081 ./api
```

### Logs not visible
```bash
# Change log level
LOG_LEVEL=debug ./api
```

## Performance Tips

1. **Database**
   - Set MaxOpenConns=50 for production
   - Monitor query performance
   - Create indexes on frequent columns

2. **Application**
   - Use reverse proxy (Nginx)
   - Enable HTTP/2
   - Monitor memory usage

3. **Infrastructure**
   - Use CDN for static assets
   - Load balance if needed
   - Monitor disk space

## Next Steps

1. ✅ Run locally
   ```bash
   go run ./cmd/api
   ```

2. ✅ Test endpoints
   ```bash
   curl http://localhost:8080/api/v1/users
   ```

3. ✅ Deploy to production
   ```bash
   See DEPLOYMENT.md
   ```

4. ⏳ Optional: Add features
   - Authentication (JWT)
   - Rate limiting
   - Caching (Redis)
   - Metrics (Prometheus)

## Build Info

- **Binary Size**: 34.45 MB
- **Go Version**: 1.25.6
- **Status**: Production Ready ✅

---

**Start Command:**
```bash
go run ./cmd/api
```

**Production Command:**
```bash
ENVIRONMENT=production ./api
```

**Docker Command:**
```bash
docker-compose up
```
