package authkit

import (
	"github.com/gofiber/fiber/v2"
)

// RegisterHandlerFiber handles user registration for Fiber
func (a *AuthKit) RegisterHandlerFiber(c *fiber.Ctx) error {
	var req RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	user, err := a.RegisterUser(req)
	if err != nil {
		status := fiber.StatusBadRequest
		if err == ErrUserAlreadyExists {
			status = fiber.StatusConflict
		}
		return c.Status(status).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User registered successfully",
		"user":    user,
	})
}

// LoginHandlerFiber handles user login for Fiber
func (a *AuthKit) LoginHandlerFiber(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	tokenResponse, err := a.LoginUser(req.Email, req.Password)
	if err != nil {
		status := fiber.StatusUnauthorized
		if err == ErrUserNotFound {
			status = fiber.StatusNotFound
		}
		return c.Status(status).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(tokenResponse)
}

// RefreshHandlerFiber handles token refresh for Fiber
func (a *AuthKit) RefreshHandlerFiber(c *fiber.Ctx) error {
	var req RefreshRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	tokenResponse, err := a.RefreshToken(req.RefreshToken)
	if err != nil {
		status := fiber.StatusUnauthorized
		if err == ErrTokenExpired {
			status = fiber.StatusUnauthorized
		}
		return c.Status(status).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(tokenResponse)
}

// ProfileHandlerFiber returns current user profile for Fiber
func (a *AuthKit) ProfileHandlerFiber(c *fiber.Ctx) error {
	claims, exists := GetUserFromFiberContext(c)
	if !exists {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not found in context",
		})
	}

	user, err := a.GetUserByID(claims.UserID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	return c.JSON(fiber.Map{
		"user": a.userToUserInfo(user),
	})
}

// UpdateProfileHandlerFiber updates current user profile for Fiber
func (a *AuthKit) UpdateProfileHandlerFiber(c *fiber.Ctx) error {
	claims, exists := GetUserFromFiberContext(c)
	if !exists {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not found in context",
		})
	}

	var updates map[string]interface{}
	if err := c.BodyParser(&updates); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Remove sensitive fields that shouldn't be updated via this endpoint
	delete(updates, "id")
	delete(updates, "password")
	delete(updates, "email")
	delete(updates, "created_at")
	delete(updates, "updated_at")

	updatedUser, err := a.UpdateUser(claims.UserID, updates)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Profile updated successfully",
		"user":    updatedUser,
	})
}

// LogoutHandlerFiber handles user logout for Fiber (for completeness - JWT is stateless)
func (a *AuthKit) LogoutHandlerFiber(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "Logged out successfully",
	})
}
