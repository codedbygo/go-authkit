package authkit

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes a password using bcrypt
func (a *AuthKit) HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), a.config.BCryptCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// ComparePassword compares a hashed password with a plaintext password
func (a *AuthKit) ComparePassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

// HashPasswordStatic is a static method for hashing passwords without AuthKit instance
func HashPasswordStatic(password string, cost int) (string, error) {
	if cost == 0 {
		cost = 12 // default cost
	}
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// ComparePasswordStatic is a static method for comparing passwords without AuthKit instance
func ComparePasswordStatic(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
