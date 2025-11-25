# AWS ECS Deployment Plan with Terraform

## Overview
This document outlines the complete plan to deploy the gRPC Go service to AWS ECS using Terraform, including infrastructure setup, containerization, CI/CD, and monitoring.

---

## 📋 Table of Contents
1. [Prerequisites & Setup](#phase-1-prerequisites--setup)
2. [Containerization](#phase-2-containerization)
3. [Terraform Infrastructure](#phase-3-terraform-infrastructure)
4. [CI/CD Pipeline](#phase-4-cicd-pipeline)
5. [Deployment Steps](#phase-5-deployment-steps)
6. [Monitoring & Observability](#phase-6-monitoring--observability)
7. [Cost Estimate](#cost-estimate)

---

## Phase 1: Prerequisites & Setup (30 min)

### 1.1 Project Structure
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
│   ├── alb/              # Network Load Balancer (for gRPC)
│   │   ├── main.tf
│   │   ├── variables.tf
│   │   └── outputs.tf
│   └── ecr/              # Container registry
│       ├── main.tf
│       ├── variables.tf
│       └── outputs.tf
└── environments/
    ├── dev.tfvars
    ├── staging.tfvars
    └── prod.tfvars
```

### 1.2 Tools to Install
```bash
# Terraform CLI
brew install terraform

# AWS CLI
brew install awscli

# Docker (if not already installed)
brew install --cask docker

# grpcurl (for testing)
brew install grpcurl
```

### 1.3 AWS Setup
- Configure AWS credentials: `aws configure`
- Create S3 bucket for Terraform state
- Create DynamoDB table for state locking
- Set up AWS IAM user with appropriate permissions

### 1.4 Required AWS Permissions
- EC2 (VPC, Subnets, Security Groups)
- ECS (Cluster, Service, Task Definition)
- ECR (Repository)
- RDS (Database Instance)
- ELB (Network Load Balancer)
- IAM (Roles, Policies)
- CloudWatch (Logs, Metrics)
- Secrets Manager (Database credentials)

---

## Phase 2: Containerization (30 min)

### 2.1 Create Dockerfile
Multi-stage build for optimal image size and security.

**Key Features:**
- Multi-stage build (builder + runtime)
- Minimal Alpine Linux base image
- Non-root user execution
- Health check endpoint
- Optimized layer caching

### 2.2 Create .dockerignore
Exclude unnecessary files from Docker context to speed up builds.

### 2.3 Build & Test Locally
```bash
# Build the image
docker build -t grpc-go:latest .

# Run locally
docker run -p 50051:50051 \
  -e DB_HOST=host.docker.internal \
  -e DB_PORT=5432 \
  -e DB_USER=superuser \
  -e DB_PASSWORD=superpass \
  -e DB_NAME=super_salary_db \
  grpc-go:latest

# Test with grpcurl
grpcurl -plaintext localhost:50051 list
```

---

## Phase 3: Terraform Infrastructure (2-3 hours)

### 3.1 Backend Configuration
- **S3 Bucket**: Store Terraform state
- **DynamoDB Table**: State locking to prevent concurrent modifications
- **Encryption**: Enable encryption at rest

### 3.2 Networking Module
**Components:**
- VPC with CIDR block (e.g., 10.0.0.0/16)
- Public subnets (2 AZs) for Load Balancer
- Private subnets (2 AZs) for ECS tasks and RDS
- Internet Gateway for public subnets
- NAT Gateways (1 per AZ) for private subnet internet access
- Route tables for public and private subnets

**Security Groups:**
- **ECS Tasks SG**: Allow inbound 50051 from NLB
- **RDS SG**: Allow inbound 5432 from ECS Tasks
- **NLB SG**: Allow inbound 50051 from internet

### 3.3 ECR Module
**Components:**
- ECR repository for Docker images
- Image scanning on push
- Lifecycle policy to clean up old images
- Repository policy for cross-account access (if needed)

### 3.4 RDS Module
**Components:**
- PostgreSQL RDS instance
- Multi-AZ deployment (for production)
- Automated backups
- Encryption at rest
- DB subnet group in private subnets
- Parameter group for PostgreSQL tuning
- Secrets Manager for database credentials

**Instance Sizing:**
- Dev: db.t3.micro (1 vCPU, 1 GB RAM)
- Staging: db.t3.small (2 vCPU, 2 GB RAM)
- Prod: db.t3.medium or larger

### 3.5 ECS Module
**Components:**

**ECS Cluster:**
- Fargate launch type (serverless)
- Container Insights enabled

**Task Definition:**
- CPU: 256 (0.25 vCPU) for dev, scale up for prod
- Memory: 512 MB for dev, scale up for prod
- Network mode: awsvpc
- Container port: 50051
- Environment variables from Terraform
- Secrets from AWS Secrets Manager
- CloudWatch Logs integration

**ECS Service:**
- Desired count: 1 for dev, 2+ for prod
- Auto-scaling policies (optional)
- Load balancer integration
- Health check grace period
- Deployment configuration (rolling update)

### 3.6 Load Balancer Module
**Network Load Balancer (NLB):**
- Better for gRPC than Application Load Balancer
- TCP protocol support
- Low latency
- Static IP addresses

**Components:**
- NLB in public subnets
- Target group for ECS tasks (IP target type)
- Listener on port 50051
- Health checks (TCP)

### 3.7 IAM Roles & Policies
**ECS Task Execution Role:**
- Pull images from ECR
- Write logs to CloudWatch
- Read secrets from Secrets Manager

**ECS Task Role:**
- Application-specific permissions
- Access to AWS services (if needed)

---

## Phase 4: CI/CD Pipeline (1 hour)

### 4.1 GitHub Actions Workflow
**Triggers:**
- Push to `main` branch
- Manual workflow dispatch

**Steps:**
1. Checkout code
2. Configure AWS credentials
3. Login to Amazon ECR
4. Run tests (optional)
5. Build Docker image
6. Tag image with git SHA and `latest`
7. Push image to ECR
8. Update ECS service (force new deployment)

### 4.2 Environment Secrets
Store in GitHub Secrets:
- `AWS_ACCESS_KEY_ID`
- `AWS_SECRET_ACCESS_KEY`
- `AWS_REGION`

### 4.3 Multi-Environment Support
- Separate workflows for dev/staging/prod
- Environment-specific variables
- Manual approval for production deployments

---

## Phase 5: Deployment Steps (30 min)

### 5.1 Initial Setup
```bash
# Create S3 bucket for Terraform state
aws s3 mb s3://your-terraform-state-bucket --region us-east-1

# Create DynamoDB table for state locking
aws dynamodb create-table \
  --table-name terraform-state-lock \
  --attribute-definitions AttributeName=LockID,AttributeType=S \
  --key-schema AttributeName=LockID,KeyType=HASH \
  --provisioned-throughput ReadCapacityUnits=5,WriteCapacityUnits=5 \
  --region us-east-1
```

### 5.2 Initialize Terraform
```bash
cd terraform
terraform init
```

### 5.3 Plan Infrastructure
```bash
terraform plan -var-file=environments/dev.tfvars -out=tfplan
```

### 5.4 Apply Infrastructure
```bash
terraform apply tfplan
```

### 5.5 Build & Push Docker Image
```bash
# Get ECR repository URL from Terraform output
ECR_REPO=$(terraform output -raw ecr_repository_url)

# Login to ECR
aws ecr get-login-password --region us-east-1 | \
  docker login --username AWS --password-stdin $ECR_REPO

# Build and push
docker build -t grpc-go:latest .
docker tag grpc-go:latest $ECR_REPO:latest
docker push $ECR_REPO:latest
```

### 5.6 Verify Deployment
```bash
# Get NLB DNS name
NLB_DNS=$(terraform output -raw nlb_dns_name)

# Test gRPC service
grpcurl -plaintext $NLB_DNS:50051 list

# Create a user
grpcurl -plaintext -d '{
  "email": "test@example.com",
  "first_name": "Test",
  "last_name": "User"
}' $NLB_DNS:50051 user.UserService/CreateUser
```

---

## Phase 6: Monitoring & Observability (1 hour)

### 6.1 CloudWatch Logs
- Log group: `/ecs/grpc-go`
- Log retention: 7 days (dev), 30 days (prod)
- View logs: `aws logs tail /ecs/grpc-go --follow`

### 6.2 CloudWatch Metrics
**ECS Metrics:**
- CPUUtilization
- MemoryUtilization
- Running task count

**RDS Metrics:**
- DatabaseConnections
- CPUUtilization
- FreeableMemory
- ReadLatency / WriteLatency

**NLB Metrics:**
- ActiveFlowCount
- ProcessedBytes
- HealthyHostCount
- UnHealthyHostCount

### 6.3 CloudWatch Alarms
- High CPU utilization (> 80%)
- High memory utilization (> 80%)
- Unhealthy target count (> 0)
- RDS connection count (> 80% of max)

### 6.4 AWS X-Ray (Optional)
- Distributed tracing
- Service map visualization
- Performance insights
- Error analysis

### 6.5 Application Logging Best Practices
- Structured logging (JSON format)
- Log levels: DEBUG, INFO, WARN, ERROR
- Include request IDs for tracing
- Log database query performance

---

## Cost Estimate

### Monthly Cost Breakdown (Dev Environment)

| Resource | Specification | Monthly Cost |
|----------|--------------|--------------|
| **ECS Fargate** | 1 task, 0.25 vCPU, 512 MB, 24/7 | ~$15 |
| **RDS PostgreSQL** | db.t3.micro, 20 GB storage | ~$15 |
| **Network Load Balancer** | 1 NLB | ~$16 |
| **NAT Gateway** | 1 NAT Gateway | ~$32 |
| **Data Transfer** | Estimated 50 GB/month | ~$5 |
| **CloudWatch Logs** | 5 GB ingestion, 7 days retention | ~$3 |
| **ECR Storage** | 5 GB images | ~$0.50 |
| **Secrets Manager** | 1 secret | ~$0.40 |
| **Total** | | **~$87/month** |

### Cost Optimization Tips
1. **Use Fargate Spot** for non-critical workloads (70% savings)
2. **Single NAT Gateway** for dev (not recommended for prod)
3. **RDS Reserved Instances** for production (up to 60% savings)
4. **Auto-scaling** to scale down during off-hours
5. **S3 Lifecycle Policies** for old Docker images in ECR

### Production Environment Estimate
- ECS: 2-4 tasks with auto-scaling: ~$60-120
- RDS: db.t3.medium Multi-AZ: ~$120
- NLB: ~$16
- NAT Gateway (2 AZs): ~$64
- Data Transfer: ~$20
- **Total: ~$280-340/month**

---

## Security Best Practices

### 6.1 Network Security
- ✅ Private subnets for ECS tasks and RDS
- ✅ Security groups with least privilege
- ✅ No public IP for ECS tasks
- ✅ VPC endpoints for AWS services (optional, reduces NAT costs)

### 6.2 Application Security
- ✅ Secrets in AWS Secrets Manager (not environment variables)
- ✅ IAM roles with least privilege
- ✅ Container image scanning
- ✅ Non-root container user
- ✅ Read-only root filesystem (optional)

### 6.3 Data Security
- ✅ RDS encryption at rest
- ✅ RDS encryption in transit (SSL/TLS)
- ✅ Automated backups
- ✅ Multi-AZ deployment for production

### 6.4 Compliance
- ✅ CloudTrail for audit logging
- ✅ Config for compliance monitoring
- ✅ GuardDuty for threat detection (optional)

---

## Disaster Recovery

### 7.1 Backup Strategy
- **RDS Automated Backups**: 7 days retention (prod)
- **RDS Snapshots**: Weekly manual snapshots
- **Terraform State**: Versioned in S3
- **Docker Images**: Retained in ECR with lifecycle policy

### 7.2 Recovery Procedures
1. **Database Restore**: Restore from RDS snapshot
2. **Infrastructure Restore**: `terraform apply` from version control
3. **Application Rollback**: Deploy previous Docker image tag

### 7.3 High Availability
- **Multi-AZ RDS**: Automatic failover
- **Multi-AZ ECS**: Tasks distributed across AZs
- **NLB**: Cross-zone load balancing enabled

---

## Scaling Strategy

### 8.1 Horizontal Scaling (ECS Auto Scaling)
**Target Tracking Policies:**
- CPU utilization target: 70%
- Memory utilization target: 80%
- Custom metric: Request count per task

**Scaling Limits:**
- Min tasks: 2 (prod), 1 (dev)
- Max tasks: 10 (prod), 2 (dev)

### 8.2 Vertical Scaling
**ECS Task Sizing:**
- Dev: 0.25 vCPU, 512 MB
- Staging: 0.5 vCPU, 1 GB
- Prod: 1 vCPU, 2 GB (adjust based on load testing)

**RDS Scaling:**
- Start with db.t3.micro (dev)
- Monitor performance metrics
- Upgrade instance class as needed
- Consider read replicas for read-heavy workloads

---

## Next Steps

### Immediate Actions
1. ✅ Review this plan
2. ✅ Set up AWS account and credentials
3. ✅ Create Terraform directory structure
4. ✅ Create Dockerfile

### Phase-by-Phase Implementation
1. **Week 1**: Containerization + Local Testing
2. **Week 2**: Terraform Infrastructure (Networking + ECR)
3. **Week 3**: Terraform Infrastructure (ECS + RDS + NLB)
4. **Week 4**: CI/CD Pipeline + Monitoring
5. **Week 5**: Testing + Documentation + Optimization

### Learning Resources
- [AWS ECS Best Practices](https://docs.aws.amazon.com/AmazonECS/latest/bestpracticesguide/intro.html)
- [Terraform AWS Provider Docs](https://registry.terraform.io/providers/hashicorp/aws/latest/docs)
- [gRPC on AWS](https://aws.amazon.com/blogs/compute/load-balancing-grpc-traffic-with-aws-app-mesh/)

---

## Questions to Consider

Before proceeding, decide on:
1. **AWS Region**: Which region to deploy to? (e.g., us-east-1, ap-southeast-1)
2. **Environment Strategy**: Single account or multi-account setup?
3. **Domain Name**: Do you need a custom domain for the gRPC endpoint?
4. **TLS/SSL**: Do you need TLS termination at the load balancer?
5. **Budget**: What's your monthly budget for this service?
6. **Scaling Requirements**: Expected traffic and growth projections?

---

## Appendix

### A. Useful Commands

```bash
# Terraform
terraform init
terraform plan -var-file=environments/dev.tfvars
terraform apply -var-file=environments/dev.tfvars
terraform destroy -var-file=environments/dev.tfvars
terraform output

# AWS CLI
aws ecs list-clusters
aws ecs list-services --cluster grpc-go-cluster
aws ecs describe-services --cluster grpc-go-cluster --services grpc-go-service
aws logs tail /ecs/grpc-go --follow

# Docker
docker build -t grpc-go .
docker run -p 50051:50051 grpc-go
docker ps
docker logs <container-id>

# grpcurl
grpcurl -plaintext <endpoint>:50051 list
grpcurl -plaintext <endpoint>:50051 describe user.UserService
```

### B. Troubleshooting

**Issue: ECS tasks failing to start**
- Check CloudWatch Logs for error messages
- Verify security group rules
- Ensure ECR image exists and is accessible
- Check IAM role permissions

**Issue: Cannot connect to RDS**
- Verify security group allows traffic from ECS tasks
- Check RDS endpoint and port
- Verify database credentials in Secrets Manager
- Ensure RDS is in the same VPC as ECS

**Issue: High costs**
- Review NAT Gateway usage (consider VPC endpoints)
- Check for idle resources
- Review CloudWatch Logs retention
- Consider Fargate Spot for dev environments

---

**Document Version**: 1.0  
**Last Updated**: 2025-11-25  
**Author**: Deployment Planning Assistant
