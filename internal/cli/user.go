package cli

import (
	"fmt"

	"github.com/codedbygo/go-authkit"
	"github.com/spf13/cobra"
)

var userCmd = &cobra.Command{
	Use:   "user",
	Short: "User management commands",
	Long:  "Commands for managing users: register, login, list, update, delete",
}

var userRegisterCmd = &cobra.Command{
	Use:   "register",
	Short: "Register a new user",
	Long:  "Register a new user with email, password, and optional metadata",
	Run:   runUserRegister,
}

var userLoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login a user",
	Long:  "Authenticate a user and get access tokens",
	Run:   runUserLogin,
}

var userListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all users",
	Long:  "List all registered users in the system",
	Run:   runUserList,
}

var userDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a user",
	Long:  "Delete a user by ID",
	Run:   runUserDelete,
}

// Flags for user commands
var (
	userEmail    string
	userPassword string
	userName     string
	userRole     string
	userID       string
)

func init() {
	// Add user command to root
	rootCmd.AddCommand(userCmd)

	// Add subcommands to user
	userCmd.AddCommand(userRegisterCmd)
	userCmd.AddCommand(userLoginCmd)
	userCmd.AddCommand(userListCmd)
	userCmd.AddCommand(userDeleteCmd)

	// Register flags
	userRegisterCmd.Flags().StringVarP(&userEmail, "email", "e", "", "User email (required)")
	userRegisterCmd.Flags().StringVarP(&userPassword, "password", "p", "", "User password (required)")
	userRegisterCmd.Flags().StringVarP(&userName, "name", "n", "", "User name (required)")
	userRegisterCmd.Flags().StringVarP(&userRole, "role", "r", "user", "User role")
	userRegisterCmd.MarkFlagRequired("email")
	userRegisterCmd.MarkFlagRequired("password")
	userRegisterCmd.MarkFlagRequired("name")

	// Login flags
	userLoginCmd.Flags().StringVarP(&userEmail, "email", "e", "", "User email (required)")
	userLoginCmd.Flags().StringVarP(&userPassword, "password", "p", "", "User password (required)")
	userLoginCmd.MarkFlagRequired("email")
	userLoginCmd.MarkFlagRequired("password")

	// Delete flags
	userDeleteCmd.Flags().StringVarP(&userID, "id", "i", "", "User ID (required)")
	userDeleteCmd.MarkFlagRequired("id")
}

func runUserRegister(cmd *cobra.Command, args []string) {
	auth := authkit.New(authkit.Config{
		JWTSecret:   secretKey,
		TokenExpiry: "24h",
		BCryptCost:  12,
	})

	req := authkit.RegisterRequest{
		Email:    userEmail,
		Password: userPassword,
		Name:     userName,
		Role:     userRole,
	}

	user, err := auth.RegisterUser(req)
	checkError(err)

	fmt.Printf("User registered successfully!\n")
	printOutput(map[string]interface{}{
		"user_id": user.ID,
		"email":   user.Email,
		"name":    user.Name,
		"role":    user.Role,
	})
}

func runUserLogin(cmd *cobra.Command, args []string) {
	auth := authkit.New(authkit.Config{
		JWTSecret:   secretKey,
		TokenExpiry: "24h",
		BCryptCost:  12,
	})

	tokenResponse, err := auth.LoginUser(userEmail, userPassword)
	checkError(err)

	fmt.Printf("Login successful!\n")
	printOutput(map[string]interface{}{
		"access_token":  tokenResponse.AccessToken,
		"refresh_token": tokenResponse.RefreshToken,
		"token_type":    tokenResponse.TokenType,
		"expires_in":    tokenResponse.ExpiresIn,
		"user":          tokenResponse.User,
	})
}

func runUserList(cmd *cobra.Command, args []string) {
	auth := authkit.New(authkit.Config{
		JWTSecret:   secretKey,
		TokenExpiry: "24h",
		BCryptCost:  12,
	})

	users := auth.ListUsers()

	fmt.Printf("Found %d users:\n", len(users))
	printOutput(map[string]interface{}{
		"count": len(users),
		"users": users,
	})
}

func runUserDelete(cmd *cobra.Command, args []string) {
	auth := authkit.New(authkit.Config{
		JWTSecret:   secretKey,
		TokenExpiry: "24h",
		BCryptCost:  12,
	})

	err := auth.DeleteUser(userID)
	checkError(err)

	fmt.Printf("User deleted successfully!\n")
	printOutput(map[string]interface{}{
		"message": "User deleted",
		"user_id": userID,
	})
}
