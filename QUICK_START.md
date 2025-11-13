# Quick Start Guide

Get the Cruder application up and running quickly.

## Prerequisites

- Go 1.25.0+
- Docker & Docker Compose
- PostgreSQL client tools (optional)
- Make

## Local Development

### 1. Setup Environment

```bash
# Create .env file from example
cp .env.example .env

# Edit .env with your values (or use defaults)
cat .env
```

### 2. Start Database

```bash
# Start PostgreSQL with Docker Compose
make db

# Or manually
docker-compose up -d db
```

### 3. Run Migrations

```bash
# Install goose if not already installed
go install github.com/pressly/goose/v3/cmd/goose@latest

# Run migrations
make migrate-up
```

### 4. Run Application

```bash
# Start application
make run

# Application will be available at http://localhost:8080
```

### 5. Test API

```bash
# Get all users
curl http://localhost:8080/api/v1/users

# With API key
curl -H "X-API-Key: your-api-key" http://localhost:8080/api/v1/users

# Get specific user
curl http://localhost:8080/api/v1/users/username/jdoe

# Create user
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-api-key" \
  -d '{
    "username": "newuser",
    "email": "newuser@example.com",
    "full_name": "New User"
  }'

# Update user (get UUID from GET request first)
curl -X PATCH http://localhost:8080/api/v1/users/{uuid} \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-api-key" \
  -d '{
    "username": "updated-username",
    "email": "updated@example.com"
  }'

# Delete user
curl -X DELETE http://localhost:8080/api/v1/users/{uuid} \
  -H "X-API-Key: your-api-key"
```

## Local with Docker

```bash
# Build Docker image
docker build -t cruder:latest .

# Run with Docker Compose
docker-compose up --build

# Or manually with database
docker-compose up -d db
docker run -p 8080:8080 \
  -e POSTGRES_HOST=localhost \
  -e POSTGRES_PASSWORD=postgres \
  -e API_KEY=dev-key \
  cruder:latest

# Access application
curl http://localhost:8080/api/v1/users
```

## AWS ECS Deployment

### Prerequisites

- AWS Account
- AWS CLI configured
- Docker image in ECR

### Deploy

```bash
# 1. Setup remote Terraform state
cd terraform
bash remote-state-setup.sh

# 2. Create terraform.tfvars
cp terraform.tfvars.example terraform.tfvars
vim terraform.tfvars

# 3. Deploy infrastructure
terraform init
terraform plan
terraform apply

# 4. Get ALB URL
aws elbv2 describe-load-balancers \
  --names cruder-alb \
  --query 'LoadBalancers[0].DNSName'

# 5. Test
curl http://alb-dns-name/api/v1/users
```

## Kubernetes Deployment

### Prerequisites

- Kubernetes cluster access
- kubectl configured
- Helm 3.x installed

### Deploy

```bash
# 1. Update values
vim helm/values.yaml

# 2. Create namespace
kubectl create namespace cruder

# 3. Deploy with Helm
helm install cruder ./helm \
  --namespace cruder \
  --set image.repository=your-registry/cruder \
  --set image.tag=latest \
  --set database.password=your-password \
  --set apiKey=your-api-key \
  --set ingress.host=cruder.example.com

# 4. Check status
kubectl get all -n cruder

# 5. Test
kubectl port-forward svc/cruder 8080:80 -n cruder
curl http://localhost:8080/api/v1/users
```

## Testing

### Run Tests

```bash
# Run all tests
make test

# Run with coverage
make coverage

# View HTML coverage report
make coverage-html
open coverage.html  # macOS
```

### Code Quality

```bash
# Format code
go fmt ./...

# Lint code
make lint

# Security check
make security

# Validate all
make validate
```

## Database Operations

### Create Migration

```bash
# Create new migration
make create-migration

# Follow the prompt to enter migration name
# Example: add_user_roles
```

### Check Migration Status

```bash
make migrate-status
```

### Rollback

```bash
make migrate-down
```

### Reset Database

```bash
make migrate-reset
```

## Logs

### Local Logs

```bash
# View application logs
# Logs will be printed to stdout in JSON format

# Example log entry:
# 2025/01/15 10:30:45 Incoming request: {
#   "timestamp":"2025-01-15T10:30:45.123456789+00:00",
#   "http.request.method":"GET",
#   "http.route":"/api/v1/users",
#   "http.response.status_code":200,
#   "http.server.request.duration":15,
#   "http.log.level":"info"
# }
```

### CloudWatch Logs (AWS)

```bash
# View logs
aws logs tail /ecs/cruder --follow

# Filter by log level
aws logs filter-log-events \
  --log-group-name /ecs/cruder \
  --filter-pattern "[... , http_log_level = \"error\", ...]"
```

### Kubernetes Logs

```bash
# View logs
kubectl logs deployment/cruder -n cruder -f

# Follow specific pod
kubectl logs <pod-name> -n cruder -f

# View logs from all pods
kubectl logs -l app=cruder -n cruder -f
```

## Common Tasks

### Scale Application

#### ECS
```bash
aws ecs update-service \
  --cluster cruder-cluster \
  --service cruder-service \
  --desired-count 5
```

#### Kubernetes
```bash
kubectl scale deployment/cruder --replicas=5 -n cruder
```

### Update Application

#### Docker
```bash
docker build -t your-registry/cruder:v1.1.0 .
docker push your-registry/cruder:v1.1.0
```

#### ECS
```bash
# Force new deployment with latest image
aws ecs update-service \
  --cluster cruder-cluster \
  --service cruder-service \
  --force-new-deployment
```

#### Kubernetes
```bash
helm upgrade cruder ./helm \
  --namespace cruder \
  --set image.tag=v1.1.0
```

### View Metrics

#### ECS
```bash
aws cloudwatch get-metric-statistics \
  --namespace AWS/ECS \
  --metric-name CPUUtilization \
  --dimensions Name=ServiceName,Value=cruder-service \
  --start-time 2025-01-15T00:00:00Z \
  --end-time 2025-01-16T00:00:00Z \
  --period 300 \
  --statistics Average,Maximum
```

#### Kubernetes
```bash
kubectl top nodes
kubectl top pods -n cruder
```

## Troubleshooting

### Connection Refused

```bash
# Check if application is running
curl http://localhost:8080/api/v1/users

# Check logs
docker logs <container-id>
kubectl logs <pod-name> -n cruder

# Check database connection
psql -h localhost -U postgres -d postgres -c "SELECT 1"
```

### Database Errors

```bash
# Check database status
docker ps | grep db

# Check database logs
docker logs database

# Connect to database
docker exec -it database psql -U postgres

# Check migrations
make migrate-status
```

### API Key Errors

```bash
# Without API key - should get 401
curl http://localhost:8080/api/v1/users

# With wrong API key - should get 403
curl -H "X-API-Key: wrong-key" http://localhost:8080/api/v1/users

# With correct API key - should work
curl -H "X-API-Key: your-api-key" http://localhost:8080/api/v1/users
```

## Useful Make Commands

```bash
make run              # Run application
make test             # Run tests
make coverage         # Show test coverage
make coverage-html    # Generate HTML coverage report
make lint             # Run linter
make security         # Run security scan
make validate         # Run all validations
make db               # Start database
make migrate-up       # Apply migrations
make migrate-down     # Rollback migrations
make migrate-status   # Check migration status
make migrate-reset    # Reset database
make create-migration # Create new migration
make up               # Start full docker-compose stack
make down             # Stop docker-compose stack
make restart          # Restart docker-compose stack
```

## Next Steps

1. Read [DEPLOYMENT.md](DEPLOYMENT.md) for production deployment
2. Read [INFRASTRUCTURE.md](INFRASTRUCTURE.md) for IaC details
3. Review [Helm Chart README](helm/README.md) for K8s deployment
4. Check CI/CD pipelines in `.github/workflows/`
5. Explore the [project README](README.md) for overview

## Support

For issues:
1. Check logs: `make logs` or `kubectl logs`
2. Review error messages
3. Check GitHub issues
4. Contact your ops team

## Security Notes

- **API Key**: Set a strong API key in production
- **Database Password**: Use strong passwords, store in secrets manager
- **HTTPS**: Use HTTPS in production
- **Secrets**: Never commit secrets to git
- **Images**: Use specific versions, scan for vulnerabilities