# Terraform Backend Configuration
# 
# OPTIONAL: Uncomment this to use S3 backend for remote state storage
# For learning, you can use local state (default)
#
# To use S3 backend:
# 1. Create S3 bucket: aws s3 mb s3://your-terraform-state-bucket
# 2. Create DynamoDB table for locking:
#    aws dynamodb create-table \
#      --table-name terraform-state-lock \
#      --attribute-definitions AttributeName=LockID,AttributeType=S \
#      --key-schema AttributeName=LockID,KeyType=HASH \
#      --provisioned-throughput ReadCapacityUnits=5,WriteCapacityUnits=5
# 3. Uncomment the block below and update bucket name
# 4. Run: terraform init -migrate-state

# terraform {
#   backend "s3" {
#     bucket         = "your-terraform-state-bucket"
#     key            = "grpc-go/terraform.tfstate"
#     region         = "us-east-1"
#     dynamodb_table = "terraform-state-lock"
#     encrypt        = true
#   }
# }

# For now, using local backend (default)
# State file will be stored in terraform.tfstate locally
