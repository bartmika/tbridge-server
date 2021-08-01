package dtos

import (
	"time"
)

type TimeSeriesDatum struct {
	Metric string `json:"metric,omitempty"`
	Labels  []*Label `json:"labels,omitempty"`
	Value float64 `json:"value,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}

type TimeSeriesData struct {
	Metric string `json:"metric,omitempty"`
	Labels  []*Label `json:"labels,omitempty"`
	Data   []*DataPoint `json:"data,omitempty"`
}

type DataPoint struct {
	Value float64 `json:"value,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}

type Label struct {
	Name string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

type Filter struct {
	Metric string `json:"metric,omitempty"`
	Labels  []*Label `json:"labels,omitempty"`
	Start time.Time `json:"start,omitempty"`
	End time.Time `json:"end,omitempty"`

}
