# Cruder Helm Chart

A Kubernetes Helm chart for deploying the Cruder application.

## Prerequisites

- Kubernetes 1.20+
- Helm 3.0+

## Installation

### Add Repository

```bash
helm repo add cruder https://charts.example.com
helm repo update
```

### Install Chart

```bash
helm install cruder cruder/cruder \
  --namespace default \
  --create-namespace \
  -f values.yaml
```

### From Local Directory

```bash
helm install cruder . \
  --namespace default \
  --create-namespace
```

## Configuration

### Basic Configuration

```bash
helm install cruder . \
  --namespace default \
  --set image.repository=your-registry/cruder \
  --set image.tag=v1.0.0 \
  --set database.password=your-password \
  --set apiKey=your-api-key \
  --set ingress.host=cruder.example.com
```

### With Custom Values File

```bash
helm install cruder . \
  --namespace default \
  --values values.yaml \
  --values values-prod.yaml
```

## Values

Key values and their defaults:

| Parameter | Description | Default |
|-----------|-------------|---------|
| `replicaCount` | Number of replicas | `3` |
| `image.repository` | Container image repository | `cruder` |
| `image.tag` | Container image tag | `latest` |
| `image.pullPolicy` | Pull policy | `IfNotPresent` |
| `appEnv` | Application environment | `production` |
| `database.host` | Database host | `postgres.default.svc.cluster.local` |
| `database.port` | Database port | `5432` |
| `database.name` | Database name | `cruder` |
| `database.user` | Database user | `postgres` |
| `database.password` | Database password | `` |
| `apiKey` | API key for authentication | `` |
| `ingress.enabled` | Enable ingress | `true` |
| `ingress.host` | Ingress hostname | `cruder.example.com` |
| `resources.requests.cpu` | CPU request | `250m` |
| `resources.requests.memory` | Memory request | `256Mi` |
| `resources.limits.cpu` | CPU limit | `500m` |
| `resources.limits.memory` | Memory limit | `512Mi` |
| `autoscaling.enabled` | Enable HPA | `true` |
| `autoscaling.minReplicas` | HPA min replicas | `2` |
| `autoscaling.maxReplicas` | HPA max replicas | `5` |
| `autoscaling.targetCPUUtilizationPercentage` | CPU target | `70` |
| `autoscaling.targetMemoryUtilizationPercentage` | Memory target | `80` |

## Upgrade

```bash
helm upgrade cruder . \
  --namespace default \
  --values values.yaml \
  --set image.tag=v1.1.0
```

## Rollback

```bash
# View history
helm history cruder --namespace default

# Rollback to previous version
helm rollback cruder --namespace default

# Rollback to specific revision
helm rollback cruder 1 --namespace default
```

## Uninstall

```bash
helm uninstall cruder --namespace default
```

## Chart Structure

```
helm/
├── Chart.yaml           # Chart metadata
├── values.yaml          # Default values
├── templates/           # Kubernetes manifests
│   ├── deployment.yaml
│   ├── service.yaml
│   ├── ingress.yaml
│   ├── hpa.yaml
│   ├── configmap.yaml
│   ├── secrets.yaml
│   └── serviceaccount.yaml
└── README.md
```

## Secrets Management

### Using Sealed Secrets

```bash
# Install sealed-secrets controller
kubectl apply -f https://github.com/bitnami-labs/sealed-secrets/releases/download/v0.18.0/controller.yaml

# Create and seal secret
echo -n 'your-password' | kubectl create secret generic db-password \
  --dry-run=client \
  --from-file=- \
  -o yaml | kubeseal -f - > sealed-secrets.yaml

# Apply sealed secret
kubectl apply -f sealed-secrets.yaml
```

### Using External Secrets

```bash
# Install external-secrets operator
helm repo add external-secrets https://charts.external-secrets.io
helm install external-secrets external-secrets/external-secrets \
  --namespace external-secrets-system \
  --create-namespace

# Create SecretStore
kubectl apply -f - <<EOF
apiVersion: external-secrets.io/v1beta1
kind: SecretStore
metadata:
  name: cruder-secret-store
  namespace: default
spec:
  provider:
    aws:
      service: SecretsManager
      region: us-east-1
      auth:
        jwt:
          serviceAccountRef:
            name: external-secrets-sa
EOF

# Create ExternalSecret
kubectl apply -f - <<EOF
apiVersion: external-secrets.io/v1beta1
kind: ExternalSecret
metadata:
  name: cruder-secrets
  namespace: default
spec:
  refreshInterval: 1h
  secretStoreRef:
    name: cruder-secret-store
    kind: SecretStore
  target:
    name: cruder-secrets
    creationPolicy: Owner
  data:
  - secretKey: db-password
    remoteRef:
      key: cruder/db/password
  - secretKey: api-key
    remoteRef:
      key: cruder/api/key
EOF
```

## Monitoring and Observability

### Prometheus Metrics

Add to values.yaml:

```yaml
prometheus:
  enabled: true
  serviceMonitor:
    enabled: true
    interval: 30s
```

### Loki Logging

Install Loki and configure promtail:

```bash
helm repo add grafana https://grafana.github.io/helm-charts
helm install loki grafana/loki-stack --namespace monitoring --create-namespace
```

### Jaeger Tracing

```bash
helm repo add jaegertracing https://jaegertracing.github.io/helm-charts
helm install jaeger jaegertracing/jaeger
```

## Troubleshooting

### Check Deployment

```bash
# Get status
helm status cruder --namespace default

# Get values
helm get values cruder --namespace default

# Get manifest
helm get manifest cruder --namespace default
```

### Debug Pods

```bash
# Get pod status
kubectl get pods -n default -l app=cruder -o wide

# Describe pod
kubectl describe pod <pod-name> -n default

# View logs
kubectl logs <pod-name> -n default

# Execute commands
kubectl exec -it <pod-name> -n default -- /bin/sh
```

### Common Issues

#### Pods not starting

Check pod events:
```bash
kubectl describe pod <pod-name> -n default
```

Check resource availability:
```bash
kubectl describe nodes
```

#### Connection timeouts

Check service:
```bash
kubectl get svc -n default
kubectl describe svc cruder -n default
```

Test connectivity:
```bash
kubectl run -it --rm debug --image=busybox --restart=Never -- sh
# Inside: wget -O- http://cruder/api/v1/users
```

#### Image pull errors

Check image:
```bash
kubectl describe pod <pod-name> -n default

# May need to create image pull secret:
kubectl create secret docker-registry regcred \
  --docker-server=your-registry \
  --docker-username=user \
  --docker-password=password
```

## Contributing

To modify the chart:

1. Edit templates in `templates/`
2. Update values in `values.yaml`
3. Update Chart.yaml version
4. Test with:
   ```bash
   helm lint .
   helm template cruder .
   helm install cruder . --dry-run --debug
   ```

## License

See LICENSE file in repository.