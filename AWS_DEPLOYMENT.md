# AWS Deployment Guide

## Overview

This guide covers deploying the Gin API to AWS using:
- **API Gateway** - RESTful API endpoint with throttling and WAF protection
- **Application Load Balancer** - Load balancing across ECS tasks
- **ECS Fargate** - Containerized application hosting
- **RDS PostgreSQL** - Managed database service
- **ElastiCache Redis** - Managed caching layer
- **VPC** - Private network infrastructure
- **CloudWatch** - Logging and monitoring
- **Secrets Manager** - Secure credential storage
- **ECR** - Container registry

## Architecture

```
Internet → Route53 → API Gateway (+ WAF) → VPC Link → ALB → ECS Fargate → RDS + ElastiCache
```

## Prerequisites

- AWS CLI installed and configured
- Terraform >= 1.0
- Docker
- AWS Account with appropriate permissions
- Domain name and ACM certificate (for custom domain)

## Setup

### 1. Configure AWS CLI

```bash
aws configure
# Enter: Access Key ID, Secret Access Key, Region, Output format
```

### 2. Create S3 Bucket for Terraform State

```bash
aws s3api create-bucket \
  --bucket gin-api-terraform-state \
  --region us-east-1

aws s3api put-bucket-versioning \
  --bucket gin-api-terraform-state \
  --versioning-configuration Status=Enabled

aws s3api put-bucket-encryption \
  --bucket gin-api-terraform-state \
  --server-side-encryption-configuration '{
    "Rules": [{
      "ApplyServerSideEncryptionByDefault": {
        "SSEAlgorithm": "AES256"
      }
    }]
  }'
```

### 3. Create DynamoDB Table for State Locking

```bash
aws dynamodb create-table \
  --table-name gin-api-terraform-locks \
  --attribute-definitions AttributeName=LockID,AttributeType=S \
  --key-schema AttributeName=LockID,KeyType=HASH \
  --billing-mode PAY_PER_REQUEST \
  --region us-east-1
```

### 4. Configure Variables

Create `terraform/terraform.tfvars`:

```hcl
aws_region     = "us-east-1"
environment    = "production"
project_name   = "gin-api"
api_stage_name = "prod"

# Database
db_instance_class = "db.t3.micro"
db_name           = "gin_db"
db_username       = "postgres"
db_password       = "YourSecurePassword123!"  # Change this!

# Redis
redis_node_type = "cache.t3.micro"

# ECS
ecs_task_cpu      = 256
ecs_task_memory   = 512
ecs_desired_count = 2

# API Gateway
enable_api_gateway       = true
enable_waf               = true
api_throttle_rate_limit  = 1000
api_throttle_burst_limit = 2000

# Optional: Custom Domain
# domain_name     = "api.example.com"
# certificate_arn = "arn:aws:acm:us-east-1:123456789012:certificate/..."
```

## Deployment Steps

### 1. Initialize Terraform

```bash
cd terraform
terraform init
```

### 2. Review Plan

```bash
terraform plan
```

### 3. Apply Infrastructure

```bash
terraform apply
```

This will create:
- VPC with public/private subnets
- NAT Gateways
- Application Load Balancer
- ECS Cluster and Task Definition
- RDS PostgreSQL instance
- ElastiCache Redis cluster
- API Gateway with VPC Link
- Security Groups
- IAM Roles
- CloudWatch Log Groups
- ECR Repository
- WAF Web ACL (if enabled)

### 4. Build and Push Docker Image

```bash
# Get ECR login
aws ecr get-login-password --region us-east-1 | \
  docker login --username AWS --password-stdin $(terraform output -raw ecr_repository_url)

# Build image
docker build -t gin-api:latest .

# Tag image
docker tag gin-api:latest $(terraform output -raw ecr_repository_url):latest

# Push to ECR
docker push $(terraform output -raw ecr_repository_url):latest
```

### 5. Deploy ECS Service

The ECS service will automatically pull the image and start tasks:

```bash
# Watch deployment
aws ecs describe-services \
  --cluster $(terraform output -raw ecs_cluster_name) \
  --services $(terraform output -raw ecs_service_name) \
  --region us-east-1
```

### 6. Initialize Database

```bash
# Get RDS endpoint
RDS_ENDPOINT=$(terraform output -raw rds_endpoint)

# Connect via ECS task (or bastion host)
# Run migrations or initialization scripts
```

## Access the API

### Via API Gateway

```bash
# Get API Gateway URL
API_URL=$(terraform output -raw api_gateway_url)

# Test health endpoint
curl $API_URL/health

# Test user endpoints
curl $API_URL/api/v1/users
```

### Via Custom Domain (if configured)

```bash
# Update Route53 to point to API Gateway
# Then access via custom domain
curl https://api.example.com/health
```

## Monitoring

### CloudWatch Logs

```bash
# View ECS logs
aws logs tail /ecs/gin-api --follow --region us-east-1

# View API Gateway logs
aws logs tail /aws/apigateway/gin-api --follow --region us-east-1
```

### CloudWatch Metrics

Access CloudWatch Console:
- ECS metrics: CPU, Memory, Task count
- API Gateway metrics: Request count, Latency, 4XX/5XX errors
- RDS metrics: CPU, Connections, Storage
- ElastiCache metrics: CPU, Memory, Connections

### WAF Metrics

```bash
# View blocked requests
aws wafv2 get-sampled-requests \
  --web-acl-arn $(terraform output -raw waf_web_acl_id) \
  --rule-metric-name RateLimitRule \
  --scope REGIONAL \
  --time-window StartTime=1234567890,EndTime=1234567900
```

## Scaling

### Manual Scaling

```bash
# Scale ECS service
aws ecs update-service \
  --cluster $(terraform output -raw ecs_cluster_name) \
  --service $(terraform output -raw ecs_service_name) \
  --desired-count 5 \
  --region us-east-1
```

### Auto Scaling

Auto scaling is already configured:
- **CPU**: Scales when CPU > 70%
- **Memory**: Scales when Memory > 80%
- **Min**: 2 tasks (configurable)
- **Max**: 10 tasks

## Security Best Practices

### 1. Secrets Management

Never hardcode secrets. Use Secrets Manager:

```bash
# Update database password
aws secretsmanager update-secret \
  --secret-id gin-api-db-credentials \
  --secret-string '{"username":"postgres","password":"NewPassword123!"}'
```

### 2. IAM Policies

Review and tighten IAM policies in production. Follow principle of least privilege.

### 3. Network Security

- All database and cache resources are in private subnets
- Security groups restrict traffic between resources
- NAT Gateways for outbound internet access

### 4. API Gateway Security

- **WAF**: Protects against common web exploits
- **Rate Limiting**: 1000 requests/second (configurable)
- **API Keys**: Can be added for authentication
- **Usage Plans**: Control API access per client

### 5. SSL/TLS

- Use ACM certificates for custom domains
- Enforce HTTPS in API Gateway and ALB

## Cost Optimization

### Development Environment

Reduce costs for dev/staging:

```hcl
# terraform/terraform.tfvars
db_instance_class = "db.t3.micro"
redis_node_type   = "cache.t3.micro"
ecs_task_cpu      = 256
ecs_task_memory   = 512
ecs_desired_count = 1
multi_az          = false
```

### Production Environment

For production, consider:

```hcl
db_instance_class = "db.t3.small"  # or larger
redis_node_type   = "cache.t3.small"
ecs_task_cpu      = 512
ecs_task_memory   = 1024
ecs_desired_count = 3
multi_az          = true  # High availability
```

### Cost Estimates

Monthly costs (approximate):
- **ECS Fargate** (2 tasks, 0.25 vCPU, 0.5 GB): ~$15
- **RDS db.t3.micro** (20 GB storage): ~$25
- **ElastiCache cache.t3.micro**: ~$12
- **Application Load Balancer**: ~$20
- **API Gateway** (1M requests): ~$3.50
- **NAT Gateway** (2 AZs): ~$64
- **Data Transfer**: Variable
- **Total**: ~$140-$200/month

## Backup and Restore

### RDS Backups

Automated backups are enabled (7 day retention):

```bash
# Create manual snapshot
aws rds create-db-snapshot \
  --db-instance-identifier gin-api-db \
  --db-snapshot-identifier gin-api-snapshot-$(date +%Y%m%d)

# Restore from snapshot
aws rds restore-db-instance-from-db-snapshot \
  --db-instance-identifier gin-api-db-restored \
  --db-snapshot-identifier gin-api-snapshot-20260129
```

### Redis Backups

Automated snapshots enabled (5 day retention):

```bash
# Create manual backup
aws elasticache create-snapshot \
  --cache-cluster-id gin-api-redis \
  --snapshot-name gin-api-redis-snapshot-$(date +%Y%m%d)
```

## Disaster Recovery

### Multi-Region Setup

For disaster recovery, deploy to multiple regions:

1. Use separate Terraform workspaces per region
2. Replicate RDS using read replicas
3. Configure Route53 health checks and failover
4. Use S3 for cross-region replication

### Rollback Procedure

```bash
# Rollback to previous task definition
aws ecs update-service \
  --cluster gin-api-cluster \
  --service gin-api-service \
  --task-definition gin-api:PREVIOUS_REVISION
```

## CI/CD Integration

### GitHub Actions Example

```yaml
# .github/workflows/deploy-aws.yml
name: Deploy to AWS

on:
  push:
    branches: [main]

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Configure AWS credentials
        uses: aws-actions/configure-aws-credentials@v2
        with:
          aws-access-key-id: ${{ secrets.AWS_ACCESS_KEY_ID }}
          aws-secret-access-key: ${{ secrets.AWS_SECRET_ACCESS_KEY }}
          aws-region: us-east-1
      
      - name: Login to ECR
        id: login-ecr
        uses: aws-actions/amazon-ecr-login@v1
      
      - name: Build and push image
        env:
          ECR_REGISTRY: ${{ steps.login-ecr.outputs.registry }}
          ECR_REPOSITORY: gin-api
          IMAGE_TAG: ${{ github.sha }}
        run: |
          docker build -t $ECR_REGISTRY/$ECR_REPOSITORY:$IMAGE_TAG .
          docker push $ECR_REGISTRY/$ECR_REPOSITORY:$IMAGE_TAG
      
      - name: Update ECS service
        run: |
          aws ecs update-service \
            --cluster gin-api-cluster \
            --service gin-api-service \
            --force-new-deployment
```

## Troubleshooting

### ECS Tasks Not Starting

```bash
# Check task status
aws ecs list-tasks --cluster gin-api-cluster

# Describe task for errors
aws ecs describe-tasks \
  --cluster gin-api-cluster \
  --tasks <task-arn>

# Check CloudWatch logs
aws logs tail /ecs/gin-api --follow
```

### API Gateway 502/503 Errors

- Check ALB target health
- Verify VPC Link status
- Check ECS task health checks
- Review CloudWatch logs

### Database Connection Issues

- Verify security group rules
- Check RDS instance status
- Verify credentials in Secrets Manager
- Test from ECS task

### High Costs

```bash
# Check resource usage
aws ce get-cost-and-usage \
  --time-period Start=2026-01-01,End=2026-01-31 \
  --granularity MONTHLY \
  --metrics BlendedCost \
  --group-by Type=SERVICE
```

## Cleanup

To destroy all resources:

```bash
cd terraform
terraform destroy
```

**Warning**: This will delete:
- All infrastructure
- Database (unless deletion protection is enabled)
- All data stored in RDS and ElastiCache
- Container images in ECR

## Additional Resources

- [AWS API Gateway Documentation](https://docs.aws.amazon.com/apigateway/)
- [AWS ECS Fargate Documentation](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/AWS_Fargate.html)
- [AWS RDS Documentation](https://docs.aws.amazon.com/rds/)
- [Terraform AWS Provider](https://registry.terraform.io/providers/hashicorp/aws/latest/docs)

## Support

For issues with:
- **Infrastructure**: Check Terraform documentation
- **Application**: Check application logs in CloudWatch
- **AWS Services**: Consult AWS Support or documentation
