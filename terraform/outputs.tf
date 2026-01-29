output "vpc_id" {
  description = "VPC ID"
  value       = aws_vpc.main.id
}

output "alb_dns_name" {
  description = "ALB DNS name"
  value       = aws_lb.main.dns_name
}

output "alb_zone_id" {
  description = "ALB Zone ID"
  value       = aws_lb.main.zone_id
}

output "api_gateway_url" {
  description = "API Gateway URL"
  value       = var.enable_api_gateway ? aws_api_gateway_stage.main[0].invoke_url : "API Gateway not enabled"
}

output "api_gateway_custom_domain" {
  description = "API Gateway custom domain name"
  value       = var.enable_api_gateway && var.certificate_arn != "" ? aws_api_gateway_domain_name.main[0].domain_name : "Custom domain not configured"
}

output "rds_endpoint" {
  description = "RDS endpoint"
  value       = aws_db_instance.main.endpoint
  sensitive   = true
}

output "redis_endpoint" {
  description = "Redis endpoint"
  value       = aws_elasticache_cluster.main.cache_nodes[0].address
  sensitive   = true
}

output "ecr_repository_url" {
  description = "ECR repository URL"
  value       = aws_ecr_repository.api.repository_url
}

output "ecs_cluster_name" {
  description = "ECS cluster name"
  value       = aws_ecs_cluster.main.name
}

output "ecs_service_name" {
  description = "ECS service name"
  value       = aws_ecs_service.api.name
}

output "cloudwatch_log_group" {
  description = "CloudWatch log group"
  value       = aws_cloudwatch_log_group.ecs.name
}

output "waf_web_acl_id" {
  description = "WAF Web ACL ID"
  value       = var.enable_api_gateway && var.enable_waf ? aws_wafv2_web_acl.api_gateway[0].id : "WAF not enabled"
}
