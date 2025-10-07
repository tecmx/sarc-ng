# AWS Secrets Manager with SAM

Guide for managing database credentials with AWS Secrets Manager.

## Benefits

- No hardcoded credentials
- Encryption at rest (KMS) and in transit (TLS)
- Centralized management - update without redeploying
- Audit trail via CloudTrail
- Optional automatic rotation

## Costs

- $0.40/secret/month
- $0.05/10,000 API calls
- Example: ~$0.65/month for typical usage

## Option 1: Simple Secret (Recommended)

```yaml
Resources:
  DatabaseSecret:
    Type: AWS::SecretsManager::Secret
    Properties:
      Name: !Sub "${Environment}-sarc-ng-db-credentials"
      SecretString: !Sub |
        {
          "username": "root",
          "password": "${DBPassword}",
          "host": "${SarcDatabase.Endpoint.Address}",
          "port": "${SarcDatabase.Endpoint.Port}",
          "database": "${DBName}"
        }

  SarcNgFunction:
    Type: AWS::Serverless::Function
    Properties:
      Environment:
        Variables:
          DB_SECRET_ARN: !Ref DatabaseSecret
      Policies:
        - Statement:
            - Effect: Allow
              Action: secretsmanager:GetSecretValue
              Resource: !Ref DatabaseSecret
```

## Option 2: Auto-Generated Password

```yaml
Resources:
  DatabaseSecret:
    Type: AWS::SecretsManager::Secret
    Properties:
      GenerateSecretString:
        SecretStringTemplate: '{"username": "root"}'
        GenerateStringKey: "password"
        PasswordLength: 32
        ExcludeCharacters: '"@/\\'

  SecretRDSAttachment:
    Type: AWS::SecretsManager::SecretTargetAttachment
    Properties:
      SecretId: !Ref DatabaseSecret
      TargetId: !Ref SarcDatabase
      TargetType: AWS::RDS::DBInstance
```

## Go Code Integration

### 1. Add Dependencies

```bash
go get github.com/aws/aws-sdk-go-v2/config
go get github.com/aws/aws-sdk-go-v2/service/secretsmanager
```

### 2. Create Secrets Helper

```go
// internal/adapter/secrets/manager.go
package secrets

import (
    "context"
    "encoding/json"
    "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

type DBCredentials struct {
    Username string `json:"username"`
    Password string `json:"password"`
    Host     string `json:"host"`
    Port     string `json:"port"`
    Database string `json:"database"`
}

func GetDBCredentials(ctx context.Context, secretARN string) (*DBCredentials, error) {
    cfg, err := config.LoadDefaultConfig(ctx)
    if err != nil {
        return nil, err
    }

    client := secretsmanager.NewFromConfig(cfg)
    result, err := client.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{
        SecretId: &secretARN,
    })
    if err != nil {
        return nil, err
    }

    var creds DBCredentials
    err = json.Unmarshal([]byte(*result.SecretString), &creds)
    return &creds, err
}
```

### 3. Use in Application

```go
// cmd/lambda/main.go
package main

import (
    "context"
    "os"
    "your-app/internal/adapter/secrets"
)

func main() {
    ctx := context.Background()
    secretARN := os.Getenv("DB_SECRET_ARN")

    creds, err := secrets.GetDBCredentials(ctx, secretARN)
    if err != nil {
        panic(err)
    }

    // Use credentials
    dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
        creds.Username, creds.Password, creds.Host, creds.Port, creds.Database)
}
```

## Caching Strategy

```go
type SecretCache struct {
    secret    *DBCredentials
    expiresAt time.Time
    mu        sync.RWMutex
}

func (c *SecretCache) Get(ctx context.Context, arn string) (*DBCredentials, error) {
    c.mu.RLock()
    if c.secret != nil && time.Now().Before(c.expiresAt) {
        c.mu.RUnlock()
        return c.secret, nil
    }
    c.mu.RUnlock()

    c.mu.Lock()
    defer c.mu.Unlock()

    creds, err := GetDBCredentials(ctx, arn)
    if err != nil {
        return nil, err
    }

    c.secret = creds
    c.expiresAt = time.Now().Add(5 * time.Minute)
    return creds, nil
}
```

## Testing

### Local Testing

Use env.json for local development:

```json
{
  "SarcNgFunction": {
    "DB_HOST": "db",
    "DB_USER": "root",
    "DB_PASSWORD": "example"
  }
}
```

### Fallback Logic

```go
func initDB(ctx context.Context) (*gorm.DB, error) {
    var dsn string

    if secretARN := os.Getenv("DB_SECRET_ARN"); secretARN != "" {
        // Production: Use Secrets Manager
        creds, err := secrets.GetDBCredentials(ctx, secretARN)
        if err != nil {
            return nil, err
        }
        dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
            creds.Username, creds.Password, creds.Host, creds.Port, creds.Database)
    } else {
        // Development: Use environment variables
        dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
            os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"),
            os.Getenv("DB_HOST"), os.Getenv("DB_PORT"),
            os.Getenv("DB_NAME"))
    }

    return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}
```

## Deployment

```bash
# Deploy with Secrets Manager
cd infrastructure/sam
sam build
sam deploy --parameter-overrides DBPassword=<secure-password>
```

## Best Practices

1. **Cache secrets** - Don't fetch on every request (5-minute TTL recommended)
2. **Fallback to env vars** - For local development
3. **Never log secrets** - Mask in logs
4. **Use IAM policies** - Restrict secret access per function
5. **Monitor access** - Enable CloudTrail logging
6. **Rotate regularly** - Use automatic rotation for production

 "github.com/aws/aws-sdk-go-v2/config"
 "github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

// DatabaseCredentials represents database secret structure
type DatabaseCredentials struct {
 Username string `json:"username"`
 Password string `json:"password"`
 Host     string `json:"host"`
 Port     string `json:"port"`
 Database string `json:"database"`
}

var (
 cachedCreds *DatabaseCredentials
 cacheMutex  sync.RWMutex
)

// GetDatabaseCredentials retrieves and caches database credentials
func GetDatabaseCredentials(ctx context.Context) (*DatabaseCredentials, error) {
 // Check cache first
 cacheMutex.RLock()
 if cachedCreds != nil {
  cacheMutex.RUnlock()
  return cachedCreds, nil
 }
 cacheMutex.RUnlock()

 // Get secret ARN from environment
 secretArn := os.Getenv("DB_SECRET_ARN")
 if secretArn == "" {
  return nil, fmt.Errorf("DB_SECRET_ARN environment variable not set")
 }

 // Load AWS config
 cfg, err := config.LoadDefaultConfig(ctx)
 if err != nil {
  return nil, fmt.Errorf("failed to load AWS config: %w", err)
 }

 // Create Secrets Manager client
 client := secretsmanager.NewFromConfig(cfg)

 // Retrieve secret value
 result, err := client.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{
  SecretId: &secretArn,
 })
 if err != nil {
  return nil, fmt.Errorf("failed to retrieve secret: %w", err)
 }

 // Parse secret JSON
 var creds DatabaseCredentials
 if err := json.Unmarshal([]byte(*result.SecretString), &creds); err != nil {
  return nil, fmt.Errorf("failed to parse secret: %w", err)
 }

 // Cache credentials
 cacheMutex.Lock()
 cachedCreds = &creds
 cacheMutex.Unlock()

 return &creds, nil
}

// BuildDSN builds database connection string from credentials
func BuildDSN(creds *DatabaseCredentials) string {
 return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
  creds.Username,
  creds.Password,
  creds.Host,
  creds.Port,
  creds.Database,
 )
}
```

### 3. Update Database Connection Logic

Modify your database connection to use Secrets Manager:

```go
func provideDatabaseConnection(config *config.Config) (*gorm.DB, error) {
 ctx := context.Background()

 // Check if using Secrets Manager (Lambda environment)
 secretArn := os.Getenv("DB_SECRET_ARN")

 var dsn string
 if secretArn != "" {
  // Get credentials from Secrets Manager
  creds, err := secrets.GetDatabaseCredentials(ctx)
  if err != nil {
   return nil, fmt.Errorf("failed to get database credentials: %w", err)
  }
  dsn = secrets.BuildDSN(creds)
 } else {
  // Use config file (local development)
  dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
   config.Database.User,
   config.Database.Password,
   config.Database.Host,
   config.Database.Port,
   config.Database.Name,
  )
 }

 // Open database connection
 db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
 if err != nil {
  return nil, fmt.Errorf("failed to connect to database: %w", err)
 }

 return db, nil
}
```

---

## 🔐 IAM Permissions Required

The Lambda execution role needs:

```yaml
Policies:
  - Statement:
      - Effect: Allow
        Action:
          - secretsmanager:GetSecretValue
          - secretsmanager:DescribeSecret
        Resource: !Ref DatabaseSecret
```

**For LabRole:** LabRole should already have these permissions, but verify with:

```bash
aws iam get-role-policy --role-name LabRole --policy-name LabRolePolicy
```

---

## 📝 Best Practices

### 1. **Cache Secrets in Lambda**

- ✅ Retrieve once per Lambda container initialization
- ✅ Reduces API calls and costs
- ✅ Improves performance
- ❌ Don't retrieve on every invocation

### 2. **Use Separate Secrets per Environment**

```
dev-sarc-ng-db-credentials
staging-sarc-ng-db-credentials
prod-sarc-ng-db-credentials
```

### 3. **Store Complete Connection Info**

```json
{
  "username": "root",
  "password": "secure_password",
  "host": "db.example.com",
  "port": "3306",
  "database": "sarcng"
}
```

### 4. **Tag Your Secrets**

```yaml
Tags:
  - Key: Environment
    Value: !Ref Environment
  - Key: Application
    Value: sarc-ng
  - Key: ManagedBy
    Value: SAM
```

### 5. **Enable CloudTrail Logging**

Monitor who accesses secrets:

```yaml
DatabaseSecret:
  Type: AWS::SecretsManager::Secret
  Properties:
    KmsKeyId: !Ref MyKMSKey  # Custom KMS key for audit
```

---

## 🧪 Testing

### Local Testing (without Secrets Manager)

```bash
# Use local config files
make docker-up
```

### Lambda Testing (with Secrets Manager)

```bash
# Create test secret
aws secretsmanager create-secret \
  --name dev-sarc-ng-db-credentials \
  --secret-string '{"username":"root","password":"example","host":"localhost","port":"3306","database":"sarcng"}'

# Test Lambda locally
sam local invoke -e events/test.json
```

---

## 🚦 Migration Steps

### Step 1: Add Secrets Manager to SAM Template

Update `infrastructure/sam/template.yaml` with Option 1 code above.

### Step 2: Add Go Dependencies

```bash
cd /home/miguel/devel/devel-tecmx/sarc-ng
go get github.com/aws/aws-sdk-go-v2/config
go get github.com/aws/aws-sdk-go-v2/service/secretsmanager
go mod tidy
```

### Step 3: Create Secret Helper

Create `internal/adapter/secrets/manager.go` with code above.

### Step 4: Update Database Connection

Modify `cmd/lambda/wire.go` and `cmd/server/wire.go`.

### Step 5: Deploy

```bash
cd infrastructure/sam
sam build
sam deploy
```

### Step 6: Verify

```bash
# Check secret exists
aws secretsmanager describe-secret --secret-id prod-sarc-ng-db-credentials

# Test Lambda
curl https://your-api-gateway-url/api/v1/buildings
```

---

## 📚 References

- [AWS Secrets Manager Documentation](https://docs.aws.amazon.com/secretsmanager/)
- [SAM Secrets Manager Integration](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/sam-resource-secret.html)
- [AWS SDK for Go v2](https://aws.github.io/aws-sdk-go-v2/docs/)
- [Secrets Manager Pricing](https://aws.amazon.com/secrets-manager/pricing/)

---

## ❓ Questions?

**Q: Can I use Secrets Manager with AWS Academy LabRole?**
A: Yes! LabRole typically has `secretsmanager:GetSecretValue` permissions.

**Q: What about costs?**
A: ~$0.65/month for basic usage. Minimal compared to security benefits.

**Q: Do I need to update secrets manually?**
A: Option 1: Manual updates via AWS Console/CLI
   Option 2: Auto-rotation (if LabRole supports it)

**Q: Can I test locally without Secrets Manager?**
A: Yes! The code checks for `DB_SECRET_ARN`. If not set, it uses config files.

**Q: Should I remove the DBPassword parameter?**
A: Keep it initially for backward compatibility. Remove after migration is complete.
