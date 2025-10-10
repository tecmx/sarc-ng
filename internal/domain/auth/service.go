package auth

import "context"

// TokenValidator defines the interface for JWT token validation
type TokenValidator interface {
	// ValidateToken validates a JWT token and returns claims
	ValidateToken(ctx context.Context, token string) (*Claims, error)

	// RefreshJWKS forces a refresh of the JWKS cache
	RefreshJWKS(ctx context.Context) error
}
