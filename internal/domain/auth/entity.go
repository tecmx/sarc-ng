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

// IsAdmin checks if user has admin privileges
func (u *User) IsAdmin() bool {
	return u.HasGroup("admin")
}

// IsManager checks if user has manager privileges
func (u *User) IsManager() bool {
	return u.HasGroup("manager") || u.IsAdmin()
}

// IsTeacher checks if user has teacher privileges
func (u *User) IsTeacher() bool {
	return u.HasGroup("teacher") || u.IsManager()
}
