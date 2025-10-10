# AWS Cognito Authentication - Implementation Checklist

## Quick Summary

This checklist provides a step-by-step guide to implement AWS Cognito authentication in the SARC-NG application. Each section contains specific files to create or modify.

---

## Prerequisites

- [ ] AWS Account with appropriate permissions
- [ ] Terraform installed (>= 1.0)
- [ ] Go 1.24+ installed
- [ ] Access to update environment configurations

---

## Phase 1: Dependencies & Configuration

### 1.1 Update Go Dependencies

```bash
go get github.com/golang-jwt/jwt/v5@latest
go mod tidy
```

### 1.2 Update Configuration Files

**Files to Modify:**

- [ ] `internal/config/config.go` - Add `CognitoConfig` struct
- [ ] `configs/default.yaml` - Add cognito configuration section
- [ ] `configs/development.yaml` - Add cognito environment variables
- [ ] `aws.env` - Add COGNITO_USER_POOL_ID and COGNITO_CLIENT_ID

**Key Configuration:**

```yaml
cognito:
  region: us-east-1
  user_pool_id: ${COGNITO_USER_POOL_ID}
  client_id: ${COGNITO_CLIENT_ID}
  jwks_cache_expiry: 1h
```

---

## Phase 2: Domain Layer Implementation

### 2.1 Create Authentication Domain

**New Files to Create:**

```
internal/domain/auth/
├── entity.go         # User and Claims entities
└── service.go        # TokenValidator interface
```

**Key Components:**

- [ ] `User` entity with ID, email, username, groups
- [ ] `Claims` struct matching Cognito JWT format
- [ ] `TokenValidator` interface
- [ ] Helper methods: `ToUser()`, `HasGroup()`, `HasAnyGroup()`

---

## Phase 3: Service Layer Implementation

### 3.1 JWT Validator Service

**New Files to Create:**

```
internal/service/auth/
├── jwt_validator.go           # Core JWT validation logic
└── jwt_validator_test.go      # Unit tests
```

**Key Features:**

- [ ] JWKS fetching and caching
- [ ] RSA public key parsing
- [ ] Token signature verification
- [ ] Claims validation (issuer, audience, expiration)
- [ ] Automatic JWKS refresh

**Critical Validations:**

```go
- Signature verification using JWKS
- Issuer matches Cognito user pool
- Token not expired
- Audience/client_id matches
- Token_use is "access" or "id"
```

---

## Phase 4: Middleware Implementation

### 4.1 Authentication Middleware

**New Files to Create:**

```
pkg/rest/middleware/
└── auth.go           # Authentication and authorization middleware
```

**Middleware to Implement:**

- [ ] `AuthMiddleware(validator)` - Required authentication
- [ ] `OptionalAuthMiddleware(validator)` - Optional authentication
- [ ] `RequireGroups(groups...)` - Group-based authorization
- [ ] `RequireAllGroups(groups...)` - Require all specified groups

**Context Helpers:**

- [ ] `GetUserFromContext(c)` - Retrieve authenticated user
- [ ] `GetClaimsFromContext(c)` - Retrieve JWT claims
- [ ] `MustGetUser(c)` - Retrieve user or panic

---

## Phase 5: Router Updates

### 5.1 Update Router Structure

**Files to Modify:**

- [ ] `internal/transport/rest/router.go`

**Changes Required:**

```go
// Add to Router struct
type Router struct {
    // ... existing services
    tokenValidator auth.TokenValidator  // NEW
}

// Update NewRouter constructor
func NewRouter(..., tokenValidator auth.TokenValidator) *Router

// Separate routes into public and protected
setupAPIRoutes() {
    publicV1 := router.Group("/api/v1")
    // Public endpoints

    protectedV1 := router.Group("/api/v1")
    protectedV1.Use(middleware.AuthMiddleware(r.tokenValidator))
    // Protected endpoints

    adminV1 := protectedV1.Group("")
    adminV1.Use(middleware.RequireGroups("admin", "manager"))
    // Admin endpoints
}
```

### 5.2 Update Domain Handlers

**Files to Modify (for each domain):**

```
internal/transport/rest/building/
├── handler.go          # Add user context usage
└── routes.go           # Split into public/protected routes

internal/transport/rest/class/
├── handler.go
└── routes.go

internal/transport/rest/lesson/
├── handler.go
└── routes.go

internal/transport/rest/reservation/
├── handler.go
└── routes.go

internal/transport/rest/resource/
├── handler.go
└── routes.go
```

**Pattern for routes.go:**

```go
// Register public routes (read-only)
func RegisterPublicRoutes(router *gin.RouterGroup, service domain.Usecase)

// Register protected routes (authenticated users)
func RegisterProtectedRoutes(router *gin.RouterGroup, service domain.Usecase)

// Register admin routes (admin/manager only)
func RegisterAdminRoutes(router *gin.RouterGroup, service domain.Usecase)
```

---

## Phase 6: Dependency Injection

### 6.1 Update Wire Configuration

**Files to Modify:**

- [ ] `cmd/server/wire.go`
- [ ] `cmd/lambda/wire.go`

**Changes:**

```go
//go:build wireinject

// Add token validator provider
func provideTokenValidator(cfg *config.Config) auth.TokenValidator {
    return authService.NewJWTValidator(
        cfg.Cognito.Region,
        cfg.Cognito.UserPoolID,
        cfg.Cognito.ClientID,
        cfg.Cognito.JWKSCacheExp,
    )
}

// Update initializeRouter
func initializeRouter(cfg *config.Config) (*rest.Router, error) {
    wire.Build(
        // ... existing providers
        provideTokenValidator,  // NEW
        rest.NewRouter,
    )
    return nil, nil
}
```

**Then regenerate wire code:**

```bash
cd cmd/server && wire
cd cmd/lambda && wire
```

---

## Phase 7: Infrastructure as Code

### 7.1 Create Cognito Terraform Module

**New Files to Create:**

```
infrastructure/terraform/modules/idp/cognito/
├── main.tf           # User pool and client configuration
├── variables.tf      # Input variables
├── outputs.tf        # Output values
├── versions.tf       # Provider versions
└── README.md         # Module documentation
```

**Key Resources:**

- [ ] `aws_cognito_user_pool` - User pool with security policies
- [ ] `aws_cognito_user_pool_client` - Application client
- [ ] `aws_cognito_user_pool_domain` - Hosted UI domain (optional)
- [ ] `aws_cognito_user_group` - User groups (admin, manager, user)
- [ ] `aws_ssm_parameter` - Store config in Parameter Store

### 7.2 Create Live Environment Configuration

**New Files to Create:**

```
infrastructure/terraform/live/accounts/dev/cognito/
└── terragrunt.hcl    # Dev environment Cognito config
```

**Example terragrunt.hcl:**

```hcl
terraform {
  source = "../../../../../modules/idp/cognito"
}

include "root" {
  path = find_in_parent_folders()
}

inputs = {
  user_pool_name              = "sarc-ng-dev"
  application_name            = "sarc-ng"
  environment                 = "dev"
  domain_name                 = "sarc-ng-dev"
  mfa_configuration           = "OPTIONAL"
  advanced_security_mode      = "AUDIT"
  allow_admin_create_user_only = false

  callback_urls = [
    "http://localhost:3000/callback",
    "https://dev.sarc-ng.example.com/callback"
  ]

  logout_urls = [
    "http://localhost:3000",
    "https://dev.sarc-ng.example.com"
  ]

  user_groups = {
    admin = {
      description = "Administrators with full access"
      precedence  = 1
      role_arn    = ""
    }
    manager = {
      description = "Managers with resource management access"
      precedence  = 2
      role_arn    = ""
    }
    teacher = {
      description = "Teachers who can create classes"
      precedence  = 3
      role_arn    = ""
    }
    student = {
      description = "Students who can make reservations"
      precedence  = 4
      role_arn    = ""
    }
  }

  tags = {
    Environment = "dev"
    Project     = "sarc-ng"
    ManagedBy   = "Terraform"
  }
}
```

### 7.3 Deploy Infrastructure

```bash
cd infrastructure/terraform/live/accounts/dev/cognito
terragrunt init
terragrunt plan
terragrunt apply
```

### 7.4 Retrieve Configuration Values

```bash
# Get User Pool ID
aws ssm get-parameter --name "/dev/sarc-ng/cognito/user-pool-id" --query "Parameter.Value" --output text

# Get Client ID
aws ssm get-parameter --name "/dev/sarc-ng/cognito/client-id" --query "Parameter.Value" --output text

# Update aws.env file
echo "COGNITO_USER_POOL_ID=<user-pool-id>" >> aws.env
echo "COGNITO_CLIENT_ID=<client-id>" >> aws.env
```

---

## Phase 8: Testing

### 8.1 Unit Tests

**New Test Files to Create:**

```
internal/service/auth/
└── jwt_validator_test.go

pkg/rest/middleware/
└── auth_test.go
```

**Test Coverage:**

- [ ] Valid token validation
- [ ] Expired token rejection
- [ ] Invalid signature rejection
- [ ] Invalid issuer rejection
- [ ] JWKS caching
- [ ] Group-based authorization
- [ ] Context user extraction

### 8.2 Integration Tests

**New Test Files to Create:**

```
test/integration/
├── auth_test.go              # Authentication flow tests
└── authorization_test.go     # Authorization tests
```

**Test Scenarios:**

- [ ] Public endpoint without token (200 OK)
- [ ] Protected endpoint without token (401 Unauthorized)
- [ ] Protected endpoint with valid token (200 OK)
- [ ] Protected endpoint with expired token (401 Unauthorized)
- [ ] Admin endpoint with user token (403 Forbidden)
- [ ] Admin endpoint with admin token (200 OK)

### 8.3 Run Tests

```bash
# Unit tests
go test ./internal/service/auth/... -v
go test ./pkg/rest/middleware/... -v

# Integration tests
go test ./test/integration/... -v

# Coverage
go test -cover ./...
```

---

## Phase 9: Documentation

### 9.1 Update API Documentation

**Files to Update:**

- [ ] `api/openapi.yaml` - Add security schemes and requirements
- [ ] `docs/content/api-reference/` - Document authentication
- [ ] `README.md` - Add authentication setup instructions

**OpenAPI Security Scheme:**

```yaml
components:
  securitySchemes:
    BearerAuth:
      type: http
      scheme: bearer
      bearerFormat: JWT
      description: JWT token from AWS Cognito

security:
  - BearerAuth: []
```

### 9.2 Create User Guides

**New Documentation Files:**

- [ ] `docs/content/authentication-setup.md` - Setup guide
- [ ] `docs/content/user-management.md` - User management procedures
- [ ] `docs/content/troubleshooting-auth.md` - Common issues

---

## Phase 10: Deployment

### 10.1 Development Environment

```bash
# 1. Deploy Cognito infrastructure
cd infrastructure/terraform/live/accounts/dev/cognito
terragrunt apply

# 2. Update environment variables
export COGNITO_USER_POOL_ID=<from-ssm>
export COGNITO_CLIENT_ID=<from-ssm>

# 3. Rebuild application
make build

# 4. Run application
make run

# 5. Create test users
aws cognito-idp admin-create-user \
  --user-pool-id $COGNITO_USER_POOL_ID \
  --username testuser \
  --user-attributes Name=email,Value=test@example.com \
  --temporary-password "TempPass123!"

# 6. Add user to group
aws cognito-idp admin-add-user-to-group \
  --user-pool-id $COGNITO_USER_POOL_ID \
  --username testuser \
  --group-name admin
```

### 10.2 Staging/Production

Follow the same process but:

- [ ] Use appropriate environment configs
- [ ] Enable MFA (ON instead of OPTIONAL)
- [ ] Set advanced_security_mode to ENFORCED
- [ ] Use proper callback URLs
- [ ] Review and adjust token validity periods
- [ ] Enable CloudWatch logging
- [ ] Setup alerts for auth failures

---

## Phase 11: Monitoring & Observability

### 11.1 Metrics to Track

Add to Prometheus/CloudWatch:

```go
- cognito_auth_attempts_total
- cognito_auth_failures_total
- cognito_token_validation_duration_seconds
- cognito_jwks_cache_hits_total
- cognito_jwks_cache_misses_total
- cognito_authorization_failures_total
```

### 11.2 CloudWatch Integration

Enable Cognito logging:

```bash
aws cognito-idp set-user-pool-mfa-config \
  --user-pool-id $COGNITO_USER_POOL_ID \
  --mfa-configuration ON
```

### 11.3 Alerts

Configure alerts for:

- [ ] High authentication failure rate (> 10% for 5 minutes)
- [ ] Token validation latency (p99 > 200ms)
- [ ] JWKS fetch failures
- [ ] Unusual authorization patterns

---

## Phase 12: Security Audit

### 12.1 Security Checklist

- [ ] HTTPS enforced in production
- [ ] Token expiration validated
- [ ] Signature verification working
- [ ] JWKS properly cached
- [ ] No tokens logged
- [ ] Rate limiting configured
- [ ] CORS properly configured
- [ ] Security headers set
- [ ] No sensitive data in URLs
- [ ] Proper error messages (no information leakage)

### 12.2 Penetration Testing

Test scenarios:

- [ ] Token tampering
- [ ] Expired token usage
- [ ] Invalid signature
- [ ] Missing claims
- [ ] Group escalation attempts
- [ ] Replay attacks
- [ ] CORS bypass attempts

### 12.3 Code Review

Review checklist:

- [ ] No hardcoded credentials
- [ ] Proper error handling
- [ ] Context propagation
- [ ] Timeout configurations
- [ ] Resource cleanup
- [ ] Thread safety (JWKS cache)

---

## Rollback Procedures

### If Authentication Issues Occur

1. **Quick Disable (Emergency)**

```go
// In router.go, comment out auth middleware temporarily
func (r *Router) setupAPIRoutes(router *gin.Engine) {
    protectedV1 := router.Group("/api/v1")
    // protectedV1.Use(middleware.AuthMiddleware(r.tokenValidator))  // DISABLED
    {
        // All routes work without auth
    }
}
```

2. **Gradual Rollback**

```bash
# Revert code changes
git revert <commit-hash>

# Rebuild and redeploy
make build && make deploy

# Infrastructure remains (no data loss)
# Users in Cognito are preserved
```

3. **Complete Rollback**

```bash
# Destroy Cognito resources (CAUTION: User data lost)
cd infrastructure/terraform/live/accounts/dev/cognito
terragrunt destroy
```

---

## Success Metrics

### Authentication

- [ ] Token validation latency < 100ms (p99)
- [ ] Authentication success rate > 95%
- [ ] Zero token bypass incidents
- [ ] JWKS cache hit rate > 90%

### Authorization

- [ ] Zero unauthorized access incidents
- [ ] Authorization check latency < 10ms
- [ ] Proper group enforcement (100%)

### User Experience

- [ ] Clear error messages
- [ ] Token refresh working smoothly
- [ ] No false authentication failures
- [ ] Smooth login/logout flow

### Operations

- [ ] Infrastructure deployed via Terraform
- [ ] Configuration in SSM Parameter Store
- [ ] Monitoring and alerting configured
- [ ] Documentation complete and accurate

---

## Support and Troubleshooting

### Common Issues

**Issue: "Invalid token" error**

```bash
# Check token expiration
# Verify COGNITO_USER_POOL_ID matches issuer in token
# Ensure COGNITO_CLIENT_ID matches audience/client_id in token
# Check JWKS is accessible
```

**Issue: "Insufficient permissions" error**

```bash
# Verify user is in required group
aws cognito-idp admin-list-groups-for-user \
  --user-pool-id $COGNITO_USER_POOL_ID \
  --username <username>

# Add user to group if needed
aws cognito-idp admin-add-user-to-group \
  --user-pool-id $COGNITO_USER_POOL_ID \
  --username <username> \
  --group-name <group>
```

**Issue: JWKS fetch failure**

```bash
# Verify network connectivity to Cognito
curl https://cognito-idp.us-east-1.amazonaws.com/<POOL_ID>/.well-known/jwks.json

# Check IAM permissions (if VPC endpoint used)
# Verify security group rules
```

### Getting Help

- Review authentication proposal: `docs/content/authentication-proposal.md`
- Check API documentation: `api/openapi.yaml`
- Review AWS Cognito docs: <https://docs.aws.amazon.com/cognito/>
- Check application logs for detailed error messages

---

## Completion Checklist

### Development Complete When

- [x] All code files created and implemented
- [x] All tests passing with >80% coverage
- [x] Documentation complete
- [x] Infrastructure deployed to dev
- [x] Manual testing completed
- [x] Security review passed

### Ready for Production When

- [ ] Staging deployment successful
- [ ] Performance testing passed
- [ ] Security audit completed
- [ ] Monitoring configured
- [ ] Alerts configured
- [ ] Runbook created
- [ ] Team trained on authentication
- [ ] Rollback plan tested

---

## Timeline Estimate

- **Phase 1-2 (Config & Domain)**: 1-2 days
- **Phase 3-4 (Service & Middleware)**: 2-3 days
- **Phase 5-6 (Router & DI)**: 2 days
- **Phase 7 (Infrastructure)**: 2-3 days
- **Phase 8 (Testing)**: 2-3 days
- **Phase 9-10 (Docs & Deploy)**: 2 days
- **Phase 11-12 (Monitoring & Security)**: 2 days

**Total Estimated Time: 2-3 weeks**

---

## Next Steps

1. Review the authentication proposal document
2. Get stakeholder approval
3. Set up development branch
4. Begin Phase 1 implementation
5. Deploy to dev environment after Phase 10
6. Conduct security audit before production
7. Deploy to production with gradual rollout

---

## Additional Resources

- [AWS Cognito Best Practices](https://docs.aws.amazon.com/cognito/latest/developerguide/security-best-practices.html)
- [JWT Best Practices](https://tools.ietf.org/html/rfc8725)
- [OWASP Authentication Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Authentication_Cheat_Sheet.html)
- [OAuth 2.0 Security Best Practices](https://tools.ietf.org/html/draft-ietf-oauth-security-topics)
