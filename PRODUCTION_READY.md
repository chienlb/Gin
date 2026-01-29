# Production Ready Gin API

Đây là hướng dẫn hoàn chỉnh để triển khai ứng dụng Gin API theo chuẩn production.

## 🚀 Nâng cấp Production

### Các cải thiện đã thực hiện:

#### 1. **Error Handling & Validation** ✅
- `pkg/apperror/` - Hệ thống lỗi tập trung
  - Định nghĩa các error code standard
  - AppError struct với status, code, message, details
  - Helper functions cho validation, not found, duplicate, vv.

- `pkg/validator/` - Validação input
  - ValidateCreateRequest - Kiểm tra đầy đủ dữ liệu tạo user
  - ValidateName - Kiểm tra tên (2-100 ký tự, valid chars)
  - ValidateEmail - Format và độ dài email
  - ValidatePassword - Yêu cầu mật khẩu mạnh (Uppercase, lowercase, digit, 6+ chars)
  - ValidateUpdateRequest - Kiểm tra dữ liệu update

#### 2. **Middleware** ✅
- `pkg/middleware/` - HTTP middleware tập trung
  - **LoggingMiddleware** - Ghi log request/response với thời gian xử lý
  - **CORSMiddleware** - Xử lý CORS cho cross-origin requests
  - **RecoveryMiddleware** - Bắt panic và trả về error response
  - **RequestIDMiddleware** - Thêm request ID cho tracing

#### 3. **Service Layer Enhancement** ✅
- Thay đổi return type từ `error` thành `*apperror.AppError`
- Thêm logging chi tiết cho mọi operation
- Kiểm tra validation trước khi gọi repository
- Xử lý lỗi với AppError thích hợp

#### 4. **Handler Improvement** ✅
- Response format chuẩn với status, code, message, data
- Xử lý AppError từ service layer
- Validation riêng biệt cho request format
- Chi tiết error response với field và reason

#### 5. **Server Configuration** ✅
- **Graceful Shutdown** - Tắt server một cách an toàn
  - Signal handling (SIGINT, SIGTERM)
  - Timeout 5 giây cho shutdown
  - Log các bước shutdown
  
- **Middleware Chain**
  - Request ID tracking
  - Logging tất cả requests
  - CORS support
  - Panic recovery
  
- **API Versioning**
  - Routes: `/api/v1/users`
  - Dễ mở rộng cho v2, v3...
  
- **Timeouts**
  - ReadTimeout: 15s
  - WriteTimeout: 15s
  - IdleTimeout: 60s

#### 6. **Configuration Management** ✅
- **Environment Support**
  - ENVIRONMENT (development, staging, production)
  - Helper functions IsProduction(), IsDevelopment()
  
- **Database Configuration**
  - MaxOpenConns configurable
  - MaxIdleConns configurable
  - ConnMaxLifetime setting
  
- **Server Timeouts**
  - All timeouts configurable via environment

## 📋 Sử dụng Production

### Environment Variables

```bash
# Server
ENVIRONMENT=production          # development|staging|production
SERVER_HOST=0.0.0.0
SERVER_PORT=8080
SERVER_READ_TIMEOUT=15s
SERVER_WRITE_TIMEOUT=15s
SERVER_IDLE_TIMEOUT=60s

# Database
DB_HOST=db.example.com
DB_PORT=5432
DB_USER=prod_user
DB_PASSWORD=strong_password_here
DB_NAME=gin_db_prod
DB_SSL_MODE=require              # enable for production
DB_MAX_OPEN_CONNS=25
DB_MAX_IDLE_CONNS=5
DB_CONN_MAX_LIFETIME=5m

# Logger
LOG_LEVEL=info                   # debug|info|warn|error
```

### Docker Deployment

```bash
# Build
docker build -t gin-api:1.0.0 .

# Run
docker run -d \
  --name gin-api \
  -e ENVIRONMENT=production \
  -e DB_HOST=postgres \
  -e DB_PASSWORD=secure_password \
  -p 8080:8080 \
  gin-api:1.0.0
```

### Docker Compose (Full Stack)

```bash
docker-compose up -d
```

## 🔒 Security Features

### Password Security
- SHA256 hashing
- Validation: Uppercase + Lowercase + Digit + 6+ chars
- Not returned in API responses

### Input Validation
- Email format validation
- Name length and character validation
- SQL injection prevention via GORM parameterized queries
- CORS protection

### Error Handling
- Consistent error format
- No sensitive info in error messages
- Request ID tracking for debugging

## 📊 API Response Format

### Success Response
```json
{
  "status": "success",
  "code": "CREATED",
  "message": "User created successfully",
  "data": {
    "id": 1,
    "name": "John Doe",
    "email": "john@example.com",
    "created_at": "2024-01-29T10:30:00Z",
    "updated_at": "2024-01-29T10:30:00Z"
  }
}
```

### Error Response
```json
{
  "status": "error",
  "code": "VALIDATION_ERROR",
  "message": "Validation failed",
  "details": {
    "field": "password",
    "reason": "password must contain uppercase, lowercase, and digit"
  }
}
```

### Special Error Responses

**Duplicate Email:**
```json
{
  "status": "error",
  "code": "DUPLICATE_EMAIL",
  "message": "User with this email already exists",
  "details": {
    "email": "duplicate@example.com"
  }
}
```

**Not Found:**
```json
{
  "status": "error",
  "code": "NOT_FOUND",
  "message": "User not found",
  "details": {
    "user_id": 999
  }
}
```

## 🔄 Request/Response Flow

```
1. HTTP Request
   ↓
2. Middleware Chain
   - RequestID: Tạo/lấy request ID
   - Logging: Log request
   - CORS: Xử lý cross-origin
   ↓
3. Handler
   - Validation request format (JSON)
   - Extract parameters
   ↓
4. Service Layer
   - Business logic validation (validator)
   - Database operations via repository
   - Error handling với AppError
   ↓
5. Database (GORM)
   - Create/Read/Update/Delete
   ↓
6. Handler Response
   - Format response
   - Set HTTP status
   ↓
7. HTTP Response
```

## 📈 Logging

Tất cả operations được log:

```
[INFO] HTTP Request: POST /api/v1/users
[DEBUG] HTTP Request: POST /api/v1/users
[INFO] HTTP Response: POST /api/v1/users | Status: 201 | Duration: 45ms
[INFO] User created successfully: john@example.com
```

## 🧪 Testing Flow

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

### Validation Errors

**Invalid password (no uppercase):**
```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "password": "password123"
  }'
```
Response: 400 - password must contain uppercase, lowercase, and digit

**Invalid email:**
```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "invalid-email",
    "password": "Password123"
  }'
```
Response: 400 - invalid email format

**Duplicate email:**
```bash
# Create user once, then try again
# Response: 409 - User with this email already exists
```

## 🔍 Monitoring

### Health Check
```bash
curl http://localhost:8080/health
```
Response:
```json
{
  "status": "OK",
  "timestamp": "2024-01-29T10:30:00Z"
}
```

### Root Endpoint
```bash
curl http://localhost:8080/
```
Response:
```json
{
  "name": "Gin Demo API",
  "version": "1.0.0",
  "status": "running"
}
```

## 📝 Production Checklist

- [x] Error handling chuẩn
- [x] Input validation hoàn chỉnh
- [x] Middleware logging
- [x] CORS support
- [x] Graceful shutdown
- [x] Database connection pooling
- [x] Configuration management
- [x] Request ID tracking
- [x] Password hashing
- [x] API versioning
- [x] Timeout settings
- [x] Panic recovery
- [ ] Rate limiting (có thể thêm)
- [ ] Authentication/Authorization (có thể thêm)
- [ ] Caching (có thể thêm)
- [ ] Monitoring/Metrics (có thể thêm)
- [ ] Database backup (production setup)
- [ ] SSL/TLS (reverse proxy)

## 🚀 Deployment

### Local Development
```bash
ENVIRONMENT=development go run ./cmd/api
```

### Staging
```bash
ENVIRONMENT=staging \
DB_SSL_MODE=require \
LOG_LEVEL=info \
./api
```

### Production
```bash
ENVIRONMENT=production \
DB_SSL_MODE=require \
LOG_LEVEL=warn \
DB_MAX_OPEN_CONNS=50 \
./api
```

---

**Trạng thái:** ✅ Sẵn sàng Production

Ứng dụng đã được nâng cấp với tất cả các best practices cho production!
