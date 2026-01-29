# Infrastructure Setup Guide

## Prerequisites

- Docker & Docker Compose
- Kubernetes cluster (minikube, kind, or cloud provider)
- kubectl CLI
- Helm (optional, for package management)

## Docker Compose Deployment

### Quick Start

1. **Generate SSL certificates** (for development):
```bash
cd nginx
bash generate-ssl.sh
cd ..
```

2. **Start all services**:
```bash
docker-compose up -d
```

3. **Check services status**:
```bash
docker-compose ps
```

4. **View logs**:
```bash
docker-compose logs -f api
```

5. **Stop services**:
```bash
docker-compose down
```

### Services

- **API**: http://localhost:8080 (or https://localhost via Nginx)
- **Nginx**: http://localhost:80, https://localhost:443
- **PostgreSQL**: localhost:5432
- **Redis**: localhost:6379
- **Kafka**: localhost:9093 (external), kafka:9092 (internal)
- **Kafka UI**: http://localhost:8090

## Kubernetes Deployment

### Setup Cluster

#### Using Minikube
```bash
minikube start --cpus=4 --memory=8192
minikube addons enable ingress
```

#### Using Kind
```bash
kind create cluster --name gin-api-cluster
```

### Deploy Application

1. **Build Docker image**:
```bash
docker build -t gin-api:latest .
```

2. **Load image into cluster** (for local clusters):
```bash
# Minikube
minikube image load gin-api:latest

# Kind
kind load docker-image gin-api:latest --name gin-api-cluster
```

3. **Deploy to Kubernetes**:
```bash
bash scripts/deploy-k8s.sh
```

Or manually:
```bash
kubectl apply -f k8s/namespace.yaml
kubectl apply -f k8s/secrets.yaml
kubectl apply -f k8s/configmap.yaml
kubectl apply -f k8s/persistent-volumes.yaml
kubectl apply -f k8s/postgres-deployment.yaml
kubectl apply -f k8s/redis-deployment.yaml
kubectl apply -f k8s/kafka-deployment.yaml
kubectl apply -f k8s/api-deployment.yaml
kubectl apply -f k8s/ingress.yaml
```

### Verify Deployment

```bash
# Check all pods
kubectl get pods -n gin-api

# Check services
kubectl get services -n gin-api

# Check ingress
kubectl get ingress -n gin-api

# View API logs
kubectl logs -f deployment/gin-api -n gin-api

# Describe pod (for troubleshooting)
kubectl describe pod <pod-name> -n gin-api
```

### Access Application

#### Port Forwarding
```bash
kubectl port-forward service/gin-api-service 8080:80 -n gin-api
```
Then access: http://localhost:8080

#### Using Ingress (production)
Configure DNS to point to your ingress IP:
```bash
kubectl get ingress -n gin-api
```

## Redis Usage

### Basic Commands

```bash
# Connect to Redis
docker-compose exec redis redis-cli -a redis_password

# Or in Kubernetes
kubectl exec -it deployment/redis -n gin-api -- redis-cli -a redis_password

# Common commands
SET key value
GET key
DEL key
KEYS *
FLUSHDB
```

### Cache Examples

```go
// In your code
ctx := context.Background()

// Cache user data
redisClient.SetUserCache(ctx, userID, user, 1*time.Hour)

// Get cached user
user, err := redisClient.GetUserCache(ctx, userID)

// Delete cache
redisClient.DeleteUserCache(ctx, userID)
```

## Kafka Usage

### Create Topics

```bash
# Docker Compose
docker-compose exec kafka kafka-topics --create \
  --bootstrap-server localhost:9092 \
  --replication-factor 1 \
  --partitions 3 \
  --topic user-events

# List topics
docker-compose exec kafka kafka-topics --list \
  --bootstrap-server localhost:9092

# Kubernetes
kubectl exec -it kafka-0 -n gin-api -- kafka-topics --create \
  --bootstrap-server localhost:9092 \
  --replication-factor 1 \
  --partitions 3 \
  --topic user-events
```

### Monitor with Kafka UI

Access: http://localhost:8090

### Producer/Consumer Examples

```go
// Produce message
event := &messaging.UserEvent{
    Type:      "created",
    UserID:    user.ID,
    Timestamp: time.Now().Unix(),
    Data:      map[string]interface{}{"name": user.Name, "email": user.Email},
}
producer.SendUserEvent(event)

// Consume messages
consumer, _ := messaging.NewKafkaConsumer(
    []string{"kafka:9092"},
    "gin-api-group",
    []string{"user-events"},
    messaging.DefaultUserEventHandler,
)
consumer.Start(ctx)
```

## Nginx Configuration

### SSL/TLS Setup

#### Development (self-signed)
```bash
cd nginx
bash generate-ssl.sh
```

#### Production (Let's Encrypt)
```bash
# Install certbot
apt-get install certbot python3-certbot-nginx

# Get certificate
certbot --nginx -d api.example.com

# Auto-renewal
certbot renew --dry-run
```

### Load Balancing

Edit [nginx/conf.d/api.conf](nginx/conf.d/api.conf):

```nginx
upstream api_backend {
    least_conn;
    server api1:8080 max_fails=3 fail_timeout=30s;
    server api2:8080 max_fails=3 fail_timeout=30s;
    server api3:8080 max_fails=3 fail_timeout=30s;
}
```

### Rate Limiting

Already configured in [nginx/nginx.conf](nginx/nginx.conf):
- API endpoints: 10 requests/second
- Login endpoint: 5 requests/minute

## Monitoring

### Health Checks

```bash
# API health
curl http://localhost:8080/health

# Via Nginx
curl http://localhost/health
```

### Kubernetes Health

```bash
# Check pod health
kubectl get pods -n gin-api

# View events
kubectl get events -n gin-api --sort-by='.lastTimestamp'

# Resource usage
kubectl top pods -n gin-api
kubectl top nodes
```

### Logs

```bash
# Docker Compose
docker-compose logs -f api
docker-compose logs -f nginx
docker-compose logs -f postgres
docker-compose logs -f redis
docker-compose logs -f kafka

# Kubernetes
kubectl logs -f deployment/gin-api -n gin-api
kubectl logs -f deployment/postgres -n gin-api
kubectl logs -f deployment/redis -n gin-api
kubectl logs -f statefulset/kafka -n gin-api
```

## Scaling

### Docker Compose
```bash
docker-compose up -d --scale api=3
```

### Kubernetes (Horizontal Pod Autoscaler)
```bash
# HPA is already configured in api-deployment.yaml
# Check status
kubectl get hpa -n gin-api

# Manual scaling
kubectl scale deployment gin-api --replicas=5 -n gin-api
```

## Backup & Restore

### PostgreSQL

#### Backup
```bash
# Docker Compose
docker-compose exec postgres pg_dump -U postgres gin_db > backup.sql

# Kubernetes
kubectl exec -it deployment/postgres -n gin-api -- \
  pg_dump -U postgres gin_db > backup.sql
```

#### Restore
```bash
# Docker Compose
docker-compose exec -T postgres psql -U postgres gin_db < backup.sql

# Kubernetes
kubectl exec -i deployment/postgres -n gin-api -- \
  psql -U postgres gin_db < backup.sql
```

### Redis

#### Backup
```bash
# Docker Compose
docker-compose exec redis redis-cli -a redis_password BGSAVE
docker cp gin-redis:/data/dump.rdb ./redis-backup.rdb

# Kubernetes
kubectl exec -it deployment/redis -n gin-api -- redis-cli -a redis_password BGSAVE
kubectl cp gin-api/redis-pod:/data/dump.rdb ./redis-backup.rdb
```

## Troubleshooting

### Common Issues

1. **Port already in use**
```bash
# Find process using port
netstat -ano | findstr :8080
# Kill process (Windows)
taskkill /PID <pid> /F
```

2. **Database connection failed**
- Check if PostgreSQL is running
- Verify credentials in environment variables
- Check network connectivity

3. **Redis connection failed**
- Verify Redis password
- Check if Redis is running
- Test connection: `redis-cli -a password ping`

4. **Kafka not starting**
- Ensure Zookeeper is running first
- Check Zookeeper connection
- Increase memory if needed

### Debug Commands

```bash
# Docker Compose
docker-compose ps
docker-compose logs <service>
docker-compose exec <service> sh

# Kubernetes
kubectl get all -n gin-api
kubectl describe pod <pod-name> -n gin-api
kubectl logs <pod-name> -n gin-api
kubectl exec -it <pod-name> -n gin-api -- sh
```

## Environment Variables

### Complete List

```bash
# Server
SERVER_HOST=0.0.0.0
SERVER_PORT=8080
ENVIRONMENT=production

# Database
DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=gin_db
DB_SSL_MODE=disable
DB_MAX_OPEN_CONNS=25
DB_MAX_IDLE_CONNS=5

# Redis
REDIS_HOST=redis
REDIS_PORT=6379
REDIS_PASSWORD=your_redis_password
REDIS_DB=0

# Kafka
KAFKA_BROKERS=kafka:9092,kafka2:9092
KAFKA_CONSUMER_GROUP=gin-api-group
KAFKA_TOPIC_USER_EVENTS=user-events

# Logger
LOG_LEVEL=info
```

## Production Checklist

- [ ] Use strong passwords for PostgreSQL and Redis
- [ ] Enable SSL/TLS with valid certificates
- [ ] Configure firewall rules
- [ ] Set up monitoring (Prometheus/Grafana)
- [ ] Configure backup automation
- [ ] Enable log aggregation
- [ ] Set resource limits
- [ ] Configure auto-scaling
- [ ] Set up alerting
- [ ] Document disaster recovery procedures
