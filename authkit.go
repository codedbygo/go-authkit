package authkit

import (
	//"errors"
	"sync"
	"time"

	"github.com/google/uuid"
)

// New creates a new AuthKit instance with the given configuration
func New(config Config) *AuthKit {
	// Set default values
	if config.BCryptCost == 0 {
		config.BCryptCost = 12
	}
	if config.TokenExpiry == "" {
		config.TokenExpiry = "24h"
	}
	if config.RefreshExpiry == "" {
		config.RefreshExpiry = "7d"
	}
	if config.RateLimitRPM == 0 {
		config.RateLimitRPM = 60
	}

	return &AuthKit{
		config: config,
		users:  make(map[string]*User),
		mutex:  sync.RWMutex{},
	}
}

// RegisterUser registers a new user
func (a *AuthKit) RegisterUser(req RegisterRequest) (*UserInfo, error) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	// Check if user already exists
	for _, user := range a.users {
		if user.Email == req.Email {
			return nil, ErrUserAlreadyExists
		}
	}

	// Hash password
	hashedPassword, err := a.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Create user
	userID := uuid.New().String()
	user := &User{
		ID:            userID,
		Email:         req.Email,
		Password:      hashedPassword,
		Name:          req.Name,
		Role:          req.Role,
		Permissions:   []string{},
		EmailVerified: !a.config.EmailRequired,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		Metadata:      req.Metadata,
	}

	// Set default role if not provided
	if user.Role == "" {
		user.Role = "user"
	}

	// Store user
	a.users[userID] = user

	return a.userToUserInfo(user), nil
}

// LoginUser authenticates a user and returns tokens
func (a *AuthKit) LoginUser(email, password string) (*TokenResponse, error) {
	a.mutex.RLock()
	defer a.mutex.RUnlock()

	// Find user by email
	var user *User
	for _, u := range a.users {
		if u.Email == email {
			user = u
			break
		}
	}

	if user == nil {
		return nil, ErrUserNotFound
	}

	// Check password
	if !a.ComparePassword(user.Password, password) {
		return nil, ErrInvalidPassword
	}

	// Generate tokens
	accessToken, err := a.GenerateAccessToken(user)
	if err != nil {
		return nil, err
	}

	refreshToken, err := a.GenerateRefreshToken(user)
	if err != nil {
		return nil, err
	}

	// Parse expiry duration
	duration, _ := time.ParseDuration(a.config.TokenExpiry)
	expiresIn := int64(duration.Seconds())

	return &TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    expiresIn,
		User:         a.userToUserInfo(user),
	}, nil
}

// GetUserByID retrieves a user by their ID
func (a *AuthKit) GetUserByID(userID string) (*User, error) {
	a.mutex.RLock()
	defer a.mutex.RUnlock()

	user, exists := a.users[userID]
	if !exists {
		return nil, ErrUserNotFound
	}

	return user, nil
}

// GetUserByEmail retrieves a user by their email
func (a *AuthKit) GetUserByEmail(email string) (*User, error) {
	a.mutex.RLock()
	defer a.mutex.RUnlock()

	for _, user := range a.users {
		if user.Email == email {
			return user, nil
		}
	}

	return nil, ErrUserNotFound
}

// UpdateUser updates user information
func (a *AuthKit) UpdateUser(userID string, updates map[string]interface{}) (*UserInfo, error) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	user, exists := a.users[userID]
	if !exists {
		return nil, ErrUserNotFound
	}

	// Update fields
	if name, ok := updates["name"].(string); ok {
		user.Name = name
	}
	if role, ok := updates["role"].(string); ok {
		user.Role = role
	}
	if permissions, ok := updates["permissions"].([]string); ok {
		user.Permissions = permissions
	}
	if metadata, ok := updates["metadata"].(map[string]interface{}); ok {
		user.Metadata = metadata
	}

	user.UpdatedAt = time.Now()

	return a.userToUserInfo(user), nil
}

// DeleteUser removes a user from the system
func (a *AuthKit) DeleteUser(userID string) error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	if _, exists := a.users[userID]; !exists {
		return ErrUserNotFound
	}

	delete(a.users, userID)
	return nil
}

// ListUsers returns all users (for admin purposes)
func (a *AuthKit) ListUsers() []*UserInfo {
	a.mutex.RLock()
	defer a.mutex.RUnlock()

	users := make([]*UserInfo, 0, len(a.users))
	for _, user := range a.users {
		users = append(users, a.userToUserInfo(user))
	}

	return users
}

// userToUserInfo converts User to UserInfo (without password)
func (a *AuthKit) userToUserInfo(user *User) *UserInfo {
	return &UserInfo{
		ID:            user.ID,
		Email:         user.Email,
		Name:          user.Name,
		Role:          user.Role,
		Permissions:   user.Permissions,
		EmailVerified: user.EmailVerified,
		Metadata:      user.Metadata,
	}
}
