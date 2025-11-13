# Infrastructure as Code (IaC) Guide

This guide explains the infrastructure components and how to manage them.

## Overview

The Cruder application infrastructure is managed using:
- **Terraform**: Infrastructure provisioning and management
- **Helm**: Kubernetes application deployment and management
- **Docker**: Container image management

## Infrastructure Architecture

### AWS ECS Deployment

```
┌─────────────────────────────────────────────────────────┐
│                        AWS Cloud                        │
├─────────────────────────────────────────────────────────┤
│                                                         │
│  ┌─────────────────────────────────────────────────┐   │
│  │              Application Load Balancer           │   │
│  │           (Public - Port 80/443)                │   │
│  └────────────────┬────────────────────────────────┘   │
│                   │                                     │
│  ┌────────────────▼────────────────────────────────┐   │
│  │         VPC (10.0.0.0/16)                       │   │
│  │                                                 │   │
│  │  ┌─────────────────────────────────────────┐   │   │
│  │  │  Public Subnets (10.0.1-2.0.0/24)       │   │   │
│  │  │  - AZ1 & AZ2                            │   │   │
│  │  │  - Route to IGW                         │   │   │
│  │  └─────────────────────────────────────────┘   │   │
│  │                                                 │   │
│  │  ┌─────────────────────────────────────────┐   │   │
│  │  │  Private Subnets (10.0.10-11.0.0/24)    │   │   │
│  │  │  - AZ1 & AZ2                            │   │   │
│  │  │  - ECS Tasks                            │   │   │
│  │  │  - RDS Instance                         │   │   │
│  │  └─────────────────────────────────────────┘   │   │
│  │                                                 │   │
│  └─────────────────────────────────────────────────┘   │
│                                                         │
│  ┌─────────────────────────────────────────────────┐   │
│  │         CloudWatch Logs & Metrics               │   │
│  │  - Application logs: /ecs/cruder                │   │
│  │  - Database metrics                            │   │
│  └─────────────────────────────────────────────────┘   │
│                                                         │
│  ┌─────────────────────────────────────────────────┐   │
│  │       AWS Secrets Manager                       │   │
│  │  - Database password                           │   │
│  │  - API key                                      │   │
│  └─────────────────────────────────────────────────┘   │
│                                                         │
└─────────────────────────────────────────────────────────┘
```

### Kubernetes Deployment

```
┌──────────────────────────────────────────────────────┐
│          Kubernetes Cluster                         │
├──────────────────────────────────────────────────────┤
│                                                      │
│  ┌───────────────────────────────────────────────┐  │
│  │  Ingress (nginx-ingress)                      │  │
│  │  - HTTPS/TLS (cert-manager)                   │  │
│  │  - Route /api/v1 → cruder service             │  │
│  └───────────────────────────────────────────────┘  │
│                                                      │
│  ┌───────────────────────────────────────────────┐  │
│  │  Service (ClusterIP)                          │  │
│  │  - Port 80 → Pod Port 8080                    │  │
│  └───────────────────────────────────────────────┘  │
│                                                      │
│  ┌───────────────────────────────────────────────┐  │
│  │  Deployment                                    │  │
│  │  - 3 replicas (HPA: 2-5)                      │  │
│  │  - Rolling update strategy                    │  │
│  │  - Pod Anti-affinity                          │  │
│  └───────────────────────────────────────────────┘  │
│                                                      │
│  ┌───────────────────────────────────────────────┐  │
│  │  HorizontalPodAutoscaler                      │  │
│  │  - CPU: 70% target                            │  │
│  │  - Memory: 80% target                         │  │
│  └───────────────────────────────────────────────┘  │
│                                                      │
│  ┌───────────────────────────────────────────────┐  │
│  │  ConfigMaps & Secrets                         │  │
│  │  - config.yaml                                │  │
│  │  - Database password                          │  │
│  │  - API key                                    │  │
│  └───────────────────────────────────────────────┘  │
│                                                      │
└──────────────────────────────────────────────────────┘
```

## Terraform Management

### Directory Structure

```
terraform/
├── main.tf                      # Core infrastructure
├── iam.tf                       # IAM roles and policies
├── secrets.tf                   # AWS Secrets Manager
├── variables.tf                 # Input variables
├── outputs.tf                   # Output values
├── terraform.tfvars.example     # Example values
├── remote-state-setup.sh        # Remote state initialization
└── .terraform/                  # Terraform working directory
```

### Setup Remote State

1. **Initialize S3 and DynamoDB for remote state:**

```bash
cd terraform
bash remote-state-setup.sh
```

2. **Configure backend in `main.tf`** (already done):

```hcl
backend "s3" {
  bucket         = "cruder-terraform-state"
  key            = "prod/terraform.tfstate"
  region         = "us-east-1"
  encrypt        = true
  dynamodb_table = "terraform-locks"
}
```

### Deploy Infrastructure

1. **Create `terraform.tfvars`:**

```bash
cp terraform/terraform.tfvars.example terraform/terraform.tfvars
```

2. **Edit variables:**

```bash
vim terraform/terraform.tfvars
```

Key variables to set:
- `aws_region`: AWS region
- `environment`: production/staging/development
- `db_password`: Strong database password
- `api_key`: Strong API key
- `ecr_repository_url`: ECR repository URL

3. **Initialize and deploy:**

```bash
cd terraform

# Initialize (downloads providers and state)
terraform init

# Plan changes
terraform plan -out=tfplan

# Apply changes
terraform apply tfplan
```

### Manage Existing Infrastructure

```bash
# View current state
terraform state list
terraform state show aws_db_instance.postgres

# Refresh state from AWS
terraform refresh

# Target specific resource
terraform apply -target=aws_ecs_service.app

# Destroy infrastructure
terraform destroy
```

### Troubleshooting Terraform

```bash
# Debug output
TF_LOG=DEBUG terraform apply

# Validate configuration
terraform validate

# Format code
terraform fmt

# Linting with tflint
tflint
```

## Helm Management

### Directory Structure

```
helm/
├── Chart.yaml                  # Chart metadata
├── values.yaml                 # Default values
├── values-prod.yaml            # Production values
├── values-dev.yaml             # Development values
├── README.md                    # Documentation
└── templates/                   # Kubernetes manifests
    ├── _helpers.tpl
    ├── deployment.yaml
    ├── service.yaml
    ├── ingress.yaml
    ├── hpa.yaml
    ├── secrets.yaml
    ├── configmap.yaml
    └── serviceaccount.yaml
```

### Deploy with Helm

1. **Install release:**

```bash
helm install cruder ./helm \
  --namespace default \
  --create-namespace \
  --values helm/values.yaml \
  --set image.tag=v1.0.0 \
  --set database.password=secure-password \
  --set apiKey=secure-api-key
```

2. **Update release:**

```bash
helm upgrade cruder ./helm \
  --namespace default \
  --values helm/values.yaml \
  --set image.tag=v1.1.0
```

3. **Rollback release:**

```bash
helm rollback cruder --namespace default
```

### Values Management

Create environment-specific values files:

**values-prod.yaml:**
```yaml
replicaCount: 5
autoscaling:
  minReplicas: 3
  maxReplicas: 10
  targetCPUUtilizationPercentage: 60
resources:
  limits:
    cpu: 1000m
    memory: 1024Mi
```

**values-dev.yaml:**
```yaml
replicaCount: 1
autoscaling:
  enabled: false
resources:
  limits:
    cpu: 250m
    memory: 512Mi
```

Deploy with specific values:

```bash
helm install cruder ./helm \
  --values helm/values.yaml \
  --values helm/values-prod.yaml
```

### Chart Development

1. **Validate chart:**

```bash
helm lint ./helm
```

2. **Dry run:**

```bash
helm install cruder ./helm --dry-run --debug
```

3. **Generate manifests:**

```bash
helm template cruder ./helm > manifests.yaml
```

4. **Test with local Kubernetes:**

```bash
# Using kind or minikube
helm install cruder ./helm --namespace test --create-namespace
kubectl get all -n test
helm uninstall cruder -n test
```

## Docker Image Management

### Build and Push

```bash
# Build image
docker build -t cruder:latest .
docker build -t cruder:v1.0.0 .

# Tag for registry
docker tag cruder:latest your-registry/cruder:latest
docker tag cruder:latest your-registry/cruder:v1.0.0

# Push to registry
docker push your-registry/cruder:latest
docker push your-registry/cruder:v1.0.0
```

### Multi-registry Push

```bash
# Push to multiple registries
docker tag cruder:latest your-registry/cruder:latest
docker tag cruder:latest your-ecr-url/cruder:latest
docker tag cruder:latest ghcr.io/yourorg/cruder:latest

docker push your-registry/cruder:latest
docker push your-ecr-url/cruder:latest
docker push ghcr.io/yourorg/cruder:latest
```

## State Management

### S3 Remote State

The remote state is stored in S3 with:
- **Encryption**: AES256 enabled
- **Versioning**: Enabled for rollback
- **Locking**: DynamoDB table prevents concurrent modifications
- **Access**: Restricted via IAM policies

View state:

```bash
# List state files
aws s3 ls s3://cruder-terraform-state/

# Download state file
aws s3 cp s3://cruder-terraform-state/prod/terraform.tfstate .

# View state lock
aws dynamodb scan --table-name terraform-locks
```

## Backup and Disaster Recovery

### Database Backups

```bash
# Automated backups (configured in Terraform)
# - Retention: 7 days
# - Backup window: 03:00-04:00 UTC
# - Multi-AZ: Enabled

# Manual backup
aws rds create-db-snapshot \
  --db-instance-identifier cruder-postgres \
  --db-snapshot-identifier cruder-postgres-manual-backup

# List backups
aws rds describe-db-snapshots --db-instance-identifier cruder-postgres
```

### Database Restore

```bash
# Restore from snapshot
aws rds restore-db-instance-from-db-snapshot \
  --db-instance-identifier cruder-postgres-restored \
  --db-snapshot-identifier cruder-postgres-manual-backup
```

### ECS Task Logs

```bash
# Export logs to S3
aws logs create-export-task \
  --log-group-name /ecs/cruder \
  --from $(date -d '7 days ago' +%s)000 \
  --to $(date +%s)000 \
  --destination cruder-logs-bucket \
  --destination-prefix backup/
```

## Cost Optimization

### AWS Cost Estimation

```bash
# Estimate costs
terraform plan | grep aws_

# Use AWS Cost Explorer for detailed analysis
aws ce get-cost-and-usage \
  --time-period Start=2025-01-01,End=2025-01-31 \
  --granularity DAILY \
  --metrics BlendedCost
```

### Optimization Tips

1. **Use Reserved Instances** for RDS
2. **Spot Instances** for non-critical ECS tasks
3. **Auto-scaling** for variable loads
4. **CloudFront** for static content CDN
5. **S3 Lifecycle Policies** for old logs

## Security

### Secrets Management

Secrets are stored in AWS Secrets Manager:
- Database password
- API key

Access via Terraform:

```hcl
data "aws_secretsmanager_secret_version" "db_password" {
  secret_id = aws_secretsmanager_secret.db_password.id
}
```

### IAM Policies

ECS tasks use least-privilege IAM roles:
- `ecs_task_execution_role`: Read secrets, write logs
- `ecs_task_role`: Application-specific permissions

### Network Security

- Security groups restrict traffic
- Private subnets for databases
- ALB in public subnets
- NACLs for subnet-level filtering

## Monitoring and Alerts

### CloudWatch Metrics

```bash
# View metrics
aws cloudwatch get-metric-statistics \
  --namespace AWS/ECS \
  --metric-name CPUUtilization \
  --dimensions Name=ServiceName,Value=cruder-service \
  --start-time 2025-01-15T00:00:00Z \
  --end-time 2025-01-16T00:00:00Z \
  --period 300 \
  --statistics Average,Maximum
```

### Create Alarms

```bash
# CPU alert
aws cloudwatch put-metric-alarm \
  --alarm-name cruder-high-cpu \
  --alarm-description "Alert if CPU > 80%" \
  --metric-name CPUUtilization \
  --namespace AWS/ECS \
  --statistic Average \
  --period 300 \
  --threshold 80 \
  --comparison-operator GreaterThanThreshold
```

## Documentation

- [AWS Documentation](https://docs.aws.amazon.com/)
- [Terraform AWS Provider](https://registry.terraform.io/providers/hashicorp/aws/latest/docs)
- [Helm Documentation](https://helm.sh/docs/)
- [Kubernetes Documentation](https://kubernetes.io/docs/)