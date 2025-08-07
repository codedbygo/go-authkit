package authkit

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// GinMiddleware returns a Gin middleware function for authentication
func (a *AuthKit) GinMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Check if the header starts with "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		// Extract the token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Validate the token
		claims, err := a.ValidateToken(tokenString)
		if err != nil {
			status := http.StatusUnauthorized
			message := "Invalid token"

			if err == ErrTokenExpired {
				status = http.StatusUnauthorized
				message = "Token expired"
			}

			c.JSON(status, gin.H{"error": message})
			c.Abort()
			return
		}

		// Set user information in context
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_role", claims.Role)
		c.Set("user_permissions", claims.Permissions)
		c.Set("user_claims", claims)

		c.Next()
	}
}

// RequireRole returns a Gin middleware that requires a specific role
func (a *AuthKit) RequireRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("user_role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			c.Abort()
			return
		}

		if userRole != role {
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireRoles returns a Gin middleware that requires one of the specified roles
func (a *AuthKit) RequireRoles(roles []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("user_role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			c.Abort()
			return
		}

		hasRole := false
		for _, role := range roles {
			if userRole == role {
				hasRole = true
				break
			}
		}

		if !hasRole {
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequirePermission returns a Gin middleware that requires a specific permission
func (a *AuthKit) RequirePermission(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userPermissions, exists := c.Get("user_permissions")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			c.Abort()
			return
		}

		permissions, ok := userPermissions.([]string)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid permissions format"})
			c.Abort()
			return
		}

		hasPermission := false
		for _, perm := range permissions {
			if perm == permission {
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// GetUserFromGinContext extracts user information from Gin context
func GetUserFromGinContext(c *gin.Context) (*Claims, bool) {
	claims, exists := c.Get("user_claims")
	if !exists {
		return nil, false
	}

	userClaims, ok := claims.(*Claims)
	return userClaims, ok
}
