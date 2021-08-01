package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/bartmika/tbridge-server/internal/utils"
)

func init() {
	// Load up our sub-command.
	rootCmd.AddCommand(createTokenCmd)
}

var createTokenCmd = &cobra.Command{
	Use:   "create_token",
	Short: "Generates a JWT token pair",
	Long:  `Generates a JWT token pair`,
	Run: func(cmd *cobra.Command, args []string) {
		doCreateToken()
	},
}


func doCreateToken() {
	sessionExpiryTime := time.Hour * 24 * 7 // 1 week

	b := []byte(hmacSecret)
	accessToken, refreshToken, err := utils.GenerateJWTTokenPair(b, sessionUuid, sessionExpiryTime)

	fmt.Println("Run in your console:")
	fmt.Println("\nexport TBRIDGE_CLIENT_ACCESS_TOKEN="+accessToken)
	fmt.Println("\nexport TBRIDGE_CLIENT_REFRESH_TOKEN="+refreshToken)
	fmt.Println("")

	if err != nil { // For debugging purposes.
		fmt.Println("Error:\n", err, "\n")
	}
}
