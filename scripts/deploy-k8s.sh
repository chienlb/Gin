#!/bin/bash

# Deploy to Kubernetes cluster

echo "🚀 Deploying Gin API to Kubernetes..."

# Apply namespace
echo "📦 Creating namespace..."
kubectl apply -f k8s/namespace.yaml

# Apply secrets and configmaps
echo "🔐 Creating secrets and configmaps..."
kubectl apply -f k8s/secrets.yaml
kubectl apply -f k8s/configmap.yaml

# Apply persistent volumes
echo "💾 Creating persistent volumes..."
kubectl apply -f k8s/persistent-volumes.yaml

# Deploy database and cache services
echo "🗄️  Deploying PostgreSQL..."
kubectl apply -f k8s/postgres-deployment.yaml

echo "🔴 Deploying Redis..."
kubectl apply -f k8s/redis-deployment.yaml

echo "📨 Deploying Kafka and Zookeeper..."
kubectl apply -f k8s/kafka-deployment.yaml

# Wait for database and cache to be ready
echo "⏳ Waiting for services to be ready..."
kubectl wait --for=condition=ready pod -l app=postgres -n gin-api --timeout=120s
kubectl wait --for=condition=ready pod -l app=redis -n gin-api --timeout=120s
kubectl wait --for=condition=ready pod -l app=kafka -n gin-api --timeout=120s

# Deploy API
echo "🚀 Deploying API..."
kubectl apply -f k8s/api-deployment.yaml

# Deploy Ingress
echo "🌐 Creating Ingress..."
kubectl apply -f k8s/ingress.yaml

echo "✅ Deployment complete!"
echo ""
echo "Check deployment status:"
echo "  kubectl get pods -n gin-api"
echo "  kubectl get services -n gin-api"
echo ""
echo "View logs:"
echo "  kubectl logs -f deployment/gin-api -n gin-api"
echo ""
echo "Access the API:"
echo "  kubectl port-forward service/gin-api-service 8080:80 -n gin-api"
