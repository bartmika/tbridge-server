package server

import (
	"fmt"
	"log"
	"net/http"

	tstorage_pb "github.com/bartmika/tstorage-server/proto"
	"google.golang.org/grpc"
)

type ServerOperation struct {
	storageAddr    string
	tstorageConn   *grpc.ClientConn
	tstorageClient tstorage_pb.TStorageClient
	hmacSecret     []byte
	sessionUuid    string
}

func New(storageAddr string, hmacSecret string, sessionUuid string) (*ServerOperation, error) {
	s := &ServerOperation{
		storageAddr: storageAddr,
	    hmacSecret: []byte(hmacSecret),
		sessionUuid: sessionUuid,
    }

	// Connect to our time-series data storage.
	log.Println("Dialing Storage...")

	// Set up a direct connection to the gRPC server.
	conn, err := grpc.Dial(
		s.storageAddr,
		grpc.WithInsecure(),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, err
	}
	log.Printf("Storage connected via address %v\n", storageAddr)

	// Set up our protocol buffer interface.
	client := tstorage_pb.NewTStorageClient(conn)

	s.tstorageConn = conn
	s.tstorageClient = client

	return s, nil
}

// Function will consume the main runtime loop and run the business logic
// of the application.
func (s *ServerOperation) RunMainRuntimeLoop() {
	defer s.shutdown()

	c := ServerOperationImpl{
		tstorageConn: s.tstorageConn,
		tstorageClient: s.tstorageClient,
		hmacSecret: s.hmacSecret,
		sessionUuid: s.sessionUuid,
	}

    mux := http.NewServeMux()
	mux.HandleFunc("/", c.ChainMiddleware(c.HandleRequests))

    httpServer := &http.Server{
        Addr: fmt.Sprintf("%s:%s", "localhost", "5000"),
        Handler: mux,
    }

    if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
        panic(err)
    }
}

// Function will tell the application to stop the main runtime loop when
// the process has been finished.
func (s *ServerOperation) StopMainRuntimeLoop() {
	s.shutdown()
}

func (s *ServerOperation) shutdown() {
	s.tstorageConn.Close()
}
