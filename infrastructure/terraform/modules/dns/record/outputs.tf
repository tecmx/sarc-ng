/**
 * DNS Record module - Output values
 */

output "fqdn" {
  description = "The FQDN of the record"
  value       = var.alias == null ? aws_route53_record.this[0].fqdn : aws_route53_record.alias[0].fqdn
}

output "name" {
  description = "The name of the record"
  value       = var.alias == null ? aws_route53_record.this[0].name : aws_route53_record.alias[0].name
} 
