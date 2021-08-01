package cmd

import (
	// "context"
	// "fmt"
	"log"
	"time"
	"encoding/json"
	"io/ioutil"
    "net/http"
	"bytes"

	"github.com/spf13/cobra"
	// "google.golang.org/grpc"
	// // "google.golang.org/grpc/credentials"
	//
	// tspb "github.com/golang/protobuf/ptypes/timestamp"
	//
	// pb "github.com/bartmika/tstorage-server/proto"
	"github.com/bartmika/tbridge-server/internal/constants"
	"github.com/bartmika/tbridge-server/internal/dtos"
)

var (
	metric string
	value  float64
	tsv    int64
)

func init() {
	// The following are required.
	insertRowCmd.Flags().StringVarP(&metric, "metric", "m", "", "The metric to attach to the TSD.")
	insertRowCmd.MarkFlagRequired("metric")
	insertRowCmd.Flags().Float64VarP(&value, "value", "v", 0.00, "The value to attach to the TSD.")
	insertRowCmd.MarkFlagRequired("value")
	insertRowCmd.Flags().Int64VarP(&tsv, "timestamp", "t", 0, "The timestamp to attach to the TSD.")
	insertRowCmd.MarkFlagRequired("timestamp")
	// The following fields are required.
    insertRowCmd.Flags().StringVarP(&remoteAddr, "remote_addr", "o", "", "The remote address of our other bridge running in `server mode`.")
    insertRowCmd.MarkFlagRequired("remote_addr")

	// The following are optional and will have defaults placed when missing.
	// None...
	rootCmd.AddCommand(insertRowCmd)
}

func doInsertRow() {
    tm := time.Unix(tsv, 0)

	labels := []*dtos.Label{}
    dataPoints := []*dtos.DataPoint{}
	dataPoints = append(dataPoints, &dtos.DataPoint{
		Value: value,
		Timestamp: tm,
	})
	reqData := &dtos.TimeSeriesData{
		Metric: metric,
		Labels: labels,
		Data: dataPoints,
	}
	postBody, _ := json.Marshal(reqData)
	responseBody := bytes.NewBuffer(postBody)
	resp, err := http.Post(remoteAddr + constants.ListOrCreateTimeSeriesDataAPIEndpointPath, "application/json", responseBody)
	if err != nil {
	    log.Fatalln(err)
	}
	//We Read the response body on the line below.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
	    log.Fatalln(err)
	}
	//Convert the body to type string
	sb := string(body)
	log.Printf(sb)

	log.Printf("Successfully inserted")
}

var insertRowCmd = &cobra.Command{
	Use:   "insert_row",
	Short: "Insert single datum",
	Long:  `Connect to the HTTP server and sends a single time-series datum.`,
	Run: func(cmd *cobra.Command, args []string) {
		doInsertRow()
	},
}
