package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/google/uuid"
)

func init() {
	rootCmd.AddCommand(createUuidCmd)
}

var createUuidCmd = &cobra.Command{
	Use:   "create_uuid",
	Short: "Print a random uuid",
	Long:  `Print a random uuid value you can use.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Run in your console:")
		fmt.Println("\nexport TBRIDGE_SERVER_SESSION_UUID="+uuid.NewString())
		fmt.Println()
	},
}
