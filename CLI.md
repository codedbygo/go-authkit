# AuthKit CLI Tool

A command-line interface for the AuthKit authentication library.

## Installation

```bash
# Build from source
go build -o authkit cmd/main.go

# Or install globally
go install github.com/your-username/go-authkit/cmd@latest
```

## Quick Start

```bash
# Set your JWT secret (required for all operations)
export JWT_SECRET="your-super-secret-jwt-key-here"

# Register a user
authkit user register --secret $JWT_SECRET --email user@example.com --password password123 --name "John Doe"

# Login user
authkit user login --secret $JWT_SECRET --email user@example.com --password password123

# Generate a token
authkit token generate --secret $JWT_SECRET --user-id user123 --expiry 24h

# Validate a token
authkit token validate --secret $JWT_SECRET --token "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."

# Hash a password
authkit password hash --password mypassword123 --cost 12

# Start development server
authkit server start --secret $JWT_SECRET --port 8080
```

## Commands

### User Management

#### Register User
```bash
authkit user register [flags]

Flags:
  -e, --email string      User email (required)
  -n, --name string       User name (required) 
  -p, --password string   User password (required)
  -r, --role string       User role (default "user")
  -s, --secret string     JWT secret key (required)
```

#### Login User
```bash
authkit user login [flags]

Flags:
  -e, --email string      User email (required)
  -p, --password string   User password (required)
  -s, --secret string     JWT secret key (required)
```

#### List Users
```bash
authkit user list [flags]

Flags:
  -s, --secret string     JWT secret key (required)
```

#### Delete User
```bash
authkit user delete [flags]

Flags:
  -i, --id string         User ID (required)
  -s, --secret string     JWT secret key (required)
```

### Token Management

#### Generate Token
```bash
authkit token generate [flags]

Flags:
  -c, --claims strings    Custom claims (key=value format)
  -x, --expiry string     Token expiry duration (default "24h")
  -s, --secret string     JWT secret key (required)
  -u, --user-id string    User ID (required)
```

#### Validate Token
```bash
authkit token validate [flags]

Flags:
  -s, --secret string     JWT secret key (required)
  -t, --token string      JWT token to validate (required)
```

#### Refresh Token
```bash
authkit token refresh [flags]

Flags:
  -r, --refresh-token string   Refresh token (required)
  -s, --secret string          JWT secret key (required)
```

### Password Utilities

#### Hash Password
```bash
authkit password hash [flags]

Flags:
  -c, --cost int          BCrypt cost (default 12)
  -p, --password string   Password to hash (required)
```

#### Compare Password
```bash
authkit password compare [flags]

Flags:
  -H, --hash string       Hashed password (required)
  -p, --password string   Plain password (required)
```

### Server Commands

#### Start Server
```bash
authkit server start [flags]

Flags:
  -c, --cors              Enable CORS (default true)
  -H, --host string       Server host (default "localhost")
  -l, --logging           Enable request logging (default true)
  -p, --port string       Server port (default "8080")
  -s, --secret string     JWT secret key (required)
```

#### Test Server
```bash
authkit server test [flags]

Flags:
  -H, --host string       Server host (default "localhost")
  -p, --port string       Server port (default "8080")
```

## Examples

### Complete User Workflow

```bash
# 1. Register a new user
authkit user register \
  --secret "mySecretKey123" \
  --email "alice@company.com" \
  --password "securePassword456" \
  --name "Alice Johnson" \
  --role "developer"

# Output:
# User registered successfully!
# {
#   "user_id": "550e8400-e29b-41d4-a716-446655440000",
#   "email": "alice@company.com",
#   "name": "Alice Johnson", 
#   "role": "developer"
# }

# 2. Login to get tokens
authkit user login \
  --secret "mySecretKey123" \
  --email "alice@company.com" \
  --password "securePassword456"

# Output:
# Login successful!
# {
#   "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
#   "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
#   "token_type": "Bearer",
#   "expires_in": 3600,
#   "user": {...}
# }

# 3. Validate the access token
authkit token validate \
  --secret "mySecretKey123" \
  --token "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."

# Output:
# Token is valid!
# {
#   "valid": true,
#   "user_id": "550e8400-e29b-41d4-a716-446655440000",
#   "email": "alice@company.com",
#   "role": "developer"
# }
```

### Password Management

```bash
# Hash a password
authkit password hash --password "mySecretPassword" --cost 12

# Output:
# Password hashed successfully!
# {
#   "original": "mySecretPassword",
#   "hashed": "$2a$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LewSuqGOSWPQtk3S6",
#   "cost": 12
# }

# Verify password against hash
authkit password compare \
  --password "mySecretPassword" \
  --hash "$2a$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LewSuqGOSWPQtk3S6"

# Output:
# Password matches hash!
# {
#   "valid": true,
#   "password": "mySecretPassword",
#   "hash": "$2a$12$..."
# }
```

### Development Server

```bash
# Start development server
authkit server start \
  --secret "mySecretKey123" \
  --port 8080 \
  --host localhost \
  --cors \
  --logging

# Output:
# Starting AuthKit Server...
# Host: localhost
# Port: 8080
# JWT Secret: mySecretKey123
# 
# Available endpoints:
#   POST /localhost:8080/api/v1/register    - User registration
#   POST /localhost:8080/api/v1/login       - User login
#   GET  /localhost:8080/api/v1/profile     - User profile (protected)
#   GET  /localhost:8080/api/v1/health      - Health check

# Test server endpoints
authkit server test --port 8080 --host localhost
```

## Global Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--secret` | `-s` | *required* | JWT secret key for signing tokens |
| `--output` | `-o` | `json` | Output format (json, table, yaml) |

## Output Formats

### JSON (default)
```bash
authkit user list --secret "key" --output json
```

### Table
```bash  
authkit user list --secret "key" --output table
```

### YAML
```bash
authkit user list --secret "key" --output yaml  
```

## Environment Variables

You can set commonly used values as environment variables:

```bash
export JWT_SECRET="your-super-secret-jwt-key"
export AUTHKIT_OUTPUT="json"
export AUTHKIT_HOST="localhost"
export AUTHKIT_PORT="8080"

# Now you can use commands without repetitive flags
authkit user register --email user@example.com --password pass123 --name "User"
authkit server start
```

## Integration with Scripts

The CLI tool is perfect for automation scripts:

```bash
#!/bin/bash

SECRET="your-secret-key"

# Register admin user
echo "Registering admin user..."
authkit user register \
  --secret "$SECRET" \
  --email "admin@company.com" \
  --password "adminPassword123" \
  --name "System Admin" \
  --role "admin"

# Generate API token for service
echo "Generating API token..."
TOKEN=$(authkit token generate \
  --secret "$SECRET" \
  --user-id "service-account" \
  --expiry "720h" \
  --output json | jq -r '.token')

echo "Service token: $TOKEN"

# Start development server
echo "Starting server..."
authkit server start --secret "$SECRET" --port 3000
```

## Troubleshooting

### Common Issues

1. **Missing secret key**
   ```bash
   Error: required flag(s) "secret" not set
   ```
   Solution: Always provide `--secret` or set `JWT_SECRET` environment variable.

2. **Invalid token format**
   ```bash
   Error: Invalid token format
   ```
   Solution: Ensure the token starts with `eyJ` and is a valid JWT.

3. **User already exists**
   ```bash
   Error: user already exists
   ```
   Solution: Use a different email or delete the existing user first.

4. **Token expired**
   ```bash
   Error: token expired
   ```
   Solution: Generate a new token or use refresh token functionality.

### Debug Mode

Enable verbose logging:

```bash
authkit --verbose user login --email user@example.com --password pass123
```

## Building from Source

```bash
# Clone repository
git clone https://github.com/your-username/go-authkit.git
cd go-authkit

# Install dependencies
go mod tidy

# Build CLI tool
go build -o authkit cmd/main.go

# Install globally
go install ./cmd
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Submit a pull request

## License

MIT License - see LICENSE file for details.
