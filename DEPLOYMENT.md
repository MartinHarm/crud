# Deployment Guide for Cruder Application

This guide covers deploying the Cruder application to various environments.

## Table of Contents

1. [Configuration](#configuration)
2. [Docker](#docker)
3. [AWS ECS Deployment](#aws-ecs-deployment)
4. [Kubernetes Deployment](#kubernetes-deployment)
5. [CI/CD Pipeline](#cicd-pipeline)

## Configuration

### Environment Variables

The application uses both environment variables and configuration files:

```bash
# Application Settings
APP_ENV=production              # development, staging, production
SERVER_HOST=0.0.0.0            # Server bind address
SERVER_PORT=8080               # Server port

# Database
POSTGRES_HOST=localhost         # Database host
POSTGRES_PORT=5432             # Database port
POSTGRES_DB=cruder             # Database name
POSTGRES_USER=postgres          # Database user
POSTGRES_PASSWORD=password      # Database password
POSTGRES_SSL_MODE=disable       # SSL mode

# API
API_KEY=your-secure-api-key    # API key for authentication
```

### Configuration File

Create `config.yaml` in the application root:

```yaml
server:
  port: 8080
  host: 0.0.0.0
  env: production

database:
  host: localhost
  port: 5432
  user: postgres
  password: ${POSTGRES_PASSWORD}
  dbname: cruder
  sslmode: disable

api:
  key: ${API_KEY}
```

## Docker

### Build Docker Image

```bash
docker build -t cruder:latest .
```

### Run with Docker

```bash
docker run -d \
  --name cruder \
  -p 8080:8080 \
  -e POSTGRES_HOST=db \
  -e POSTGRES_PASSWORD=password \
  -e API_KEY=your-api-key \
  cruder:latest
```

### Docker Compose

```bash
docker-compose up --build
```

## AWS ECS Deployment

### Prerequisites

- AWS Account with appropriate permissions
- AWS CLI configured
- Docker image pushed to ECR

### Setup Remote Terraform State

```bash
cd terraform
bash remote-state-setup.sh
```

### Deploy Infrastructure

```bash
cd terraform

# Create terraform.tfvars from example
cp terraform.tfvars.example terraform.tfvars

# Edit terraform.tfvars with your values
vim terraform.tfvars

# Initialize and deploy
terraform init
terraform plan
terraform apply
```

### Manage Deployment

```bash
# View deployment status
aws ecs describe-services \
  --cluster cruder-cluster \
  --services cruder-service

# View logs
aws logs tail /ecs/cruder --follow

# Update task definition
aws ecs update-service \
  --cluster cruder-cluster \
  --service cruder-service \
  --task-definition cruder:latest \
  --force-new-deployment

# Scale service
aws ecs update-service \
  --cluster cruder-cluster \
  --service cruder-service \
  --desired-count 5
```

### Access Application

The application is accessible via the ALB:

```bash
# Get ALB DNS name
aws elbv2 describe-load-balancers \
  --names cruder-alb \
  --query 'LoadBalancers[0].DNSName'

# Make API request
curl -H "X-API-Key: your-api-key" \
  http://alb-dns-name/api/v1/users
```

## Kubernetes Deployment

### Prerequisites

- Kubernetes cluster (EKS, AKS, GKE, etc.)
- kubectl configured
- Helm 3.x installed
- Docker image in container registry

### Deploy with Helm

```bash
# Add and update repository (if using Helm repo)
helm repo add cruder https://charts.example.com
helm repo update

# Or deploy from local directory
cd helm

# Install release
helm install cruder . \
  --namespace default \
  --create-namespace \
  --values values.yaml \
  --set image.repository=your-registry/cruder \
  --set image.tag=latest \
  --set database.password=your-password \
  --set apiKey=your-api-key \
  --set ingress.host=cruder.example.com

# Or upgrade existing release
helm upgrade cruder . \
  --namespace default \
  --values values.yaml \
  --set image.tag=v1.0.0
```

### Verify Deployment

```bash
# Check deployment status
kubectl rollout status deployment/cruder -n default

# View pods
kubectl get pods -n default -l app=cruder

# View service
kubectl get svc -n default

# View ingress
kubectl get ingress -n default

# Check logs
kubectl logs -n default -l app=cruder -f

# Describe deployment
kubectl describe deployment cruder -n default
```

### Configure Ingress

For external access, configure your Ingress controller:

```bash
# For nginx-ingress with cert-manager:
helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
helm repo add jetstack https://charts.jetstack.io
helm repo update

helm install nginx-ingress ingress-nginx/ingress-nginx \
  --namespace ingress-nginx \
  --create-namespace

helm install cert-manager jetstack/cert-manager \
  --namespace cert-manager \
  --create-namespace \
  --set installCRDs=true
```

### Scale Application

```bash
# Manual scaling
kubectl scale deployment cruder --replicas=5 -n default

# Using HPA (Horizontal Pod Autoscaler)
kubectl autoscale deployment cruder \
  --min=2 --max=10 \
  --cpu-percent=70 \
  -n default
```

### Rollback Deployment

```bash
# View rollout history
kubectl rollout history deployment/cruder -n default

# Rollback to previous version
kubectl rollout undo deployment/cruder -n default

# Rollback to specific revision
kubectl rollout undo deployment/cruder --to-revision=2 -n default
```

## CI/CD Pipeline

### GitHub Actions Workflow

The project includes automated CI/CD pipelines:

1. **CI Pipeline** (`.github/workflows/ci.yml`)
   - Code formatting check
   - Go vet analysis
   - Go linting
   - Security scanning
   - Unit tests
   - Docker image build

2. **CD Pipeline** (`.github/workflows/cd.yml`)
   - Build and push Docker image
   - Deploy to AWS ECS
   - Deploy to Kubernetes
   - Run smoke tests
   - Slack notifications

### Configure Secrets

Add these secrets to your GitHub repository:

```
AWS_ACCESS_KEY_ID              - AWS access key
AWS_SECRET_ACCESS_KEY          - AWS secret key
DB_PASSWORD                    - Database password
API_KEY                        - API key
KUBE_CONFIG                    - Base64-encoded kubeconfig
K8S_INGRESS_HOST              - Kubernetes ingress hostname
SLACK_WEBHOOK_URL             - Slack webhook for notifications
```

### Trigger Deployments

Deployments are automatically triggered on:
- Push to `main` branch
- Pull requests (CI only)
- Manual workflow dispatch

### Monitor Pipeline

- View workflow runs in GitHub Actions
- Check deployment status in AWS Console or kubectl
- Review logs in CloudWatch or kubectl
- Receive notifications via Slack

## Monitoring and Logging

### CloudWatch Logs (AWS ECS)

```bash
# View logs
aws logs tail /ecs/cruder --follow

# Create metric filters
aws logs put-metric-filter \
  --log-group-name /ecs/cruder \
  --filter-name ErrorCount \
  --filter-pattern "[... , http_log_level = \"error\", ...]" \
  --metric-transformations metricName=ErrorCount,metricValue=1
```

### Kubernetes Logs

```bash
# View application logs
kubectl logs deployment/cruder -n default -f

# View logs with labels
kubectl logs -l app=cruder -n default --tail=100

# Stream logs from multiple pods
kubectl logs -f deployment/cruder -n default --all-containers=true
```

### JSON Structured Logging

The application outputs JSON-formatted logs:

```json
{
  "timestamp": "2025-01-15T10:30:45.123456789+00:00",
  "http.request.method": "GET",
  "http.route": "/api/v1/users",
  "http.request.host": "localhost",
  "http.request.remote_addr": "127.0.0.1",
  "http.response.status_code": 200,
  "http.server.request.duration": 15,
  "http.log.level": "info",
  "http.request.message": "Incoming request:",
  "http.user_agent": "curl/7.64.1",
  "request_id": "20250115103045-abc123def456"
}
```

## Troubleshooting

### ECS Issues

```bash
# Check task logs
aws logs tail /ecs/cruder --follow

# Describe task
aws ecs describe-tasks \
  --cluster cruder-cluster \
  --tasks <task-arn>

# Check service events
aws ecs describe-services \
  --cluster cruder-cluster \
  --services cruder-service \
  --query 'services[0].events'
```

### Kubernetes Issues

```bash
# Get pod status
kubectl get pods -n default -l app=cruder

# Describe pod for events
kubectl describe pod <pod-name> -n default

# Check resource usage
kubectl top nodes
kubectl top pods -n default

# Debug pod
kubectl exec -it <pod-name> -n default -- /bin/sh
```

### Network Issues

```bash
# Test service connectivity
kubectl run -it --rm debug --image=busybox --restart=Never -- sh
# Inside pod: wget -O- http://cruder.default.svc.cluster.local/api/v1/users

# Check DNS
kubectl run -it --rm debug --image=busybox --restart=Never -- nslookup cruder.default
```

## Cleanup

### Remove ECS Resources

```bash
cd terraform
terraform destroy
```

### Remove Kubernetes Resources

```bash
helm uninstall cruder --namespace default
kubectl delete namespace default
```

## Support

For issues or questions, contact your operations team or create an issue in the repository.