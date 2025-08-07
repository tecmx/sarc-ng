/**
 * DNS Zone module - Output values
 */

output "zone_id" {
  description = "The ID of the hosted zone"
  value       = aws_route53_zone.this.zone_id
}

output "zone_name" {
  description = "The name of the hosted zone"
  value       = aws_route53_zone.this.name
}

output "name_servers" {
  description = "The name servers for the hosted zone"
  value       = aws_route53_zone.this.name_servers
}

output "acm_certificate_arn" {
  description = "The ARN of the ACM certificate"
  value       = var.create_acm_certificate ? aws_acm_certificate.this[0].arn : null
} 
