package authkit

import (
	"testing"
	"time"
)

func TestAuthKit(t *testing.T) {
	// Initialize AuthKit for testing
	auth := New(Config{
		JWTSecret:     "test-secret-key-for-testing-only",
		TokenExpiry:   "1h",
		RefreshExpiry: "24h",
		BCryptCost:    4, // Lower cost for faster tests
		EmailRequired: false,
	})

	t.Run("RegisterUser", func(t *testing.T) {
		req := RegisterRequest{
			Email:    "test@example.com",
			Password: "testpassword123",
			Name:     "Test User",
			Role:     "user",
		}

		user, err := auth.RegisterUser(req)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if user.Email != req.Email {
			t.Errorf("Expected email %s, got %s", req.Email, user.Email)
		}

		if user.Name != req.Name {
			t.Errorf("Expected name %s, got %s", req.Name, user.Name)
		}

		// Try to register the same user again (should fail)
		_, err = auth.RegisterUser(req)
		if err != ErrUserAlreadyExists {
			t.Errorf("Expected ErrUserAlreadyExists, got %v", err)
		}
	})

	t.Run("LoginUser", func(t *testing.T) {
		// First register a user
		req := RegisterRequest{
			Email:    "login@example.com",
			Password: "loginpassword123",
			Name:     "Login Test User",
		}
		_, err := auth.RegisterUser(req)
		if err != nil {
			t.Fatalf("Failed to register user: %v", err)
		}

		// Test successful login
		tokenResponse, err := auth.LoginUser(req.Email, req.Password)
		if err != nil {
			t.Fatalf("Expected successful login, got error: %v", err)
		}

		if tokenResponse.AccessToken == "" {
			t.Error("Expected access token, got empty string")
		}

		if tokenResponse.RefreshToken == "" {
			t.Error("Expected refresh token, got empty string")
		}

		if tokenResponse.User.Email != req.Email {
			t.Errorf("Expected user email %s, got %s", req.Email, tokenResponse.User.Email)
		}

		// Test login with wrong password
		_, err = auth.LoginUser(req.Email, "wrongpassword")
		if err != ErrInvalidPassword {
			t.Errorf("Expected ErrInvalidPassword, got %v", err)
		}

		// Test login with non-existent user
		_, err = auth.LoginUser("nonexistent@example.com", "password")
		if err != ErrUserNotFound {
			t.Errorf("Expected ErrUserNotFound, got %v", err)
		}
	})

	t.Run("ValidateToken", func(t *testing.T) {
		// Register and login a user
		req := RegisterRequest{
			Email:    "validate@example.com",
			Password: "validatepassword123",
			Name:     "Validate Test User",
		}
		user, _ := auth.RegisterUser(req)
		tokenResponse, _ := auth.LoginUser(req.Email, req.Password)

		// Test valid token
		claims, err := auth.ValidateToken(tokenResponse.AccessToken)
		if err != nil {
			t.Fatalf("Expected valid token, got error: %v", err)
		}

		if claims.UserID != user.ID {
			t.Errorf("Expected user ID %s, got %s", user.ID, claims.UserID)
		}

		if claims.Email != req.Email {
			t.Errorf("Expected email %s, got %s", req.Email, claims.Email)
		}

		// Test invalid token
		_, err = auth.ValidateToken("invalid-token")
		if err != ErrInvalidToken {
			t.Errorf("Expected ErrInvalidToken, got %v", err)
		}
	})

	t.Run("RefreshToken", func(t *testing.T) {
		// Register and login a user
		req := RegisterRequest{
			Email:    "refresh@example.com",
			Password: "refreshpassword123",
			Name:     "Refresh Test User",
		}
		_, _ = auth.RegisterUser(req)
		tokenResponse, _ := auth.LoginUser(req.Email, req.Password)

		// Add a small delay to ensure different timestamps
		time.Sleep(time.Millisecond * 10)

		// Test valid refresh token
		newTokens, err := auth.RefreshToken(tokenResponse.RefreshToken)
		if err != nil {
			t.Fatalf("Expected successful refresh, got error: %v", err)
		}

		if newTokens.AccessToken == "" {
			t.Error("Expected new access token, got empty string")
		}

		if newTokens.AccessToken == tokenResponse.AccessToken {
			t.Error("Expected new access token to be different from original")
		}

		// Test invalid refresh token
		_, err = auth.RefreshToken("invalid-refresh-token")
		if err != ErrInvalidToken {
			t.Errorf("Expected ErrInvalidToken, got %v", err)
		}
	})

	t.Run("PasswordUtilities", func(t *testing.T) {
		password := "testpassword123"

		// Test password hashing
		hashedPassword, err := auth.HashPassword(password)
		if err != nil {
			t.Fatalf("Expected no error hashing password, got %v", err)
		}

		if hashedPassword == password {
			t.Error("Expected hashed password to be different from original")
		}

		// Test password comparison
		if !auth.ComparePassword(hashedPassword, password) {
			t.Error("Expected password comparison to be true")
		}

		if auth.ComparePassword(hashedPassword, "wrongpassword") {
			t.Error("Expected password comparison to be false for wrong password")
		}
	})

	t.Run("UserManagement", func(t *testing.T) {
		// Register a user
		req := RegisterRequest{
			Email:    "management@example.com",
			Password: "managementpassword123",
			Name:     "Management Test User",
		}
		user, _ := auth.RegisterUser(req)

		// Test get user by ID
		retrievedUser, err := auth.GetUserByID(user.ID)
		if err != nil {
			t.Fatalf("Expected to find user, got error: %v", err)
		}

		if retrievedUser.Email != req.Email {
			t.Errorf("Expected email %s, got %s", req.Email, retrievedUser.Email)
		}

		// Test get user by email
		retrievedUser, err = auth.GetUserByEmail(req.Email)
		if err != nil {
			t.Fatalf("Expected to find user by email, got error: %v", err)
		}

		if retrievedUser.ID != user.ID {
			t.Errorf("Expected user ID %s, got %s", user.ID, retrievedUser.ID)
		}

		// Test update user
		updates := map[string]interface{}{
			"name": "Updated Name",
			"role": "admin",
		}
		updatedUser, err := auth.UpdateUser(user.ID, updates)
		if err != nil {
			t.Fatalf("Expected successful update, got error: %v", err)
		}

		if updatedUser.Name != "Updated Name" {
			t.Errorf("Expected updated name 'Updated Name', got %s", updatedUser.Name)
		}

		if updatedUser.Role != "admin" {
			t.Errorf("Expected updated role 'admin', got %s", updatedUser.Role)
		}

		// Test list users
		users := auth.ListUsers()
		if len(users) == 0 {
			t.Error("Expected at least one user in the list")
		}

		// Test delete user
		err = auth.DeleteUser(user.ID)
		if err != nil {
			t.Fatalf("Expected successful deletion, got error: %v", err)
		}

		// Verify user is deleted
		_, err = auth.GetUserByID(user.ID)
		if err != ErrUserNotFound {
			t.Errorf("Expected ErrUserNotFound after deletion, got %v", err)
		}
	})

	t.Run("CustomToken", func(t *testing.T) {
		userID := "test-user-123"
		customClaims := map[string]interface{}{
			"role":        "admin",
			"permissions": []string{"read", "write", "delete"},
		}

		token, err := auth.GenerateCustomToken(userID, customClaims, time.Hour)
		if err != nil {
			t.Fatalf("Expected successful token generation, got error: %v", err)
		}

		if token == "" {
			t.Error("Expected token string, got empty")
		}
	})
}

func TestStaticPasswordUtilities(t *testing.T) {
	password := "testpassword123"

	// Test static password hashing
	hashedPassword, err := HashPasswordStatic(password, 4)
	if err != nil {
		t.Fatalf("Expected no error hashing password, got %v", err)
	}

	if hashedPassword == password {
		t.Error("Expected hashed password to be different from original")
	}

	// Test static password comparison
	if !ComparePasswordStatic(hashedPassword, password) {
		t.Error("Expected password comparison to be true")
	}

	if ComparePasswordStatic(hashedPassword, "wrongpassword") {
		t.Error("Expected password comparison to be false for wrong password")
	}
}
