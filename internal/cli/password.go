package cli

import (
	"fmt"

	"github.com/codedbygo/go-authkit"
	"github.com/spf13/cobra"
)

var passwordCmd = &cobra.Command{
	Use:   "password",
	Short: "Password utility commands",
	Long:  "Commands for hashing and comparing passwords",
}

var passwordHashCmd = &cobra.Command{
	Use:   "hash",
	Short: "Hash a password",
	Long:  "Hash a password using bcrypt",
	Run:   runPasswordHash,
}

var passwordCompareCmd = &cobra.Command{
	Use:   "compare",
	Short: "Compare a password with its hash",
	Long:  "Compare a plain password with its bcrypt hash",
	Run:   runPasswordCompare,
}

// Flags for password commands
var (
	plainPassword  string
	hashedPassword string
	bcryptCost     int
)

func init() {
	// Add password command to root
	rootCmd.AddCommand(passwordCmd)

	// Add subcommands to password
	passwordCmd.AddCommand(passwordHashCmd)
	passwordCmd.AddCommand(passwordCompareCmd)

	// Hash flags
	passwordHashCmd.Flags().StringVarP(&plainPassword, "password", "p", "", "Password to hash (required)")
	passwordHashCmd.Flags().IntVarP(&bcryptCost, "cost", "c", 12, "BCrypt cost (4-31)")
	passwordHashCmd.MarkFlagRequired("password")

	// Compare flags
	passwordCompareCmd.Flags().StringVarP(&plainPassword, "password", "p", "", "Plain password (required)")
	passwordCompareCmd.Flags().StringVarP(&hashedPassword, "hash", "H", "", "Hashed password (required)")
	passwordCompareCmd.MarkFlagRequired("password")
	passwordCompareCmd.MarkFlagRequired("hash")
}

func runPasswordHash(cmd *cobra.Command, args []string) {
	hashed, err := authkit.HashPasswordStatic(plainPassword, bcryptCost)
	checkError(err)

	fmt.Printf("Password hashed successfully!\n")
	printOutput(map[string]interface{}{
		"original": plainPassword,
		"hashed":   hashed,
		"cost":     bcryptCost,
	})
}

func runPasswordCompare(cmd *cobra.Command, args []string) {
	isValid := authkit.ComparePasswordStatic(hashedPassword, plainPassword)

	if isValid {
		fmt.Printf("Password matches hash!\n")
	} else {
		fmt.Printf("Password does not match hash!\n")
	}

	printOutput(map[string]interface{}{
		"valid":    isValid,
		"password": plainPassword,
		"hash":     hashedPassword,
	})
}
