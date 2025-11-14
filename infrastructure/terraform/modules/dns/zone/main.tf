/**
 * DNS Zone module - Creates a Route53 hosted zone and optional ACM certificate
 */

# Route53 Hosted Zone
resource "aws_route53_zone" "this" {
  name = var.zone_name

  dynamic "vpc" {
    for_each = var.create_public_zone ? [] : [1]
    content {
      vpc_id = var.vpc_id
    }
  }

  tags = local.tags
}

# ACM Certificate with DNS validation
resource "aws_acm_certificate" "this" {
  count = var.create_acm_certificate ? 1 : 0

  domain_name               = var.zone_name
  subject_alternative_names = local.cert_domains
  validation_method         = "DNS"

  lifecycle {
    create_before_destroy = true
  }

  tags = local.tags
}

# DNS Validation Records
resource "aws_route53_record" "validation" {
  for_each = var.create_acm_certificate ? {
    for dvo in aws_acm_certificate.this[0].domain_validation_options : dvo.domain_name => {
      name   = dvo.resource_record_name
      record = dvo.resource_record_value
      type   = dvo.resource_record_type
    }
  } : {}

  zone_id = aws_route53_zone.this.zone_id
  name    = each.value.name
  type    = each.value.type
  ttl     = 60
  records = [each.value.record]
}

# Certificate validation
resource "aws_acm_certificate_validation" "this" {
  count = var.create_acm_certificate ? 1 : 0

  certificate_arn         = aws_acm_certificate.this[0].arn
  validation_record_fqdns = [for record in aws_route53_record.validation : record.fqdn]
}
