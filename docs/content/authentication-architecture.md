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
    "start_time": "2025-10-15T10:00:00Z",
    "end_time": "2025-10-15T11:00:00Z",
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
  "start_time": "2025-10-15T10:00:00Z",
  "end_time": "2025-10-15T11:00:00Z",
  "purpose": "Team meeting",
  "status": "confirmed",
  "created_at": "2025-10-10T14:23:45Z"
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
