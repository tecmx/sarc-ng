/**
 * DNS Record module - Creates Route53 records for services
 */

# Standard DNS record
resource "aws_route53_record" "this" {
  count = var.alias == null ? 1 : 0

  zone_id = var.zone_id
  name    = local.full_record_name
  type    = var.record_type
  ttl     = var.record_ttl
  records = var.records
}

# Alias record
resource "aws_route53_record" "alias" {
  count = var.alias != null ? 1 : 0

  zone_id = var.zone_id
  name    = local.full_record_name
  type    = var.record_type

  alias {
    name                   = var.alias.name
    zone_id                = var.alias.zone_id
    evaluate_target_health = var.alias.evaluate_target_health
  }
}
