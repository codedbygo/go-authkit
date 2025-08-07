// Package main demonstrates AuthKit with Gin web framework
package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/your-username/go-authkit"
)

func main() {
	// Initialize AuthKit
	auth := authkit.New(authkit.Config{
		JWTSecret:     "your-super-secret-jwt-key-here",
		TokenExpiry:   "24h",
		RefreshExpiry: "7d",
		BCryptCost:    12,
		EmailRequired: false,
	})

	// Create Gin router
	r := gin.Default()

	// Add CORS middleware (optional)
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Public routes (no authentication required)
	api := r.Group("/api/v1")
	{
		api.POST("/register", auth.RegisterHandler)
		api.POST("/login", auth.LoginHandler)
		api.POST("/refresh", auth.RefreshHandler)

		// Health check
		api.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status":  "ok",
				"message": "AuthKit API is running",
				"version": "1.0.0",
			})
		})
	}

	// Protected routes (authentication required)
	protected := api.Group("")
	protected.Use(auth.GinMiddleware())
	{
		// User profile routes
		protected.GET("/profile", auth.ProfileHandler)
		protected.PUT("/profile", auth.UpdateProfileHandler)
		protected.POST("/logout", auth.LogoutHandler)

		// Protected resource examples
		protected.GET("/posts", getPostsHandler)
		protected.POST("/posts", createPostHandler)
		protected.GET("/posts/:id", getPostHandler)

		// Dashboard route
		protected.GET("/dashboard", func(c *gin.Context) {
			claims, _ := authkit.GetUserFromGinContext(c)
			c.JSON(http.StatusOK, gin.H{
				"message":     "Welcome to your dashboard",
				"user":        claims.Email,
				"role":        claims.Role,
				"permissions": claims.Permissions,
			})
		})
	}

	// Admin routes (admin role required)
	admin := protected.Group("/admin")
	admin.Use(auth.RequireRole("admin"))
	{
		admin.GET("/users", listUsersHandler(auth))
		admin.DELETE("/users/:id", deleteUserHandler(auth))
		admin.PUT("/users/:id/role", updateUserRoleHandler(auth))
	}

	// Moderator and Admin routes (multiple roles allowed)
	modAdmin := protected.Group("/moderate")
	modAdmin.Use(auth.RequireRoles([]string{"admin", "moderator"}))
	{
		modAdmin.POST("/posts/:id/approve", approvePostHandler)
		modAdmin.DELETE("/posts/:id", deletePostHandler)
	}

	log.Println("AuthKit Gin Server starting on :8080")
	log.Println("Available endpoints:")
	log.Println("   POST /api/v1/register    - User registration")
	log.Println("   POST /api/v1/login       - User login")
	log.Println("   POST /api/v1/refresh     - Refresh token")
	log.Println("   GET  /api/v1/health      - Health check")
	log.Println("   GET  /api/v1/profile     - User profile (protected)")
	log.Println("   GET  /api/v1/dashboard   - Dashboard (protected)")
	log.Println("   GET  /api/v1/admin/users - Admin only")
	log.Println("")
	log.Println("Example login request:")
	log.Println(`   curl -X POST http://localhost:8080/api/v1/login \
     -H "Content-Type: application/json" \
     -d '{"email":"test@example.com","password":"password123"}'`)

	// Start server
	r.Run(":8080")
}

// Example handlers for demonstration
func getPostsHandler(c *gin.Context) {
	claims, _ := authkit.GetUserFromGinContext(c)
	c.JSON(http.StatusOK, gin.H{
		"message": "Posts retrieved successfully",
		"user":    claims.Email,
		"posts": []map[string]interface{}{
			{"id": 1, "title": "First Post", "content": "Hello World"},
			{"id": 2, "title": "Second Post", "content": "AuthKit is awesome!"},
		},
	})
}

func createPostHandler(c *gin.Context) {
	var post struct {
		Title   string `json:"title" binding:"required"`
		Content string `json:"content" binding:"required"`
	}

	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	claims, _ := authkit.GetUserFromGinContext(c)
	c.JSON(http.StatusCreated, gin.H{
		"message": "Post created successfully",
		"post": gin.H{
			"id":      123,
			"title":   post.Title,
			"content": post.Content,
			"author":  claims.Email,
		},
	})
}

func getPostHandler(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"message": "Post retrieved successfully",
		"post": gin.H{
			"id":      id,
			"title":   "Sample Post",
			"content": "This is a sample post content",
		},
	})
}

func listUsersHandler(auth *authkit.AuthKit) gin.HandlerFunc {
	return func(c *gin.Context) {
		users := auth.ListUsers()
		c.JSON(http.StatusOK, gin.H{
			"message": "Users retrieved successfully",
			"count":   len(users),
			"users":   users,
		})
	}
}

func deleteUserHandler(auth *authkit.AuthKit) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Param("id")
		err := auth.DeleteUser(userID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "User deleted successfully",
		})
	}
}

func updateUserRoleHandler(auth *authkit.AuthKit) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Param("id")
		var req struct {
			Role string `json:"role" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		updates := map[string]interface{}{"role": req.Role}
		updatedUser, err := auth.UpdateUser(userID, updates)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "User role updated successfully",
			"user":    updatedUser,
		})
	}
}

func approvePostHandler(c *gin.Context) {
	postID := c.Param("id")
	claims, _ := authkit.GetUserFromGinContext(c)
	c.JSON(http.StatusOK, gin.H{
		"message":     "Post approved successfully",
		"post_id":     postID,
		"approved_by": claims.Email,
	})
}

func deletePostHandler(c *gin.Context) {
	postID := c.Param("id")
	claims, _ := authkit.GetUserFromGinContext(c)
	c.JSON(http.StatusOK, gin.H{
		"message":    "Post deleted successfully",
		"post_id":    postID,
		"deleted_by": claims.Email,
	})
}
