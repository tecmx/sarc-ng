/**
 * Example usage of the DNS record module
 */

provider "aws" {
  region = "us-east-1"
}

# Assume hosted zone exists
data "aws_route53_zone" "main" {
  name = "dev.sarc-ng.example.com"
}

# Create various DNS records
module "api_record" {
  source = "../.."

  project_name = "sarc-ng"
  environment  = "dev"

  zone_id = data.aws_route53_zone.main.zone_id
  name    = "api.dev.sarc-ng.example.com"
  type    = "A"

  # ALB alias record
  alias = {
    name                   = "my-alb-123456.us-east-1.elb.amazonaws.com"
    zone_id                = "Z35SXDOTRQ7X7K" # ALB zone ID
    evaluate_target_health = true
  }

  additional_tags = {
    Example = "complete"
    Type    = "api"
  }
}

# CNAME record example
module "www_record" {
  source = "../.."

  project_name = "sarc-ng"
  environment  = "dev"

  zone_id = data.aws_route53_zone.main.zone_id
  name    = "www.dev.sarc-ng.example.com"
  type    = "CNAME"
  ttl     = 300
  records = ["dev.sarc-ng.example.com"]

  additional_tags = {
    Example = "complete"
    Type    = "www"
  }
}

