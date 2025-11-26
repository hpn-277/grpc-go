# Networking Module - VPC and Network Resources

locals {
  name_prefix = "${var.project_name}-${var.environment}"
  
  common_tags = {
    Project     = var.project_name
    Environment = var.environment
    ManagedBy   = "Terraform"
  }
}

# ============================================
# VPC
# ============================================

resource "aws_vpc" "main" {
  cidr_block           = var.vpc_cidr
  enable_dns_hostnames = true
  enable_dns_support   = true

  tags = merge(
    local.common_tags,
    {
      Name = "${local.name_prefix}-vpc"
    }
  )
}

# ============================================
# Public Subnets (for ECS tasks)
# ============================================

# Calculate subnet CIDR blocks
# VPC: 10.0.0.0/16 (65,536 IPs)
# Public Subnet 1: 10.0.1.0/24 (256 IPs)
# Public Subnet 2: 10.0.2.0/24 (256 IPs)

resource "aws_subnet" "public" {
  count = length(var.availability_zones)

  vpc_id                  = aws_vpc.main.id
  cidr_block              = cidrsubnet(var.vpc_cidr, 8, count.index + 1)
  availability_zone       = var.availability_zones[count.index]
  map_public_ip_on_launch = true  # Auto-assign public IPs to instances

  tags = merge(
    local.common_tags,
    {
      Name = "${local.name_prefix}-public-subnet-${count.index + 1}"
      Type = "Public"
    }
  )
}

# ============================================
# Private Subnets (for RDS)
# ============================================

# Calculate subnet CIDR blocks
# Private Subnet 1: 10.0.11.0/24 (256 IPs)
# Private Subnet 2: 10.0.12.0/24 (256 IPs)

resource "aws_subnet" "private" {
  count = length(var.availability_zones)

  vpc_id                  = aws_vpc.main.id
  cidr_block              = cidrsubnet(var.vpc_cidr, 8, count.index + 11)
  availability_zone       = var.availability_zones[count.index]
  map_public_ip_on_launch = false  # No public IPs for private subnets

  tags = merge(
    local.common_tags,
    {
      Name = "${local.name_prefix}-private-subnet-${count.index + 1}"
      Type = "Private"
    }
  )
}

# ============================================
# Internet Gateway (for public subnet internet access)
# ============================================

resource "aws_internet_gateway" "main" {
  vpc_id = aws_vpc.main.id

  tags = merge(
    local.common_tags,
    {
      Name = "${local.name_prefix}-igw"
    }
  )
}

# ============================================
# Route Tables
# ============================================

# Public Route Table (routes to Internet Gateway)
resource "aws_route_table" "public" {
  vpc_id = aws_vpc.main.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.main.id
  }

  tags = merge(
    local.common_tags,
    {
      Name = "${local.name_prefix}-public-rt"
      Type = "Public"
    }
  )
}

# Associate public subnets with public route table
resource "aws_route_table_association" "public" {
  count = length(var.availability_zones)

  subnet_id      = aws_subnet.public[count.index].id
  route_table_id = aws_route_table.public.id
}

# ============================================
# Security Groups
# ============================================

# 1. ECS Task Security Group
resource "aws_security_group" "ecs_tasks" {
  name        = "${local.name_prefix}-ecs-tasks-sg"
  description = "Allow inbound traffic to ECS tasks"
  vpc_id      = aws_vpc.main.id

  # Inbound: Allow gRPC traffic from anywhere
  ingress {
    description = "Allow gRPC traffic from anywhere"
    from_port   = var.container_port
    to_port     = var.container_port
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  # Outbound: Allow all traffic (needed for pulling images, reaching RDS, etc.)
  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = merge(
    local.common_tags,
    {
      Name = "${local.name_prefix}-ecs-tasks-sg"
    }
  )
}

# 2. RDS Security Group
resource "aws_security_group" "rds" {
  name        = "${local.name_prefix}-rds-sg"
  description = "Allow inbound traffic to RDS from ECS tasks only"
  vpc_id      = aws_vpc.main.id

  # Inbound: Allow PostgreSQL traffic ONLY from ECS tasks
  ingress {
    description     = "Allow PostgreSQL from ECS tasks"
    from_port       = 5432
    to_port         = 5432
    protocol        = "tcp"
    security_groups = [aws_security_group.ecs_tasks.id]
  }

  tags = merge(
    local.common_tags,
    {
      Name = "${local.name_prefix}-rds-sg"
    }
  )
}
