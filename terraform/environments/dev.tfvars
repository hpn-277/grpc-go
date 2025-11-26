# Development Environment Configuration

aws_region         = "ap-southeast-2"
environment        = "dev"
project_name       = "grpc-go"

# Network Configuration
vpc_cidr           = "10.0.0.0/16"
availability_zones = ["ap-southeast-2a", "ap-southeast-2b"]

# ECS Configuration (Minimal for cost savings)
ecs_task_cpu       = 256   # 0.25 vCPU
ecs_task_memory    = 512   # 512 MB
ecs_desired_count  = 1     # Single task for dev
use_fargate_spot   = true  # 70% cost savings!

# RDS Configuration (Minimal for cost savings)
db_instance_class    = "db.t3.micro"  # Smallest instance
db_allocated_storage = 20             # Minimum storage
db_name              = "super_salary_db"
db_username          = "superuser"
db_password          = "temporary-password-change-later"  # Temporary default, will be required when RDS module is enabled

# Application Configuration
container_port      = 50051
log_retention_days  = 7  # Short retention for dev

# Cost Optimization Notes:
# - Using Fargate Spot saves 70% on compute costs
# - db.t3.micro is the smallest RDS instance
# - Single task deployment
# - Remember to stop RDS when not in use: aws rds stop-db-instance --db-instance-identifier grpc-go-dev-postgres
# - Scale ECS to 0 when not testing: aws ecs update-service --desired-count 0
