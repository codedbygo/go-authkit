package authkit

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

// FiberMiddleware returns a Fiber middleware function for authentication
func (a *AuthKit) FiberMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get token from Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Authorization header required",
			})
		}

		// Check if the header starts with "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid authorization header format",
			})
		}

		// Extract the token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Validate the token
		claims, err := a.ValidateToken(tokenString)
		if err != nil {
			status := fiber.StatusUnauthorized
			message := "Invalid token"

			if err == ErrTokenExpired {
				status = fiber.StatusUnauthorized
				message = "Token expired"
			}

			return c.Status(status).JSON(fiber.Map{
				"error": message,
			})
		}

		// Set user information in context
		c.Locals("user_id", claims.UserID)
		c.Locals("user_email", claims.Email)
		c.Locals("user_role", claims.Role)
		c.Locals("user_permissions", claims.Permissions)
		c.Locals("user_claims", claims)

		return c.Next()
	}
}

// RequireRoleFiber returns a Fiber middleware that requires a specific role
func (a *AuthKit) RequireRoleFiber(role string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRole := c.Locals("user_role")
		if userRole == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "User not authenticated",
			})
		}

		if userRole != role {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Insufficient permissions",
			})
		}

		return c.Next()
	}
}

// RequireRolesFiber returns a Fiber middleware that requires one of the specified roles
func (a *AuthKit) RequireRolesFiber(roles []string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRole := c.Locals("user_role")
		if userRole == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "User not authenticated",
			})
		}

		hasRole := false
		for _, role := range roles {
			if userRole == role {
				hasRole = true
				break
			}
		}

		if !hasRole {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Insufficient permissions",
			})
		}

		return c.Next()
	}
}

// RequirePermissionFiber returns a Fiber middleware that requires a specific permission
func (a *AuthKit) RequirePermissionFiber(permission string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userPermissions := c.Locals("user_permissions")
		if userPermissions == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "User not authenticated",
			})
		}

		permissions, ok := userPermissions.([]string)
		if !ok {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Invalid permissions format",
			})
		}

		hasPermission := false
		for _, perm := range permissions {
			if perm == permission {
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Insufficient permissions",
			})
		}

		return c.Next()
	}
}

// GetUserFromFiberContext extracts user information from Fiber context
func GetUserFromFiberContext(c *fiber.Ctx) (*Claims, bool) {
	claims := c.Locals("user_claims")
	if claims == nil {
		return nil, false
	}

	userClaims, ok := claims.(*Claims)
	return userClaims, ok
}
