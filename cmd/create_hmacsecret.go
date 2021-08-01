package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/bartmika/tbridge-server/internal/utils"
)

var (
	hmacSecretLength int
)

func init() {
	createHmacSecretCmd.Flags().IntVarP(&hmacSecretLength, "hmac_secret_length", "s", 51, "The HMAC secret length.")
	rootCmd.AddCommand(createHmacSecretCmd)
}

var createHmacSecretCmd = &cobra.Command{
	Use:   "create_hmacsecret",
	Short: "Print a random HMAC secret",
	Long:  `Print a random HMAC secret value you can use.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Run in your console:\n")
		fmt.Println("export TBRIDGE_SERVER_HMAC_SECRET="+utils.NewRandomString(hmacSecretLength))
		fmt.Println()
	},
}
