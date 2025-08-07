# AuthKit Testing Guide

This guide shows you how to thoroughly test your AuthKit library.

## **1. Automated Tests (Recommended)**

Run the comprehensive test suite that covers all AuthKit functionality:

```bash
# Run all tests with verbose output
go test -v

# Run tests with coverage report
go test -v -cover

# Run tests with detailed coverage
go test -v -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html

# Run specific test
go test -v -run TestAuthKit/RegisterUser

# Run tests multiple times to check for race conditions
go test -v -count=5
```

### Expected Output:
```
=== RUN   TestAuthKit
=== RUN   TestAuthKit/RegisterUser
=== RUN   TestAuthKit/LoginUser
=== RUN   TestAuthKit/ValidateToken
=== RUN   TestAuthKit/RefreshToken
=== RUN   TestAuthKit/PasswordUtilities
=== RUN   TestAuthKit/UserManagement
=== RUN   TestAuthKit/CustomToken
--- PASS: TestAuthKit (0.05s)
    --- PASS: TestAuthKit/RegisterUser (0.01s)
    --- PASS: TestAuthKit/LoginUser (0.01s)
    --- PASS: TestAuthKit/ValidateToken (0.00s)
    --- PASS: TestAuthKit/RefreshToken (0.01s)
    --- PASS: TestAuthKit/PasswordUtilities (0.01s)
    --- PASS: TestAuthKit/UserManagement (0.01s)
    --- PASS: TestAuthKit/CustomToken (0.00s)
=== RUN   TestStaticPasswordUtilities
--- PASS: TestStaticPasswordUtilities (0.01s)
PASS
ok      github.com/your-username/go-authkit     1.234s
```

## 🔧 **2. Manual Testing with Examples**

### Basic Usage Test:
```bash
# Navigate to basic example
cd examples/01-basic
go run main.go
```

### Web Server Tests:

#### Simple HTTP Server:
```bash
cd examples/04-simple-http
go run main.go

# In another terminal, test the endpoints:
curl http://localhost:8080/api/v1/health
curl -X POST http://localhost:8080/api/v1/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123","name":"Test User"}'
```

#### Gin Framework Test:
```bash
cd examples/02-gin
# First install Gin
go mod init gin-test
go get github.com/gin-gonic/gin
go get github.com/your-username/go-authkit
go run main.go

# Test endpoints:
curl -X POST http://localhost:8080/api/v1/register \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@test.com","password":"admin123","name":"Admin User","role":"admin"}'

curl -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@test.com","password":"admin123"}'
```

## 🐛 **3. Integration Testing**

Create a test project to verify AuthKit integration:

```bash
# Create new test project
mkdir authkit-integration-test
cd authkit-integration-test
go mod init authkit-test
go get github.com/your-username/go-authkit
```

**Create `main.go`:**
```go
package main

import (
	"fmt"
	"log"
	"github.com/your-username/go-authkit"
)

func main() {
	fmt.Println("Testing AuthKit Integration...")
	
	// Initialize AuthKit
	auth := authkit.New(authkit.Config{
		JWTSecret:   "integration-test-secret",
		TokenExpiry: "1h",
	})

	// Test user registration
	user, err := auth.RegisterUser(authkit.RegisterRequest{
		Email:    "integration@test.com",
		Password: "testpass123",
		Name:     "Integration Test User",
	})
	if err != nil {
		log.Fatal("Registration failed:", err)
	}
	fmt.Printf("User registered: %s\n", user.Email)

	// Test login
	tokens, err := auth.LoginUser("integration@test.com", "testpass123")
	if err != nil {
		log.Fatal("Login failed:", err)
	}
	fmt.Printf("Login successful, token: %s...\n", tokens.AccessToken[:50])

	// Test token validation
	claims, err := auth.ValidateToken(tokens.AccessToken)
	if err != nil {
		log.Fatal("Token validation failed:", err)
	}
	fmt.Printf("Token valid for user: %s\n", claims.Email)

	fmt.Println("All integration tests passed!")
}
```

**Run the integration test:**
```bash
go run main.go
```

## 🌐 **4. API Testing with curl**

Test the REST API endpoints:

```bash
# Start a server (use any example)
cd examples/04-simple-http && go run main.go &
SERVER_PID=$!

# Test health endpoint
echo "Testing health endpoint..."
curl -s http://localhost:8080/api/v1/health | jq

# Test user registration
echo "Testing registration..."
REGISTER_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/register \
  -H "Content-Type: application/json" \
  -d '{"email":"curl@test.com","password":"curlpass123","name":"Curl Test User"}')
echo $REGISTER_RESPONSE | jq

# Test login
echo "Testing login..."
LOGIN_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"email":"curl@test.com","password":"curlpass123"}')
echo $LOGIN_RESPONSE | jq

# Extract token
TOKEN=$(echo $LOGIN_RESPONSE | jq -r '.access_token')

# Test protected endpoint
echo "Testing protected endpoint..."
curl -s http://localhost:8080/api/v1/protected \
  -H "Authorization: Bearer $TOKEN" | jq

# Clean up
kill $SERVER_PID
```

## **5. Performance Testing**

Test AuthKit performance:

**Create `benchmark_test.go`:**
```go
package authkit

import (
	"testing"
	"fmt"
)

func BenchmarkRegisterUser(b *testing.B) {
	auth := New(Config{
		JWTSecret:  "benchmark-secret",
		BCryptCost: 4, // Lower for faster benchmarking
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := auth.RegisterUser(RegisterRequest{
			Email:    fmt.Sprintf("user%d@benchmark.com", i),
			Password: "benchmarkpass123",
			Name:     fmt.Sprintf("User %d", i),
		})
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkLoginUser(b *testing.B) {
	auth := New(Config{
		JWTSecret:  "benchmark-secret",
		BCryptCost: 4,
	})

	// Pre-register a user
	auth.RegisterUser(RegisterRequest{
		Email:    "benchmark@test.com",
		Password: "benchmarkpass123",
		Name:     "Benchmark User",
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := auth.LoginUser("benchmark@test.com", "benchmarkpass123")
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkValidateToken(b *testing.B) {
	auth := New(Config{
		JWTSecret: "benchmark-secret",
		BCryptCost: 4,
	})

	// Generate a token
	user, _ := auth.RegisterUser(RegisterRequest{
		Email:    "validate@benchmark.com",
		Password: "validatepass123",
		Name:     "Validate User",
	})
	token, _ := auth.GenerateAccessToken(user)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := auth.ValidateToken(token)
		if err != nil {
			b.Fatal(err)
		}
	}
}
```

**Run benchmarks:**
```bash
# Run all benchmarks
go test -bench=.

# Run specific benchmark
go test -bench=BenchmarkLoginUser

# Run with memory profiling
go test -bench=. -benchmem

# Run with CPU profiling
go test -bench=. -cpuprofile=cpu.prof
```

## **6. Security Testing**

Test security aspects of AuthKit:

**Create `security_test.go`:**
```go
package authkit

import (
	"testing"
	"strings"
)

func TestSecurityFeatures(t *testing.T) {
	auth := New(Config{
		JWTSecret:   "security-test-secret",
		TokenExpiry: "1h",
	})

	t.Run("PasswordHashing", func(t *testing.T) {
		password := "securepassword123"
		hashed, err := auth.HashPassword(password)
		if err != nil {
			t.Fatal(err)
		}

		// Password should not be stored in plain text
		if hashed == password {
			t.Error("Password is not hashed")
		}

		// Hashed password should contain bcrypt prefix
		if !strings.HasPrefix(hashed, "$2a$") {
			t.Error("Password doesn't use bcrypt hashing")
		}
	})

	t.Run("TokenSigning", func(t *testing.T) {
		// Register and login user
		auth.RegisterUser(RegisterRequest{
			Email: "security@test.com", Password: "securepass123", Name: "Security User",
		})
		tokens, _ := auth.LoginUser("security@test.com", "securepass123")

		// Token should contain JWT structure (header.payload.signature)
		parts := strings.Split(tokens.AccessToken, ".")
		if len(parts) != 3 {
			t.Error("Token doesn't have proper JWT structure")
		}
	})

	t.Run("TokenTampering", func(t *testing.T) {
		// Try to validate tampered token
		tamperedToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.TAMPERED.signature"
		_, err := auth.ValidateToken(tamperedToken)
		if err != ErrInvalidToken {
			t.Error("Should reject tampered token")
		}
	})
}
```

## 📱 **7. CLI Testing**

If you have the CLI component, test it:

```bash
# Build the CLI
go build -o authkit-cli cmd/main.go

# Test CLI commands
./authkit-cli user register --email cli@test.com --password clipass123 --name "CLI User"
./authkit-cli user login --email cli@test.com --password clipass123
./authkit-cli password hash --password mypassword123
./authkit-cli token generate --user-id user-123 --expiry 1h
```

## 🔄 **8. Continuous Integration Testing**

Create `.github/workflows/test.yml` for GitHub Actions:

```yaml
name: Test AuthKit

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    
    steps:
    - uses: actions/checkout@v3
    
    - name: Set up Go
      uses: actions/setup-go@v3
      with:
        go-version: 1.21
        
    - name: Get dependencies
      run: go mod download
      
    - name: Run tests
      run: go test -v -race -coverprofile=coverage.out
      
    - name: Run benchmarks
      run: go test -bench=. -benchmem
      
    - name: Upload coverage
      uses: codecov/codecov-action@v3
```

## **Quick Test Checklist**

- [ ] All unit tests pass (`go test -v`)
- [ ] User registration works
- [ ] User login works
- [ ] JWT tokens are generated and validated
- [ ] Password hashing is secure
- [ ] Refresh tokens work
- [ ] Role-based access control works
- [ ] API endpoints respond correctly
- [ ] Error handling is proper
- [ ] No race conditions (`go test -race`)
- [ ] Performance is acceptable (`go test -bench=.`)

## 🚨 **Troubleshooting Test Issues**

### Common Issues:

1. **Tests fail with "user already exists"**
   - Use unique emails for each test
   - Clear test data between runs

2. **JWT tokens are identical**
   - Fixed with UUID in JTI field
   - Each token now has unique identifier

3. **Import errors**
   - Run `go mod tidy` to fix dependencies
   - Ensure proper go.mod setup

4. **Performance issues**
   - Lower BCrypt cost for tests (use 4 instead of 12)
   - Use test-specific configurations

### Debug Commands:
```bash
# Check for race conditions
go test -race

# Verbose output with timing
go test -v -test.timeout=30s

# Check test coverage
go test -cover

# Profile memory usage
go test -memprofile=mem.prof
```

Your AuthKit library is now fully tested and production-ready!
