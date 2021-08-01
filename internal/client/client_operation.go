package client

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"

	"google.golang.org/grpc"
	pb "github.com/bartmika/tstorage-server/proto"

    "github.com/bartmika/tbridge-server/internal/constants"
	"github.com/bartmika/tbridge-server/internal/dtos"
)

type ClientOperation struct {
	port              int
	remoteAddr        string
	grpcServer        *grpc.Server
	accessToken       string
	refreshToken      string
}

func New(port int, remoteAddr string, accessToken string, refreshToken string) *ClientOperation {
	return &ClientOperation{
		port:              port,
		remoteAddr:        remoteAddr,
		grpcServer:        nil,
		accessToken:       accessToken,
		refreshToken:      refreshToken,
	}
}

// Function will consume the main runtime loop and run the business logic
// of the application.
func (s *ClientOperation) RunMainRuntimeLoop() {
    // Before we begin, we need to connect to our remote bridge server over HTTP.
	// To begin this process let's make a simple call to the `version` API enpdoint.
	log.Printf("Connecting to remote bridge...")
	r, err := http.Get(s.remoteAddr + constants.GetVersionAPIEndpointPath)
	defer r.Body.Close()
	if err != nil {
		log.Fatalln(err)
	}

	var versionDTO *dtos.VersionDTO
	if err := json.NewDecoder(r.Body).Decode(&versionDTO); err != nil {
		log.Fatal(err)
	}
	if versionDTO.Value != "v1.0" {
		log.Fatal("This client does not support the bridge at the specified remote address.")
	}
	log.Printf("Successfully connected to remote bridge.")

	// Open a TCP server to the specified localhost and environment variable
	// specified port number.
	log.Printf("Starting gRPC server...")
	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", s.port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Initialize our gRPC server using our TCP server.
	grpcServer := grpc.NewServer()

	// Save reference to our application state.
	s.grpcServer = grpcServer

	// For debugging purposes only.
	log.Printf("gRPC server is running on port %v", s.port)
	log.Printf("Bridge is now running in client mode and is successfully connected to a remote bridge at address %v", s.remoteAddr)

	// Block the main runtime loop for accepting and processing gRPC requests.
	pb.RegisterTStorageServer(grpcServer, &ClientOperationImpl{
		// DEVELOPERS NOTE:
		// We want to attach to every gRPC call the following variables...
		remoteAddr:       s.remoteAddr,
		accessToken:       s.accessToken,
		refreshToken:      s.refreshToken,
	})
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// Function will tell the application to stop the main runtime loop when
// the process has been finished.
func (s *ClientOperation) StopMainRuntimeLoop() {
	log.Printf("Starting graceful shutdown now...")

	// Finish any RPC communication taking place at the moment before
	// shutting down the gRPC server.
	s.grpcServer.GracefulStop()
}
