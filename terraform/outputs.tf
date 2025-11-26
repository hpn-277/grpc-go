output "vpc_id" {
  description = "VPC ID"
  value       = module.networking.vpc_id
}

output "vpc_cidr" {
  description = "VPC CIDR block"
  value       = module.networking.vpc_cidr
}


output "public_subnet_ids" {
  description = "Public subnet IDs"
  value       = module.networking.public_subnet_ids
}


# output "private_subnet_ids" {
#   description = "Private subnet IDs"
#   value       = module.networking.private_subnet_ids
# }

# TODO: Uncomment after creating ECR module
# output "ecr_repository_url" {
#   description = "ECR repository URL"
#   value       = module.ecr.repository_url
# }

# output "ecr_repository_name" {
#   description = "ECR repository name"
#   value       = module.ecr.repository_name
# }

# TODO: Uncomment after creating RDS module
# output "rds_endpoint" {
#   description = "RDS endpoint"
#   value       = module.rds.db_endpoint
# }

# output "rds_database_name" {
#   description = "RDS database name"
#   value       = module.rds.db_name
# }

# TODO: Uncomment after creating ECS module
# output "ecs_cluster_name" {
#   description = "ECS cluster name"
#   value       = module.ecs.cluster_name
# }

# output "ecs_service_name" {
#   description = "ECS service name"
#   value       = module.ecs.service_name
# }

# output "ecs_task_security_group_id" {
#   description = "ECS task security group ID"
#   value       = module.ecs.task_security_group_id
# }

# output "cloudwatch_log_group" {
#   description = "CloudWatch log group name"
#   value       = module.ecs.log_group_name
# }

# Helper outputs for testing
# TODO: Uncomment after all modules are created
# output "quick_start_commands" {
#   description = "Quick start commands for testing"
#   value = <<-EOT
#     # 1. Build and push Docker image:
#     aws ecr get-login-password --region ${var.aws_region} | docker login --username AWS --password-stdin ${module.ecr.repository_url}
#     docker build -t ${var.project_name}:latest .
#     docker tag ${var.project_name}:latest ${module.ecr.repository_url}:latest
#     docker push ${module.ecr.repository_url}:latest
#
#     # 2. Get ECS task public IP (after deployment):
#     ./scripts/get-task-ip.sh
#
#     # 3. Test gRPC service:
#     grpcurl -plaintext <TASK_IP>:50051 list
#
#     # 4. View logs:
#     aws logs tail ${module.ecs.log_group_name} --follow
#
#     # 5. Stop environment (save costs):
#     ./scripts/stop-dev.sh
#   EOT
# }
