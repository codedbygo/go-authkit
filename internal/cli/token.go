package cli

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/your-username/go-authkit"
)

var tokenCmd = &cobra.Command{
	Use:   "token",
	Short: "Token management commands",
	Long:  "Commands for generating, validating, and refreshing JWT tokens",
}

var tokenGenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate a JWT token",
	Long:  "Generate a JWT token for a user with optional custom claims",
	Run:   runTokenGenerate,
}

var tokenValidateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate a JWT token",
	Long:  "Validate and decode a JWT token",
	Run:   runTokenValidate,
}

var tokenRefreshCmd = &cobra.Command{
	Use:   "refresh",
	Short: "Refresh an access token",
	Long:  "Refresh an access token using a refresh token",
	Run:   runTokenRefresh,
}

// Flags for token commands
var (
	tokenString  string
	tokenUserID  string
	tokenExpiry  string
	refreshToken string
	customClaims []string
)

func init() {
	// Add token command to root
	rootCmd.AddCommand(tokenCmd)

	// Add subcommands to token
	tokenCmd.AddCommand(tokenGenerateCmd)
	tokenCmd.AddCommand(tokenValidateCmd)
	tokenCmd.AddCommand(tokenRefreshCmd)

	// Generate flags
	tokenGenerateCmd.Flags().StringVarP(&tokenUserID, "user-id", "u", "", "User ID (required)")
	tokenGenerateCmd.Flags().StringVarP(&tokenExpiry, "expiry", "x", "24h", "Token expiry duration")
	tokenGenerateCmd.Flags().StringSliceVarP(&customClaims, "claims", "c", []string{}, "Custom claims (key=value format)")
	tokenGenerateCmd.MarkFlagRequired("user-id")

	// Validate flags
	tokenValidateCmd.Flags().StringVarP(&tokenString, "token", "t", "", "JWT token to validate (required)")
	tokenValidateCmd.MarkFlagRequired("token")

	// Refresh flags
	tokenRefreshCmd.Flags().StringVarP(&refreshToken, "refresh-token", "r", "", "Refresh token (required)")
	tokenRefreshCmd.MarkFlagRequired("refresh-token")
}

func runTokenGenerate(cmd *cobra.Command, args []string) {
	auth := authkit.New(authkit.Config{
		JWTSecret:   secretKey,
		TokenExpiry: tokenExpiry,
		BCryptCost:  12,
	})

	// Parse expiry duration
	duration, err := time.ParseDuration(tokenExpiry)
	checkError(err)

	// Parse custom claims
	claims := make(map[string]interface{})
	for _, claim := range customClaims {
		// Simple key=value parsing (could be enhanced)
		fmt.Printf("Custom claim: %s\n", claim)
		// For demo, just add as string
		claims[claim] = "custom-value"
	}

	token, err := auth.GenerateCustomToken(tokenUserID, claims, duration)
	checkError(err)

	fmt.Printf("Token generated successfully!\n")
	printOutput(map[string]interface{}{
		"token":   token,
		"user_id": tokenUserID,
		"expiry":  tokenExpiry,
		"claims":  claims,
	})
}

func runTokenValidate(cmd *cobra.Command, args []string) {
	auth := authkit.New(authkit.Config{
		JWTSecret:   secretKey,
		TokenExpiry: "24h",
		BCryptCost:  12,
	})

	claims, err := auth.ValidateToken(tokenString)
	if err != nil {
		fmt.Printf("Token validation failed: %v\n", err)
		return
	}

	fmt.Printf("Token is valid!\n")
	printOutput(map[string]interface{}{
		"valid":       true,
		"user_id":     claims.UserID,
		"email":       claims.Email,
		"role":        claims.Role,
		"permissions": claims.Permissions,
		"issued_at":   claims.IssuedAt,
		"expires_at":  claims.ExpiresAt,
	})
}

func runTokenRefresh(cmd *cobra.Command, args []string) {
	auth := authkit.New(authkit.Config{
		JWTSecret:     secretKey,
		TokenExpiry:   "24h",
		RefreshExpiry: "7d",
		BCryptCost:    12,
	})

	newTokens, err := auth.RefreshToken(refreshToken)
	checkError(err)

	fmt.Printf("Token refreshed successfully!\n")
	printOutput(map[string]interface{}{
		"access_token":  newTokens.AccessToken,
		"refresh_token": newTokens.RefreshToken,
		"token_type":    newTokens.TokenType,
		"expires_in":    newTokens.ExpiresIn,
		"user":          newTokens.User,
	})
}
