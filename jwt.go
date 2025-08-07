package authkit

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// GenerateAccessToken generates a JWT access token for the user
func (a *AuthKit) GenerateAccessToken(user *User) (string, error) {
	duration, err := time.ParseDuration(a.config.TokenExpiry)
	if err != nil {
		duration = 24 * time.Hour // default to 24 hours
	}

	claims := &Claims{
		UserID:      user.ID,
		Email:       user.Email,
		Role:        user.Role,
		Permissions: user.Permissions,
		Metadata:    user.Metadata,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.New().String(), // Add unique JTI (JWT ID)
			Subject:   user.ID,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "authkit",
			Audience:  []string{"authkit-users"},
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(a.config.JWTSecret))
}

// GenerateRefreshToken generates a JWT refresh token
func (a *AuthKit) GenerateRefreshToken(user *User) (string, error) {
	duration, err := time.ParseDuration(a.config.RefreshExpiry)
	if err != nil {
		duration = 7 * 24 * time.Hour // default to 7 days
	}

	claims := &jwt.RegisteredClaims{
		ID:        uuid.New().String(), // Add unique JTI (JWT ID)
		Subject:   user.ID,
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		NotBefore: jwt.NewNumericDate(time.Now()),
		Issuer:    "authkit-refresh",
		Audience:  []string{"authkit-refresh"},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(a.config.JWTSecret))
}

// ValidateToken validates and parses a JWT token
func (a *AuthKit) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(a.config.JWTSecret), nil
	})

	if err != nil {
		return nil, ErrInvalidToken
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrInvalidToken
}

// RefreshToken validates a refresh token and generates new access token
func (a *AuthKit) RefreshToken(refreshTokenString string) (*TokenResponse, error) {
	// Parse the refresh token
	token, err := jwt.ParseWithClaims(refreshTokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(a.config.JWTSecret), nil
	})

	if err != nil {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	// Get user from claims
	user, err := a.GetUserByID(claims.Subject)
	if err != nil {
		return nil, err
	}

	// Generate new tokens
	accessToken, err := a.GenerateAccessToken(user)
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := a.GenerateRefreshToken(user)
	if err != nil {
		return nil, err
	}

	// Parse expiry duration
	duration, _ := time.ParseDuration(a.config.TokenExpiry)
	expiresIn := int64(duration.Seconds())

	return &TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    expiresIn,
		User:         a.userToUserInfo(user),
	}, nil
}

// GenerateCustomToken generates a token with custom claims
func (a *AuthKit) GenerateCustomToken(userID string, customClaims map[string]interface{}, expiry time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"jti":     uuid.New().String(), // Add unique JTI
		"user_id": userID,
		"iss":     "authkit",
		"aud":     "authkit-users",
		"iat":     time.Now().Unix(),
		"exp":     time.Now().Add(expiry).Unix(),
		"nbf":     time.Now().Unix(),
	}

	// Add custom claims
	for key, value := range customClaims {
		claims[key] = value
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(a.config.JWTSecret))
}
