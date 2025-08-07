# AuthKit Usage Guide

This guide demonstrates how to use the AuthKit library in your Go applications.

## Installation

```bash
go get github.com/codedbygo/go-authkit
```

## Quick Start

### 1. Initialize AuthKit

```go
package main

import "github.com/codedbygo/go-authkit"

func main() {
    // Initialize AuthKit with configuration
    auth := authkit.New(authkit.Config{
        JWTSecret:     "your-super-secret-jwt-key-here",
        TokenExpiry:   "24h",        // Access token expiry
        RefreshExpiry: "7d",         // Refresh token expiry
        BCryptCost:    12,           // Password hashing cost
        RateLimitRPM:  100,          // Rate limit per minute
        EmailRequired: false,        // Email verification required
    })
}
```

### 2. User Registration

```go
// Register a new user
registerReq := authkit.RegisterRequest{
    Email:    "user@example.com",
    Password: "securepassword123",
    Name:     "John Doe",
    Role:     "user",
    Metadata: map[string]interface{}{
        "department": "engineering",
        "level":      "junior",
    },
}

user, err := auth.RegisterUser(registerReq)
if err != nil {
    // Handle registration error
    if err == authkit.ErrUserAlreadyExists {
        log.Println("User already exists")
    }
    return
}

log.Printf("User registered: %+v", user)
```

### 3. User Login

```go
// Login user and get tokens
tokenResponse, err := auth.LoginUser("user@example.com", "securepassword123")
if err != nil {
    // Handle login error
    if err == authkit.ErrUserNotFound {
        log.Println("User not found")
    } else if err == authkit.ErrInvalidPassword {
        log.Println("Invalid password")
    }
    return
}

log.Printf("Login successful!")
log.Printf("Access Token: %s", tokenResponse.AccessToken)
log.Printf("Refresh Token: %s", tokenResponse.RefreshToken)
log.Printf("Expires In: %d seconds", tokenResponse.ExpiresIn)
log.Printf("User: %+v", tokenResponse.User)
```

### 4. Token Validation

```go
// Validate a JWT token
claims, err := auth.ValidateToken(tokenString)
if err != nil {
    // Handle validation error
    if err == authkit.ErrTokenExpired {
        log.Println("Token expired")
    } else if err == authkit.ErrInvalidToken {
        log.Println("Invalid token")
    }
    return
}

log.Printf("Token is valid!")
log.Printf("User ID: %s", claims.UserID)
log.Printf("Email: %s", claims.Email)
log.Printf("Role: %s", claims.Role)
log.Printf("Permissions: %+v", claims.Permissions)
```

### 5. Token Refresh

```go
// Refresh access token using refresh token
newTokens, err := auth.RefreshToken(refreshTokenString)
if err != nil {
    // Handle refresh error
    log.Printf("Token refresh failed: %v", err)
    return
}

log.Printf("Token refreshed!")
log.Printf("New Access Token: %s", newTokens.AccessToken)
log.Printf("New Refresh Token: %s", newTokens.RefreshToken)
```

## Web Framework Integration

### Gin Framework

```go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/your-username/go-authkit"
)

func main() {
    auth := authkit.New(authkit.Config{
        JWTSecret: "your-secret-key",
        TokenExpiry: "24h",
    })

    r := gin.Default()

    // Public routes
    r.POST("/register", auth.RegisterHandler)
    r.POST("/login", auth.LoginHandler)
    r.POST("/refresh", auth.RefreshHandler)

    // Protected routes
    protected := r.Group("/api/v1")
    protected.Use(auth.GinMiddleware())
    {
        protected.GET("/profile", auth.ProfileHandler)
        protected.PUT("/profile", auth.UpdateProfileHandler)
        protected.GET("/posts", getPostsHandler)
    }

    // Admin only routes
    admin := protected.Group("/admin")
    admin.Use(auth.RequireRole("admin"))
    {
        admin.GET("/users", listUsersHandler)
        admin.DELETE("/users/:id", deleteUserHandler)
    }

    r.Run(":8080")
}
```

### Fiber Framework

```go
package main

import (
    "github.com/gofiber/fiber/v2"
    "github.com/your-username/go-authkit"
)

func main() {
    auth := authkit.New(authkit.Config{
        JWTSecret: "your-secret-key",
        TokenExpiry: "24h",
    })

    app := fiber.New()

    // Public routes
    app.Post("/register", auth.RegisterHandlerFiber)
    app.Post("/login", auth.LoginHandlerFiber)
    app.Post("/refresh", auth.RefreshHandlerFiber)

    // Protected routes
    protected := app.Group("/api/v1")
    protected.Use(auth.FiberMiddleware())
    protected.Get("/profile", auth.ProfileHandlerFiber)
    protected.Put("/profile", auth.UpdateProfileHandlerFiber)

    // Admin routes
    admin := protected.Group("/admin")
    admin.Use(auth.RequireRoleFiber("admin"))
    admin.Get("/users", listUsersHandlerFiber)

    app.Listen(":8080")
}
```

## Advanced Features

### Role-Based Access Control

```go
// Require specific role
router.Use(auth.RequireRole("admin"))

// Require one of multiple roles
router.Use(auth.RequireRoles([]string{"admin", "moderator"}))

// Require specific permission
router.Use(auth.RequirePermission("posts:write"))
```

### Custom Claims

```go
// Generate token with custom claims
customClaims := map[string]interface{}{
    "department": "engineering",
    "clearance":  "level-3",
    "scopes":     []string{"api:read", "api:write"},
}

token, err := auth.GenerateCustomToken(userID, customClaims, time.Hour*2)
```

### User Management

```go
// Update user information
updates := map[string]interface{}{
    "name": "John Smith",
    "role": "senior",
    "permissions": []string{"read", "write", "admin"},
}
updatedUser, err := auth.UpdateUser(userID, updates)

// Get user by ID
user, err := auth.GetUserByID(userID)

// Get user by email
user, err := auth.GetUserByEmail("user@example.com")

// List all users
users := auth.ListUsers()

// Delete user
err := auth.DeleteUser(userID)
```

### Password Utilities

```go
// Hash password
hashedPassword, err := auth.HashPassword("plainPassword")

// Compare password
isValid := auth.ComparePassword(hashedPassword, "plainPassword")

// Static methods (without AuthKit instance)
hashedPassword, err := authkit.HashPasswordStatic("plainPassword", 12)
isValid := authkit.ComparePasswordStatic(hashedPassword, "plainPassword")
```

## Error Handling

AuthKit provides specific error types for better error handling:

```go
import "errors"

tokenResponse, err := auth.LoginUser(email, password)
if err != nil {
    switch {
    case errors.Is(err, authkit.ErrUserNotFound):
        // User doesn't exist
    case errors.Is(err, authkit.ErrInvalidPassword):
        // Wrong password
    case errors.Is(err, authkit.ErrUserAlreadyExists):
        // User already registered
    case errors.Is(err, authkit.ErrTokenExpired):
        // Token has expired
    case errors.Is(err, authkit.ErrInvalidToken):
        // Token is invalid
    case errors.Is(err, authkit.ErrUnauthorized):
        // User not authorized
    case errors.Is(err, authkit.ErrInsufficientRole):
        // User doesn't have required role
    default:
        // Other error
    }
}
```

## Context Helpers

Extract user information from request context:

### Gin Context

```go
func protectedHandler(c *gin.Context) {
    // Get user claims from context
    claims, exists := authkit.GetUserFromGinContext(c)
    if !exists {
        c.JSON(401, gin.H{"error": "User not found in context"})
        return
    }

    c.JSON(200, gin.H{
        "user_id":     claims.UserID,
        "email":       claims.Email,
        "role":        claims.Role,
        "permissions": claims.Permissions,
    })
}
```

### Fiber Context

```go
func protectedHandler(c *fiber.Ctx) error {
    // Get user claims from context
    claims, exists := authkit.GetUserFromFiberContext(c)
    if !exists {
        return c.Status(401).JSON(fiber.Map{
            "error": "User not found in context",
        })
    }

    return c.JSON(fiber.Map{
        "user_id":     claims.UserID,
        "email":       claims.Email,
        "role":        claims.Role,
        "permissions": claims.Permissions,
    })
}
```

## Best Practices

### 1. Secure JWT Secret

```go
// Use environment variables for secrets
import "os"

auth := authkit.New(authkit.Config{
    JWTSecret: os.Getenv("JWT_SECRET"), // Set via environment
    TokenExpiry: "24h",
})
```

### 2. Handle Errors Properly

```go
tokenResponse, err := auth.LoginUser(email, password)
if err != nil {
    // Log the error for debugging
    log.Printf("Login failed for %s: %v", email, err)
    
    // Return generic error to client for security
    c.JSON(401, gin.H{"error": "Invalid credentials"})
    return
}
```

### 3. Rate Limiting

```go
// Configure rate limiting
auth := authkit.New(authkit.Config{
    JWTSecret:    "your-secret",
    RateLimitRPM: 60, // 60 requests per minute
})
```

### 4. Token Expiry

```go
// Use appropriate token expiry times
auth := authkit.New(authkit.Config{
    JWTSecret:     "your-secret",
    TokenExpiry:   "15m",  // Short-lived access tokens
    RefreshExpiry: "7d",   // Longer refresh tokens
})
```

### 5. Database Integration

For production use, replace the in-memory storage:

```go
// Implement your own user storage
type DatabaseUserStore struct {
    db *sql.DB
}

func (d *DatabaseUserStore) CreateUser(user *User) error {
    // Implement database user creation
}

func (d *DatabaseUserStore) GetUserByEmail(email string) (*User, error) {
    // Implement database user retrieval
}

// Use dependency injection to provide your storage implementation
```

## Testing

AuthKit includes comprehensive tests. Run them with:

```bash
go test -v ./...
```

## Configuration Options

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `JWTSecret` | `string` | **required** | Secret key for signing JWT tokens |
| `TokenExpiry` | `string` | `"24h"` | Access token expiry duration |
| `RefreshExpiry` | `string` | `"7d"` | Refresh token expiry duration |
| `BCryptCost` | `int` | `12` | BCrypt hashing cost (4-31) |
| `RateLimitRPM` | `int` | `60` | Rate limit requests per minute |
| `EmailRequired` | `bool` | `false` | Require email verification |

## Examples

Check the `/examples` folder for complete working examples:

- `basic_usage.go` - Core AuthKit functionality
- `gin_example.go` - Gin web framework integration  
- `fiber_example.go` - Fiber web framework integration
- `simple_http.go` - Standard HTTP server integration

## Support

For questions, issues, or contributions, please visit the GitHub repository.
