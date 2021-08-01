package cmd

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	serv "github.com/bartmika/tbridge-server/internal/server"
)

func init() {
	// The following are required.
	serverModeCmd.Flags().StringVarP(&storageAddr, "storage_addr", "o", "localhost:50051", "The time-series data storage gRPC server address.")
	serverModeCmd.MarkFlagRequired("storage_addr")

	// The following are optional and will have defaults placed when missing.
	serverModeCmd.Flags().IntVarP(&port, "port", "p", 50054, "The port to run this bridge server on in server mode. Default: 50053.")

	// Make this sub-command part of our application.
	rootCmd.AddCommand(serverModeCmd)
}

func doServerMode() {
	log.Print("Starting server mode...")

	// Setup our server.
	server, err := serv.New(storageAddr, hmacSecret, sessionUuid)
	if err != nil {
		log.Fatal(err)
	}

	done := make(chan os.Signal, 1)
    signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

    go server.RunMainRuntimeLoop()

    log.Printf("Bridge is now running in server mode on port %v.", port)

    // Run the main loop blocking code.
    <-done

    server.StopMainRuntimeLoop()
}

var serverModeCmd = &cobra.Command{
	Use:   "server_mode",
	Short: "Run the bridge in server mode",
	Long:  `Run the bridge in server mode by running a gRPC server to allow other services to access so the bridge will make HTTP requests to a remote tstorage-server.`,
	Run: func(cmd *cobra.Command, args []string) {
		doServerMode()
	},
}
