/**
 * DNS Zone module - SSM parameters
 */

# SSM Parameters to store zone information
resource "aws_ssm_parameter" "zone_id" {
  name        = "/${var.project_name}/${var.environment}/dns/zone_id"
  description = "Route53 zone ID"
  type        = "String"
  value       = aws_route53_zone.this.zone_id
  tags        = local.tags
}

resource "aws_ssm_parameter" "zone_name" {
  name        = "/${var.project_name}/${var.environment}/dns/zone_name"
  description = "Route53 zone name"
  type        = "String"
  value       = aws_route53_zone.this.name
  tags        = local.tags
}

resource "aws_ssm_parameter" "acm_cert_arn" {
  count       = var.create_acm_certificate ? 1 : 0
  name        = "/${var.project_name}/${var.environment}/dns/acm_cert_arn"
  description = "ACM certificate ARN"
  type        = "String"
  value       = aws_acm_certificate.this[0].arn
  tags        = local.tags
} 
