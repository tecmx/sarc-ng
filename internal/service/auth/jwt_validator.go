package auth

import (
	"context"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"
	"sync"
	"time"

	"sarc-ng/internal/domain/auth"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken     = errors.New("invalid token")
	ErrExpiredToken     = errors.New("token has expired")
	ErrInvalidSignature = errors.New("invalid token signature")
	ErrInvalidIssuer    = errors.New("invalid token issuer")
	ErrInvalidAudience  = errors.New("invalid token audience")
	ErrInvalidTokenUse  = errors.New("invalid token use")
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

	validator := &JWTValidator{
		region:      region,
		userPoolID:  userPoolID,
		clientID:    clientID,
		jwksURL:     jwksURL,
		issuer:      issuer,
		jwksCache:   make(map[string]*rsa.PublicKey),
		cacheExpiry: cacheExpiry,
		httpClient: &http.Client{
			Timeout: 3 * time.Second, // Reduced timeout
		},
	}

	// Pre-load JWKS from environment variable if available (for VPC Lambda without NAT Gateway)
	if jwksJSON := os.Getenv("COGNITO_JWKS"); jwksJSON != "" {
		validator.loadJWKSFromJSON([]byte(jwksJSON))
	}

	return validator
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

// loadJWKSFromJSON loads JWKS from JSON data (for pre-fetched JWKS from env var)
func (v *JWTValidator) loadJWKSFromJSON(data []byte) {
	var jwks JWKS
	if err := json.Unmarshal(data, &jwks); err != nil {
		log.Printf("Failed to unmarshal JWKS from environment variable: %v", err)
		return
	}

	v.cacheMutex.Lock()
	defer v.cacheMutex.Unlock()

	for _, key := range jwks.Keys {
		publicKey, err := v.parseRSAPublicKey(key)
		if err != nil {
			log.Printf("Failed to convert JWK to public key for kid %s: %v", key.Kid, err)
			continue
		}
		v.jwksCache[key.Kid] = publicKey
	}

	v.lastCacheUpdate = time.Now()
	log.Printf("Successfully loaded %d keys from COGNITO_JWKS environment variable", len(v.jwksCache))
}
