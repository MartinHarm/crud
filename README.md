### Start Application

**Local Development**:
```bash
make run
```

**On Kubernetes**:
```bash
helm install cruder ./helm
```

**On AWS ECS**:
```bash
cd terraform && terraform apply
```

### Test API

```bash
# Get all users
curl -H "X-API-Key: your-key" http://localhost:8080/api/v1/users

# Create user
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-key" \
  -d '{
    "username": "newuser",
    "email": "newuser@example.com",
    "full_name": "New User"
  }'

# Update user
curl -X PATCH http://localhost:8080/api/v1/users/{uuid} \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-key" \
  -d '{"username": "updated"}'

# Delete user
curl -X DELETE http://localhost:8080/api/v1/users/{uuid} \
  -H "X-API-Key: your-key"
```

### View Logs

**Local**:
```bash
docker-compose logs -f app
```

**Kubernetes**:
```bash
kubectl logs deployment/cruder -f
```

**AWS ECS**:
```bash
aws logs tail /ecs/cruder --follow
```

---

## 🔧 Key Commands

### Development
```bash
make run              # Run application
make test             # Run tests
make lint             # Run linter
make security         # Security scan
make migrate-up       # Apply migrations
make db               # Start database
```

### Docker
```bash
docker build -t cruder:latest .
docker-compose up --build
docker-compose down
```

### Kubernetes
```bash
helm lint ./helm
helm install cruder ./helm
helm upgrade cruder ./helm
helm uninstall cruder
kubectl logs deployment/cruder -f
```

### Terraform
```bash
cd terraform
terraform init
terraform plan
terraform apply
terraform destroy
```

## 📊 Architecture Diagrams

### Request Flow & Business Logic
```mermaid
graph TD
    A["📱 HTTP Request<br/>GET/POST/PATCH/DELETE<br/>X-API-Key Header"] 
    
    B["🔑 Request ID Middleware<br/>Generate/Extract X-Request-ID<br/>Attach to Context"]
    
    C["🔐 API Key Middleware<br/>Check X-API-Key Header"]
    
    C1["❌ Missing Key<br/>Return 401 Unauthorized"]
    C2["❌ Invalid Key<br/>Return 403 Forbidden"]
    C3["✅ Valid Key<br/>Proceed"]
    
    D["📋 JSON Logger Middleware<br/>Start Request Timer<br/>Log Request Details"]
    
    E["🛣️ Router<br/>Match Endpoint<br/>Route to Controller"]
    
    F["🎮 Controller Layer<br/>Validate Request DTO<br/>Call Service Layer"]
    
    G["⚙️ Service Layer<br/>Business Logic<br/>Data Validation<br/>Authorization Checks"]
    
    H{{"Operation<br/>Type"}}
    
    I["📖 GET: Fetch Users<br/>GetAll, GetByUsername<br/>GetByID, GetByUUID"]
    
    J["✏️ POST: Create User<br/>Validate Email/Username<br/>Generate UUID<br/>Set Timestamps"]
    
    K["🔄 PATCH: Update User<br/>Find by UUID<br/>Merge Fields<br/>Update Timestamps"]
    
    L["🗑️ DELETE: Remove User<br/>Find by UUID<br/>Delete from DB"]
    
    M["🗄️ Repository Layer<br/>SQL Query Execution<br/>PostgreSQL Interaction"]
    
    N["📊 Database Response<br/>User Records<br/>Affected Rows"]
    
    O{{"Success?"}}
    
    P["✅ Success Response<br/>200 OK / 201 Created<br/>204 No Content"]
    
    Q["❌ Error Response<br/>404 Not Found<br/>500 Internal Server Error"]
    
    R["📝 JSON Logger Middleware<br/>Log Response Status<br/>Log Duration<br/>Log Request ID<br/>Output to stdout"]
    
    S["📤 HTTP Response<br/>Status Code<br/>Body/Headers<br/>Request ID"]
    
    T["📊 Log Aggregation<br/>CloudWatch/ELK/Splunk<br/>Request Tracing<br/>Performance Metrics"]
    
    A --> B
    B --> C
    C -->|Missing| C1
    C -->|Invalid| C2
    C -->|Valid| C3
    C1 --> R
    C2 --> R
    C3 --> D
    D --> E
    E --> F
    F --> G
    G --> H
    H -->|Read| I
    H -->|Create| J
    H -->|Update| K
    H -->|Delete| L
    I --> M
    J --> M
    K --> M
    L --> M
    M --> N
    N --> O
    O -->|Yes| P
    O -->|No| Q
    P --> R
    Q --> R
    R --> S
    S --> T
    
    style A fill:#e1f5ff
    style C1 fill:#ffcdd2
    style C2 fill:#ffcdd2
    style C3 fill:#c8e6c9
    style P fill:#c8e6c9
    style Q fill:#ffcdd2
    style T fill:#fff3e0
```

### Service Layer Details
```mermaid
graph LR
    subgraph Input["📥 Input Validation"]
        V1["Check Required Fields"]
        V2["Validate Email Format"]
        V3["Validate Username<br/>Not Empty, No Spaces"]
    end
    
    subgraph Business["⚙️ Business Logic"]
        B1["Check for Duplicates<br/>Username/Email"]
        B2["Generate UUID<br/>for New Records"]
        B3["Set Timestamps<br/>created_at/updated_at"]
        B4["Handle Soft Deletes<br/>if Applicable"]
    end
    
    subgraph DB["🗄️ Database Operations"]
        D1["Query Validation"]
        D2["SQL Execution"]
        D3["Transaction Handling"]
        D4["Error Recovery"]
    end
    
    subgraph Output["📤 Output & Response"]
        O1["Format Response DTO"]
        O2["Set HTTP Status Code"]
        O3["Attach Request ID"]
        O4["Send to Logger"]
    end
    
    V1 --> V2 --> V3
    V3 --> B1
    B1 --> B2 --> B3 --> B4
    B4 --> D1 --> D2 --> D3 --> D4
    D4 --> O1 --> O2 --> O3 --> O4
    
    style Input fill:#e3f2fd
    style Business fill:#f3e5f5
    style DB fill:#fce4ec
    style Output fill:#e0f2f1
```

### Data Flow Through Layers
```mermaid
sequenceDiagram
    participant Client as 🖥️ Client
    participant Middleware as 🔐 Middleware
    participant Handler as 🎮 Handler
    participant Service as ⚙️ Service
    participant Repository as 🗄️ Repository
    participant DB as 💾 PostgreSQL

    Client->>Middleware: HTTP Request + X-API-Key
    Middleware->>Middleware: Generate Request ID
    Middleware->>Middleware: Validate API Key
    Middleware->>Middleware: Log Request Start
    
    Middleware->>Handler: Pass to Router
    Handler->>Service: Call Business Logic
    
    Service->>Service: Validate Input
    Service->>Service: Apply Business Rules
    Service->>Repository: Call Repository Method
    
    Repository->>Repository: Build SQL Query
    Repository->>DB: Execute Query
    DB->>DB: Process Query
    DB-->>Repository: Return Result Set
    
    Repository-->>Service: Return Domain Model
    Service->>Service: Format Response
    Service-->>Handler: Return Result/Error
    
    Handler->>Middleware: Response Ready
    Middleware->>Middleware: Log Response
    Middleware->>Middleware: Attach Request ID
    Middleware->>Client: HTTP Response + Status
    
    Client->>Client: Receive Response
```

### Error Handling Flow
```mermaid
graph TD
    E["❌ Error Occurs<br/>at Any Layer"]
    
    E1{"Error<br/>Type?"}
    
    E2["❌ Validation Error<br/>- Invalid Input<br/>- Missing Required Fields"]
    E2R["🔴 400 Bad Request"]
    
    E3["❌ Not Found<br/>- User UUID Not Found<br/>- Resource Missing"]
    E3R["🟡 404 Not Found"]
    
    E4["❌ Conflict<br/>- Duplicate Username<br/>- Duplicate Email"]
    E4R["🟠 409 Conflict"]
    
    E5["❌ Auth Error<br/>- Missing API Key<br/>- Invalid API Key"]
    E5R["🔐 401/403"]
    
    E6["❌ Server Error<br/>- DB Connection Error<br/>- Query Execution Error"]
    E6R["🔴 500 Internal Server Error"]
    
    RESP["📝 Error Response<br/>JSON Body<br/>Error Message<br/>Request ID"]
    
    LOG["📊 Structured Log<br/>Error Details<br/>Stack Trace<br/>Timestamp"]
    
    E --> E1
    E1 -->|Validation| E2 --> E2R
    E1 -->|Not Found| E3 --> E3R
    E1 -->|Conflict| E4 --> E4R
    E1 -->|Auth| E5 --> E5R
    E1 -->|Server| E6 --> E6R
    
    E2R --> RESP
    E3R --> RESP
    E4R --> RESP
    E5R --> RESP
    E6R --> RESP
    
    RESP --> LOG
    LOG --> Client["📱 Client Receives<br/>Error Response"]
    
    style E fill:#ffcdd2
    style E2R fill:#ffcdd2
    style E3R fill:#ffcdd2
    style E4R fill:#ffcdd2
    style E5R fill:#ffcdd2
    style E6R fill:#ffcdd2
```

### AWS ECS Architecture
```
Internet (HTTP/HTTPS)
    ↓
Application Load Balancer (Public)
    ↓
ECS Fargate Tasks (Private Subnets, Multi-AZ)
    ↓
RDS PostgreSQL (Private, Multi-AZ)
    ↓
CloudWatch Logs & Metrics
```

### Kubernetes Architecture
```
External Users
    ↓
Ingress (HTTPS/TLS)
    ↓
Service (ClusterIP)
    ↓
Pods (3+ replicas, HPA)
    ↓
Database (External or In-cluster)
```

---

## 📈 Monitoring

### Local Development
- JSON logs to stdout
- Request ID correlation
- Status-based log levels

### AWS ECS
- CloudWatch Logs: `/ecs/cruder`
- CloudWatch Metrics: CPU, memory, task count
- Alarms for high utilization
- RDS monitoring and backups

### Kubernetes
- Kubectl logs and events
- Pod metrics via kubectl top
- Prometheus integration (optional)
- Loki log aggregation (optional)
- Grafana dashboards (optional)

---

## 🚢 Deployment Options

### Local
```bash
make run
# Application at http://localhost:8080
```

### Docker Compose
```bash
docker-compose up --build
# Application at http://localhost:8080
```

### AWS ECS
```bash
cd terraform && terraform apply
# Application at ALB URL (from Terraform output)
```

### Kubernetes
```bash
helm install cruder ./helm
# Application at Ingress hostname
```

### Cloud Providers
- AWS: ECS (Terraform)
- GCP: GKE (Kubernetes)
- Azure: AKS (Kubernetes)
- DigitalOcean: App Platform or DOKS
- Self-hosted: Docker Compose or Kubernetes

---

## 🔄 CI/CD Pipeline Flow

```
┌─────────────────────┐
│  Push to main       │
└──────────┬──────────┘
           ↓
┌──────────────────────────┐
│  CI Pipeline Runs:       │
│  - Format check          │
│  - Lint/Vet             │
│  - Security scan         │
│  - Tests                │
│  - Docker build         │
└──────────┬───────────────┘
           ↓ (if all pass)
┌──────────────────────────┐
│  CD Pipeline Runs:       │
│  - Push image            │
│  - Deploy to ECS         │
│  - Deploy to K8s         │
│  - Smoke tests           │
│  - Notify Slack          │
└──────────┬───────────────┘
           ↓
┌─────────────────────┐
│  Production Live!   │
└─────────────────────┘
```

---

## ❓ FAQ

### Q: How do I configure the API key?
**A**: Set `API_KEY` environment variable or in `config.yaml`

### Q: Can I disable authentication?
**A**: Yes, leave `API_KEY` empty (not recommended for production)

### Q: How do I scale the application?
**A**: Use HPA in K8s or AWS auto-scaling in ECS

### Q: What's the recommended database size?
**A**: Start with `db.t3.micro`, use RDS Performance Insights for monitoring

### Q: How do I backup the database?
**A**: Automated in Terraform (7-day retention), manual snapshots available

### Q: Can I use a different cloud provider?
**A**: Yes, adapt Terraform code or use Kubernetes

### Q: How do I monitor the application?
**A**: CloudWatch for AWS, kubectl for K8s, Prometheus for both

---

## 🆘 Troubleshooting Guide

### Application won't start
1. Check config file exists
2. Verify database connectivity
3. Check environment variables
4. Review logs

### Database connection fails
1. Verify credentials
2. Check network connectivity
3. Ensure database is running
4. Check security groups

### Deployment fails
1. Review Terraform errors
2. Check AWS credentials
3. Verify resource quotas
4. Review CloudWatch logs

### Tests fail
1. Run locally first
2. Check database state
3. Review log output
4. Run with `-v` flag

---