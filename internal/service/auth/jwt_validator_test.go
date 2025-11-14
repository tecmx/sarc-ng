package auth

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testKeyID = "test-key-id"
)

// generateTestRSAKey generates an RSA key pair for testing
func generateTestRSAKey() (*rsa.PrivateKey, error) {
	return rsa.GenerateKey(rand.Reader, 2048)
}

// createTestJWKS creates a test JWKS server
func createTestJWKSServer(privateKey *rsa.PrivateKey, kid string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jwks := JWKS{
			Keys: []JWK{
				{
					Kid: kid,
					Kty: "RSA",
					Alg: "RS256",
					Use: "sig",
					N:   base64.RawURLEncoding.EncodeToString(privateKey.PublicKey.N.Bytes()),
					E:   base64.RawURLEncoding.EncodeToString(big.NewInt(int64(privateKey.PublicKey.E)).Bytes()),
				},
			},
		}
		json.NewEncoder(w).Encode(jwks)
	}))
}

// createTestToken creates a test JWT token
func createTestToken(privateKey *rsa.PrivateKey, claims map[string]interface{}, kid string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims(claims))
	token.Header["kid"] = kid
	return token.SignedString(privateKey)
}

// testValidatorSetup holds test validator setup components
type testValidatorSetup struct {
	Validator  *JWTValidator
	PrivateKey *rsa.PrivateKey
	Server     *httptest.Server
}

// createTestValidatorSetup creates a complete test validator setup
func createTestValidatorSetup(t *testing.T) *testValidatorSetup {
	t.Helper()

	privateKey, err := generateTestRSAKey()
	require.NoError(t, err)

	server := createTestJWKSServer(privateKey, testKeyID)

	validator := &JWTValidator{
		region:      "us-east-1",
		userPoolID:  "us-east-1_TestPool",
		clientID:    "test-client-id",
		jwksURL:     server.URL,
		issuer:      "https://cognito-idp.us-east-1.amazonaws.com/us-east-1_TestPool",
		jwksCache:   make(map[string]*rsa.PublicKey),
		cacheExpiry: time.Hour,
		httpClient:  &http.Client{Timeout: 3 * time.Second},
	}

	return &testValidatorSetup{
		Validator:  validator,
		PrivateKey: privateKey,
		Server:     server,
	}
}

// createStandardClaims creates a standard set of test claims
func createStandardClaims() map[string]interface{} {
	now := time.Now().Unix()
	return map[string]interface{}{
		"sub":        "test-user-id",
		"token_use":  "access",
		"client_id":  "test-client-id",
		"iss":        "https://cognito-idp.us-east-1.amazonaws.com/us-east-1_TestPool",
		"auth_time":  now,
		"iat":        now,
		"exp":        now + 3600,
	}
}

func TestNewJWTValidator(t *testing.T) {
	validator := NewJWTValidator("us-east-1", "test-pool-id", "test-client-id", time.Hour)

	assert.NotNil(t, validator)
	assert.Equal(t, "us-east-1", validator.region)
	assert.Equal(t, "test-pool-id", validator.userPoolID)
	assert.Equal(t, "test-client-id", validator.clientID)
	assert.Equal(t, time.Hour, validator.cacheExpiry)
	assert.NotNil(t, validator.jwksCache)
	assert.NotNil(t, validator.httpClient)
}

func TestValidateToken_Success_AccessToken(t *testing.T) {
	ctx := context.Background()

	// Generate test key
	privateKey, err := generateTestRSAKey()
	require.NoError(t, err)

	kid := "test-key-id"

	// Create test JWKS server
	server := createTestJWKSServer(privateKey, kid)
	defer server.Close()

	// Create validator pointing to test server
	validator := &JWTValidator{
		region:      "us-east-1",
		userPoolID:  "us-east-1_TestPool",
		clientID:    "test-client-id",
		jwksURL:     server.URL,
		issuer:      "https://cognito-idp.us-east-1.amazonaws.com/us-east-1_TestPool",
		jwksCache:   make(map[string]*rsa.PublicKey),
		cacheExpiry: time.Hour,
		httpClient:  &http.Client{Timeout: 3 * time.Second},
	}

	// Create test token (access token)
	now := time.Now().Unix()
	claims := map[string]interface{}{
		"sub":        "test-user-id",
		"token_use":  "access",
		"client_id":  "test-client-id",
		"iss":        "https://cognito-idp.us-east-1.amazonaws.com/us-east-1_TestPool",
		"scope":      "openid profile email",
		"auth_time":  now,
		"iat":        now,
		"exp":        now + 3600,
	}

	tokenString, err := createTestToken(privateKey, claims, kid)
	require.NoError(t, err)

	// Validate token
	result, err := validator.ValidateToken(ctx, tokenString)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "test-user-id", result.Sub)
	assert.Equal(t, "access", result.TokenUse)
	assert.Equal(t, "test-client-id", result.ClientID)
}

func TestValidateToken_Success_IDToken(t *testing.T) {
	ctx := context.Background()

	// Generate test key
	privateKey, err := generateTestRSAKey()
	require.NoError(t, err)

	kid := "test-key-id"

	// Create test JWKS server
	server := createTestJWKSServer(privateKey, kid)
	defer server.Close()

	// Create validator
	validator := &JWTValidator{
		region:      "us-east-1",
		userPoolID:  "us-east-1_TestPool",
		clientID:    "test-client-id",
		jwksURL:     server.URL,
		issuer:      "https://cognito-idp.us-east-1.amazonaws.com/us-east-1_TestPool",
		jwksCache:   make(map[string]*rsa.PublicKey),
		cacheExpiry: time.Hour,
		httpClient:  &http.Client{Timeout: 3 * time.Second},
	}

	// Create test token (ID token)
	now := time.Now().Unix()
	claims := map[string]interface{}{
		"sub":            "test-user-id",
		"email":          "test@example.com",
		"email_verified": true,
		"cognito:username": "testuser",
		"cognito:groups":   []string{"admin", "user"},
		"token_use":      "id",
		"aud":            "test-client-id",
		"iss":            "https://cognito-idp.us-east-1.amazonaws.com/us-east-1_TestPool",
		"auth_time":      now,
		"iat":            now,
		"exp":            now + 3600,
	}

	tokenString, err := createTestToken(privateKey, claims, kid)
	require.NoError(t, err)

	// Validate token
	result, err := validator.ValidateToken(ctx, tokenString)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "test-user-id", result.Sub)
	assert.Equal(t, "test@example.com", result.Email)
	assert.True(t, result.EmailVerified)
	assert.Equal(t, "testuser", result.Username)
	assert.Equal(t, []string{"admin", "user"}, result.Groups)
	assert.Equal(t, "id", result.TokenUse)
	assert.Equal(t, "test-client-id", result.Audience)
}

func TestValidateToken_InvalidIssuer(t *testing.T) {
	ctx := context.Background()

	setup := createTestValidatorSetup(t)
	defer setup.Server.Close()

	claims := createStandardClaims()
	claims["iss"] = "https://evil.com/fake-issuer" // Wrong issuer

	tokenString, err := createTestToken(setup.PrivateKey, claims, testKeyID)
	require.NoError(t, err)

	result, err := setup.Validator.ValidateToken(ctx, tokenString)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.ErrorIs(t, err, ErrInvalidIssuer)
}

func TestValidateToken_InvalidTokenUse(t *testing.T) {
	ctx := context.Background()

	setup := createTestValidatorSetup(t)
	defer setup.Server.Close()

	claims := createStandardClaims()
	claims["token_use"] = "refresh" // Invalid token_use

	tokenString, err := createTestToken(setup.PrivateKey, claims, testKeyID)
	require.NoError(t, err)

	result, err := setup.Validator.ValidateToken(ctx, tokenString)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.ErrorIs(t, err, ErrInvalidTokenUse)
}

func TestValidateToken_InvalidAudience_AccessToken(t *testing.T) {
	ctx := context.Background()

	setup := createTestValidatorSetup(t)
	defer setup.Server.Close()

	claims := createStandardClaims()
	claims["client_id"] = "wrong-client-id" // Wrong client_id

	tokenString, err := createTestToken(setup.PrivateKey, claims, testKeyID)
	require.NoError(t, err)

	result, err := setup.Validator.ValidateToken(ctx, tokenString)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.ErrorIs(t, err, ErrInvalidAudience)
}

func TestValidateToken_ExpiredToken(t *testing.T) {
	ctx := context.Background()

	privateKey, err := generateTestRSAKey()
	require.NoError(t, err)

	kid := testKeyID
	server := createTestJWKSServer(privateKey, kid)
	defer server.Close()

	validator := &JWTValidator{
		region:      "us-east-1",
		userPoolID:  "us-east-1_TestPool",
		clientID:    "test-client-id",
		jwksURL:     server.URL,
		issuer:      "https://cognito-idp.us-east-1.amazonaws.com/us-east-1_TestPool",
		jwksCache:   make(map[string]*rsa.PublicKey),
		cacheExpiry: time.Hour,
		httpClient:  &http.Client{Timeout: 3 * time.Second},
	}

	now := time.Now().Unix()
	claims := map[string]interface{}{
		"sub":        "test-user-id",
		"token_use":  "access",
		"client_id":  "test-client-id",
		"iss":        "https://cognito-idp.us-east-1.amazonaws.com/us-east-1_TestPool",
		"auth_time":  now - 7200,
		"iat":        now - 7200,
		"exp":        now - 3600,  // Expired 1 hour ago
	}

	tokenString, err := createTestToken(privateKey, claims, kid)
	require.NoError(t, err)

	result, err := validator.ValidateToken(ctx, tokenString)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.ErrorIs(t, err, ErrExpiredToken)
}

func TestValidateToken_MissingKid(t *testing.T) {
	ctx := context.Background()

	privateKey, err := generateTestRSAKey()
	require.NoError(t, err)

	kid := testKeyID
	server := createTestJWKSServer(privateKey, kid)
	defer server.Close()

	validator := &JWTValidator{
		region:      "us-east-1",
		userPoolID:  "us-east-1_TestPool",
		clientID:    "test-client-id",
		jwksURL:     server.URL,
		issuer:      "https://cognito-idp.us-east-1.amazonaws.com/us-east-1_TestPool",
		jwksCache:   make(map[string]*rsa.PublicKey),
		cacheExpiry: time.Hour,
		httpClient:  &http.Client{Timeout: 3 * time.Second},
	}

	now := time.Now().Unix()
	claims := map[string]interface{}{
		"sub":        "test-user-id",
		"token_use":  "access",
		"client_id":  "test-client-id",
		"iss":        "https://cognito-idp.us-east-1.amazonaws.com/us-east-1_TestPool",
		"auth_time":  now,
		"iat":        now,
		"exp":        now + 3600,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims(claims))
	// Don't set kid in header
	tokenString, err := token.SignedString(privateKey)
	require.NoError(t, err)

	result, err := validator.ValidateToken(ctx, tokenString)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "kid")
}

func TestJWKSCaching(t *testing.T) {
	ctx := context.Background()

	privateKey, err := generateTestRSAKey()
	require.NoError(t, err)

	kid := "test-key-id"

	// Track number of JWKS requests
	requestCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		jwks := JWKS{
			Keys: []JWK{
				{
					Kid: kid,
					Kty: "RSA",
					Alg: "RS256",
					Use: "sig",
					N:   base64.RawURLEncoding.EncodeToString(privateKey.PublicKey.N.Bytes()),
					E:   base64.RawURLEncoding.EncodeToString(big.NewInt(int64(privateKey.PublicKey.E)).Bytes()),
				},
			},
		}
		json.NewEncoder(w).Encode(jwks)
	}))
	defer server.Close()

	validator := &JWTValidator{
		region:      "us-east-1",
		userPoolID:  "us-east-1_TestPool",
		clientID:    "test-client-id",
		jwksURL:     server.URL,
		issuer:      "https://cognito-idp.us-east-1.amazonaws.com/us-east-1_TestPool",
		jwksCache:   make(map[string]*rsa.PublicKey),
		cacheExpiry: time.Hour,
		httpClient:  &http.Client{Timeout: 3 * time.Second},
	}

	now := time.Now().Unix()
	claims := map[string]interface{}{
		"sub":        "test-user-id",
		"token_use":  "access",
		"client_id":  "test-client-id",
		"iss":        "https://cognito-idp.us-east-1.amazonaws.com/us-east-1_TestPool",
		"auth_time":  now,
		"iat":        now,
		"exp":        now + 3600,
	}

	tokenString, err := createTestToken(privateKey, claims, kid)
	require.NoError(t, err)

	// First validation - should fetch JWKS
	_, err = validator.ValidateToken(ctx, tokenString)
	assert.NoError(t, err)
	assert.Equal(t, 1, requestCount)

	// Second validation - should use cache
	_, err = validator.ValidateToken(ctx, tokenString)
	assert.NoError(t, err)
	assert.Equal(t, 1, requestCount, "JWKS should be cached")

	// Third validation - should still use cache
	_, err = validator.ValidateToken(ctx, tokenString)
	assert.NoError(t, err)
	assert.Equal(t, 1, requestCount, "JWKS should still be cached")
}

func TestRefreshJWKS(t *testing.T) {
	ctx := context.Background()

	privateKey, err := generateTestRSAKey()
	require.NoError(t, err)

	kid := testKeyID
	server := createTestJWKSServer(privateKey, kid)
	defer server.Close()

	validator := &JWTValidator{
		region:      "us-east-1",
		userPoolID:  "us-east-1_TestPool",
		clientID:    "test-client-id",
		jwksURL:     server.URL,
		issuer:      "https://cognito-idp.us-east-1.amazonaws.com/us-east-1_TestPool",
		jwksCache:   make(map[string]*rsa.PublicKey),
		cacheExpiry: time.Hour,
		httpClient:  &http.Client{Timeout: 3 * time.Second},
	}

	// Refresh JWKS
	err = validator.RefreshJWKS(ctx)
	assert.NoError(t, err)

	// Verify cache is populated
	assert.NotEmpty(t, validator.jwksCache)
	assert.Contains(t, validator.jwksCache, kid)
}

func TestLoadJWKSFromJSON(t *testing.T) {
	privateKey, err := generateTestRSAKey()
	require.NoError(t, err)

	kid := "test-key-id"

	jwksJSON := fmt.Sprintf(`{
		"keys": [
			{
				"kid": "%s",
				"kty": "RSA",
				"alg": "RS256",
				"use": "sig",
				"n": "%s",
				"e": "%s"
			}
		]
	}`, kid,
		base64.RawURLEncoding.EncodeToString(privateKey.PublicKey.N.Bytes()),
		base64.RawURLEncoding.EncodeToString(big.NewInt(int64(privateKey.PublicKey.E)).Bytes()))

	validator := &JWTValidator{
		jwksCache:   make(map[string]*rsa.PublicKey),
		cacheExpiry: time.Hour,
	}

	validator.loadJWKSFromJSON([]byte(jwksJSON))

	// Verify cache is populated
	assert.NotEmpty(t, validator.jwksCache)
	assert.Contains(t, validator.jwksCache, kid)
}

