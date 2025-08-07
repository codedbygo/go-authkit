package authkit

import (
	"errors"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// AuthKit is the main struct that holds configuration and methods
type AuthKit struct {
	config Config
	users  map[string]*User // In-memory storage for demo (use database in production)
	mutex  sync.RWMutex     // For thread-safe operations
}

// Config holds the configuration for AuthKit
type Config struct {
	JWTSecret     string
	TokenExpiry   string // e.g., "24h", "1h", "30m"
	RefreshExpiry string // e.g., "7d", "30d"
	BCryptCost    int    // bcrypt cost (default: 12)
	RateLimitRPM  int    // Rate limit per minute
	EmailRequired bool   // Require email verification
}

// User represents a user in the system
type User struct {
	ID            string                 `json:"id"`
	Email         string                 `json:"email"`
	Password      string                 `json:"password,omitempty"` // Hashed password
	Name          string                 `json:"name"`
	Role          string                 `json:"role"`
	Permissions   []string               `json:"permissions"`
	EmailVerified bool                   `json:"email_verified"`
	CreatedAt     time.Time              `json:"created_at"`
	UpdatedAt     time.Time              `json:"updated_at"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
}

// Claims represents JWT claims
type Claims struct {
	UserID      string                 `json:"user_id"`
	Email       string                 `json:"email"`
	Role        string                 `json:"role"`
	Permissions []string               `json:"permissions"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	jwt.RegisteredClaims
}

// TokenResponse represents the response after successful login
type TokenResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	TokenType    string    `json:"token_type"`
	ExpiresIn    int64     `json:"expires_in"`
	User         *UserInfo `json:"user"`
}

// UserInfo represents safe user information (without password)
type UserInfo struct {
	ID            string                 `json:"id"`
	Email         string                 `json:"email"`
	Name          string                 `json:"name"`
	Role          string                 `json:"role"`
	Permissions   []string               `json:"permissions"`
	EmailVerified bool                   `json:"email_verified"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
}

// LoginRequest represents login request payload
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// RegisterRequest represents registration request payload
type RegisterRequest struct {
	Email    string                 `json:"email" binding:"required,email"`
	Password string                 `json:"password" binding:"required,min=8"`
	Name     string                 `json:"name" binding:"required"`
	Role     string                 `json:"role,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// RefreshRequest represents refresh token request
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// Common errors
var (
	ErrUserNotFound      = errors.New("user not found")
	ErrInvalidPassword   = errors.New("invalid password")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrInvalidToken      = errors.New("invalid token")
	ErrTokenExpired      = errors.New("token expired")
	ErrUnauthorized      = errors.New("unauthorized")
	ErrInsufficientRole  = errors.New("insufficient role permissions")
)
