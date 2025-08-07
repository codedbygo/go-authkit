// Package main demonstrates AuthKit with Fiber web framework
package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/codedbygo/go-authkit"
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

	// Create Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// Add CORS middleware
	app.Use(cors.New())

	// API routes
	api := app.Group("/api/v1")

	// Public routes
	api.Post("/register", auth.RegisterHandlerFiber)
	api.Post("/login", auth.LoginHandlerFiber)
	api.Post("/refresh", auth.RefreshHandlerFiber)

	// Health check
	api.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "AuthKit Fiber API is running",
			"version": "1.0.0",
		})
	})

	// Protected routes
	protected := api.Group("")
	protected.Use(auth.FiberMiddleware())

	// User profile routes
	protected.Get("/profile", auth.ProfileHandlerFiber)
	protected.Put("/profile", auth.UpdateProfileHandlerFiber)
	protected.Post("/logout", auth.LogoutHandlerFiber)

	// Protected resource examples
	protected.Get("/posts", getPostsHandlerFiber)
	protected.Post("/posts", createPostHandlerFiber)
	protected.Get("/posts/:id", getPostHandlerFiber)

	// Dashboard
	protected.Get("/dashboard", func(c *fiber.Ctx) error {
		claims, _ := authkit.GetUserFromFiberContext(c)
		return c.JSON(fiber.Map{
			"message":     "Welcome to your dashboard",
			"user":        claims.Email,
			"role":        claims.Role,
			"permissions": claims.Permissions,
		})
	})

	// Admin routes
	admin := protected.Group("/admin")
	admin.Use(auth.RequireRoleFiber("admin"))

	admin.Get("/users", listUsersHandlerFiber(auth))
	admin.Delete("/users/:id", deleteUserHandlerFiber(auth))
	admin.Put("/users/:id/role", updateUserRoleHandlerFiber(auth))

	// Moderator and Admin routes
	modAdmin := protected.Group("/moderate")
	modAdmin.Use(auth.RequireRolesFiber([]string{"admin", "moderator"}))

	modAdmin.Post("/posts/:id/approve", approvePostHandlerFiber)
	modAdmin.Delete("/posts/:id", deletePostHandlerFiber)

	log.Println("AuthKit Fiber Server starting on :8080")
	log.Println("Available endpoints:")
	log.Println("   POST /api/v1/register    - User registration")
	log.Println("   POST /api/v1/login       - User login")
	log.Println("   POST /api/v1/refresh     - Refresh token")
	log.Println("   GET  /api/v1/health      - Health check")
	log.Println("   GET  /api/v1/profile     - User profile (protected)")
	log.Println("   GET  /api/v1/dashboard   - Dashboard (protected)")
	log.Println("   GET  /api/v1/admin/users - Admin only")

	// Start server
	log.Fatal(app.Listen(":8080"))
}

// Example handlers for Fiber
func getPostsHandlerFiber(c *fiber.Ctx) error {
	claims, _ := authkit.GetUserFromFiberContext(c)
	return c.JSON(fiber.Map{
		"message": "Posts retrieved successfully",
		"user":    claims.Email,
		"posts": []fiber.Map{
			{"id": 1, "title": "First Post", "content": "Hello World"},
			{"id": 2, "title": "Second Post", "content": "AuthKit is awesome!"},
		},
	})
}

func createPostHandlerFiber(c *fiber.Ctx) error {
	var post struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}

	if err := c.BodyParser(&post); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	claims, _ := authkit.GetUserFromFiberContext(c)
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Post created successfully",
		"post": fiber.Map{
			"id":      123,
			"title":   post.Title,
			"content": post.Content,
			"author":  claims.Email,
		},
	})
}

func getPostHandlerFiber(c *fiber.Ctx) error {
	id := c.Params("id")
	return c.JSON(fiber.Map{
		"message": "Post retrieved successfully",
		"post": fiber.Map{
			"id":      id,
			"title":   "Sample Post",
			"content": "This is a sample post content",
		},
	})
}

func listUsersHandlerFiber(auth *authkit.AuthKit) fiber.Handler {
	return func(c *fiber.Ctx) error {
		users := auth.ListUsers()
		return c.JSON(fiber.Map{
			"message": "Users retrieved successfully",
			"count":   len(users),
			"users":   users,
		})
	}
}

func deleteUserHandlerFiber(auth *authkit.AuthKit) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := c.Params("id")
		err := auth.DeleteUser(userID)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.JSON(fiber.Map{
			"message": "User deleted successfully",
		})
	}
}

func updateUserRoleHandlerFiber(auth *authkit.AuthKit) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := c.Params("id")
		var req struct {
			Role string `json:"role"`
		}

		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		updates := map[string]interface{}{"role": req.Role}
		updatedUser, err := auth.UpdateUser(userID, updates)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"message": "User role updated successfully",
			"user":    updatedUser,
		})
	}
}

func approvePostHandlerFiber(c *fiber.Ctx) error {
	postID := c.Params("id")
	claims, _ := authkit.GetUserFromFiberContext(c)
	return c.JSON(fiber.Map{
		"message":     "Post approved successfully",
		"post_id":     postID,
		"approved_by": claims.Email,
	})
}

func deletePostHandlerFiber(c *fiber.Ctx) error {
	postID := c.Params("id")
	claims, _ := authkit.GetUserFromFiberContext(c)
	return c.JSON(fiber.Map{
		"message":    "Post deleted successfully",
		"post_id":    postID,
		"deleted_by": claims.Email,
	})
}
