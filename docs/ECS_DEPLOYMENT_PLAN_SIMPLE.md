# AWS ECS Deployment Plan (Simplified - ECS + RDS Only)

## Overview
Simplified deployment plan for gRPC service to AWS ECS with RDS PostgreSQL. **No load balancer** - ECS tasks will have public IPs for direct access.

---

## 📋 Table of Contents
1. [Prerequisites & Setup](#phase-1-prerequisites--setup)
2. [Containerization](#phase-2-containerization)
3. [Terraform Infrastructure](#phase-3-terraform-infrastructure)
4. [Deployment Steps](#phase-4-deployment-steps)
5. [Testing & Verification](#phase-5-testing--verification)
6. [Cost Estimate](#cost-estimate)

---

## Phase 1: Prerequisites & Setup (20 min)

### 1.1 Simplified Project Structure
```
terraform/
├── main.tf                 # Main configuration
├── variables.tf            # Input variables
├── outputs.tf             # Output values
├── versions.tf            # Provider versions
├── backend.tf             # S3 backend configuration
├── modules/
│   ├── networking/        # VPC, subnets, security groups
│   │   ├── main.tf
│   │   ├── variables.tf
│   │   └── outputs.tf
│   ├── ecs/              # ECS cluster, service, task definition
│   │   ├── main.tf
│   │   ├── variables.tf
│   │   └── outputs.tf
│   ├── rds/              # PostgreSQL RDS instance
│   │   ├── main.tf
│   │   ├── variables.tf
│   │   └── outputs.tf
│   └── ecr/              # Container registry
│       ├── main.tf
│       ├── variables.tf
│       └── outputs.tf
└── environments/
    ├── dev.tfvars
    └── prod.tfvars
```

### 1.2 Tools to Install
```bash
# Terraform CLI
brew install terraform

# AWS CLI
brew install awscli

# Docker
brew install --cask docker

# grpcurl (for testing)
brew install grpcurl
```

### 1.3 AWS Setup
```bash
# Configure AWS credentials
aws configure

# Set your AWS region (e.g., us-east-1)
# Set your AWS access key and secret key
```

### 1.4 Required AWS Permissions
- EC2 (VPC, Subnets, Security Groups)
- ECS (Cluster, Service, Task Definition)
- ECR (Repository)
- RDS (Database Instance)
- IAM (Roles, Policies)
- CloudWatch (Logs)
- Secrets Manager (Database credentials)

---

## Phase 2: Containerization (20 min)

### 2.1 Create Dockerfile

Create `Dockerfile` in project root:

```dockerfile
# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o server cmd/server/main.go

# Runtime stage
FROM alpine:latest

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/server .

# Expose gRPC port
EXPOSE 50051

# Run the server
CMD ["./server"]
```

### 2.2 Create .dockerignore

Create `.dockerignore` in project root:

```
.git
.gitignore
.env
*.md
docs/
terraform/
migrations/
.DS_Store
*.log
server
```

### 2.3 Build & Test Locally

```bash
# Build the image
docker build -t grpc-go:latest .

# Run with local PostgreSQL
docker run -p 50051:50051 \
  -e DB_HOST=host.docker.internal \
  -e DB_PORT=5432 \
  -e DB_USER=superuser \
  -e DB_PASSWORD=superpass \
  -e DB_NAME=super_salary_db \
  -e DB_SSLMODE=disable \
  grpc-go:latest

# Test in another terminal
grpcurl -plaintext localhost:50051 list
```

---

## Phase 3: Terraform Infrastructure (1-2 hours)

### 3.1 Architecture Overview

```
┌─────────────────────────────────────────────────────────┐
│                         VPC                              │
│  ┌────────────────────┐      ┌────────────────────┐    │
│  │  Public Subnet     │      │  Public Subnet     │    │
│  │  (AZ-1)            │      │  (AZ-2)            │    │
│  │                    │      │                    │    │
│  │  ┌──────────────┐  │      │  ┌──────────────┐  │    │
│  │  │ ECS Task     │  │      │  │ ECS Task     │  │    │
│  │  │ (Public IP)  │  │      │  │ (Public IP)  │  │    │
│  │  └──────────────┘  │      │  └──────────────┘  │    │
│  └────────────────────┘      └────────────────────┘    │
│                                                          │
│  ┌────────────────────┐      ┌────────────────────┐    │
│  │  Private Subnet    │      │  Private Subnet    │    │
│  │  (AZ-1)            │      │  (AZ-2)            │    │
│  │                    │      │                    │    │
│  │  ┌──────────────┐  │      │                    │    │
│  │  │ RDS Primary  │  │      │                    │    │
│  │  │ PostgreSQL   │  │      │                    │    │
│  │  └──────────────┘  │      │                    │    │
│  └────────────────────┘      └────────────────────┘    │
└─────────────────────────────────────────────────────────┘
```

**Key Points:**
- ECS tasks in **public subnets** with **public IPs** (no NAT Gateway needed!)
- RDS in **private subnets** (more secure)
- No load balancer (cost savings)
- Direct access to ECS tasks via public IP

### 3.2 Module Breakdown

#### Module 1: Networking

**Creates:**
- VPC (10.0.0.0/16)
- 2 Public subnets (for ECS tasks)
- 2 Private subnets (for RDS)
- Internet Gateway (for public subnet internet access)
- Route tables
- Security Groups:
  - **ECS Tasks SG**: Allow inbound 50051 from anywhere (0.0.0.0/0)
  - **RDS SG**: Allow inbound 5432 from ECS Tasks SG only

**No NAT Gateway needed** since ECS tasks have public IPs!

#### Module 2: ECR

**Creates:**
- ECR repository for Docker images
- Image scanning on push
- Lifecycle policy (keep last 10 images)

#### Module 3: RDS

**Creates:**
- PostgreSQL RDS instance
- DB subnet group (private subnets)
- Security group (allow 5432 from ECS)
- Secrets Manager secret for password
- Parameter group (optional tuning)

**Sizing:**
- Dev: `db.t3.micro` (1 vCPU, 1 GB RAM, 20 GB storage)
- Prod: `db.t3.small` or larger

#### Module 4: ECS

**Creates:**
- ECS Cluster (Fargate)
- Task Definition:
  - CPU: 256 (0.25 vCPU)
  - Memory: 512 MB
  - Container: gRPC service
  - Environment variables
  - Secrets from Secrets Manager
  - CloudWatch Logs
- ECS Service:
  - Desired count: 1
  - Public IP: **Enabled**
  - No load balancer
  - Auto-assign public IP

**IAM Roles:**
- Task Execution Role (pull images, logs, secrets)
- Task Role (application permissions)

---

## Phase 4: Deployment Steps (30 min)

### 4.1 Create Terraform Directory Structure

```bash
cd /Users/nguyenphuoc/Desktop/personal-repos/go-grpc

# Create directory structure
mkdir -p terraform/{modules/{networking,ecs,rds,ecr},environments}
cd terraform
```

### 4.2 Create Terraform Files

You'll create these files in the next step:
- `terraform/versions.tf` - Provider versions
- `terraform/backend.tf` - S3 backend (optional for now)
- `terraform/variables.tf` - Input variables
- `terraform/outputs.tf` - Output values
- `terraform/main.tf` - Main configuration
- `terraform/environments/dev.tfvars` - Dev environment variables

And modules:
- `terraform/modules/networking/` - VPC, subnets, security groups
- `terraform/modules/ecr/` - ECR repository
- `terraform/modules/rds/` - RDS instance
- `terraform/modules/ecs/` - ECS cluster, service, task definition

### 4.3 Initialize Terraform

```bash
cd terraform
terraform init
```

### 4.4 Plan Infrastructure

```bash
terraform plan -var-file=environments/dev.tfvars -out=tfplan
```

### 4.5 Apply Infrastructure

```bash
terraform apply tfplan
```

### 4.6 Build & Push Docker Image

```bash
# Get ECR repository URL from Terraform output
ECR_REPO=$(terraform output -raw ecr_repository_url)

# Login to ECR
aws ecr get-login-password --region us-east-1 | \
  docker login --username AWS --password-stdin $ECR_REPO

# Build image
docker build -t grpc-go:latest .

# Tag and push
docker tag grpc-go:latest $ECR_REPO:latest
docker push $ECR_REPO:latest
```

### 4.7 Update ECS Service

After pushing the image, ECS will automatically pull and deploy it.

```bash
# Force new deployment (if needed)
aws ecs update-service \
  --cluster grpc-go-cluster \
  --service grpc-go-service \
  --force-new-deployment
```

---

## Phase 5: Testing & Verification (15 min)

### 5.1 Get ECS Task Public IP

```bash
# Get task ARN
TASK_ARN=$(aws ecs list-tasks \
  --cluster grpc-go-cluster \
  --service-name grpc-go-service \
  --query 'taskArns[0]' \
  --output text)

# Get task details and extract public IP
PUBLIC_IP=$(aws ecs describe-tasks \
  --cluster grpc-go-cluster \
  --tasks $TASK_ARN \
  --query 'tasks[0].attachments[0].details[?name==`networkInterfaceId`].value' \
  --output text | xargs -I {} aws ec2 describe-network-interfaces \
  --network-interface-ids {} \
  --query 'NetworkInterfaces[0].Association.PublicIp' \
  --output text)

echo "ECS Task Public IP: $PUBLIC_IP"
```

Or use Terraform output:
```bash
terraform output ecs_task_public_ip
```

### 5.2 Test gRPC Service

```bash
# List services
grpcurl -plaintext $PUBLIC_IP:50051 list

# Create a user
grpcurl -plaintext -d '{
  "email": "test@example.com",
  "first_name": "Test",
  "last_name": "User"
}' $PUBLIC_IP:50051 user.UserService/CreateUser

# Get user (use ID from create response)
grpcurl -plaintext -d '{
  "user_id": "YOUR_USER_ID"
}' $PUBLIC_IP:50051 user.UserService/GetUser

# List users
grpcurl -plaintext -d '{
  "offset": 0,
  "limit": 10
}' $PUBLIC_IP:50051 user.UserService/ListUsers
```

### 5.3 Check Logs

```bash
# View CloudWatch Logs
aws logs tail /ecs/grpc-go --follow

# Or in AWS Console:
# CloudWatch > Log groups > /ecs/grpc-go
```

### 5.4 Verify Database Connection

```bash
# Connect to RDS from your local machine (if allowed)
# Or check ECS logs for successful database connection
aws logs tail /ecs/grpc-go --follow | grep -i "database\|postgres"
```

---

## Cost Estimate

### Monthly Cost Breakdown (Dev Environment)

| Resource | Specification | Monthly Cost |
|----------|--------------|--------------|
| **ECS Fargate** | 1 task, 0.25 vCPU, 512 MB, 24/7 | ~$15 |
| **RDS PostgreSQL** | db.t3.micro, 20 GB storage | ~$15 |
| **Data Transfer** | Estimated 10 GB/month | ~$1 |
| **CloudWatch Logs** | 2 GB ingestion, 7 days retention | ~$1 |
| **ECR Storage** | 2 GB images | ~$0.20 |
| **Secrets Manager** | 1 secret | ~$0.40 |
| **Total** | | **~$33/month** |

### Cost Savings vs Full Setup
- ✅ **No NAT Gateway**: Save ~$32/month
- ✅ **No Load Balancer**: Save ~$16/month
- ✅ **Total Savings**: ~$48/month (~60% cheaper!)

### Production Environment Estimate
- ECS: 2 tasks: ~$30
- RDS: db.t3.small: ~$30
- Data Transfer: ~$5
- **Total: ~$65-70/month**

---

## 💰 Cost Minimization for Learning (IMPORTANT!)

### Overview
If you're deploying this **just for learning** and plan to **stop it after testing**, you can reduce costs to almost zero!

### 💵 Cost Comparison

| Scenario | Cost |
|----------|------|
| **Learning session (3 hours)** | ~$0.10 |
| **5 learning sessions** | ~$0.50 |
| **Storage between sessions** | ~$0.10/month |
| **If left running 24/7** | ~$33/month ❌ |

---

### 🎯 Cost-Saving Strategies

#### 1. Use Fargate Spot (70% cheaper!)

Update your ECS service to use Fargate Spot:

```hcl
# In terraform/modules/ecs/main.tf
resource "aws_ecs_service" "grpc" {
  # ... other config ...
  
  capacity_provider_strategy {
    capacity_provider = "FARGATE_SPOT"
    weight           = 100
  }
}
```

**Savings:**
- Regular Fargate: $0.04048/hour
- Fargate Spot: $0.01214/hour
- **Save 70%!**

#### 2. Stop RDS When Not in Use

```bash
# Stop RDS after testing (can be stopped for up to 7 days)
aws rds stop-db-instance --db-instance-identifier grpc-go-postgres

# Start when needed
aws rds start-db-instance --db-instance-identifier grpc-go-postgres
```

**Savings:** RDS costs $0.017/hour when running, $0/hour when stopped (only pay ~$0.10/month for storage)

#### 3. Scale ECS to 0 When Not Testing

```bash
# Stop all tasks
aws ecs update-service \
  --cluster grpc-go-cluster \
  --service grpc-go-service \
  --desired-count 0

# Start when needed
aws ecs update-service \
  --cluster grpc-go-cluster \
  --service grpc-go-service \
  --desired-count 1
```

#### 4. Use Smallest Possible Sizes

```hcl
# ECS Task (already minimal)
cpu    = "256"   # 0.25 vCPU
memory = "512"   # 512 MB

# RDS (already minimal)
instance_class = "db.t3.micro"
storage        = 20  # GB
```

#### 5. Destroy Everything When Done Learning

```bash
# When completely done, destroy all infrastructure
terraform destroy -var-file=environments/dev.tfvars
```

---

### 🚀 Recommended Workflow for Learning

#### Setup (One-time, ~30 min)
```bash
# 1. Create infrastructure with Terraform
cd terraform
terraform init
terraform apply -var-file=environments/dev.tfvars

# 2. Build and push Docker image
docker build -t grpc-go:latest .
ECR_REPO=$(terraform output -raw ecr_repository_url)
aws ecr get-login-password --region us-east-1 | docker login --username AWS --password-stdin $ECR_REPO
docker tag grpc-go:latest $ECR_REPO:latest
docker push $ECR_REPO:latest
```

#### Learning Session - Start
```bash
# 1. Start RDS (takes ~5 min)
aws rds start-db-instance --db-instance-identifier grpc-go-postgres

# 2. Wait for RDS to be available
aws rds wait db-instance-available --db-instance-identifier grpc-go-postgres

# 3. Start ECS service
aws ecs update-service \
  --cluster grpc-go-cluster \
  --service grpc-go-service \
  --desired-count 1

# 4. Get public IP and test (see Phase 5)
```

#### Learning Session - Stop ⚠️ **CRITICAL**
```bash
# 1. Stop ECS tasks
aws ecs update-service \
  --cluster grpc-go-cluster \
  --service grpc-go-service \
  --desired-count 0

# 2. Stop RDS
aws rds stop-db-instance --db-instance-identifier grpc-go-postgres
```

---

### 🔧 Helper Scripts

Create these scripts to make starting/stopping easier:

#### `scripts/start-dev.sh`
```bash
#!/bin/bash
set -e

echo "🚀 Starting dev environment..."

echo "1️⃣ Starting RDS..."
aws rds start-db-instance --db-instance-identifier grpc-go-postgres

echo "⏳ Waiting for RDS to be available (this takes ~5 min)..."
aws rds wait db-instance-available --db-instance-identifier grpc-go-postgres

echo "2️⃣ Starting ECS service..."
aws ecs update-service \
  --cluster grpc-go-cluster \
  --service grpc-go-service \
  --desired-count 1 \
  --region us-east-1

echo "⏳ Waiting for ECS task to start..."
sleep 30

echo "3️⃣ Getting ECS task public IP..."
TASK_ARN=$(aws ecs list-tasks \
  --cluster grpc-go-cluster \
  --service-name grpc-go-service \
  --query 'taskArns[0]' \
  --output text \
  --region us-east-1)

if [ "$TASK_ARN" != "None" ]; then
  ENI_ID=$(aws ecs describe-tasks \
    --cluster grpc-go-cluster \
    --tasks $TASK_ARN \
    --query 'tasks[0].attachments[0].details[?name==`networkInterfaceId`].value' \
    --output text \
    --region us-east-1)
  
  PUBLIC_IP=$(aws ec2 describe-network-interfaces \
    --network-interface-ids $ENI_ID \
    --query 'NetworkInterfaces[0].Association.PublicIp' \
    --output text \
    --region us-east-1)
  
  echo ""
  echo "✅ Dev environment is ready!"
  echo "📍 ECS Task Public IP: $PUBLIC_IP"
  echo ""
  echo "Test with:"
  echo "  grpcurl -plaintext $PUBLIC_IP:50051 list"
else
  echo "⚠️ No tasks running yet. Wait a bit and check:"
  echo "  aws ecs list-tasks --cluster grpc-go-cluster --service-name grpc-go-service"
fi
```

#### `scripts/stop-dev.sh`
```bash
#!/bin/bash
set -e

echo "🛑 Stopping dev environment..."

echo "1️⃣ Stopping ECS tasks..."
aws ecs update-service \
  --cluster grpc-go-cluster \
  --service grpc-go-service \
  --desired-count 0 \
  --region us-east-1

echo "2️⃣ Stopping RDS..."
aws rds stop-db-instance \
  --db-instance-identifier grpc-go-postgres \
  --region us-east-1

echo "✅ Dev environment stopped!"
echo "💰 You're no longer being charged for compute (only storage: ~$0.10/month)"
```

#### Make scripts executable
```bash
mkdir -p scripts
chmod +x scripts/start-dev.sh
chmod +x scripts/stop-dev.sh
```

---

### ⚠️ Set Up Billing Alarms

**Create a $5 billing alarm:**
```bash
aws cloudwatch put-metric-alarm \
  --alarm-name billing-alarm-5-dollars \
  --alarm-description "Alert when charges exceed $5" \
  --metric-name EstimatedCharges \
  --namespace AWS/Billing \
  --statistic Maximum \
  --period 21600 \
  --threshold 5 \
  --comparison-operator GreaterThanThreshold \
  --evaluation-periods 1 \
  --dimensions Name=Currency,Value=USD
```

**Set up AWS Budget in Console:**
1. Go to AWS Billing Console
2. Budgets > Create budget
3. Cost budget > $10/month
4. Set email alerts at 80% and 100%

---

### 📊 Detailed Cost Breakdown

#### 3-Hour Learning Session

| Resource | Cost/Hour | Hours | Total |
|----------|-----------|-------|-------|
| ECS Fargate Spot (0.25 vCPU, 512 MB) | $0.01214 | 3 | $0.036 |
| RDS db.t3.micro | $0.017 | 3 | $0.051 |
| Data Transfer (minimal) | ~$0.001 | 3 | $0.003 |
| CloudWatch Logs | ~$0.0001 | 3 | $0.0003 |
| ECR Storage (2 GB) | $0.10/month | - | $0.01 |
| **Total** | | | **~$0.10** |

#### If You Forget to Stop

| Resource | Daily Cost | Monthly Cost |
|----------|------------|--------------|
| ECS Fargate Spot | $0.29 | $8.75 |
| RDS db.t3.micro | $0.41 | $12.30 |
| Other | $0.05 | $1.50 |
| **Total** | **$0.75/day** | **~$22/month** |

---

### 💡 Additional Cost-Saving Tips

#### 1. Use AWS Free Tier (First 12 months)
If your account is < 12 months old:
- ✅ 750 hours/month of db.t2.micro RDS (use t2.micro instead of t3.micro)
- ✅ 5 GB of CloudWatch Logs
- ⚠️ Fargate is NOT included in free tier

#### 2. Delete ECR Images When Done
```bash
# Delete all images to save storage costs
aws ecr batch-delete-image \
  --repository-name grpc-go \
  --image-ids imageTag=latest
```

#### 3. Use Terraform Cloud Free Tier
- Free remote state storage (instead of S3)
- No S3 bucket costs
- Sign up at https://app.terraform.io

---

### 📝 Cost Tracking Checklist

**Before each learning session:**
- [ ] Check current AWS bill in console
- [ ] Ensure billing alarm is set ($5 threshold)
- [ ] Run `./scripts/start-dev.sh`

**After each learning session:**
- [ ] Run `./scripts/stop-dev.sh`
- [ ] Verify ECS tasks stopped: `aws ecs list-tasks --cluster grpc-go-cluster`
- [ ] Verify RDS stopped: `aws rds describe-db-instances --db-instance-identifier grpc-go-postgres | grep DBInstanceStatus`

**When completely done learning:**
- [ ] Run `terraform destroy -var-file=environments/dev.tfvars`
- [ ] Delete ECR images
- [ ] Check AWS bill after 24 hours to confirm no charges

---

### 🎯 Summary: Minimize Costs for Learning

1. ✅ **Use Fargate Spot** (70% cheaper)
2. ✅ **Stop RDS** when not in use
3. ✅ **Scale ECS to 0** when not in use
4. ✅ **Use helper scripts** to start/stop easily
5. ✅ **Set billing alarms** ($5 threshold)
6. ✅ **Destroy everything** when done learning

**Expected Total Cost:**
- Per 3-hour session: ~$0.10
- 5 learning sessions: ~$0.50
- Storage between sessions: ~$0.10/month
- **Total for learning: < $2** 🎉

---

## Security Considerations

### ⚠️ Important Notes

**Public IP for ECS Tasks:**
- ✅ **Pros**: Simple, no NAT Gateway cost, easy to test
- ⚠️ **Cons**: Tasks are directly exposed to internet
- 🔒 **Mitigation**: 
  - Use security groups to restrict access
  - Consider IP whitelisting if you know client IPs
  - Add TLS/SSL for production
  - Consider adding authentication/authorization

**Recommended Security Groups:**
```hcl
# ECS Tasks Security Group
ingress {
  from_port   = 50051
  to_port     = 50051
  protocol    = "tcp"
  cidr_blocks = ["0.0.0.0/0"]  # For dev/testing
  # cidr_blocks = ["YOUR_IP/32"]  # For production (whitelist)
}

egress {
  from_port   = 0
  to_port     = 0
  protocol    = "-1"
  cidr_blocks = ["0.0.0.0/0"]
}
```

### Best Practices
1. ✅ RDS in private subnet (not directly accessible)
2. ✅ Use Secrets Manager for database password
3. ✅ Enable CloudWatch Logs
4. ✅ Use IAM roles (no hardcoded credentials)
5. ✅ Enable RDS encryption at rest
6. ⚠️ Consider adding TLS for gRPC in production
7. ⚠️ Consider VPN or bastion host for production

---

## Scaling Strategy

### Horizontal Scaling
```bash
# Scale up to 2 tasks
aws ecs update-service \
  --cluster grpc-go-cluster \
  --service grpc-go-service \
  --desired-count 2

# Scale down to 1 task
aws ecs update-service \
  --cluster grpc-go-cluster \
  --service grpc-go-service \
  --desired-count 1
```

**Note**: Without a load balancer, clients need to know all task IPs. For production with multiple tasks, you'd want:
- Service discovery (AWS Cloud Map)
- Or add a load balancer later
- Or use DNS round-robin

### Vertical Scaling
Update task definition with more CPU/memory:
```hcl
# In terraform/modules/ecs/main.tf
cpu    = "512"   # 0.5 vCPU
memory = "1024"  # 1 GB
```

Then apply:
```bash
terraform apply -var-file=environments/dev.tfvars
```

---

## Monitoring

### CloudWatch Metrics to Watch
1. **ECS Metrics**:
   - CPUUtilization
   - MemoryUtilization
   - RunningTaskCount

2. **RDS Metrics**:
   - DatabaseConnections
   - CPUUtilization
   - FreeableMemory

### CloudWatch Alarms (Optional)
```bash
# Create alarm for high CPU
aws cloudwatch put-metric-alarm \
  --alarm-name ecs-high-cpu \
  --alarm-description "Alert when CPU exceeds 80%" \
  --metric-name CPUUtilization \
  --namespace AWS/ECS \
  --statistic Average \
  --period 300 \
  --threshold 80 \
  --comparison-operator GreaterThanThreshold \
  --evaluation-periods 2
```

---

## Troubleshooting

### Issue: ECS task fails to start
**Check:**
```bash
# View task logs
aws logs tail /ecs/grpc-go --follow

# Describe task to see stopped reason
aws ecs describe-tasks \
  --cluster grpc-go-cluster \
  --tasks <task-arn>
```

**Common causes:**
- Image not found in ECR
- Invalid environment variables
- Insufficient IAM permissions
- Port already in use

### Issue: Cannot connect to gRPC service
**Check:**
```bash
# Verify task is running
aws ecs list-tasks --cluster grpc-go-cluster

# Check security group allows port 50051
aws ec2 describe-security-groups --group-ids <sg-id>

# Verify public IP is assigned
aws ecs describe-tasks --cluster grpc-go-cluster --tasks <task-arn>
```

### Issue: Database connection failed
**Check:**
```bash
# Verify RDS endpoint
terraform output rds_endpoint

# Check security group allows ECS -> RDS
# Check database credentials in Secrets Manager
aws secretsmanager get-secret-value --secret-id grpc-go-db-password
```

### Issue: High costs
**Check:**
```bash
# View cost breakdown in AWS Cost Explorer
# Common culprits:
# - Data transfer (use CloudWatch to monitor)
# - Running tasks when not needed (scale to 0 for dev)
# - RDS instance running 24/7 (consider stopping for dev)
```

---

## Useful Commands

### Terraform
```bash
# Initialize
terraform init

# Plan
terraform plan -var-file=environments/dev.tfvars

# Apply
terraform apply -var-file=environments/dev.tfvars

# Destroy (cleanup)
terraform destroy -var-file=environments/dev.tfvars

# Show outputs
terraform output

# Format code
terraform fmt -recursive
```

### AWS ECS
```bash
# List clusters
aws ecs list-clusters

# List services
aws ecs list-services --cluster grpc-go-cluster

# List tasks
aws ecs list-tasks --cluster grpc-go-cluster --service-name grpc-go-service

# Describe service
aws ecs describe-services --cluster grpc-go-cluster --services grpc-go-service

# Update service (force new deployment)
aws ecs update-service --cluster grpc-go-cluster --service grpc-go-service --force-new-deployment

# Scale service
aws ecs update-service --cluster grpc-go-cluster --service grpc-go-service --desired-count 2

# Stop task
aws ecs stop-task --cluster grpc-go-cluster --task <task-arn>
```

### AWS RDS
```bash
# Describe DB instance
aws rds describe-db-instances --db-instance-identifier grpc-go-postgres

# Stop DB (dev only, saves cost)
aws rds stop-db-instance --db-instance-identifier grpc-go-postgres

# Start DB
aws rds start-db-instance --db-instance-identifier grpc-go-postgres
```

### Docker & ECR
```bash
# Build
docker build -t grpc-go:latest .

# Login to ECR
aws ecr get-login-password --region us-east-1 | docker login --username AWS --password-stdin <account-id>.dkr.ecr.us-east-1.amazonaws.com

# Tag
docker tag grpc-go:latest <ecr-repo>:latest

# Push
docker push <ecr-repo>:latest

# List images
aws ecr list-images --repository-name grpc-go
```

### CloudWatch Logs
```bash
# Tail logs
aws logs tail /ecs/grpc-go --follow

# Get logs for specific time range
aws logs tail /ecs/grpc-go --since 1h

# Search logs
aws logs filter-log-events --log-group-name /ecs/grpc-go --filter-pattern "ERROR"
```

---

## Next Steps

### Ready to implement?

**Option 1: I create all Terraform files for you**
- I'll create all the modules and configuration files
- You just need to run `terraform apply`

**Option 2: Step-by-step guidance**
- I guide you through creating each file
- You learn Terraform as you go

**Option 3: Hybrid approach**
- I create the boilerplate
- You customize for your needs

Which approach would you prefer? 🚀

---

## Migration Path (Future)

When you're ready to add more features:

### Add Load Balancer
1. Create ALB/NLB module
2. Move ECS tasks to private subnets
3. Add NAT Gateway
4. Update security groups
5. Point DNS to load balancer

### Add Auto Scaling
1. Create auto-scaling policies
2. Define target tracking metrics
3. Set min/max task counts

### Add CI/CD
1. Create GitHub Actions workflow
2. Automate Docker build and push
3. Automate ECS deployment

### Add Monitoring
1. Set up CloudWatch Alarms
2. Add AWS X-Ray tracing
3. Create CloudWatch Dashboard

---

**Document Version**: 1.0 (Simplified)  
**Last Updated**: 2025-11-25  
**Focus**: ECS + RDS only, no load balancer
