package middleware

import (
	"context"
	"net/http"
	"strings"

	"sarc-ng/internal/domain/auth"

	"github.com/gin-gonic/gin"
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
				"error":           "Insufficient permissions",
				"code":            "INSUFFICIENT_PERMISSIONS",
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
					"error":           "Insufficient permissions",
					"code":            "INSUFFICIENT_PERMISSIONS",
					"required_groups": groups,
				})
				return
			}
		}

		c.Next()
	}
}

// RequireAdmin middleware ensures user has admin privileges
func RequireAdmin() gin.HandlerFunc {
	return RequireGroups("admin")
}

// RequireManager middleware ensures user has manager or admin privileges
func RequireManager() gin.HandlerFunc {
	return RequireGroups("admin", "manager")
}

// RequireTeacher middleware ensures user has teacher, manager, or admin privileges
func RequireTeacher() gin.HandlerFunc {
	return RequireGroups("admin", "manager", "teacher")
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

// GetUserIDFromContext retrieves the authenticated user ID from context
func GetUserIDFromContext(c *gin.Context) (string, bool) {
	user, exists := GetUserFromContext(c)
	if !exists {
		return "", false
	}
	return user.ID, true
}

// GetUserEmailFromContext retrieves the authenticated user email from context
func GetUserEmailFromContext(c *gin.Context) (string, bool) {
	user, exists := GetUserFromContext(c)
	if !exists {
		return "", false
	}
	return user.Email, true
}

// WithUser adds a user to the context (useful for testing)
func WithUser(ctx context.Context, user *auth.User) context.Context {
	return context.WithValue(ctx, ContextKeyUser, user)
}

// WithClaims adds claims to the context (useful for testing)
func WithClaims(ctx context.Context, claims *auth.Claims) context.Context {
	return context.WithValue(ctx, ContextKeyClaims, claims)
}
