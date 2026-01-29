# Advanced Features Setup Guide

## Storage (S3/MinIO)

### Using AWS S3

```bash
# Set environment variables
export STORAGE_TYPE=s3
export STORAGE_REGION=us-east-1
export STORAGE_ACCESS_KEY_ID=your-access-key
export STORAGE_SECRET_ACCESS_KEY=your-secret-key
export STORAGE_BUCKET=gin-demo-uploads
```

### Using MinIO (Local Development)

```bash
# Start MinIO with Docker
docker run -d \
  -p 9000:9000 -p 9001:9001 \
  -e MINIO_ROOT_USER=minioadmin \
  -e MINIO_ROOT_PASSWORD=minioadmin123 \
  --name gin-minio \
  minio/minio server /data --console-address ":9001"

# Set environment variables
export STORAGE_TYPE=minio
export STORAGE_ENDPOINT=http://localhost:9000
export STORAGE_REGION=us-east-1
export STORAGE_ACCESS_KEY_ID=minioadmin
export STORAGE_SECRET_ACCESS_KEY=minioadmin123
export STORAGE_BUCKET=gin-demo-uploads
export STORAGE_USE_PATH_STYLE=true

# Create bucket
mc alias set local http://localhost:9000 minioadmin minioadmin123
mc mb local/gin-demo-uploads
```

### Using with Docker Compose

```bash
# Merge MinIO service
docker-compose -f docker-compose.yml -f docker-compose.minio.yml up -d
```

### CDN Integration (CloudFront)

1. Create CloudFront distribution pointing to your S3 bucket
2. Set `STORAGE_CDN_DOMAIN=https://d123456.cloudfront.net`
3. URLs will use CDN domain for public access

### API Endpoints

```bash
# Upload file
curl -X POST http://localhost:8080/api/v1/files/upload \
  -F "file=@/path/to/file.jpg"

# Get presigned URL
curl http://localhost:8080/api/v1/files/presigned-url?key=uploads/123_file.jpg

# List files
curl http://localhost:8080/api/v1/files?prefix=uploads/

# Delete file
curl -X DELETE http://localhost:8080/api/v1/files/uploads/123_file.jpg
```

## Background Workers

### Starting Worker Pool

```go
// In main.go or separate worker service
workerPool := worker.NewWorkerPool(5)

// Register handlers
workerPool.RegisterHandler("email", &worker.EmailJobHandler{})
workerPool.RegisterHandler("data_processing", &worker.DataProcessingJobHandler{})
workerPool.RegisterHandler("user_cleanup", &worker.UserCleanupJobHandler{})

// Start workers
workerPool.Start()
defer workerPool.Stop()
```

### Submitting Jobs

```go
job := &worker.Job{
	ID:        uuid.New().String(),
	Type:      "email",
	Payload: map[string]interface{}{
		"to":      "user@example.com",
		"subject": "Welcome",
		"body":    "Welcome to Gin Demo!",
	},
	CreatedAt: time.Now(),
	Status:    "pending",
}

workerPool.Submit(job)
```

### Kubernetes CronJobs

```bash
# Apply CronJob manifests
kubectl apply -f k8s/cronjob.yaml

# List CronJobs
kubectl get cronjobs -n gin-demo

# View job runs
kubectl get jobs -n gin-demo

# View logs
kubectl logs -n gin-demo job/user-cleanup-job-xxxxx
```

## Service Discovery

### Kubernetes Service Discovery

```bash
# Apply service discovery manifests
kubectl apply -f k8s/service-discovery.yaml

# Services are accessible via DNS:
# - postgres-service.gin-demo.svc.cluster.local:5432
# - redis-service.gin-demo.svc.cluster.local:6379
# - kafka-service.gin-demo.svc.cluster.local:9092
```

### Using Consul

```bash
# Start Consul
docker run -d \
  -p 8500:8500 \
  -p 8600:8600/udp \
  --name=consul \
  consul agent -server -ui -bootstrap-expect=1 -client=0.0.0.0

# Or with Kubernetes
kubectl apply -f k8s/service-discovery.yaml

# Set environment variables
export DISCOVERY_ENABLED=true
export DISCOVERY_TYPE=consul
export CONSUL_ADDR=localhost:8500
```

### Service Registration

```go
// Initialize Consul client
consulClient, err := discovery.NewConsulClient(discovery.ConsulConfig{
	Address:     cfg.Discovery.ConsulAddr,
	ServiceID:   cfg.Discovery.ServiceID,
	ServiceName: cfg.Discovery.ServiceName,
	ServicePort: cfg.Discovery.ServicePort,
	Tags:        cfg.Discovery.Tags,
})

// Register service
consulClient.Register("http://localhost:8080/health")

// Deregister on shutdown
defer consulClient.Deregister()
```

### Service Discovery

```go
// Discover services
address, err := consulClient.GetServiceAddress("gin-demo-api")

// Watch for service changes
consulClient.WatchService("gin-demo-api", func(services []*api.ServiceEntry) {
	log.Printf("Service instances: %d", len(services))
})
```

### Consul KV Store

```go
// Set configuration
consulClient.SetKV("config/database/host", []byte("localhost"))

// Get configuration
value, err := consulClient.GetKV("config/database/host")

// Delete configuration
consulClient.DeleteKV("config/database/host")
```

## Database Read Replicas

### Configuration

```bash
# Set read replica URLs
export DB_READ_REPLICAS=replica1.example.com:5432,replica2.example.com:5432
```

### AWS RDS Read Replicas

1. Create read replica in AWS Console or Terraform
2. Update environment variables with replica endpoints
3. Repository layer will automatically use replicas for read operations

### Terraform Configuration

See `terraform/aws/rds_elasticache.tf` for read replica configuration:

```hcl
resource "aws_db_instance" "postgres_read_replica" {
  replicate_source_db = aws_db_instance.postgres.id
  instance_class      = "db.t3.micro"
  identifier          = "gin-demo-postgres-replica"
  skip_final_snapshot = true
}
```

## Testing

### Unit Tests

```bash
go test ./internal/service/... -v
go test ./pkg/validator/... -v
```

### Integration Tests

```bash
# Set up test database
export DB_NAME=gin_db_test
go test ./tests/integration/... -v
```

### E2E Tests

```bash
# Start API server first
go run cmd/api/main.go

# Run E2E tests
go test ./tests/e2e/... -v
```

## Monitoring & Health Checks

### Health Check Endpoint

```bash
curl http://localhost:8080/health
```

### Kubernetes Health Checks

```yaml
livenessProbe:
  httpGet:
    path: /health
    port: 8080
  initialDelaySeconds: 30
  periodSeconds: 10

readinessProbe:
  httpGet:
    path: /health
    port: 8080
  initialDelaySeconds: 5
  periodSeconds: 5
```

### Circuit Breaker Monitoring

Check circuit breaker state via middleware or custom metrics endpoint.

### Rate Limiting

Global and per-IP rate limiting is automatically applied via middleware.

## Production Checklist

- [ ] Configure S3/MinIO storage
- [ ] Set up CDN (CloudFront/CloudFlare)
- [ ] Enable worker pool for background jobs
- [ ] Deploy Kubernetes CronJobs for scheduled tasks
- [ ] Configure service discovery (Consul or Kubernetes)
- [ ] Set up database read replicas
- [ ] Configure backup automation
- [ ] Enable rate limiting
- [ ] Configure circuit breakers
- [ ] Set up monitoring and alerting
- [ ] Configure feature flags
- [ ] Test A/B testing implementation
- [ ] Review security settings (SSL, authentication)
- [ ] Load test with production-like data
