package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	port        int
	storageAddr string
	remoteAddr  string
	hmacSecret  string
	sessionUuid string
)

// Initialize function will be called when every command gets called.
func init() {
	// Get our environment variables which will used to configure our application and save across all the sub-commands.
	rootCmd.PersistentFlags().StringVar(&hmacSecret, "hmacSecret", os.Getenv("TBRIDGE_SERVER_HMAC_SECRET"), "The bridge signing key.")
	rootCmd.PersistentFlags().StringVar(&sessionUuid, "sessionUuid", os.Getenv("TBRIDGE_SERVER_SESSION_UUID"), "The bridge session uuid.")
}

var rootCmd = &cobra.Command{
	Use:   "tbridge-server",
	Short: "Bridge time-series data across a network",
	Long:  `Bridge time-series data across a network using HTTP in client and serve mode.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Do nothing...
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
