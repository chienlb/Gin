#!/bin/bash

# AWS Deployment Script for Gin API

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}🚀 AWS Deployment Script${NC}\n"

# Check prerequisites
command -v aws >/dev/null 2>&1 || { echo -e "${RED}❌ AWS CLI is not installed${NC}"; exit 1; }
command -v terraform >/dev/null 2>&1 || { echo -e "${RED}❌ Terraform is not installed${NC}"; exit 1; }
command -v docker >/dev/null 2>&1 || { echo -e "${RED}❌ Docker is not installed${NC}"; exit 1; }

echo -e "${GREEN}✓ Prerequisites check passed${NC}\n"

# Variables
PROJECT_NAME="gin-api"
AWS_REGION=${AWS_REGION:-us-east-1}
ENVIRONMENT=${ENVIRONMENT:-production}

# Step 1: Initialize Terraform
echo -e "${YELLOW}Step 1: Initializing Terraform...${NC}"
cd terraform
terraform init
echo -e "${GREEN}✓ Terraform initialized${NC}\n"

# Step 2: Plan infrastructure
echo -e "${YELLOW}Step 2: Planning infrastructure...${NC}"
terraform plan -out=tfplan
echo -e "${GREEN}✓ Terraform plan created${NC}\n"

# Confirm deployment
read -p "Do you want to apply this plan? (yes/no): " confirm
if [ "$confirm" != "yes" ]; then
    echo -e "${YELLOW}Deployment cancelled${NC}"
    exit 0
fi

# Step 3: Apply infrastructure
echo -e "${YELLOW}Step 3: Applying infrastructure...${NC}"
terraform apply tfplan
rm tfplan
echo -e "${GREEN}✓ Infrastructure deployed${NC}\n"

# Get outputs
ECR_REPO=$(terraform output -raw ecr_repository_url)
ECS_CLUSTER=$(terraform output -raw ecs_cluster_name)
ECS_SERVICE=$(terraform output -raw ecs_service_name)

echo -e "${GREEN}ECR Repository: $ECR_REPO${NC}"
echo -e "${GREEN}ECS Cluster: $ECS_CLUSTER${NC}"
echo -e "${GREEN}ECS Service: $ECS_SERVICE${NC}\n"

# Step 4: Build Docker image
echo -e "${YELLOW}Step 4: Building Docker image...${NC}"
cd ..
docker build -t $PROJECT_NAME:latest .
echo -e "${GREEN}✓ Docker image built${NC}\n"

# Step 5: Login to ECR
echo -e "${YELLOW}Step 5: Logging in to ECR...${NC}"
aws ecr get-login-password --region $AWS_REGION | \
  docker login --username AWS --password-stdin $ECR_REPO
echo -e "${GREEN}✓ Logged in to ECR${NC}\n"

# Step 6: Tag and push image
echo -e "${YELLOW}Step 6: Pushing image to ECR...${NC}"
docker tag $PROJECT_NAME:latest $ECR_REPO:latest
docker tag $PROJECT_NAME:latest $ECR_REPO:$(git rev-parse --short HEAD 2>/dev/null || echo "manual")
docker push $ECR_REPO:latest
docker push $ECR_REPO:$(git rev-parse --short HEAD 2>/dev/null || echo "manual")
echo -e "${GREEN}✓ Image pushed to ECR${NC}\n"

# Step 7: Force new deployment
echo -e "${YELLOW}Step 7: Updating ECS service...${NC}"
aws ecs update-service \
  --cluster $ECS_CLUSTER \
  --service $ECS_SERVICE \
  --force-new-deployment \
  --region $AWS_REGION > /dev/null
echo -e "${GREEN}✓ ECS service updated${NC}\n"

# Step 8: Wait for deployment
echo -e "${YELLOW}Step 8: Waiting for deployment to complete...${NC}"
aws ecs wait services-stable \
  --cluster $ECS_CLUSTER \
  --services $ECS_SERVICE \
  --region $AWS_REGION

echo -e "${GREEN}✓ Deployment completed successfully!${NC}\n"

# Show API Gateway URL
cd terraform
API_URL=$(terraform output -raw api_gateway_url)
echo -e "${GREEN}🎉 Deployment Complete!${NC}"
echo -e "${GREEN}API Gateway URL: $API_URL${NC}"
echo -e "${GREEN}Test: curl $API_URL/health${NC}\n"

cd ..
