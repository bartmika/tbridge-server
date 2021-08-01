package server

import (
	"context"
	"encoding/json"
	"log"
	// "time"
	"net/http"
	"strings"

	tstorage_pb "github.com/bartmika/tstorage-server/proto"
	"google.golang.org/grpc"
	// "google.golang.org/grpc/credentials"
	tspb "github.com/golang/protobuf/ptypes/timestamp"
    pb "github.com/bartmika/tstorage-server/proto"

	"github.com/bartmika/tbridge-server/internal/dtos"
)

type ServerOperationImpl struct {
	storageAddr    string
	tstorageConn   *grpc.ClientConn
	tstorageClient tstorage_pb.TStorageClient
	hmacSecret     []byte
	sessionUuid    string
}

func (c *ServerOperationImpl) HandleRequests(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

    // Split path into slash-separated parts, for example, path "/foo/bar"
    // gives p==["foo", "bar"] and path "/" gives p==[""]. Our API starts with
    // "/api/v1", as a result we will start the array slice at "3".
    p := strings.Split(r.URL.Path, "/")[3:]
    n := len(p)

    // fmt.Println(p, n) // For debugging purposes only.

    switch {
    case n == 1 && p[0] == "version" && r.Method == http.MethodGet:
        c.getVersion(w, r)
	case n == 1 && p[0] == "time-series-data" && r.Method == http.MethodGet:
        c.getTimeSeriesData(w, r)
	case n == 1 && p[0] == "time-series-data" && r.Method == http.MethodPost:
        c.postTimeSeriesData(w, r)
	// case n == 2 && p[0] == "time-series-datum" && r.Method == http.MethodDelete:
    //     c.deleteTimeSeriesDatum(w, r, p[1])
    default:
        http.NotFound(w, r)
    }
}


func (c *ServerOperationImpl) getVersion(w http.ResponseWriter, req *http.Request) {
    w.Header().Set("Content-Type", "application/json")

	versionDTO := &dtos.VersionDTO{Value:"v1.0"}
	if err := json.NewEncoder(w).Encode(&versionDTO); err != nil { // [2]
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (c *ServerOperationImpl) getTimeSeriesData(w http.ResponseWriter, req *http.Request) {
    w.Write([]byte("TODO: List Time Series Data")) //TODO: IMPLEMENT.
}

func (c *ServerOperationImpl) postTimeSeriesData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

    // Take our raw data from our request and serialize it into our golang struct.
	var requestData *dtos.TimeSeriesData
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Generate our `labels`, this is part of the gRPC service definition requirements.
	labels := []*pb.Label{}
	labels = append(labels, &pb.Label{Name: "Source", Value: "Command"})

	// DEVELOPERS NOTE:
	// To stream from a client to a server using gRPC, the following documentation
	// will help explain how it works. Please visit it if the code below does
	// not make any sense.
	// https://grpc.io/docs/languages/go/basics/#client-side-streaming-rpc-1
	stream, err := c.tstorageClient.InsertRows(context.Background())
	if err != nil {
		log.Fatalf("%v.InsertRows(_) = _, %v", c.tstorageClient, err)
	}

    // Iterate through all the data points that were sent through the API call
	// and stream them to our gRPC server.
	for _, dataPoint := range requestData.Data {
		ts := &tspb.Timestamp{
			Seconds: dataPoint.Timestamp.Unix(),
			Nanos:   0,
		}
		tsd := &pb.TimeSeriesDatum{
			Labels: labels,
			Metric: requestData.Metric,
			Value: dataPoint.Value,
			Timestamp: ts,
		}
		if err := stream.Send(tsd); err != nil {
			log.Fatalf("%v.Send(%v) = %v", stream, tsd, err)
		}
	}

	_, err = stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("%v.CloseAndRecv() got error %v, want %v", stream, err, nil)
	}
	log.Printf("Successfully inserted")
}

// func (c *ServerOperationImpl) deleteTimeSeriesDatum(w http.ResponseWriter, req *http.Request, uuid string) {
//     w.Write([]byte("TODO: Delete Series Datum with UUID: " + uuid)) //TODO: IMPLEMENT.
// }
