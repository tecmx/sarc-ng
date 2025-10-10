# AWS Cognito Authentication Implementation Proposal

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

**File: `pkg/rest/middleware/auth.go`**

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
    buildingRest "sarc-ng/internal/transport/rest/building"
    classRest "sarc-ng/internal/transport/rest/class"
    lessonRest "sarc-ng/internal/transport/rest/lesson"
    reservationRest "sarc-ng/internal/transport/rest/reservation"
    resourceRest "sarc-ng/internal/transport/rest/resource"
    "sarc-ng/pkg/rest/middleware"

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
  -d '{"resource_id": 1, "start_time": "2025-10-15T10:00:00Z"}'
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
