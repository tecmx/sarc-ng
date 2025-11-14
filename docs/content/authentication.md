# AWS Cognito Authentication - Documentation Index

## 📚 Documentation Overview

This is the complete documentation suite for implementing AWS Cognito authentication in the SARC-NG application. The documentation is organized into three complementary documents, each serving a specific purpose.

---

## 📖 Document Guide

### 1. 📄 [Authentication Proposal](./authentication-proposal.md)

**Purpose:** Comprehensive technical specification
**Audience:** Architects, senior developers, technical reviewers
**Reading Time:** 45 minutes

**What's Inside:**

- Current state analysis
- Detailed architecture design
- Complete code specifications
- Security best practices
- Authorization model
- Migration strategy
- Monitoring and compliance
- Cost analysis

**Use This For:**

- Technical review and approval
- Architecture decisions
- Security audit preparation
- Understanding design rationale
- Reference during implementation
- Business case and cost analysis
- Executive overview

---

### 2. ✅ [Authentication Implementation Checklist](./authentication-implementation-checklist.md)

**Purpose:** Step-by-step implementation guide
**Audience:** Developers actively implementing the solution
**Reading Time:** 30 minutes (reference throughout implementation)

**What's Inside:**

- Phase-by-phase implementation steps
- Complete file listing (new and modified)
- Code snippets for each phase
- Testing procedures
- Deployment instructions
- Troubleshooting guide
- Success criteria

**Use This For:**

- Daily implementation tasks
- Tracking progress
- Code review checklist
- Deployment procedures
- Troubleshooting issues

---

### 3. 🏗️ [Authentication Architecture](./authentication-architecture.md)

**Purpose:** Visual architecture and examples
**Audience:** Developers, DevOps, technical support
**Reading Time:** 20 minutes

**What's Inside:**

- System architecture diagrams
- Authentication flow diagrams
- Token structure examples
- Authorization model visualizations
- Code structure overview
- Request/response examples
- Performance characteristics
- Monitoring metrics

**Use This For:**

- Understanding system flow
- Debugging authentication issues
- Performance optimization
- Setting up monitoring
- Training new team members

---

## 🚀 Quick Start Guide

### For Stakeholders/Management

1. Read: [Authentication Proposal](./authentication-proposal.md) - Cost analysis and timeline sections (15 min)
2. Review: Key benefits and decision rationale
3. Decision: Approve or request clarifications

### For Technical Reviewers

1. Read: [Authentication Proposal](./authentication-proposal.md) (45 min)
2. Review: [Authentication Architecture](./authentication-architecture.md) (20 min)
3. Provide: Technical feedback and approval

### For Developers

1. Skim: [Authentication Proposal](./authentication-proposal.md) (15 min)
2. Study: [Authentication Architecture](./authentication-architecture.md) (20 min)
3. Follow: [Authentication Implementation Checklist](./authentication-implementation-checklist.md) (ongoing)

### For DevOps/Infrastructure

1. Review: Infrastructure sections in [Proposal](./authentication-proposal.md)
2. Follow: Infrastructure deployment in [Checklist](./authentication-implementation-checklist.md)
3. Reference: [Architecture](./authentication-architecture.md) for monitoring

---

## 📋 Implementation Phases Summary

### Phase 1-2: Foundation & Core (Week 1)

**Documents:** Checklist (Phases 1-3), Proposal (Sections 1-3), Architecture (Code Structure)

**Tasks:**

- Update configuration
- Create domain entities
- Implement JWT validator
- Create middleware
- Write unit tests

**Deliverable:** Authentication logic complete and tested

### Phase 3-4: Integration & Infrastructure (Week 2)

**Documents:** Checklist (Phases 4-5), Proposal (Sections 5-6), Architecture (Deployment)

**Tasks:**

- Update router and handlers
- Deploy Cognito infrastructure
- Configure user groups
- Integration testing
- Deploy to dev environment

**Deliverable:** Fully integrated and deployed to dev

### Phase 5-6: Testing & Production (Week 3)

**Documents:** Checklist (Phases 6-12), Proposal (Sections 8-12)

**Tasks:**

- Security testing
- Documentation updates
- Monitoring setup
- Production deployment
- Post-deployment validation

**Deliverable:** Production-ready authentication system

---

## 🎯 Key Decision Points

### Before Starting Implementation

| Question | Answer | Document Reference |
|----------|--------|-------------------|
| Why AWS Cognito? | Managed, scalable, secure, cost-effective | Proposal → Decision Points |
| What's the cost? | Free up to 50K users | Proposal → Cost Analysis |
| How long to implement? | 2-3 weeks | Proposal → Timeline |
| What are the risks? | Low - well-established technology | Proposal → Risk Assessment |
| Is it reversible? | Yes - stateless architecture | Proposal → Rollback Plan |

### During Implementation

| Question | Answer | Document Reference |
|----------|--------|-------------------|
| Which files to create? | ~15 new files | Checklist → Files to Create |
| Which files to modify? | ~10 existing files | Checklist → Files to Modify |
| How to structure code? | See directory layout | Architecture → Code Structure |
| How to test? | Unit + integration tests | Checklist → Phase 8 |
| How to deploy? | Terraform + manual steps | Checklist → Phase 10 |

### After Deployment

| Question | Answer | Document Reference |
|----------|--------|-------------------|
| How to monitor? | Prometheus metrics | Architecture → Monitoring |
| What to alert on? | Auth failures, latency | Architecture → Alert Conditions |
| How to troubleshoot? | Check logs, verify config | Checklist → Troubleshooting |
| How to rollback? | Remove middleware | Checklist → Rollback |
| How to add users? | AWS CLI or Console | Checklist → User Management |

---

## 🔧 Common Tasks Quick Reference

### Create a New User

```bash
aws cognito-idp admin-create-user \
  --user-pool-id $COGNITO_USER_POOL_ID \
  --username newuser \
  --user-attributes Name=email,Value=user@example.com \
  --temporary-password "TempPass123!"
```

**Reference:** Checklist → Phase 10.1

### Add User to Group

```bash
aws cognito-idp admin-add-user-to-group \
  --user-pool-id $COGNITO_USER_POOL_ID \
  --username newuser \
  --group-name admin
```

**Reference:** Checklist → Phase 10.1

### Test Authentication

```bash
# Public endpoint (should work)
curl http://localhost:8080/api/v1/buildings

# Protected endpoint (should fail without token)
curl http://localhost:8080/api/v1/reservations

# Protected endpoint (with token)
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/api/v1/reservations
```

**Reference:** Architecture → Request/Response Examples

### Check Metrics

```bash
curl http://localhost:8080/metrics | grep cognito
```

**Reference:** Architecture → Monitoring Dashboard

---

## 🛠️ Development Workflow

### Daily Development Cycle

```
1. Check Checklist → Identify current phase
2. Review Proposal → Understand requirements
3. Reference Architecture → See examples
4. Implement → Write code
5. Test → Run unit/integration tests
6. Review → Code review with team
7. Update Checklist → Mark tasks complete
```

### Code Review Checklist

- [ ] Follows clean architecture principles
- [ ] Single responsibility per function
- [ ] Proper error handling
- [ ] No hardcoded credentials
- [ ] Unit tests written
- [ ] Documentation updated
- [ ] Security best practices followed

**Reference:** Proposal → Code Quality Standards

---

## 📊 Progress Tracking

### Implementation Progress

Track your progress through the phases:

```
☐ Phase 1: Configuration & Dependencies
☐ Phase 2: Domain Layer
☐ Phase 3: Service Layer
☐ Phase 4: Middleware Layer
☐ Phase 5: Router Updates
☐ Phase 6: Dependency Injection
☐ Phase 7: Infrastructure
☐ Phase 8: Testing
☐ Phase 9: Documentation
☐ Phase 10: Deployment
☐ Phase 11: Monitoring
☐ Phase 12: Security Audit
```

### Completion Criteria

```
Development Complete:
☐ All code implemented
☐ All tests passing (>80% coverage)
☐ Documentation complete
☐ Dev deployment successful
☐ Security review passed

Production Ready:
☐ Staging deployment successful
☐ Performance testing passed
☐ Security audit completed
☐ Monitoring configured
☐ Team trained
☐ Rollback plan tested
```

---

## 🆘 Need Help?

### Finding Information

| If You Need... | Look In... |
|----------------|------------|
| Business justification | Proposal → Key Benefits |
| Technical details | Proposal → Complete specification |
| Step-by-step guide | Checklist → Phase-by-phase |
| Code examples | Architecture → Request/Response |
| Security info | Proposal → Security Best Practices |
| Cost information | Proposal → Cost Analysis |
| Troubleshooting | Checklist → Troubleshooting |
| Monitoring setup | Architecture → Monitoring |

### Common Questions

**Q: Where do I start?**
A: Read the Proposal for overview, then follow the Checklist phase by phase.

**Q: How do I understand the architecture?**
A: Study the Architecture document's diagrams and flows.

**Q: What if I get stuck?**
A: Check the Checklist troubleshooting section, review the Proposal for details.

**Q: How do I know it's working?**
A: Follow the testing procedures in the Checklist, check metrics in Architecture.

**Q: Is it secure?**
A: Review the Security section in the Proposal, complete the security audit checklist.

---

## 📦 Deliverables

### Documentation

- [x] Authentication Proposal
- [x] Authentication Implementation Checklist
- [x] Authentication Architecture
- [x] Authentication Index (this document)

### Code (To Be Implemented)

- [ ] Domain entities and interfaces
- [ ] JWT validator service
- [ ] Authentication middleware
- [ ] Router updates
- [ ] Handler modifications
- [ ] Dependency injection setup
- [ ] Unit tests
- [ ] Integration tests

### Infrastructure (To Be Deployed)

- [ ] Cognito User Pool
- [ ] User Groups
- [ ] Terraform modules
- [ ] SSM parameters
- [ ] Monitoring configuration
- [ ] Alert rules

---

## 🎓 Learning Path

### For New Team Members

1. **Day 1: Understand the System**
   - Read: Proposal (business context and technical overview)
   - Read: Architecture (technical overview)
   - Goal: Understand what we're building and why

2. **Day 2-3: Deep Dive**
   - Read: Proposal (detailed design)
   - Study: Code examples in Architecture
   - Goal: Understand how it works

3. **Day 4-5: Hands-On**
   - Follow: Checklist (implementation steps)
   - Practice: Setting up dev environment
   - Goal: Ready to contribute

### For Experienced Developers

1. **Quick Start (1 hour)**
   - Skim: Proposal → Architecture → Checklist
   - Jump: Directly to implementation phase needed
   - Reference: Proposal for specific details

---

## 📅 Document Maintenance

### When to Update

- **Proposal**: When architecture decisions or business requirements change
- **Checklist**: When implementation steps change
- **Architecture**: When system design changes
- **Index**: When any document is added/modified

### Version History

| Version | Date | Changes | Author |
|---------|------|---------|--------|
| 1.0 | 2024-10-10 | Initial documentation suite | AI Assistant |

---

## 🔗 Related Resources

### Internal

- `README.md` - Project overview
- `CONTRIBUTING.md` - Development guidelines
- `api/openapi.yaml` - API specification
- `configs/` - Configuration files

### External

- [AWS Cognito Documentation](https://docs.aws.amazon.com/cognito/)
- [JWT Introduction](https://jwt.io/introduction)
- [OWASP Authentication Guide](https://cheatsheetseries.owasp.org/cheatsheets/Authentication_Cheat_Sheet.html)
- [OAuth 2.0 Specification](https://oauth.net/2/)

---

## ✅ Review Checklist

Before proceeding with implementation, ensure:

- [ ] All stakeholders have reviewed the Proposal
- [ ] Technical team has reviewed the Proposal
- [ ] Security team has reviewed security sections
- [ ] DevOps team has reviewed infrastructure sections
- [ ] Timeline and budget approved
- [ ] Success criteria agreed upon
- [ ] Rollback plan understood
- [ ] Team has access to all documentation

---

**Ready to Begin?** Start with the [Authentication Proposal](./authentication-proposal.md) document!

**Need Implementation Details?** Jump to the [Authentication Implementation Checklist](./authentication-implementation-checklist.md)!

**Want to Understand the Architecture?** Review the [Authentication Architecture](./authentication-architecture.md) document!

**Looking for Complete Specifications?** Read the [Authentication Proposal](./authentication-proposal.md)!

---

*This documentation suite provides everything you need to successfully implement AWS Cognito authentication in SARC-NG. Good luck with your implementation!*


---

# AWS Cognito Authentication Implementation Proposal

## Table of Contents

- [Executive Summary](#executive-summary)
- [Current State Analysis](#current-state-analysis)
- [Proposed Solution Architecture](#proposed-solution-architecture)
  - [Configuration Layer](#1-configuration-layer)
  - [Domain Layer](#2-domain-layer)
  - [Service Layer](#3-service-layer)
  - [Middleware Layer](#4-middleware-layer)
  - [Transport Layer Updates](#5-transport-layer-updates)
  - [Infrastructure as Code](#6-infrastructure-as-code)
  - [Dependency Injection with Wire](#7-dependency-injection-with-wire)
  - [Updated Dependencies](#8-updated-dependencies)
  - [Testing Strategy](#9-testing-strategy)
  - [Implementation Phases](#10-implementation-phases)
  - [Security Best Practices](#11-security-best-practices)
  - [Authorization Model](#12-authorization-model)
  - [Environment Variables](#13-environment-variables)
  - [API Request Examples](#14-api-request-examples)
  - [Migration Strategy](#15-migration-strategy)
  - [Monitoring and Logging](#16-monitoring-and-logging)
  - [Cost Estimation](#17-cost-estimation)
  - [Compliance and Standards](#18-compliance-and-standards)
  - [Rollback Plan](#19-rollback-plan)
  - [Success Criteria](#20-success-criteria)
- [Conclusion](#conclusion)

---

## Executive Summary

This document proposes implementing AWS Cognito as the Identity Provider (IDP) for the SARC-NG application following industry best practices, security standards, and clean architecture principles.

## Current State Analysis

### Existing Infrastructure

- Clean architecture with domain-driven design (DDD)
- REST API using Gin framework with middleware support
- Configuration management using Viper
- AWS SDK v2 dependencies already present
- Infrastructure as Code with Terraform modules
- Middleware stack: CORS, Logger, Recovery, Metrics

### Missing Components

- JWT token validation
- Authentication middleware
- User context management
- Cognito configuration
- Authorization/RBAC capabilities
- Protected vs public route differentiation

---

## Proposed Solution Architecture

### 1. Configuration Layer

#### Update Configuration Structure

**File: `internal/config/config.go`**

```go
// CognitoConfig holds AWS Cognito-related configuration
type CognitoConfig struct {
    Region       string        `mapstructure:"region"`
    UserPoolID   string        `mapstructure:"user_pool_id"`
    ClientID     string        `mapstructure:"client_id"`
    JWKSCacheExp time.Duration `mapstructure:"jwks_cache_expiry"`
}

// Add to Config struct
type Config struct {
    Server   ServerConfig   `mapstructure:"server"`
    Database DatabaseConfig `mapstructure:"database"`
    Redis    RedisConfig    `mapstructure:"redis"`
    Cognito  CognitoConfig  `mapstructure:"cognito"`  // NEW
    Logging  LoggingConfig  `mapstructure:"logging"`
    API      APIConfig      `mapstructure:"api"`
}
```

#### Configuration Files

**File: `configs/default.yaml`**

```yaml
cognito:
  region: us-east-1
  user_pool_id: ""  # To be provided via environment
  client_id: ""     # To be provided via environment
  jwks_cache_expiry: 1h
```

**File: `configs/development.yaml`**

```yaml
cognito:
  region: us-east-1
  user_pool_id: ${COGNITO_USER_POOL_ID}
  client_id: ${COGNITO_CLIENT_ID}
  jwks_cache_expiry: 1h
```

---

### 2. Domain Layer

#### Authentication Domain Entities

**File: `internal/domain/auth/entity.go`**

```go
package auth

import "time"

// User represents an authenticated user from Cognito
type User struct {
    ID         string
    Email      string
    Username   string
    Groups     []string
    Attributes map[string]string
    AuthTime   time.Time
}

// Claims represents JWT token claims from Cognito
type Claims struct {
    Sub            string   `json:"sub"`
    Email          string   `json:"email"`
    EmailVerified  bool     `json:"email_verified"`
    Username       string   `json:"cognito:username"`
    Groups         []string `json:"cognito:groups"`
    TokenUse       string   `json:"token_use"`
    Scope          string   `json:"scope"`
    AuthTime       int64    `json:"auth_time"`
    IssuedAt       int64    `json:"iat"`
    ExpirationTime int64    `json:"exp"`
    Issuer         string   `json:"iss"`
    ClientID       string   `json:"client_id"`
    Audience       string   `json:"aud"`
}

// ToUser converts claims to User entity
func (c *Claims) ToUser() *User {
    return &User{
        ID:         c.Sub,
        Email:      c.Email,
        Username:   c.Username,
        Groups:     c.Groups,
        AuthTime:   time.Unix(c.AuthTime, 0),
        Attributes: make(map[string]string),
    }
}

// HasGroup checks if user belongs to a specific group
func (u *User) HasGroup(group string) bool {
    for _, g := range u.Groups {
        if g == group {
            return true
        }
    }
    return false
}

// HasAnyGroup checks if user belongs to any of the specified groups
func (u *User) HasAnyGroup(groups []string) bool {
    for _, group := range groups {
        if u.HasGroup(group) {
            return true
        }
    }
    return false
}
```

#### Authentication Service Interface

**File: `internal/domain/auth/service.go`**

```go
package auth

import "context"

// TokenValidator defines the interface for JWT token validation
type TokenValidator interface {
    // ValidateToken validates a JWT token and returns claims
    ValidateToken(ctx context.Context, token string) (*Claims, error)

    // RefreshJWKS forces a refresh of the JWKS cache
    RefreshJWKS(ctx context.Context) error
}
```

---

### 3. Service Layer

#### JWT Validator Implementation

**File: `internal/service/auth/jwt_validator.go`**

```go
package auth

import (
    "context"
    "crypto/rsa"
    "encoding/base64"
    "encoding/json"
    "errors"
    "fmt"
    "math/big"
    "net/http"
    "strings"
    "sync"
    "time"

    "github.com/golang-jwt/jwt/v5"
    "sarc-ng/internal/domain/auth"
)

var (
    ErrInvalidToken      = errors.New("invalid token")
    ErrExpiredToken      = errors.New("token has expired")
    ErrInvalidSignature  = errors.New("invalid token signature")
    ErrInvalidIssuer     = errors.New("invalid token issuer")
    ErrInvalidAudience   = errors.New("invalid token audience")
    ErrInvalidTokenUse   = errors.New("invalid token use")
)

// JWKS represents the JSON Web Key Set structure
type JWKS struct {
    Keys []JWK `json:"keys"`
}

// JWK represents a JSON Web Key
type JWK struct {
    Kid string `json:"kid"`
    Kty string `json:"kty"`
    Alg string `json:"alg"`
    Use string `json:"use"`
    N   string `json:"n"`
    E   string `json:"e"`
}

// JWTValidator validates JWT tokens from AWS Cognito
type JWTValidator struct {
    region          string
    userPoolID      string
    clientID        string
    jwksURL         string
    issuer          string
    jwksCache       map[string]*rsa.PublicKey
    cacheMutex      sync.RWMutex
    cacheExpiry     time.Duration
    lastCacheUpdate time.Time
    httpClient      *http.Client
}

// NewJWTValidator creates a new JWT validator for Cognito
func NewJWTValidator(region, userPoolID, clientID string, cacheExpiry time.Duration) *JWTValidator {
    jwksURL := fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/%s/.well-known/jwks.json", region, userPoolID)
    issuer := fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/%s", region, userPoolID)

    return &JWTValidator{
        region:      region,
        userPoolID:  userPoolID,
        clientID:    clientID,
        jwksURL:     jwksURL,
        issuer:      issuer,
        jwksCache:   make(map[string]*rsa.PublicKey),
        cacheExpiry: cacheExpiry,
        httpClient: &http.Client{
            Timeout: 10 * time.Second,
        },
    }
}

// ValidateToken validates a JWT token and returns claims
func (v *JWTValidator) ValidateToken(ctx context.Context, tokenString string) (*auth.Claims, error) {
    // Parse token without verification first to get the kid
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        // Verify signing algorithm
        if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }

        // Get key ID from token header
        kid, ok := token.Header["kid"].(string)
        if !ok {
            return nil, errors.New("missing kid in token header")
        }

        // Get public key for this kid
        publicKey, err := v.getPublicKey(ctx, kid)
        if err != nil {
            return nil, err
        }

        return publicKey, nil
    })

    if err != nil {
        if errors.Is(err, jwt.ErrTokenExpired) {
            return nil, ErrExpiredToken
        }
        return nil, fmt.Errorf("%w: %v", ErrInvalidToken, err)
    }

    if !token.Valid {
        return nil, ErrInvalidToken
    }

    // Extract claims
    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok {
        return nil, ErrInvalidToken
    }

    // Validate issuer
    iss, ok := claims["iss"].(string)
    if !ok || iss != v.issuer {
        return nil, ErrInvalidIssuer
    }

    // Validate token_use
    tokenUse, ok := claims["token_use"].(string)
    if !ok || (tokenUse != "access" && tokenUse != "id") {
        return nil, ErrInvalidTokenUse
    }

    // Validate client_id (for access tokens) or aud (for id tokens)
    if tokenUse == "access" {
        clientID, ok := claims["client_id"].(string)
        if !ok || clientID != v.clientID {
            return nil, ErrInvalidAudience
        }
    } else {
        aud, ok := claims["aud"].(string)
        if !ok || aud != v.clientID {
            return nil, ErrInvalidAudience
        }
    }

    // Convert to our Claims structure
    authClaims := &auth.Claims{}
    if err := v.mapClaimsToDomain(claims, authClaims); err != nil {
        return nil, err
    }

    return authClaims, nil
}

// RefreshJWKS forces a refresh of the JWKS cache
func (v *JWTValidator) RefreshJWKS(ctx context.Context) error {
    return v.fetchJWKS(ctx)
}

// getPublicKey retrieves the public key for a given kid
func (v *JWTValidator) getPublicKey(ctx context.Context, kid string) (*rsa.PublicKey, error) {
    // Check cache first
    v.cacheMutex.RLock()
    key, exists := v.jwksCache[kid]
    cacheExpired := time.Since(v.lastCacheUpdate) > v.cacheExpiry
    v.cacheMutex.RUnlock()

    if exists && !cacheExpired {
        return key, nil
    }

    // Fetch JWKS
    if err := v.fetchJWKS(ctx); err != nil {
        return nil, err
    }

    // Try again from cache
    v.cacheMutex.RLock()
    key, exists = v.jwksCache[kid]
    v.cacheMutex.RUnlock()

    if !exists {
        return nil, fmt.Errorf("key with kid %s not found in JWKS", kid)
    }

    return key, nil
}

// fetchJWKS fetches and caches the JWKS from Cognito
func (v *JWTValidator) fetchJWKS(ctx context.Context) error {
    req, err := http.NewRequestWithContext(ctx, http.MethodGet, v.jwksURL, nil)
    if err != nil {
        return fmt.Errorf("failed to create request: %w", err)
    }

    resp, err := v.httpClient.Do(req)
    if err != nil {
        return fmt.Errorf("failed to fetch JWKS: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("failed to fetch JWKS: status code %d", resp.StatusCode)
    }

    var jwks JWKS
    if err := json.NewDecoder(resp.Body).Decode(&jwks); err != nil {
        return fmt.Errorf("failed to decode JWKS: %w", err)
    }

    // Update cache
    v.cacheMutex.Lock()
    defer v.cacheMutex.Unlock()

    v.jwksCache = make(map[string]*rsa.PublicKey)
    for _, key := range jwks.Keys {
        if key.Kty != "RSA" {
            continue
        }

        publicKey, err := v.parseRSAPublicKey(key)
        if err != nil {
            continue
        }

        v.jwksCache[key.Kid] = publicKey
    }

    v.lastCacheUpdate = time.Now()
    return nil
}

// parseRSAPublicKey converts a JWK to an RSA public key
func (v *JWTValidator) parseRSAPublicKey(jwk JWK) (*rsa.PublicKey, error) {
    nBytes, err := base64.RawURLEncoding.DecodeString(jwk.N)
    if err != nil {
        return nil, err
    }

    eBytes, err := base64.RawURLEncoding.DecodeString(jwk.E)
    if err != nil {
        return nil, err
    }

    n := new(big.Int).SetBytes(nBytes)
    e := new(big.Int).SetBytes(eBytes)

    return &rsa.PublicKey{
        N: n,
        E: int(e.Int64()),
    }, nil
}

// mapClaimsToDomain maps jwt.MapClaims to our domain Claims structure
func (v *JWTValidator) mapClaimsToDomain(mapClaims jwt.MapClaims, authClaims *auth.Claims) error {
    authClaims.Sub, _ = mapClaims["sub"].(string)
    authClaims.Email, _ = mapClaims["email"].(string)
    authClaims.EmailVerified, _ = mapClaims["email_verified"].(bool)
    authClaims.Username, _ = mapClaims["cognito:username"].(string)
    authClaims.TokenUse, _ = mapClaims["token_use"].(string)
    authClaims.Scope, _ = mapClaims["scope"].(string)
    authClaims.Issuer, _ = mapClaims["iss"].(string)
    authClaims.ClientID, _ = mapClaims["client_id"].(string)
    authClaims.Audience, _ = mapClaims["aud"].(string)

    // Handle numeric claims
    if authTime, ok := mapClaims["auth_time"].(float64); ok {
        authClaims.AuthTime = int64(authTime)
    }
    if iat, ok := mapClaims["iat"].(float64); ok {
        authClaims.IssuedAt = int64(iat)
    }
    if exp, ok := mapClaims["exp"].(float64); ok {
        authClaims.ExpirationTime = int64(exp)
    }

    // Handle groups array
    if groups, ok := mapClaims["cognito:groups"].([]interface{}); ok {
        authClaims.Groups = make([]string, 0, len(groups))
        for _, group := range groups {
            if groupStr, ok := group.(string); ok {
                authClaims.Groups = append(authClaims.Groups, groupStr)
            }
        }
    }

    return nil
}
```

---

### 4. Middleware Layer

#### Authentication Middleware

**File: `internal/transport/middleware/auth.go`**

```go
package middleware

import (
    "context"
    "net/http"
    "strings"

    "github.com/gin-gonic/gin"
    "sarc-ng/internal/domain/auth"
)

const (
    // ContextKeyUser is the key for user in context
    ContextKeyUser = "user"

    // ContextKeyClaims is the key for claims in context
    ContextKeyClaims = "claims"
)

// AuthMiddleware validates JWT tokens from request headers
func AuthMiddleware(validator auth.TokenValidator) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Extract token from Authorization header
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
                "error": "Authorization header required",
                "code":  "AUTH_HEADER_MISSING",
            })
            return
        }

        // Extract Bearer token
        tokenString := extractBearerToken(authHeader)
        if tokenString == "" {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
                "error": "Invalid authorization header format. Expected: Bearer <token>",
                "code":  "AUTH_HEADER_INVALID",
            })
            return
        }

        // Validate token
        claims, err := validator.ValidateToken(c.Request.Context(), tokenString)
        if err != nil {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
                "error": "Invalid or expired token",
                "code":  "TOKEN_INVALID",
            })
            return
        }

        // Set claims and user in context
        c.Set(ContextKeyClaims, claims)
        c.Set(ContextKeyUser, claims.ToUser())

        c.Next()
    }
}

// OptionalAuthMiddleware validates JWT if present but doesn't require it
func OptionalAuthMiddleware(validator auth.TokenValidator) gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.Next()
            return
        }

        tokenString := extractBearerToken(authHeader)
        if tokenString == "" {
            c.Next()
            return
        }

        claims, err := validator.ValidateToken(c.Request.Context(), tokenString)
        if err == nil {
            c.Set(ContextKeyClaims, claims)
            c.Set(ContextKeyUser, claims.ToUser())
        }

        c.Next()
    }
}

// RequireGroups middleware ensures user belongs to required groups
func RequireGroups(groups ...string) gin.HandlerFunc {
    return func(c *gin.Context) {
        user, exists := GetUserFromContext(c)
        if !exists {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
                "error": "User not authenticated",
                "code":  "USER_NOT_AUTHENTICATED",
            })
            return
        }

        if !user.HasAnyGroup(groups) {
            c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
                "error": "Insufficient permissions",
                "code":  "INSUFFICIENT_PERMISSIONS",
                "required_groups": groups,
            })
            return
        }

        c.Next()
    }
}

// RequireAllGroups middleware ensures user belongs to all specified groups
func RequireAllGroups(groups ...string) gin.HandlerFunc {
    return func(c *gin.Context) {
        user, exists := GetUserFromContext(c)
        if !exists {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
                "error": "User not authenticated",
                "code":  "USER_NOT_AUTHENTICATED",
            })
            return
        }

        for _, group := range groups {
            if !user.HasGroup(group) {
                c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
                    "error": "Insufficient permissions",
                    "code":  "INSUFFICIENT_PERMISSIONS",
                    "required_groups": groups,
                })
                return
            }
        }

        c.Next()
    }
}

// extractBearerToken extracts the token from "Bearer <token>" format
func extractBearerToken(authHeader string) string {
    parts := strings.SplitN(authHeader, " ", 2)
    if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
        return ""
    }
    return strings.TrimSpace(parts[1])
}

// GetUserFromContext retrieves the authenticated user from context
func GetUserFromContext(c *gin.Context) (*auth.User, bool) {
    value, exists := c.Get(ContextKeyUser)
    if !exists {
        return nil, false
    }

    user, ok := value.(*auth.User)
    return user, ok
}

// GetClaimsFromContext retrieves the JWT claims from context
func GetClaimsFromContext(c *gin.Context) (*auth.Claims, bool) {
    value, exists := c.Get(ContextKeyClaims)
    if !exists {
        return nil, false
    }

    claims, ok := value.(*auth.Claims)
    return claims, ok
}

// MustGetUser retrieves user from context or panics (use in handlers after auth middleware)
func MustGetUser(c *gin.Context) *auth.User {
    user, exists := GetUserFromContext(c)
    if !exists {
        panic("user not found in context - ensure AuthMiddleware is applied")
    }
    return user
}
```

---

### 5. Transport Layer Updates

#### Update Router Structure

**File: `internal/transport/rest/router.go`**

```go
package rest

import (
    "sarc-ng/internal/domain/auth"
    "sarc-ng/internal/domain/building"
    "sarc-ng/internal/domain/class"
    "sarc-ng/internal/domain/lesson"
    "sarc-ng/internal/domain/reservation"
    "sarc-ng/internal/domain/resource"
    "sarc-ng/internal/transport/middleware"
    buildingRest "sarc-ng/internal/transport/rest/building"
    classRest "sarc-ng/internal/transport/rest/class"
    lessonRest "sarc-ng/internal/transport/rest/lesson"
    reservationRest "sarc-ng/internal/transport/rest/reservation"
    resourceRest "sarc-ng/internal/transport/rest/resource"

    "github.com/gin-gonic/gin"
    "github.com/prometheus/client_golang/prometheus/promhttp"
    swaggerFiles "github.com/swaggo/files"
    ginSwagger "github.com/swaggo/gin-swagger"

    _ "sarc-ng/api/swagger"
)

// Router contains all the dependencies for setting up routes
type Router struct {
    buildingService    building.Usecase
    classService       class.Usecase
    lessonService      lesson.Usecase
    reservationService reservation.Usecase
    resourceService    resource.Usecase
    tokenValidator     auth.TokenValidator
}

// NewRouter creates a new router with all dependencies
func NewRouter(
    buildingService building.Usecase,
    classService class.Usecase,
    lessonService lesson.Usecase,
    reservationService reservation.Usecase,
    resourceService resource.Usecase,
    tokenValidator auth.TokenValidator,
) *Router {
    return &Router{
        buildingService:    buildingService,
        classService:       classService,
        lessonService:      lessonService,
        reservationService: reservationService,
        resourceService:    resourceService,
        tokenValidator:     tokenValidator,
    }
}

// SetupRoutes configures all the routes for the API
func (r *Router) SetupRoutes(router *gin.Engine) {
    // Apply global middleware
    r.setupMiddleware(router)

    // Setup system routes
    r.setupSystemRoutes(router)

    // Setup API routes
    r.setupAPIRoutes(router)
}

// setupMiddleware applies global middleware to the router
func (r *Router) setupMiddleware(router *gin.Engine) {
    router.Use(middleware.Logger())
    router.Use(middleware.Recovery())
    router.Use(middleware.CORS())
    router.Use(middleware.Metrics())
}

// setupSystemRoutes configures system routes like health and swagger
func (r *Router) setupSystemRoutes(router *gin.Engine) {
    // Swagger routes
    router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

    // Health check route
    router.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "status":  "healthy",
            "service": "sarc-ng",
        })
    })

    // Prometheus metrics endpoint
    router.GET("/metrics", gin.WrapH(promhttp.Handler()))
}

// setupAPIRoutes configures API v1 routes with authentication
func (r *Router) setupAPIRoutes(router *gin.Engine) {
    // Public API routes (no authentication required)
    publicV1 := router.Group("/api/v1")
    {
        // Public read-only endpoints
        buildingRest.RegisterPublicRoutes(publicV1, r.buildingService)
        resourceRest.RegisterPublicRoutes(publicV1, r.resourceService)
        classRest.RegisterPublicRoutes(publicV1, r.classService)
        lessonRest.RegisterPublicRoutes(publicV1, r.lessonService)
    }

    // Protected API routes (authentication required)
    protectedV1 := router.Group("/api/v1")
    protectedV1.Use(middleware.AuthMiddleware(r.tokenValidator))
    {
        // Protected endpoints requiring authentication
        reservationRest.RegisterProtectedRoutes(protectedV1, r.reservationService)

        // Admin endpoints requiring specific groups
        adminV1 := protectedV1.Group("")
        adminV1.Use(middleware.RequireGroups("admin", "manager"))
        {
            buildingRest.RegisterAdminRoutes(adminV1, r.buildingService)
            resourceRest.RegisterAdminRoutes(adminV1, r.resourceService)
            classRest.RegisterAdminRoutes(adminV1, r.classService)
            lessonRest.RegisterAdminRoutes(adminV1, r.lessonService)
        }
    }
}
```

---

### 6. Infrastructure as Code

#### Cognito Terraform Module

**File: `infrastructure/terraform/modules/idp/cognito/main.tf`**

```hcl
terraform {
  required_version = ">= 1.0"
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

resource "aws_cognito_user_pool" "main" {
  name = var.user_pool_name

  # Password policy
  password_policy {
    minimum_length                   = 8
    require_lowercase                = true
    require_uppercase                = true
    require_numbers                  = true
    require_symbols                  = true
    temporary_password_validity_days = 7
  }

  # Account recovery
  account_recovery_setting {
    recovery_mechanism {
      name     = "verified_email"
      priority = 1
    }
  }

  # Auto-verified attributes
  auto_verified_attributes = ["email"]

  # User attributes
  schema {
    attribute_data_type      = "String"
    name                     = "email"
    required                 = true
    mutable                  = true
    developer_only_attribute = false

    string_attribute_constraints {
      min_length = 5
      max_length = 256
    }
  }

  # MFA configuration
  mfa_configuration = var.mfa_configuration

  software_token_mfa_configuration {
    enabled = var.mfa_configuration != "OFF"
  }

  # Admin create user configuration
  admin_create_user_config {
    allow_admin_create_user_only = var.allow_admin_create_user_only

    invite_message_template {
      email_message = "Your username is {username} and temporary password is {####}."
      email_subject = "Your temporary password for ${var.application_name}"
      sms_message   = "Your username is {username} and temporary password is {####}."
    }
  }

  # User pool add-ons
  user_pool_add_ons {
    advanced_security_mode = var.advanced_security_mode
  }

  # Prevent destruction of user pool
  lifecycle {
    prevent_destroy = var.prevent_destroy
  }

  tags = merge(
    var.tags,
    {
      Name        = var.user_pool_name
      Environment = var.environment
      ManagedBy   = "Terraform"
    }
  )
}

# User Pool Domain
resource "aws_cognito_user_pool_domain" "main" {
  count        = var.domain_name != "" ? 1 : 0
  domain       = var.domain_name
  user_pool_id = aws_cognito_user_pool.main.id
}

# User Pool Client for Application
resource "aws_cognito_user_pool_client" "app_client" {
  name         = "${var.application_name}-client"
  user_pool_id = aws_cognito_user_pool.main.id

  generate_secret = var.generate_client_secret

  # OAuth flows
  allowed_oauth_flows_user_pool_client = true
  allowed_oauth_flows                  = var.oauth_flows
  allowed_oauth_scopes                 = var.oauth_scopes

  # Callback URLs
  callback_urls         = var.callback_urls
  logout_urls           = var.logout_urls
  default_redirect_uri  = length(var.callback_urls) > 0 ? var.callback_urls[0] : null

  # Supported identity providers
  supported_identity_providers = var.identity_providers

  # Token validity
  access_token_validity  = var.access_token_validity
  id_token_validity      = var.id_token_validity
  refresh_token_validity = var.refresh_token_validity

  token_validity_units {
    access_token  = "minutes"
    id_token      = "minutes"
    refresh_token = "days"
  }

  # Read/Write attributes
  read_attributes  = var.read_attributes
  write_attributes = var.write_attributes

  # Prevent user existence errors
  prevent_user_existence_errors = "ENABLED"

  # Enable token revocation
  enable_token_revocation = true

  # Explicit auth flows
  explicit_auth_flows = var.explicit_auth_flows
}

# User Pool Groups
resource "aws_cognito_user_group" "groups" {
  for_each = var.user_groups

  name         = each.key
  user_pool_id = aws_cognito_user_pool.main.id
  description  = each.value.description
  precedence   = each.value.precedence
  role_arn     = each.value.role_arn
}

# Identity Pool (optional - for AWS resource access)
resource "aws_cognito_identity_pool" "main" {
  count                            = var.create_identity_pool ? 1 : 0
  identity_pool_name               = "${var.application_name}-identity-pool"
  allow_unauthenticated_identities = var.allow_unauthenticated_identities

  cognito_identity_providers {
    client_id               = aws_cognito_user_pool_client.app_client.id
    provider_name           = aws_cognito_user_pool.main.endpoint
    server_side_token_check = true
  }

  tags = var.tags
}

# SSM Parameters for application configuration
resource "aws_ssm_parameter" "user_pool_id" {
  name        = "/${var.environment}/${var.application_name}/cognito/user-pool-id"
  description = "Cognito User Pool ID"
  type        = "String"
  value       = aws_cognito_user_pool.main.id

  tags = var.tags
}

resource "aws_ssm_parameter" "user_pool_client_id" {
  name        = "/${var.environment}/${var.application_name}/cognito/client-id"
  description = "Cognito User Pool Client ID"
  type        = "String"
  value       = aws_cognito_user_pool_client.app_client.id

  tags = var.tags
}

resource "aws_ssm_parameter" "user_pool_client_secret" {
  count       = var.generate_client_secret ? 1 : 0
  name        = "/${var.environment}/${var.application_name}/cognito/client-secret"
  description = "Cognito User Pool Client Secret"
  type        = "SecureString"
  value       = aws_cognito_user_pool_client.app_client.client_secret

  tags = var.tags
}
```

**File: `infrastructure/terraform/modules/idp/cognito/variables.tf`**

```hcl
variable "user_pool_name" {
  description = "Name of the Cognito User Pool"
  type        = string
}

variable "application_name" {
  description = "Name of the application"
  type        = string
}

variable "environment" {
  description = "Environment name (dev, staging, prod)"
  type        = string
}

variable "domain_name" {
  description = "Domain name for the Cognito hosted UI"
  type        = string
  default     = ""
}

variable "mfa_configuration" {
  description = "MFA configuration (OFF, ON, OPTIONAL)"
  type        = string
  default     = "OPTIONAL"
  validation {
    condition     = contains(["OFF", "ON", "OPTIONAL"], var.mfa_configuration)
    error_message = "MFA configuration must be OFF, ON, or OPTIONAL"
  }
}

variable "advanced_security_mode" {
  description = "Advanced security mode (OFF, AUDIT, ENFORCED)"
  type        = string
  default     = "AUDIT"
  validation {
    condition     = contains(["OFF", "AUDIT", "ENFORCED"], var.advanced_security_mode)
    error_message = "Advanced security mode must be OFF, AUDIT, or ENFORCED"
  }
}

variable "allow_admin_create_user_only" {
  description = "Allow only admin to create users"
  type        = bool
  default     = false
}

variable "prevent_destroy" {
  description = "Prevent destruction of the user pool"
  type        = bool
  default     = true
}

variable "generate_client_secret" {
  description = "Generate client secret for the app client"
  type        = bool
  default     = false
}

variable "oauth_flows" {
  description = "OAuth flows to enable"
  type        = list(string)
  default     = ["code", "implicit"]
}

variable "oauth_scopes" {
  description = "OAuth scopes to enable"
  type        = list(string)
  default     = ["email", "openid", "profile", "aws.cognito.signin.user.admin"]
}

variable "callback_urls" {
  description = "List of allowed callback URLs"
  type        = list(string)
  default     = []
}

variable "logout_urls" {
  description = "List of allowed logout URLs"
  type        = list(string)
  default     = []
}

variable "identity_providers" {
  description = "List of supported identity providers"
  type        = list(string)
  default     = ["COGNITO"]
}

variable "access_token_validity" {
  description = "Access token validity in minutes"
  type        = number
  default     = 60
}

variable "id_token_validity" {
  description = "ID token validity in minutes"
  type        = number
  default     = 60
}

variable "refresh_token_validity" {
  description = "Refresh token validity in days"
  type        = number
  default     = 30
}

variable "read_attributes" {
  description = "List of user pool attributes the app client can read"
  type        = list(string)
  default     = ["email", "email_verified"]
}

variable "write_attributes" {
  description = "List of user pool attributes the app client can write"
  type        = list(string)
  default     = ["email"]
}

variable "explicit_auth_flows" {
  description = "List of authentication flows"
  type        = list(string)
  default = [
    "ALLOW_USER_SRP_AUTH",
    "ALLOW_REFRESH_TOKEN_AUTH",
    "ALLOW_USER_PASSWORD_AUTH"
  ]
}

variable "user_groups" {
  description = "Map of user groups to create"
  type = map(object({
    description = string
    precedence  = number
    role_arn    = string
  }))
  default = {}
}

variable "create_identity_pool" {
  description = "Create a Cognito Identity Pool"
  type        = bool
  default     = false
}

variable "allow_unauthenticated_identities" {
  description = "Allow unauthenticated identities in identity pool"
  type        = bool
  default     = false
}

variable "tags" {
  description = "Tags to apply to resources"
  type        = map(string)
  default     = {}
}
```

**File: `infrastructure/terraform/modules/idp/cognito/outputs.tf`**

```hcl
output "user_pool_id" {
  description = "The ID of the Cognito User Pool"
  value       = aws_cognito_user_pool.main.id
}

output "user_pool_arn" {
  description = "The ARN of the Cognito User Pool"
  value       = aws_cognito_user_pool.main.arn
}

output "user_pool_endpoint" {
  description = "The endpoint of the Cognito User Pool"
  value       = aws_cognito_user_pool.main.endpoint
}

output "user_pool_client_id" {
  description = "The ID of the Cognito User Pool Client"
  value       = aws_cognito_user_pool_client.app_client.id
}

output "user_pool_client_secret" {
  description = "The secret of the Cognito User Pool Client"
  value       = aws_cognito_user_pool_client.app_client.client_secret
  sensitive   = true
}

output "user_pool_domain" {
  description = "The domain of the Cognito User Pool"
  value       = var.domain_name != "" ? aws_cognito_user_pool_domain.main[0].domain : null
}

output "identity_pool_id" {
  description = "The ID of the Cognito Identity Pool"
  value       = var.create_identity_pool ? aws_cognito_identity_pool.main[0].id : null
}

output "user_groups" {
  description = "Map of created user groups"
  value = {
    for k, v in aws_cognito_user_group.groups : k => {
      name        = v.name
      description = v.description
      precedence  = v.precedence
    }
  }
}

output "issuer_url" {
  description = "The issuer URL for JWT validation"
  value       = "https://cognito-idp.${data.aws_region.current.name}.amazonaws.com/${aws_cognito_user_pool.main.id}"
}

output "jwks_uri" {
  description = "The JWKS URI for token validation"
  value       = "https://cognito-idp.${data.aws_region.current.name}.amazonaws.com/${aws_cognito_user_pool.main.id}/.well-known/jwks.json"
}

data "aws_region" "current" {}
```

---

### 7. Dependency Injection with Wire

**File: `cmd/server/wire.go`**

```go
//go:build wireinject
// +build wireinject

package main

import (
    "sarc-ng/internal/adapter/db"
    gormBuilding "sarc-ng/internal/adapter/gorm/building"
    gormClass "sarc-ng/internal/adapter/gorm/class"
    gormLesson "sarc-ng/internal/adapter/gorm/lesson"
    gormReservation "sarc-ng/internal/adapter/gorm/reservation"
    gormResource "sarc-ng/internal/adapter/gorm/resource"
    "sarc-ng/internal/config"
    buildingService "sarc-ng/internal/service/building"
    classService "sarc-ng/internal/service/class"
    lessonService "sarc-ng/internal/service/lesson"
    reservationService "sarc-ng/internal/service/reservation"
    resourceService "sarc-ng/internal/service/resource"
    authService "sarc-ng/internal/service/auth"
    "sarc-ng/internal/transport/rest"

    "github.com/google/wire"
)

// initializeRouter creates and initializes the router with all dependencies
func initializeRouter(cfg *config.Config) (*rest.Router, error) {
    wire.Build(
        // Database connection
        db.NewConnection,

        // Token validator
        provideTokenValidator,

        // Repositories
        gormBuilding.NewAdapter,
        wire.Bind(new(building.Repository), new(*gormBuilding.Adapter)),

        gormClass.NewAdapter,
        wire.Bind(new(class.Repository), new(*gormClass.Adapter)),

        gormLesson.NewAdapter,
        wire.Bind(new(lesson.Repository), new(*gormLesson.Adapter)),

        gormReservation.NewAdapter,
        wire.Bind(new(reservation.Repository), new(*gormReservation.Adapter)),

        gormResource.NewAdapter,
        wire.Bind(new(resource.Repository), new(*gormResource.Adapter)),

        // Services
        buildingService.NewService,
        wire.Bind(new(building.Usecase), new(*buildingService.Service)),

        classService.NewService,
        wire.Bind(new(class.Usecase), new(*classService.Service)),

        lessonService.NewService,
        wire.Bind(new(lesson.Usecase), new(*lessonService.Service)),

        reservationService.NewService,
        wire.Bind(new(reservation.Usecase), new(*reservationService.Service)),

        resourceService.NewService,
        wire.Bind(new(resource.Usecase), new(*resourceService.Service)),

        // Router
        rest.NewRouter,
    )

    return nil, nil
}

// provideTokenValidator creates a new JWT token validator
func provideTokenValidator(cfg *config.Config) auth.TokenValidator {
    return authService.NewJWTValidator(
        cfg.Cognito.Region,
        cfg.Cognito.UserPoolID,
        cfg.Cognito.ClientID,
        cfg.Cognito.JWKSCacheExp,
    )
}
```

---

### 8. Updated Dependencies

**Add to `go.mod`:**

```go
require (
    github.com/golang-jwt/jwt/v5 v5.2.0
    // ... existing dependencies
)
```

---

### 9. Testing Strategy

#### Unit Tests for JWT Validator

**File: `internal/service/auth/jwt_validator_test.go`**

```go
package auth

import (
    "context"
    "crypto/rand"
    "crypto/rsa"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
    "time"

    "github.com/golang-jwt/jwt/v5"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestJWTValidator_ValidateToken(t *testing.T) {
    // Generate test RSA keys
    privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
    require.NoError(t, err)

    publicKey := &privateKey.PublicKey

    // Setup mock JWKS server
    jwksServer := setupMockJWKSServer(t, publicKey)
    defer jwksServer.Close()

    // Create validator with mock JWKS URL
    validator := createTestValidator(t, jwksServer.URL)

    tests := []struct {
        name        string
        token       string
        expectError bool
        errorType   error
    }{
        {
            name:        "valid token",
            token:       createValidToken(t, privateKey),
            expectError: false,
        },
        {
            name:        "expired token",
            token:       createExpiredToken(t, privateKey),
            expectError: true,
            errorType:   ErrExpiredToken,
        },
        {
            name:        "invalid signature",
            token:       createInvalidSignatureToken(t),
            expectError: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            claims, err := validator.ValidateToken(context.Background(), tt.token)

            if tt.expectError {
                assert.Error(t, err)
                assert.Nil(t, claims)
            } else {
                assert.NoError(t, err)
                assert.NotNil(t, claims)
            }
        })
    }
}
```

#### Integration Tests

**File: `test/integration/auth_test.go`**

```go
package integration

import (
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestAuthenticationFlow(t *testing.T) {
    // Setup test router with auth middleware
    router := setupTestRouter(t)

    tests := []struct {
        name           string
        endpoint       string
        method         string
        token          string
        expectedStatus int
    }{
        {
            name:           "public endpoint without token",
            endpoint:       "/api/v1/buildings",
            method:         http.MethodGet,
            token:          "",
            expectedStatus: http.StatusOK,
        },
        {
            name:           "protected endpoint without token",
            endpoint:       "/api/v1/reservations",
            method:         http.MethodPost,
            token:          "",
            expectedStatus: http.StatusUnauthorized,
        },
        {
            name:           "protected endpoint with valid token",
            endpoint:       "/api/v1/reservations",
            method:         http.MethodPost,
            token:          getValidTestToken(),
            expectedStatus: http.StatusOK,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            req := httptest.NewRequest(tt.method, tt.endpoint, nil)
            if tt.token != "" {
                req.Header.Set("Authorization", "Bearer "+tt.token)
            }

            w := httptest.NewRecorder()
            router.ServeHTTP(w, req)

            assert.Equal(t, tt.expectedStatus, w.Code)
        })
    }
}
```

---

### 10. Implementation Phases

#### Phase 1: Foundation (Week 1)

- [ ] Update configuration structures
- [ ] Add Cognito configuration to YAML files
- [ ] Create domain entities for authentication
- [ ] Add golang-jwt/jwt dependency

#### Phase 2: Core Implementation (Week 1-2)

- [ ] Implement JWT validator service
- [ ] Create authentication middleware
- [ ] Add authorization middleware (groups, scopes)
- [ ] Update dependency injection with Wire

#### Phase 3: Router Integration (Week 2)

- [ ] Update router to support protected routes
- [ ] Separate public and protected endpoints
- [ ] Add admin-only routes with group requirements
- [ ] Update existing handlers to use user context

#### Phase 4: Infrastructure (Week 2-3)

- [ ] Create Cognito Terraform module
- [ ] Setup user pool with proper security settings
- [ ] Configure user groups (admin, manager, user)
- [ ] Create SSM parameters for configuration

#### Phase 5: Testing (Week 3)

- [ ] Write unit tests for JWT validator
- [ ] Create integration tests for auth flow
- [ ] Test authorization middleware
- [ ] Perform security testing

#### Phase 6: Documentation & Deployment (Week 3-4)

- [ ] Update API documentation
- [ ] Create authentication setup guide
- [ ] Document user management procedures
- [ ] Deploy to development environment
- [ ] Conduct security audit

---

### 11. Security Best Practices

#### Token Validation

- ✅ Verify token signature using JWKS
- ✅ Validate issuer matches Cognito user pool
- ✅ Check token expiration
- ✅ Verify audience/client_id
- ✅ Validate token_use claim
- ✅ Cache JWKS with expiration

#### Password Policy

- ✅ Minimum 8 characters
- ✅ Require uppercase, lowercase, numbers, symbols
- ✅ Temporary password validity (7 days)

#### MFA Configuration

- ✅ Optional MFA by default
- ✅ Support TOTP (software tokens)
- ✅ Can be enforced for specific groups

#### Advanced Security

- ✅ Enable advanced security mode (AUDIT/ENFORCED)
- ✅ Detect compromised credentials
- ✅ Risk-based adaptive authentication
- ✅ Account takeover protection

#### Session Management

- ✅ Short-lived access tokens (60 minutes)
- ✅ Short-lived ID tokens (60 minutes)
- ✅ Long-lived refresh tokens (30 days)
- ✅ Token revocation support

#### HTTPS Only

- ✅ Enforce HTTPS in production
- ✅ Secure cookie flags
- ✅ HSTS headers

---

### 12. Authorization Model

#### User Groups

```
admin
  - Full system access
  - User management
  - All CRUD operations

manager
  - Manage resources in their domain
  - Create/update classes and lessons
  - View all reservations

teacher
  - Create classes and lessons
  - View their own classes
  - Manage their own schedule

student
  - Create reservations
  - View available resources
  - Cancel own reservations
```

#### Example Route Protection

```go
// Public - anyone
GET /api/v1/buildings
GET /api/v1/resources

// Authenticated - any logged in user
POST /api/v1/reservations
GET /api/v1/reservations/my

// Manager or Admin only
POST /api/v1/classes
PUT /api/v1/classes/:id

// Admin only
POST /api/v1/buildings
DELETE /api/v1/buildings/:id
POST /api/v1/users
```

---

### 13. Environment Variables

```bash
# AWS Configuration
AWS_REGION=us-east-1

# Cognito Configuration
COGNITO_USER_POOL_ID=us-east-1_XXXXXXXXX
COGNITO_CLIENT_ID=xxxxxxxxxxxxxxxxxxxx
COGNITO_JWKS_CACHE_EXPIRY=1h

# Application Configuration
SERVER_PORT=8080
DATABASE_HOST=localhost
DATABASE_PORT=3306
```

---

### 14. API Request Examples

#### Public Endpoint (No Auth)

```bash
curl -X GET http://localhost:8080/api/v1/buildings
```

#### Protected Endpoint (With Auth)

```bash
curl -X POST http://localhost:8080/api/v1/reservations \
  -H "Authorization: Bearer <JWT_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{"resource_id": 1, "start_time": "2024-10-15T10:00:00Z"}'
```

#### Get User Info from Token

```bash
# The user information is automatically extracted
# from the token and available in request context
```

---

### 15. Migration Strategy

#### For Existing Deployments

1. **Phase 1: Deploy with Optional Auth**
   - Deploy authentication code
   - Make all routes public initially
   - Test authentication with test users

2. **Phase 2: Enable Auth for New Features**
   - Require auth for new endpoints
   - Keep existing endpoints public

3. **Phase 3: Migrate Users**
   - Create Cognito users for existing users
   - Notify users of new authentication

4. **Phase 4: Enforce Authentication**
   - Make all sensitive endpoints protected
   - Deprecate old authentication method

---

### 16. Monitoring and Logging

#### Metrics to Track

- Authentication success/failure rate
- Token validation latency
- JWKS cache hit rate
- Failed authorization attempts
- Token expiration rate

#### Logging Requirements

- Log authentication attempts (success/failure)
- Log authorization failures with user context
- Log JWKS refresh operations
- Sanitize sensitive data (never log tokens)

---

### 17. Cost Estimation

#### AWS Cognito Pricing (US East)

- First 50,000 MAUs: Free
- 50,001 - 100,000 MAUs: $0.0055/MAU
- Advanced Security: +$0.05/MAU (if enabled)

#### Example Monthly Costs

- 1,000 users: $0 (free tier)
- 10,000 users: $0 (free tier)
- 100,000 users: ~$275/month
- With Advanced Security (100K): ~$5,275/month

---

### 18. Compliance and Standards

- ✅ OWASP Top 10 compliance
- ✅ OAuth 2.0 / OpenID Connect standards
- ✅ JWT best practices (RFC 7519)
- ✅ GDPR-ready (user data management)
- ✅ SOC 2 Type II (AWS Cognito)
- ✅ HIPAA eligible (with BAA)

---

### 19. Rollback Plan

#### If Issues Occur

1. Remove authentication middleware from router
2. Revert to public endpoints
3. Roll back infrastructure changes
4. Investigate and fix issues
5. Redeploy with fixes

#### Database Rollback

- No database changes required
- Authentication is stateless (JWT)
- User data in Cognito (separate system)

---

### 20. Success Criteria

- [ ] All protected endpoints require valid JWT
- [ ] Public endpoints remain accessible
- [ ] User groups properly enforced
- [ ] Token validation < 100ms p99
- [ ] Zero authentication bypasses
- [ ] Comprehensive test coverage (>80%)
- [ ] Security audit passed
- [ ] Documentation complete
- [ ] Successfully deployed to dev/staging/prod

---

## Conclusion

This proposal provides a comprehensive, production-ready implementation of AWS Cognito authentication following industry best practices and clean architecture principles. The solution is:

- **Secure**: Uses AWS Cognito with advanced security features
- **Scalable**: Stateless JWT authentication
- **Maintainable**: Clean architecture with clear separation of concerns
- **Flexible**: Supports multiple user groups and authorization models
- **Testable**: Comprehensive unit and integration tests
- **Infrastructure as Code**: Terraform modules for repeatable deployments
- **Cost-Effective**: Free tier covers most small-medium applications

The phased implementation approach allows for gradual rollout with minimal risk, while the monitoring and rollback plans ensure production stability.


---

# AWS Cognito Authentication - Implementation Checklist

## Table of Contents

- [Quick Summary](#quick-summary)
- [Prerequisites](#prerequisites)
- [Phase 1: Dependencies & Configuration](#phase-1-dependencies--configuration)
- [Phase 2: Domain Layer Implementation](#phase-2-domain-layer-implementation)
- [Phase 3: Service Layer Implementation](#phase-3-service-layer-implementation)
- [Phase 4: Middleware Implementation](#phase-4-middleware-implementation)
- [Phase 5: Router Updates](#phase-5-router-updates)
- [Phase 6: Dependency Injection](#phase-6-dependency-injection)
- [Phase 7: Infrastructure as Code](#phase-7-infrastructure-as-code)
- [Phase 8: Testing](#phase-8-testing)
- [Phase 9: Documentation](#phase-9-documentation)
- [Phase 10: Deployment](#phase-10-deployment)
- [Phase 11: Monitoring & Observability](#phase-11-monitoring--observability)
- [Phase 12: Security Audit](#phase-12-security-audit)
- [Rollback Procedures](#rollback-procedures)
- [Success Metrics](#success-metrics)
- [Support and Troubleshooting](#support-and-troubleshooting)

---

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
internal/transport/middleware/
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

internal/transport/middleware/
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
go test ./internal/transport/middleware/... -v

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
make build
cd infrastructure/sam && make deploy-all ENV=prod

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


---

# Authentication Architecture Overview

## System Architecture

```
┌─────────────────────────────────────────────────────────────────────┐
│                          Client Application                          │
│                     (Web, Mobile, CLI)                               │
└────────────────────────────┬────────────────────────────────────────┘
                             │
                             │ 1. Login Request
                             ▼
┌─────────────────────────────────────────────────────────────────────┐
│                         AWS Cognito                                  │
│  ┌────────────────┐  ┌─────────────────┐  ┌────────────────────┐  │
│  │  User Pool     │  │   User Groups   │  │   Hosted UI        │  │
│  │  - Users       │  │   - admin       │  │   (Optional)       │  │
│  │  - Passwords   │  │   - manager     │  │                    │  │
│  │  - MFA         │  │   - teacher     │  │                    │  │
│  │  - Attributes  │  │   - student     │  │                    │  │
│  └────────────────┘  └─────────────────┘  └────────────────────┘  │
└────────────────────────────┬────────────────────────────────────────┘
                             │
                             │ 2. JWT Token (Access + ID)
                             ▼
┌─────────────────────────────────────────────────────────────────────┐
│                       SARC-NG API Gateway                            │
│                                                                       │
│  ┌────────────────────────────────────────────────────────────────┐ │
│  │                 Middleware Pipeline                            │ │
│  │                                                                 │ │
│  │  1. Logger        → Log request details                        │ │
│  │  2. Recovery      → Handle panics                              │ │
│  │  3. CORS          → Handle cross-origin                        │ │
│  │  4. Metrics       → Collect metrics                            │ │
│  │  5. Auth          → Validate JWT token ◄──────┐               │ │
│  │                                                │               │ │
│  └────────────────────────────────────────────────┼───────────────┘ │
│                                                    │                 │
│  ┌────────────────────────────────────────────────┴───────────────┐ │
│  │              JWT Validator Service                             │ │
│  │                                                                 │ │
│  │  - Fetch JWKS from Cognito                                     │ │
│  │  - Cache public keys (1 hour)                                  │ │
│  │  - Verify token signature                                      │ │
│  │  - Validate claims:                                            │ │
│  │    • Issuer (iss)                                              │ │
│  │    • Audience (aud)                                            │ │
│  │    • Expiration (exp)                                          │ │
│  │    • Token use (token_use)                                     │ │
│  │  - Extract user information                                    │ │
│  └────────────────────────────────────────────────────────────────┘ │
│                                                                       │
│  ┌────────────────────────────────────────────────────────────────┐ │
│  │                    Route Groups                                │ │
│  │                                                                 │ │
│  │  ┌──────────────────────────────────────────────────────────┐ │ │
│  │  │  Public Routes (No Auth)                                 │ │ │
│  │  │  GET  /api/v1/buildings                                  │ │ │
│  │  │  GET  /api/v1/resources                                  │ │ │
│  │  │  GET  /api/v1/classes                                    │ │ │
│  │  └──────────────────────────────────────────────────────────┘ │ │
│  │                                                                 │ │
│  │  ┌──────────────────────────────────────────────────────────┐ │ │
│  │  │  Protected Routes (Auth Required)                        │ │ │
│  │  │  POST /api/v1/reservations         [student+]           │ │ │
│  │  │  GET  /api/v1/reservations/my      [student+]           │ │ │
│  │  │  POST /api/v1/classes               [teacher+]           │ │ │
│  │  └──────────────────────────────────────────────────────────┘ │ │
│  │                                                                 │ │
│  │  ┌──────────────────────────────────────────────────────────┐ │ │
│  │  │  Admin Routes (Admin/Manager Only)                       │ │ │
│  │  │  POST   /api/v1/buildings          [admin, manager]     │ │ │
│  │  │  DELETE /api/v1/buildings/:id      [admin]              │ │ │
│  │  │  POST   /api/v1/resources          [admin, manager]     │ │ │
│  │  └──────────────────────────────────────────────────────────┘ │ │
│  └────────────────────────────────────────────────────────────────┘ │
└───────────────────────────┬───────────────────────────────────────────┘
                            │
                            │ 4. Response
                            ▼
┌─────────────────────────────────────────────────────────────────────┐
│                      Application Services                            │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐             │
│  │  Building    │  │  Reservation │  │   Resource   │             │
│  │  Service     │  │  Service     │  │   Service    │             │
│  └──────────────┘  └──────────────┘  └──────────────┘             │
└───────────────────────────┬───────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────────────┐
│                         Database Layer                               │
│                        (MySQL/RDS)                                   │
└─────────────────────────────────────────────────────────────────────┘
```

---

## Authentication Flow

### 1. User Login Flow

```
┌────────┐         ┌─────────┐         ┌──────────────┐
│ Client │         │ Cognito │         │  SARC-NG API │
└───┬────┘         └────┬────┘         └──────┬───────┘
    │                   │                      │
    │ 1. Login Request  │                      │
    ├──────────────────►│                      │
    │ (username/pass)   │                      │
    │                   │                      │
    │ 2. JWT Tokens     │                      │
    │◄──────────────────┤                      │
    │ (access + id)     │                      │
    │                   │                      │
    │ 3. API Request + Bearer Token            │
    ├──────────────────────────────────────────►
    │                   │                      │
    │                   │ 4. Fetch JWKS        │
    │                   │◄─────────────────────┤
    │                   │                      │
    │                   │ 5. JWKS Response     │
    │                   ├─────────────────────►│
    │                   │                      │
    │                   │ 6. Validate Token    │
    │                   │      (in memory)     │
    │                   │                      │
    │ 7. API Response with User Context        │
    │◄──────────────────────────────────────────┤
    │                   │                      │
```

### 2. Token Validation Flow

```
┌──────────────────────────────────────────────────────────────┐
│                    Token Validation                           │
└───────────────────────────────┬──────────────────────────────┘
                                │
                    ┌───────────▼────────────┐
                    │  Extract Bearer Token  │
                    └───────────┬────────────┘
                                │
                    ┌───────────▼────────────┐
                    │  Parse JWT (no verify) │
                    │  Extract 'kid' header  │
                    └───────────┬────────────┘
                                │
                    ┌───────────▼────────────┐
                    │  Get Public Key        │
                    │  (from JWKS cache)     │
                    └───────────┬────────────┘
                                │
                        ┌───────┴────────┐
                        │ Key in cache?  │
                        └───┬────────┬───┘
                            │ No     │ Yes
                  ┌─────────▼──┐     │
                  │ Fetch JWKS │     │
                  │from Cognito│     │
                  └─────────┬──┘     │
                            │        │
                  ┌─────────▼────────▼───┐
                  │  Verify Signature    │
                  │  using Public Key    │
                  └─────────┬────────────┘
                            │
                  ┌─────────▼────────────┐
                  │  Validate Claims     │
                  │  - Issuer (iss)      │
                  │  - Audience (aud)    │
                  │  - Expiration (exp)  │
                  │  - Token Use         │
                  └─────────┬────────────┘
                            │
                  ┌─────────▼────────────┐
                  │  Extract User Info   │
                  │  - ID (sub)          │
                  │  - Email             │
                  │  - Username          │
                  │  - Groups            │
                  └─────────┬────────────┘
                            │
                  ┌─────────▼────────────┐
                  │  Set in Context      │
                  │  - User object       │
                  │  - Claims object     │
                  └──────────────────────┘
```

---

## Token Structure

### Access Token Claims

```json
{
  "sub": "12345678-1234-1234-1234-123456789012",
  "event_id": "abc123",
  "token_use": "access",
  "scope": "aws.cognito.signin.user.admin",
  "auth_time": 1697001600,
  "iss": "https://cognito-idp.us-east-1.amazonaws.com/us-east-1_XXXXXXXXX",
  "exp": 1697005200,
  "iat": 1697001600,
  "jti": "xyz789",
  "client_id": "abcdefghijklmnopqrstuvwx",
  "username": "john.doe",
  "cognito:groups": ["admin", "manager"]
}
```

### ID Token Claims

```json
{
  "sub": "12345678-1234-1234-1234-123456789012",
  "email_verified": true,
  "iss": "https://cognito-idp.us-east-1.amazonaws.com/us-east-1_XXXXXXXXX",
  "cognito:username": "john.doe",
  "cognito:groups": ["admin", "manager"],
  "aud": "abcdefghijklmnopqrstuvwx",
  "token_use": "id",
  "auth_time": 1697001600,
  "exp": 1697005200,
  "iat": 1697001600,
  "email": "john.doe@example.com"
}
```

---

## Authorization Model

### Role Hierarchy

```
┌──────────────────────────────────────────────────────────────┐
│                         Admin                                 │
│  - Full system access                                         │
│  - User management                                            │
│  - All CRUD operations                                        │
│  - System configuration                                       │
└───────────────────────────┬──────────────────────────────────┘
                            │
            ┌───────────────┴────────────────┐
            │                                │
┌───────────▼──────────────┐    ┌───────────▼──────────────┐
│        Manager           │    │        Teacher           │
│  - Resource management   │    │  - Create classes        │
│  - Create/edit classes   │    │  - Manage own schedule   │
│  - View all reservations │    │  - View own classes      │
│  - Generate reports      │    │  - Student management    │
└───────────┬──────────────┘    └───────────┬──────────────┘
            │                                │
            └────────────────┬───────────────┘
                             │
                 ┌───────────▼──────────────┐
                 │        Student           │
                 │  - Create reservations   │
                 │  - View resources        │
                 │  - Cancel own bookings   │
                 │  - View own schedule     │
                 └──────────────────────────┘
```

### Permission Matrix

| Resource      | Public | Student | Teacher | Manager | Admin |
|---------------|--------|---------|---------|---------|-------|
| Buildings     |        |         |         |         |       |
| - List        | ✓      | ✓       | ✓       | ✓       | ✓     |
| - View        | ✓      | ✓       | ✓       | ✓       | ✓     |
| - Create      |        |         |         | ✓       | ✓     |
| - Update      |        |         |         | ✓       | ✓     |
| - Delete      |        |         |         |         | ✓     |
|               |        |         |         |         |       |
| Resources     |        |         |         |         |       |
| - List        | ✓      | ✓       | ✓       | ✓       | ✓     |
| - View        | ✓      | ✓       | ✓       | ✓       | ✓     |
| - Create      |        |         |         | ✓       | ✓     |
| - Update      |        |         |         | ✓       | ✓     |
| - Delete      |        |         |         |         | ✓     |
|               |        |         |         |         |       |
| Classes       |        |         |         |         |       |
| - List        | ✓      | ✓       | ✓       | ✓       | ✓     |
| - View        | ✓      | ✓       | ✓       | ✓       | ✓     |
| - Create      |        |         | ✓       | ✓       | ✓     |
| - Update (own)|        |         | ✓       | ✓       | ✓     |
| - Update (all)|        |         |         | ✓       | ✓     |
| - Delete      |        |         |         | ✓       | ✓     |
|               |        |         |         |         |       |
| Reservations  |        |         |         |         |       |
| - List (all)  |        |         |         | ✓       | ✓     |
| - List (own)  |        | ✓       | ✓       | ✓       | ✓     |
| - Create      |        | ✓       | ✓       | ✓       | ✓     |
| - Cancel (own)|        | ✓       | ✓       | ✓       | ✓     |
| - Cancel (all)|        |         |         | ✓       | ✓     |
|               |        |         |         |         |       |
| Users         |        |         |         |         |       |
| - List        |        |         |         |         | ✓     |
| - Create      |        |         |         |         | ✓     |
| - Update      |        |         |         |         | ✓     |
| - Delete      |        |         |         |         | ✓     |
| - View (self) |        | ✓       | ✓       | ✓       | ✓     |

---

## Code Structure

### Directory Layout

```
sarc-ng/
├── internal/
│   ├── domain/
│   │   └── auth/                    # NEW - Authentication domain
│   │       ├── entity.go            # User, Claims entities
│   │       └── service.go           # TokenValidator interface
│   │
│   ├── service/
│   │   └── auth/                    # NEW - Auth service implementation
│   │       ├── jwt_validator.go     # JWT validation logic
│   │       └── jwt_validator_test.go
│   │
│   ├── config/
│   │   └── config.go                # UPDATED - Add CognitoConfig
│   │
│   └── transport/
│       └── rest/
│           └── router.go            # UPDATED - Add auth middleware
│
├── pkg/
│   └── rest/
│       └── middleware/
│           └── auth.go              # NEW - Auth middleware
│
├── infrastructure/
│   └── terraform/
│       └── modules/
│           └── idp/
│               └── cognito/         # NEW - Cognito Terraform module
│                   ├── main.tf
│                   ├── variables.tf
│                   ├── outputs.tf
│                   └── versions.tf
│
└── configs/
    ├── default.yaml                 # UPDATED - Add cognito config
    └── development.yaml             # UPDATED - Add cognito config
```

---

## Request/Response Examples

### 1. Login to Cognito (External)

**Request:**

```bash
curl -X POST https://cognito-idp.us-east-1.amazonaws.com/ \
  -H "Content-Type: application/x-amz-json-1.1" \
  -H "X-Amz-Target: AWSCognitoIdentityProviderService.InitiateAuth" \
  -d '{
    "AuthFlow": "USER_PASSWORD_AUTH",
    "ClientId": "abcdefghijklmnopqrstuvwx",
    "AuthParameters": {
      "USERNAME": "john.doe",
      "PASSWORD": "SecurePassword123!"
    }
  }'
```

**Response:**

```json
{
  "AuthenticationResult": {
    "AccessToken": "eyJraWQiOiI...",
    "IdToken": "eyJraWQiOiJ...",
    "RefreshToken": "eyJjdHkiOi...",
    "ExpiresIn": 3600,
    "TokenType": "Bearer"
  }
}
```

### 2. Access Protected Endpoint

**Request:**

```bash
curl -X POST http://localhost:8080/api/v1/reservations \
  -H "Authorization: Bearer eyJraWQiOiI..." \
  -H "Content-Type: application/json" \
  -d '{
    "resource_id": 5,
    "start_time": "2024-10-15T10:00:00Z",
    "end_time": "2024-10-15T11:00:00Z",
    "purpose": "Team meeting"
  }'
```

**Success Response (200):**

```json
{
  "id": 123,
  "resource_id": 5,
  "user_id": "12345678-1234-1234-1234-123456789012",
  "user_email": "john.doe@example.com",
  "start_time": "2024-10-15T10:00:00Z",
  "end_time": "2024-10-15T11:00:00Z",
  "purpose": "Team meeting",
  "status": "confirmed",
  "created_at": "2024-10-10T14:23:45Z"
}
```

**Error Response - No Token (401):**

```json
{
  "error": "Authorization header required",
  "code": "AUTH_HEADER_MISSING"
}
```

**Error Response - Invalid Token (401):**

```json
{
  "error": "Invalid or expired token",
  "code": "TOKEN_INVALID"
}
```

**Error Response - Insufficient Permissions (403):**

```json
{
  "error": "Insufficient permissions",
  "code": "INSUFFICIENT_PERMISSIONS",
  "required_groups": ["admin", "manager"]
}
```

### 3. Get User Info from Context

**In Handler:**

```go
func (h *Handler) CreateReservation(c *gin.Context) {
    // Get authenticated user from context
    user := middleware.MustGetUser(c)

    // User information is available
    log.Printf("User ID: %s", user.ID)
    log.Printf("Email: %s", user.Email)
    log.Printf("Groups: %v", user.Groups)

    // Check permissions
    if user.HasGroup("admin") {
        // Admin-specific logic
    }

    // Continue with business logic
    // ...
}
```

---

## Security Considerations

### Token Validation Checklist

✅ **Signature Verification**

- Token signed with RSA private key from Cognito
- Signature verified using public key from JWKS
- Keys cached for performance (1 hour)

✅ **Claim Validation**

- `iss` (issuer) matches Cognito user pool URL
- `aud` or `client_id` matches application client ID
- `exp` (expiration) is in the future
- `token_use` is either "access" or "id"

✅ **Transport Security**

- HTTPS enforced in production
- Tokens never logged
- Tokens never stored in database

✅ **Token Lifecycle**

- Access tokens: 60 minutes
- ID tokens: 60 minutes
- Refresh tokens: 30 days
- Token revocation supported

### Common Vulnerabilities Prevented

| Vulnerability | Mitigation |
|---------------|------------|
| Token Tampering | RSA signature verification |
| Token Replay | Short expiration times |
| Man-in-the-Middle | HTTPS only |
| Token Theft | Secure storage, HTTPS |
| Privilege Escalation | Group-based authorization |
| Brute Force | Cognito rate limiting |
| Account Enumeration | Generic error messages |
| Information Disclosure | No sensitive data in tokens |

---

## Performance Characteristics

### Latency Targets

| Operation | Target | Notes |
|-----------|--------|-------|
| Token Validation (cached JWKS) | < 50ms | p99 |
| Token Validation (fetch JWKS) | < 200ms | p99 |
| Authorization Check | < 10ms | p99 |
| User Context Extraction | < 1ms | p99 |

### Caching Strategy

```
JWKS Cache:
├── Cache Duration: 1 hour (configurable)
├── Cache Key: Key ID (kid)
├── Cache Invalidation: Time-based
└── Cache Size: ~10 keys (typical)

Benefits:
├── Reduces latency by 150ms per request
├── Reduces load on Cognito
└── Provides resilience if Cognito temporarily unavailable
```

### Scalability

- **Stateless Authentication**: No server-side sessions
- **Horizontal Scaling**: Each instance validates independently
- **No Database Dependency**: Token validation doesn't hit database
- **Cognito SLA**: 99.99% availability

---

## Monitoring Dashboard

### Key Metrics

```
Authentication Metrics:
├── auth_requests_total
├── auth_failures_total  (by reason)
├── token_validation_duration_seconds
└── jwks_cache_hit_rate

Authorization Metrics:
├── authorization_checks_total
├── authorization_failures_total (by endpoint)
└── group_requirement_failures_total

Performance Metrics:
├── jwks_fetch_duration_seconds
├── jwks_cache_size
└── active_authenticated_sessions (approx)

Error Metrics:
├── invalid_token_count
├── expired_token_count
├── unauthorized_access_attempts
└── jwks_fetch_failures_total
```

### Alert Conditions

```yaml
- alert: HighAuthenticationFailureRate
  expr: rate(auth_failures_total[5m]) > 0.1
  severity: warning

- alert: JWKSFetchFailure
  expr: jwks_fetch_failures_total > 0
  severity: critical

- alert: HighTokenValidationLatency
  expr: histogram_quantile(0.99, token_validation_duration_seconds) > 0.2
  severity: warning

- alert: UnauthorizedAccessSpike
  expr: rate(authorization_failures_total[5m]) > 10
  severity: warning
```

---

## Deployment Checklist

### Pre-Deployment

- [ ] Cognito User Pool created
- [ ] User groups configured
- [ ] Test users created
- [ ] Configuration values in SSM Parameter Store
- [ ] Environment variables configured
- [ ] HTTPS certificates in place
- [ ] Monitoring configured
- [ ] Alerts configured

### Post-Deployment Verification

```bash
# 1. Health check
curl https://api.example.com/health

# 2. Public endpoint (no auth)
curl https://api.example.com/api/v1/buildings

# 3. Protected endpoint (should fail)
curl https://api.example.com/api/v1/reservations

# 4. Protected endpoint (with token)
curl -H "Authorization: Bearer $TOKEN" \
  https://api.example.com/api/v1/reservations

# 5. Check metrics
curl https://api.example.com/metrics | grep cognito
```

---

## Future Enhancements

### Phase 2 Features

- **Social Identity Providers**: Google, Facebook, Apple Sign-In
- **SAML Federation**: Enterprise SSO integration
- **Custom Authentication Flow**: Lambda triggers
- **User Attributes**: Extended profile information
- **Email Templates**: Customized invitation/verification emails
- **SMS MFA**: Alternative to TOTP
- **Password Reset Flow**: Self-service password reset
- **Account Linking**: Merge multiple identity providers

### Phase 3 Features

- **API Key Authentication**: For service-to-service communication
- **Fine-Grained Permissions**: Resource-level authorization
- **Audit Logging**: Detailed access logs
- **Session Management**: Active session viewing/termination
- **Rate Limiting Per User**: Prevent abuse
- **Geographic Restrictions**: Location-based access control

---

## References

- [AWS Cognito Documentation](https://docs.aws.amazon.com/cognito/)
- [JWT Best Practices RFC 8725](https://tools.ietf.org/html/rfc8725)
- [OAuth 2.0 RFC 6749](https://tools.ietf.org/html/rfc6749)
- [OpenID Connect Core 1.0](https://openid.net/specs/openid-connect-core-1_0.html)
- [OWASP Authentication Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Authentication_Cheat_Sheet.html)
