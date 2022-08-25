package API

import (
	"encoding/json"
	"fmt"
	"os"

	Util "github.com/artziel/go-utilities"
)

type Result struct {
	MaxDuration int64
	MinDuration int64
	Durations   []int64
	MaxTransfer int
	Success     int
	Fails       int
}

func (r *Result) String() string {
	return fmt.Sprintf(
		"Success: %d\tFails: %d\nAverage Transfer: %s\nDuration [ Min %s\tMax %s\tAvg %s ]",
		r.Success,
		r.Fails,
		Util.HummanReadSize(int64(r.MaxTransfer)),
		Util.HumanReadDuration(r.MinDuration),
		Util.HumanReadDuration(r.MaxDuration),
		Util.HumanReadDuration(r.Average()),
	)
}

func (r *Result) Average() int64 {

	if len(r.Durations) < 1 {
		return -1
	}

	var avg int64
	for _, d := range r.Durations {
		avg = avg + d
	}

	return avg / int64(len(r.Durations))
}

type RequestResult struct {
	Duration   int64
	StatusCode int
	Body       string
	BodySize   int
}

type JsonFile struct {
	Endpoints []Endpoint `json:"endpoints"`
}

func ReadJSON(fileName string) (JsonFile, error) {
	jsonFile := JsonFile{}

	content, err := os.ReadFile(fileName)
	if err != nil {
		return jsonFile, err
	}

	// Unmarshal or Decode the JSON to the interface.
	err = json.Unmarshal(content, &jsonFile)
	if err != nil {
		return jsonFile, err
	}

	return jsonFile, err
}
