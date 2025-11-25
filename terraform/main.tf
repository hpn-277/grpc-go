# Main Terraform configuration
# This file orchestrates all modules

locals {
  name_prefix = "${var.project_name}-${var.environment}"
  
  common_tags = {
    Project     = var.project_name
    Environment = var.environment
  }
}

# Networking Module - VPC, Subnets, Security Groups
module "networking" {
  source = "./modules/networking"

  project_name       = var.project_name
  environment        = var.environment
  vpc_cidr           = var.vpc_cidr
  availability_zones = var.availability_zones
  container_port     = var.container_port
}

# ECR Module - Container Registry
module "ecr" {
  source = "./modules/ecr"

  project_name = var.project_name
  environment  = var.environment
}

# RDS Module - PostgreSQL Database
module "rds" {
  source = "./modules/rds"

  project_name         = var.project_name
  environment          = var.environment
  vpc_id               = module.networking.vpc_id
  private_subnet_ids   = module.networking.private_subnet_ids
  db_instance_class    = var.db_instance_class
  db_allocated_storage = var.db_allocated_storage
  db_name              = var.db_name
  db_username          = var.db_username
  db_password          = var.db_password
  
  # Allow access from ECS tasks
  allowed_security_group_ids = [module.ecs.task_security_group_id]
}

# ECS Module - Fargate Service
module "ecs" {
  source = "./modules/ecs"

  project_name        = var.project_name
  environment         = var.environment
  vpc_id              = module.networking.vpc_id
  public_subnet_ids   = module.networking.public_subnet_ids
  container_port      = var.container_port
  ecs_task_cpu        = var.ecs_task_cpu
  ecs_task_memory     = var.ecs_task_memory
  ecs_desired_count   = var.ecs_desired_count
  use_fargate_spot    = var.use_fargate_spot
  log_retention_days  = var.log_retention_days
  
  # ECR image
  ecr_repository_url = module.ecr.repository_url
  image_tag          = "latest"
  
  # Database connection
  db_host     = module.rds.db_endpoint
  db_port     = module.rds.db_port
  db_name     = module.rds.db_name
  db_username = var.db_username
  db_password = var.db_password
}
