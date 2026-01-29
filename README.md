# Production-Ready Go API with Gin

Complete production-ready REST API with advanced features including:
- ✅ Clean Architecture (Handler → Service → Repository)
- ✅ PostgreSQL with GORM (indexes, transactions, soft deletes, read replicas)
- ✅ Redis Caching
- ✅ Kafka Messaging
- ✅ Docker & Docker Compose
- ✅ Kubernetes Manifests
- ✅ Nginx Reverse Proxy
- ✅ AWS Deployment (Terraform, ECS, API Gateway, RDS, ElastiCache)
- ✅ Comprehensive Testing (Unit, Integration, E2E)
- ✅ Resilience Patterns (Circuit Breaker, Retry, Rate Limiting, Timeout)
- ✅ Feature Flags & A/B Testing
- ✅ S3/MinIO File Storage
- ✅ Background Workers & CronJobs
- ✅ Service Discovery (Kubernetes & Consul)
- ✅ Database Backups & Restore Scripts

## Project Structure

```
gin-demo/
├── cmd/
│   └── api/
│       └── main.go                 # Application entry point
├── internal/
│   ├── app/
│   │   ├── middleware.go           # CORS, Logging, Recovery
│   │   ├── route.go                # Route definitions
│   │   └── server.go               # HTTP server setup
│   ├── config/
│   │   └── config.go               # Configuration management
│   ├── database/
│   │   ├── postgres.go             # Database connection
│   │   └── transaction.go          # Transaction helpers
│   ├── domain/
│   │   └── user.go                 # User entity with indexes
│   ├── handler/
│   │   ├── user_handler.go         # User HTTP handlers
│   │   ├── feature_flag_handler.go # Feature flag API
│   │   └── file_upload_handler.go  # File upload API
│   ├── repository/
│   │   └── user_repository.go      # Data access layer
│   ├── service/
│   │   └── user_service.go         # Business logic
│   └── worker/
│       └── worker.go               # Background job workers
├── pkg/
│   ├── cache/
│   │   └── redis.go                # Redis client wrapper
│   ├── discovery/
│   │   └── consul.go               # Consul service discovery
│   ├── feature/
│   │   └── feature_flag.go         # Feature flags & A/B testing
│   ├── logger/
│   │   └── logger.go               # Structured logging
│   ├── messaging/
│   │   └── kafka.go                # Kafka producer/consumer
│   ├── middleware/
│   │   └── rate_limit.go           # Rate limiting middleware
│   ├── resilience/
│   │   └── circuit_breaker.go      # Resilience patterns
│   ├── storage/
│   │   └── s3.go                   # S3/MinIO storage
│   └── validator/
│       └── user_validator.go       # Custom validators
├── tests/
│   ├── integration/
│   │   └── api_test.go             # Integration tests
│   └── e2e/
│       └── api_e2e_test.go         # End-to-end tests
├── k8s/
│   ├── namespace.yaml              # Kubernetes namespace
│   ├── configmap.yaml              # Configuration
│   ├── secrets.yaml                # Secrets
│   ├── postgres-pvc.yaml           # Persistent volumes
│   ├── postgres-deployment.yaml    # Postgres deployment
│   ├── redis-deployment.yaml       # Redis deployment
│   ├── kafka-deployment.yaml       # Kafka deployment
│   ├── api-deployment.yaml         # API deployment
│   ├── ingress.yaml                # Ingress controller
│   ├── cronjob.yaml                # Scheduled jobs
│   └── service-discovery.yaml      # Service discovery
├── terraform/
│   └── aws/
│       ├── main.tf                 # Terraform main config
│       ├── variables.tf            # Variables
│       ├── vpc.tf                  # VPC setup
│       ├── alb.tf                  # Application Load Balancer
│       ├── api_gateway.tf          # API Gateway
│       ├── ecs.tf                  # ECS Fargate
│       ├── rds_elasticache.tf      # RDS & ElastiCache
│       └── outputs.tf              # Terraform outputs
├── docker/
│   ├── Dockerfile                  # Multi-stage Docker build
│   ├── docker-compose.yml          # Full stack
│   ├── docker-compose.minio.yml    # MinIO service
│   └── nginx/
│       ├── nginx.conf              # Nginx config
│       ├── api.conf                # API upstream
│       └── generate-ssl.sh         # SSL certificate generation
├── scripts/
│   ├── backup-db.sh                # Database backup
│   └── restore-db.sh               # Database restore
├── docs/
│   ├── API.md                      # API documentation
│   ├── DEPLOYMENT.md               # Deployment guide
│   ├── ADVANCED_FEATURES.md        # Advanced features setup
│   └── AWS_DEPLOYMENT.md           # AWS deployment guide
├── go.mod
├── go.sum
├── .env.example
├── .env.storage
└── README.md
```

## Quick Start

### Local Development

```bash
# 1. Clone and setup
git clone <repository>
cd gin-demo

# 2. Install dependencies
go mod download

# 3. Set up environment
cp .env.example .env
# Edit .env with your configuration

# 4. Start infrastructure
docker-compose up -d postgres redis kafka

# 5. Run migrations (auto-migration enabled)
go run cmd/api/main.go

# 6. Access API
curl http://localhost:8080/health
```

### With Docker Compose (Full Stack)

```bash
# Start all services
docker-compose up -d

# With MinIO
docker-compose -f docker-compose.yml -f docker-compose.minio.yml up -d

# View logs
docker-compose logs -f api

# Stop all services
docker-compose down
```

### With Kubernetes

```bash
# Apply all manifests
kubectl apply -f k8s/

# Check status
kubectl get all -n gin-demo

# Port forward
kubectl port-forward -n gin-demo svc/gin-demo-api-service 8080:8080

# View logs
kubectl logs -n gin-demo -l app=gin-demo-api -f
```

## Features

### 1. Testing

#### Unit Tests
```bash
go test ./internal/service/... -v
go test ./pkg/validator/... -v
```

#### Integration Tests
```bash
export DB_NAME=gin_db_test
go test ./tests/integration/... -v
```

#### E2E Tests
```bash
# Start server first
go run cmd/api/main.go

# Run tests
go test ./tests/e2e/... -v
```

### 2. Database Features

- **Indexes**: Optimized queries (email, name, timestamps)
- **Transactions**: WithTransaction wrapper with auto-rollback
- **Soft Deletes**: DeletedAt field with GORM
- **Backups**: Automated scripts with S3 upload
- **Read Replicas**: Configurable via `DB_READ_REPLICAS`

### 3. Resilience Patterns

- **Circuit Breaker**: 3-state (Closed/Open/HalfOpen)
- **Retry**: Exponential backoff
- **Rate Limiting**: Token bucket (global + per-IP)
- **Timeout**: Context-aware operations

### 4. Feature Management

- **Feature Flags**: Enable/disable features dynamically
- **Targeting Rules**: User-based targeting
- **Rollout Percentage**: Gradual feature rollout
- **A/B Testing**: Consistent variant assignment

### 5. Storage

- **S3**: AWS S3 integration
- **MinIO**: S3-compatible local storage
- **Presigned URLs**: Temporary access URLs
- **CDN**: CloudFront integration

### 6. Background Jobs

- **Worker Pool**: Concurrent job processing
- **Job Handlers**: Email, data processing, cleanup
- **Kubernetes CronJobs**: Scheduled tasks
- **Job Queue**: Buffered channel with backpressure

### 7. Service Discovery

- **Kubernetes**: Built-in service DNS
- **Consul**: Dynamic service registration
- **Health Checks**: Automatic health monitoring
- **KV Store**: Configuration management

## API Endpoints

### Users
```
GET    /api/v1/users          - List all users
GET    /api/v1/users/:id      - Get user by ID
POST   /api/v1/users          - Create user
PUT    /api/v1/users/:id      - Update user
DELETE /api/v1/users/:id      - Delete user
```

### Feature Flags
```
GET    /api/v1/feature-flags          - List flags
GET    /api/v1/feature-flags/:key     - Get flag
GET    /api/v1/feature-flags/:key/check - Check if enabled
POST   /api/v1/feature-flags          - Create flag
PUT    /api/v1/feature-flags/:key     - Update flag
```

### File Upload
```
POST   /api/v1/files/upload           - Upload file
GET    /api/v1/files/presigned-url    - Get presigned URL
GET    /api/v1/files                  - List files
DELETE /api/v1/files/:key             - Delete file
```

## Configuration

See `.env.example` and `.env.storage` for all configuration options.

### Key Environment Variables

```bash
# Server
SERVER_HOST=0.0.0.0
SERVER_PORT=8080

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=gin_db
DB_READ_REPLICAS=replica1:5432,replica2:5432

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379

# Kafka
KAFKA_BROKERS=localhost:9092

# Storage
STORAGE_TYPE=s3
STORAGE_BUCKET=gin-demo-uploads
STORAGE_CDN_DOMAIN=https://cdn.example.com

# Worker
WORKER_ENABLED=true
WORKER_COUNT=5

# Service Discovery
DISCOVERY_ENABLED=true
DISCOVERY_TYPE=consul
CONSUL_ADDR=localhost:8500
```

## Deployment

### AWS (Terraform)

```bash
cd terraform/aws

# Initialize
terraform init

# Plan
terraform plan

# Apply
terraform apply

# Get outputs
terraform output
```

See [docs/AWS_DEPLOYMENT.md](docs/AWS_DEPLOYMENT.md) for detailed instructions.

### Kubernetes

```bash
# Apply manifests
kubectl apply -f k8s/

# Apply CronJobs
kubectl apply -f k8s/cronjob.yaml

# Apply service discovery
kubectl apply -f k8s/service-discovery.yaml
```

## Monitoring & Health

### Health Check
```bash
curl http://localhost:8080/health
```

### Metrics
- Circuit breaker state
- Rate limit counters
- Worker pool status
- Cache hit/miss rates

### Logs
Structured JSON logging with contextual information.

## Dependencies

```
github.com/gin-gonic/gin          v1.11.0
gorm.io/gorm                      v1.25.7
gorm.io/driver/postgres           v1.5.7
github.com/redis/go-redis/v9      v9.5.1
github.com/IBM/sarama             v1.43.0
github.com/aws/aws-sdk-go-v2      v1.32.7
github.com/hashicorp/consul/api   v1.30.0
golang.org/x/time                 v0.9.0
```

## Contributing

1. Fork the repository
2. Create feature branch (`git checkout -b feature/amazing-feature`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing-feature`)
5. Open Pull Request

## License

MIT License

## Support

For detailed documentation, see:
- [API Documentation](docs/API.md)
- [Deployment Guide](docs/DEPLOYMENT.md)
- [Advanced Features](docs/ADVANCED_FEATURES.md)
- [AWS Deployment](docs/AWS_DEPLOYMENT.md)
