package cmd

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/bartmika/tbridge-server/internal/client"
)

var (
	accessToken string
	refreshToken string
)

func init() {
	// The following are optional and will have defaults placed when missing.
	clientModeCmd.Flags().IntVarP(&port, "port", "p", 50053, "The port to run this bridge server on in client mode. Recommended: 50053.")
	clientModeCmd.Flags().StringVarP(&accessToken, "access_token", "a", os.Getenv("TBRIDGE_CLIENT_ACCESS_TOKEN"), "The access token to authenticate requests.")
	clientModeCmd.Flags().StringVarP(&refreshToken, "refresh_token", "r", os.Getenv("TBRIDGE_CLIENT_REFRESH_TOKEN"), "The refresh token to authenticate requests.")

	// The following fields are required.
    clientModeCmd.Flags().StringVarP(&remoteAddr, "remote_addr", "o", "", "The remote address of our other bridge running in `server mode`.")
    clientModeCmd.MarkFlagRequired("remote_addr")

	// Make this sub-command part of our application.
	rootCmd.AddCommand(clientModeCmd)
}

func doClientMode() {
	// Setup our client.
	c := client.New(port, remoteAddr, accessToken, refreshToken)

	// DEVELOPERS CODE:
	// The following code will create an anonymous goroutine which will have a
	// blocking chan `sigs`. This blocking chan will only unblock when the
	// golang app receives a termination command; therfore the anyomous
	// goroutine will run and terminate our running application.
	//
	// Special Thanks:
	// (1) https://gobyexample.com/signals
	// (2) https://guzalexander.com/2017/05/31/gracefully-exit-server-in-go.html
	//
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs // Block execution until signal from terminal gets triggered here.
		c.StopMainRuntimeLoop()
	}()
	c.RunMainRuntimeLoop()
}

var clientModeCmd = &cobra.Command{
	Use:   "client_mode",
	Short: "Run the bridge in client mode",
	Long:  `Run the bridge in client mode by running a gRPC server to allow other services to access so the bridge will make HTTP requests to a remote tstorage-server.`,
	Run: func(cmd *cobra.Command, args []string) {
		doClientMode()
	},
}
