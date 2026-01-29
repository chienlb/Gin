# Production Deployment Guide

## Triển khai Production

### 1. Prerequisites

- Go 1.25.6+
- PostgreSQL 15+
- Docker & Docker Compose (recommended)
- Nginx/Reverse Proxy
- SSL Certificate (for HTTPS)

### 2. Database Setup

```bash
# Create production database
createdb gin_db_prod

# Create production user
psql -c "CREATE USER prod_user WITH PASSWORD 'your_strong_password';"
psql -c "GRANT ALL PRIVILEGES ON DATABASE gin_db_prod TO prod_user;"
```

### 3. Application Build

```bash
# Build binary
go build -o /usr/local/bin/gin-api ./cmd/api

# Or use Docker
docker build -t gin-api:1.0.0 .
docker tag gin-api:1.0.0 your-registry/gin-api:1.0.0
docker push your-registry/gin-api:1.0.0
```

### 4. Systemd Service (Linux)

Create `/etc/systemd/system/gin-api.service`:

```ini
[Unit]
Description=Gin Demo API
After=network.target postgresql.service

[Service]
Type=simple
User=api
ExecStart=/usr/local/bin/gin-api
Restart=on-failure
RestartSec=5s

Environment="ENVIRONMENT=production"
Environment="DB_HOST=localhost"
Environment="DB_USER=prod_user"
Environment="DB_PASSWORD=your_strong_password"
Environment="DB_NAME=gin_db_prod"
Environment="DB_SSL_MODE=disable"
Environment="LOG_LEVEL=warn"

[Install]
WantedBy=multi-user.target
```

Start service:
```bash
sudo systemctl daemon-reload
sudo systemctl enable gin-api
sudo systemctl start gin-api
sudo systemctl status gin-api
```

### 5. Nginx Reverse Proxy

Create `/etc/nginx/sites-available/gin-api`:

```nginx
upstream gin_api {
    server localhost:8080;
    keepalive 32;
}

server {
    listen 80;
    server_name api.yourdomain.com;
    
    # Redirect to HTTPS
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name api.yourdomain.com;
    
    # SSL Configuration
    ssl_certificate /etc/ssl/certs/your-cert.pem;
    ssl_certificate_key /etc/ssl/private/your-key.pem;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;
    
    # Security Headers
    add_header Strict-Transport-Security "max-age=31536000" always;
    add_header X-Frame-Options "DENY" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header X-XSS-Protection "1; mode=block" always;
    
    # Logging
    access_log /var/log/nginx/gin-api-access.log;
    error_log /var/log/nginx/gin-api-error.log;
    
    # API Proxy
    location / {
        proxy_pass http://gin_api;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_cache_bypass $http_upgrade;
        
        # Timeouts
        proxy_connect_timeout 15s;
        proxy_send_timeout 15s;
        proxy_read_timeout 15s;
    }
}
```

Enable site:
```bash
sudo ln -s /etc/nginx/sites-available/gin-api /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl restart nginx
```

### 6. Docker Compose Production

Create `docker-compose.prod.yml`:

```yaml
version: '3.8'

services:
  postgres:
    image: postgres:15-alpine
    container_name: gin_postgres_prod
    environment:
      POSTGRES_USER: prod_user
      POSTGRES_PASSWORD: ${DB_PASSWORD}
      POSTGRES_DB: gin_db_prod
    volumes:
      - postgres_data_prod:/var/lib/postgresql/data
    networks:
      - gin_network
    restart: always

  api:
    build:
      context: .
      dockerfile: Dockerfile
    container_name: gin_api_prod
    environment:
      ENVIRONMENT: production
      SERVER_HOST: 0.0.0.0
      SERVER_PORT: 8080
      DB_HOST: postgres
      DB_USER: prod_user
      DB_PASSWORD: ${DB_PASSWORD}
      DB_NAME: gin_db_prod
      DB_SSL_MODE: disable
      LOG_LEVEL: warn
    ports:
      - "8080:8080"
    depends_on:
      - postgres
    networks:
      - gin_network
    restart: always
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 10s

volumes:
  postgres_data_prod:

networks:
  gin_network:
    driver: bridge
```

Deploy:
```bash
DB_PASSWORD=your_strong_password docker-compose -f docker-compose.prod.yml up -d
```

### 7. Monitoring

#### Health Check Script

Create `monitoring/health_check.sh`:

```bash
#!/bin/bash

API_URL="https://api.yourdomain.com/health"
CHECK_INTERVAL=30

while true; do
    RESPONSE=$(curl -s -w "\n%{http_code}" "$API_URL")
    HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
    
    if [ "$HTTP_CODE" = "200" ]; then
        echo "[$(date)] ✅ API is healthy"
    else
        echo "[$(date)] ❌ API is down (HTTP $HTTP_CODE)"
        # Send alert (email, Slack, PagerDuty, etc.)
    fi
    
    sleep $CHECK_INTERVAL
done
```

#### Logging

View logs:
```bash
# Systemd
journalctl -u gin-api -f

# Docker
docker-compose logs -f api
```

#### Metrics

Monitor with Prometheus/Grafana (future enhancement)

### 8. Backup Strategy

#### Database Backup

```bash
# Daily backup
0 2 * * * pg_dump -U prod_user gin_db_prod | gzip > /backup/gin_db_$(date +\%Y\%m\%d).sql.gz

# S3 backup
0 3 * * * aws s3 cp /backup/gin_db_$(date +\%Y\%m\%d).sql.gz s3://your-bucket/backups/
```

#### Restore from Backup

```bash
gunzip < gin_db_20240129.sql.gz | psql -U prod_user gin_db_prod
```

### 9. Security Checklist

- [x] HTTPS/SSL enabled
- [x] Strong database password
- [x] Input validation
- [x] SQL injection prevention
- [x] CORS configured
- [x] Security headers
- [x] Firewall rules (22, 80, 443 only)
- [x] SSH key authentication
- [ ] API Authentication (can add JWT)
- [ ] Rate limiting (can add)
- [ ] WAF rules (can add)
- [ ] DDoS protection (can add)

### 10. Performance Optimization

#### Database
- Set appropriate MaxOpenConns (50-100)
- Monitor query performance
- Create indexes on frequently used columns

#### Application
- Enable HTTPS/HTTP2
- Use compression
- Set appropriate timeouts
- Monitor memory usage

#### Infrastructure
- Use CDN for static assets
- Load balance if needed
- Monitor disk space
- Keep logs rotated

### 11. Rollback Plan

```bash
# Keep previous version
cp /usr/local/bin/gin-api /usr/local/bin/gin-api.backup

# Rollback if needed
cp /usr/local/bin/gin-api.backup /usr/local/bin/gin-api
systemctl restart gin-api
```

### 12. Update Procedure

```bash
# 1. Build new version
go build -o gin-api-new ./cmd/api

# 2. Test new version
./gin-api-new --help

# 3. Backup old version
cp /usr/local/bin/gin-api /usr/local/bin/gin-api.backup

# 4. Deploy
cp gin-api-new /usr/local/bin/gin-api

# 5. Restart
systemctl restart gin-api

# 6. Verify
curl https://api.yourdomain.com/health
```

---

**Deployed**: Ready for production! ✅
