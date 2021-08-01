package client

import (
	"context"
	"io"
    "net/http"
	"time"
	"encoding/json"
	"bytes"
	"log"
	"io/ioutil"

	"github.com/golang/protobuf/ptypes/empty"
	// tspb "github.com/golang/protobuf/ptypes/timestamp"
	"github.com/nakabonne/tstorage"
	pb "github.com/bartmika/tstorage-server/proto"

	"github.com/bartmika/tbridge-server/internal/dtos"
	"github.com/bartmika/tbridge-server/internal/constants"
)

type ClientOperationImpl struct {
	storage tstorage.Storage
	remoteAddr        string
	accessToken       string
	refreshToken      string
	pb.TStorageServer
}

func (s *ClientOperationImpl) InsertRow(ctx context.Context, in *pb.TimeSeriesDatum) (*empty.Empty, error) {
    panic("Not implemented yet.")
	return &empty.Empty{}, nil
}

func (s *ClientOperationImpl) InsertRows(stream pb.TStorage_InsertRowsServer) error {
	// DEVELOPERS NOTE:
	// If you don't understand how server side streaming works using gRPC then
	// please visit the documentation to get an understanding:
	// https://grpc.io/docs/languages/go/basics/#server-side-streaming-rpc-1

	// Wait and receieve the stream from the client.
	for {
		datum, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&empty.Empty{})
		}
		if err != nil {
			return err
		}
		go s.processDatum(datum)
	}

	return nil
}

func (s *ClientOperationImpl) processDatum(datum *pb.TimeSeriesDatum) {
	// Generate our labels, if there are any.
	labels := []*dtos.Label{}
	for _, label := range datum.Labels {
		labels = append(labels, &dtos.Label{Name: label.Name, Value: label.Value})
	}

	// Generate our datapoint.
	tm := time.Unix(datum.Timestamp.Seconds, 0)
	dataPoint := &dtos.DataPoint{Timestamp: tm, Value: datum.Value}


	dataPoints := []*dtos.DataPoint{}

	dataPoints = append(dataPoints, dataPoint)
	postData := &dtos.TimeSeriesData{
		Metric: datum.Metric,
		Labels: labels,
		Data: dataPoints,
	}

	s.postTimeSeriesData(postData)
}

func (s *ClientOperationImpl) postTimeSeriesData(postData *dtos.TimeSeriesData) {
	postBodyBin, _ := json.Marshal(postData)
	postBodyBuff := bytes.NewBuffer(postBodyBin)

	// Create a Bearer string by appending string access token
    var bearer = "Bearer " + s.accessToken

	// Create a new request using http
	url := s.remoteAddr + constants.ListOrCreateTimeSeriesDataAPIEndpointPath
    req, err := http.NewRequest("POST", url, postBodyBuff)
	if err != nil {
        log.Println("Error on http.NewReques.\n[ERROR] -", err)
    }

	// add authorization header to the req
    req.Header.Add("Authorization", bearer)
    req.Header.Add("Accept", "application/json")

	// Send req using http Client
    client := &http.Client{}
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		// Special thanks for this block of code from stackoverflow via 
		// https://stackoverflow.com/a/58493142
        for key, val := range via[0].Header {
            req.Header[key] = val
        }
        return err
    }

    resp, err := client.Do(req)
    if err != nil {
        log.Println("Error on response.\n[ERROR] -", err)
    }
    defer resp.Body.Close()


	// resp, err := http.Post(s.remoteAddr + constants.ListOrCreateTimeSeriesDataAPIEndpointPath, "application/json", postBodyBuff)
	// if err != nil {
	//     log.Fatalln(err)
	// }

	//We Read the response body on the line below.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
	    log.Fatalln(err)
	}

	//Convert the body to type string
	sb := string(body)
	log.Println(sb)
}

func (s *ClientOperationImpl) Select(in *pb.Filter, stream pb.TStorage_SelectServer) error {
	// // Generate our labels, if there are any.
	// labels := []tstorage.Label{}
	// for _, label := range in.Labels {
	// 	labels = append(labels, tstorage.Label{Name: label.Name, Value: label.Value})
	// }
	//
	// points, err := s.storage.Select(in.Metric, labels, in.Start.Seconds, in.End.Seconds)
	// if err != nil {
	// 	return err
	// }
	//
	// for _, point := range points {
	// 	ts := &tspb.Timestamp{
	// 		Seconds: point.Timestamp,
	// 		Nanos:   0,
	// 	}
	// 	dataPoint := &pb.DataPoint{Value: point.Value, Timestamp: ts}
	// 	if err := stream.Send(dataPoint); err != nil {
	// 		return err
	// 	}
	// }

	panic("Not implemented yet.")
	return nil
}
