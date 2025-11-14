/**
 * Example usage of the DNS zone module
 */

provider "aws" {
  region = "us-east-1"
}

# Create Route53 hosted zone with ACM certificate
module "domain" {
  source = "../.."

  project_name = "sarc-ng"
  environment  = "dev"

  zone_name = "dev.sarc-ng.example.com"
  comment   = "Dev environment hosted zone"

  # Create ACM certificate with SANs
  create_certificate = true
  subject_alternative_names = [
    "*.dev.sarc-ng.example.com",
    "api.dev.sarc-ng.example.com"
  ]

  additional_tags = {
    Example = "complete"
  }
}

