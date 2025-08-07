package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"
)

// Simple HTTP server example without external dependencies
func main() {
	log.Println("AuthKit Simple HTTP Server starting on :8080")
	log.Println("This is a basic example without external web frameworks")

	// Routes
	http.HandleFunc("/api/v1/health", healthHandler)
	http.HandleFunc("/api/v1/register", corsHandler(registerHandler))
	http.HandleFunc("/api/v1/login", corsHandler(loginHandler))
	http.HandleFunc("/api/v1/protected", corsHandler(protectedHandler))

	log.Println("Available endpoints:")
	log.Println("   GET  /api/v1/health     - Health check")
	log.Println("   POST /api/v1/register   - User registration")
	log.Println("   POST /api/v1/login      - User login")
	log.Println("   GET  /api/v1/protected  - Protected route (requires Bearer token)")
	log.Println("")
	log.Println("Example requests:")
	log.Println("Register:")
	log.Println(`  curl -X POST http://localhost:8080/api/v1/register \
    -H "Content-Type: application/json" \
    -d '{"email":"test@example.com","password":"password123","name":"Test User"}'`)
	log.Println("")
	log.Println("Login:")
	log.Println(`  curl -X POST http://localhost:8080/api/v1/login \
    -H "Content-Type: application/json" \
    -d '{"email":"test@example.com","password":"password123"}'`)

	// Start server
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// CORS middleware wrapper
func corsHandler(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

// Health check handler
func healthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "ok",
		"message": "AuthKit Simple HTTP Server is running",
		"time":    time.Now().Format(time.RFC3339),
	})
}

// Register handler
func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Name     string `json:"name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Simulate user registration (in real app, use AuthKit)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "User registered successfully (demo)",
		"user": map[string]interface{}{
			"id":    "user-123",
			"email": req.Email,
			"name":  req.Name,
		},
	})
}

// Login handler
func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Simulate login (in real app, use AuthKit)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":      "Login successful (demo)",
		"access_token": "demo-jwt-token-here",
		"token_type":   "Bearer",
		"expires_in":   3600,
		"user": map[string]interface{}{
			"id":    "user-123",
			"email": req.Email,
			"name":  "Demo User",
			"role":  "user",
		},
	})
}

// Protected handler
func protectedHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Check Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Authorization header required", http.StatusUnauthorized)
		return
	}

	if !strings.HasPrefix(authHeader, "Bearer ") {
		http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
		return
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")

	// In real app, validate token with AuthKit
	if token != "demo-jwt-token-here" {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Access granted to protected resource",
		"user":    "demo-user@example.com",
		"data":    "This is protected data",
		"time":    time.Now().Format(time.RFC3339),
	})
}

// Example of how to integrate with AuthKit:
/*
import "github.com/your-username/go-authkit"

func realAuthExample() {
	// Initialize AuthKit
	auth := authkit.New(authkit.Config{
		JWTSecret:   "your-secret-key",
		TokenExpiry: "24h",
	})

	// In your login handler:
	tokenResponse, err := auth.LoginUser(req.Email, req.Password)
	if err != nil {
		// Handle error
		return
	}

	// In your protected handler:
	claims, err := auth.ValidateToken(tokenString)
	if err != nil {
		// Handle invalid token
		return
	}

	// Use claims.UserID, claims.Email, claims.Role, etc.
}
*/
