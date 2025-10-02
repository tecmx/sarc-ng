# AWS Secrets Manager Integration with SAM

Guide for integrating AWS Secrets Manager with your SARC-NG Lambda function.

## ✅ Yes, It's Fully Supported

AWS Secrets Manager is natively supported in SAM templates and is the **recommended approach** for managing sensitive data like database credentials.

---

## 🎯 Benefits

### Security

- ✅ **No hardcoded credentials** in code or environment variables
- ✅ **Encryption at rest** (AWS KMS)
- ✅ **Encryption in transit** (TLS)
- ✅ **Audit trail** (CloudTrail logs all secret access)
- ✅ **Automatic rotation** (optional)

### Operations

- ✅ **Centralized management** - Update secrets without redeploying Lambda
- ✅ **Version control** - Track secret changes
- ✅ **Fine-grained IAM permissions** - Control who can access secrets
- ✅ **Multi-region replication** (optional)

---

## 📊 Costs

**AWS Secrets Manager Pricing:**

- **$0.40 per secret per month**
- **$0.05 per 10,000 API calls**

**Example for SARC-NG:**

- 1 secret (database credentials): **$0.40/month**
- ~50,000 Lambda invocations/month: **$0.25/month** (secret access)
- **Total: ~$0.65/month**

**Note:** Free tier includes 30-day trial for secrets rotation.

---

## 🚀 Implementation Options

### Option 1: Simple Secret (Recommended for Your Case)

**Best for:** Storing database credentials that rarely change.

```yaml
Resources:
  # Create the secret
  DatabaseSecret:
    Type: AWS::SecretsManager::Secret
    Properties:
      Name: !Sub "${Environment}-sarc-ng-db-credentials"
      Description: Database credentials for SARC-NG
      SecretString: !Sub |
        {
          "username": "root",
          "password": "${DBPassword}",
          "host": "${SarcDatabase.Endpoint.Address}",
          "port": "${SarcDatabase.Endpoint.Port}",
          "database": "${DBName}"
        }
      Tags:
        - Key: Environment
          Value: !Ref Environment

  # Grant Lambda permission to read the secret
  SarcNgFunction:
    Type: AWS::Serverless::Function
    Properties:
      # ... existing properties ...
      Environment:
        Variables:
          DB_SECRET_ARN: !Ref DatabaseSecret
          ENVIRONMENT: !Ref Environment
      Policies:
        - Statement:
            - Effect: Allow
              Action:
                - secretsmanager:GetSecretValue
              Resource: !Ref DatabaseSecret
```

---

### Option 2: Auto-Generated Password with RDS

**Best for:** Maximum security with automatic password generation.

```yaml
Resources:
  # Generate random password
  DatabaseSecret:
    Type: AWS::SecretsManager::Secret
    Properties:
      Name: !Sub "${Environment}-sarc-ng-db-credentials"
      GenerateSecretString:
        SecretStringTemplate: '{"username": "root"}'
        GenerateStringKey: "password"
        PasswordLength: 32
        ExcludeCharacters: '"@/\\'
        RequireEachIncludedType: true

  # Attach secret to RDS (CloudFormation manages RDS password)
  SecretRDSAttachment:
    Type: AWS::SecretsManager::SecretTargetAttachment
    Properties:
      SecretId: !Ref DatabaseSecret
      TargetId: !Ref SarcDatabase
      TargetType: AWS::RDS::DBInstance

  # Update RDS to use secret
  SarcDatabase:
    Type: AWS::RDS::DBInstance
    Properties:
      # Remove MasterUserPassword parameter
      MasterUsername: !GetAtt DatabaseSecret.SecretString:username
      MasterUserPassword: !Sub '{{resolve:secretsmanager:${DatabaseSecret}:SecretString:password}}'
      # ... other properties ...
```

---

### Option 3: Automatic Secret Rotation

**Best for:** Production environments requiring periodic password rotation.

```yaml
Resources:
  DatabaseSecret:
    Type: AWS::SecretsManager::Secret
    Properties:
      Name: !Sub "${Environment}-sarc-ng-db-credentials"
      GenerateSecretString:
        SecretStringTemplate: '{"username": "root"}'
        GenerateStringKey: "password"
        PasswordLength: 32
        ExcludeCharacters: '"@/\\'

  # Attach to RDS
  SecretRDSAttachment:
    Type: AWS::SecretsManager::SecretTargetAttachment
    Properties:
      SecretId: !Ref DatabaseSecret
      TargetId: !Ref SarcDatabase
      TargetType: AWS::RDS::DBInstance

  # Rotate every 30 days
  SecretRotationSchedule:
    Type: AWS::SecretsManager::RotationSchedule
    Properties:
      SecretId: !Ref DatabaseSecret
      HostedRotationLambda:
        RotationType: MySQLSingleUser
        RotationLambdaName: !Sub "${Environment}-sarc-ng-db-rotation"
      RotationRules:
        AutomaticallyAfterDays: 30
```

**Note:** LabRole may not support automatic rotation. Use Option 1 or 2 for AWS Academy.

---

## 💻 Go Code Integration

### 1. Add AWS SDK Dependencies

```bash
go get github.com/aws/aws-sdk-go-v2/config
go get github.com/aws/aws-sdk-go-v2/service/secretsmanager
```

Update `go.mod`:

```go
require (
    github.com/aws/aws-sdk-go-v2/config v1.27.0
    github.com/aws/aws-sdk-go-v2/service/secretsmanager v1.28.0
)
```

### 2. Create Secret Retrieval Helper

Create `internal/adapter/secrets/manager.go`:

```go
package secrets

import (
 "context"
 "encoding/json"
 "fmt"
 "os"
 "sync"

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
