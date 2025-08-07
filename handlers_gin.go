package authkit

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RegisterHandler handles user registration for Gin
func (a *AuthKit) RegisterHandler(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := a.RegisterUser(req)
	if err != nil {
		status := http.StatusBadRequest
		if err == ErrUserAlreadyExists {
			status = http.StatusConflict
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"user":    user,
	})
}

// LoginHandler handles user login for Gin
func (a *AuthKit) LoginHandler(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tokenResponse, err := a.LoginUser(req.Email, req.Password)
	if err != nil {
		status := http.StatusUnauthorized
		if err == ErrUserNotFound {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tokenResponse)
}

// RefreshHandler handles token refresh for Gin
func (a *AuthKit) RefreshHandler(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tokenResponse, err := a.RefreshToken(req.RefreshToken)
	if err != nil {
		status := http.StatusUnauthorized
		if err == ErrTokenExpired {
			status = http.StatusUnauthorized
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tokenResponse)
}

// ProfileHandler returns current user profile for Gin
func (a *AuthKit) ProfileHandler(c *gin.Context) {
	claims, exists := GetUserFromGinContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}

	user, err := a.GetUserByID(claims.UserID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": a.userToUserInfo(user),
	})
}

// UpdateProfileHandler updates current user profile for Gin
func (a *AuthKit) UpdateProfileHandler(c *gin.Context) {
	claims, exists := GetUserFromGinContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Remove sensitive fields that shouldn't be updated via this endpoint
	delete(updates, "id")
	delete(updates, "password")
	delete(updates, "email")
	delete(updates, "created_at")
	delete(updates, "updated_at")

	updatedUser, err := a.UpdateUser(claims.UserID, updates)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Profile updated successfully",
		"user":    updatedUser,
	})
}

// LogoutHandler handles user logout for Gin (for completeness - JWT is stateless)
func (a *AuthKit) LogoutHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Logged out successfully",
	})
}
