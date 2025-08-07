package main

import (
	"log"
	"time"

	"github.com/your-username/go-authkit"
)

func main() {
	// Initialize AuthKit with configuration
	auth := authkit.New(authkit.Config{
		JWTSecret:     "your-super-secret-jwt-key-here",
		TokenExpiry:   "24h",
		RefreshExpiry: "7d",
		BCryptCost:    12,
		RateLimitRPM:  100,
		EmailRequired: false,
	})

	// Example 1: Register a new user
	log.Println("=== User Registration ===")
	registerReq := authkit.RegisterRequest{
		Email:    "john@example.com",
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
		log.Printf("Registration failed: %v", err)
	} else {
		log.Printf("User registered: %+v", user)
	}

	// Example 2: Login user
	log.Println("\n=== User Login ===")
	tokenResponse, err := auth.LoginUser("john@example.com", "securepassword123")
	if err != nil {
		log.Printf("Login failed: %v", err)
		return
	}
	log.Printf("Login successful!")
	log.Printf("Access Token: %s", tokenResponse.AccessToken)
	log.Printf("Expires In: %d seconds", tokenResponse.ExpiresIn)
	log.Printf("User Info: %+v", tokenResponse.User)

	// Example 3: Validate token
	log.Println("\n=== Token Validation ===")
	claims, err := auth.ValidateToken(tokenResponse.AccessToken)
	if err != nil {
		log.Printf("Token validation failed: %v", err)
	} else {
		log.Printf("Token is valid!")
		log.Printf("Claims: %+v", claims)
	}

	// Example 4: Refresh token
	log.Println("\n=== Token Refresh ===")
	newTokens, err := auth.RefreshToken(tokenResponse.RefreshToken)
	if err != nil {
		log.Printf("Token refresh failed: %v", err)
	} else {
		log.Printf("Token refreshed!")
		log.Printf("New Access Token: %s", newTokens.AccessToken)
	}

	// Example 5: Generate custom token
	log.Println("\n=== Custom Token ===")
	customClaims := map[string]interface{}{
		"role":        "admin",
		"permissions": []string{"read", "write", "delete"},
		"scope":       "api:full",
	}
	customToken, err := auth.GenerateCustomToken(user.ID, customClaims, time.Hour*2)
	if err != nil {
		log.Printf("Custom token generation failed: %v", err)
	} else {
		log.Printf("Custom token generated: %s", customToken)
	}

	// Example 6: Password utilities
	log.Println("\n=== Password Utilities ===")
	hashedPassword, _ := authkit.HashPasswordStatic("mypassword123", 12)
	log.Printf("Hashed password: %s", hashedPassword)

	isValid := authkit.ComparePasswordStatic(hashedPassword, "mypassword123")
	log.Printf("Password is valid: %v", isValid)

	// Example 7: User management
	log.Println("\n=== User Management ===")
	// Update user
	updates := map[string]interface{}{
		"name": "John Smith",
		"role": "senior",
		"metadata": map[string]interface{}{
			"department": "engineering",
			"level":      "senior",
		},
	}
	updatedUser, err := auth.UpdateUser(user.ID, updates)
	if err != nil {
		log.Printf("User update failed: %v", err)
	} else {
		log.Printf("User updated: %+v", updatedUser)
	}

	// List all users
	allUsers := auth.ListUsers()
	log.Printf("Total users: %d", len(allUsers))
	for _, u := range allUsers {
		log.Printf("User: %s (%s) - Role: %s", u.Name, u.Email, u.Role)
	}
}
