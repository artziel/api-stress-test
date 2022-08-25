package API

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	Parallel "github.com/artziel/go-parallel"
)

type Response struct {
	StatusCode   int
	ResponseSize int64
	StartAt      int64
	EndAt        int64
	Body         string
}

func (r *Response) Duration() int64 {
	diff := r.EndAt - r.StartAt
	return diff
}

type Endpoint struct {
	URL         string                 `json:"url"`
	Header      map[string][]string    `json:"header"`
	Data        map[string]interface{} `json:"data"`
	XX          map[string]interface{} `json:"xx"`
	Method      string                 `json:"method"`
	Iterations  int                    `json:"iterations"`
	Concurrents int                    `json:"concurrents"`
	Response    Response
}

func (e *Endpoint) doRequest() (RequestResult, error) {
	var req *http.Request
	var err error
	var payload *bytes.Buffer
	result := RequestResult{}

	client := http.Client{}

	if e.Data != nil {
		data, err := json.MarshalIndent(e.Data, "", "    ")
		if err != nil {
			return result, fmt.Errorf("error parsing endpoint data: %s", err)
		}
		payload = bytes.NewBuffer(data)
		req, err = http.NewRequest(e.Method, e.URL, payload)
		if err != nil {
			return result, err
		}
	} else {
		req, err = http.NewRequest(e.Method, e.URL, nil)
		if err != nil {
			return result, err
		}
	}
	req.Header = e.Header

	StartAt := time.Now().UnixNano() / int64(time.Millisecond)

	resp, err := client.Do(req)
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return result, err
	}
	result.BodySize = len(b)
	result.Body = string(b)
	EndAt := time.Now().UnixNano() / int64(time.Millisecond)
	result.Duration = EndAt - StartAt
	result.StatusCode = resp.StatusCode

	return result, nil
}

func (e *Endpoint) setDefaults() {
	// Set default headers
	header := http.Header{
		"accept-encoding": {"gzip, deflate, br"},
		"accept-language": {"es-MX,es;q=0.9,en-US;q=0.8,en;q=0.7,es-419;q=0.6"},
		"user-agent":      {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.5005.63 Safari/537.36"},
	}
	if e.Data != nil {
		header["Content-Type"] = []string{"application/json"}
	}
	// Merge Default headers
	for k, v := range e.Header {
		header[k] = v
	}
	e.Header = header

	if e.Method == "" {
		e.Method = "GET"
	} else {
		e.Method = strings.ToUpper(e.Method)
	}
	if e.Iterations == 0 {
		e.Iterations = 1
	}
	if e.Concurrents == 0 {
		e.Concurrents = 1
	}
}

func (e *Endpoint) Exec() (Result, error) {

	e.setDefaults()

	runner := Parallel.Processor{}

	result := Result{
		Durations: []int64{},
	}
	for i := 0; i < e.Iterations; i++ {
		runner.AddWorker(
			&Parallel.Worker{
				Data: i,
				Fnc: func(w *Parallel.Worker) error {
					r, _ := e.doRequest()
					w.Data = r
					return nil
				},
			})
	}

	err := runner.Run(e.Concurrents, nil)

	for _, w := range runner.Workers {
		r := w.Data.(RequestResult)
		result.Durations = append(result.Durations, r.Duration)
		if result.MaxDuration < r.Duration {
			result.MaxDuration = r.Duration
		}
		if result.MinDuration > r.Duration || result.MinDuration == 0 {
			result.MinDuration = r.Duration
		}
		if r.StatusCode < 200 || r.StatusCode > 299 {
			result.Fails = result.Fails + 1
		} else {
			result.Success = result.Success + 1
		}
		if result.MaxTransfer < r.BodySize {
			result.MaxTransfer = r.BodySize
		}
	}

	if err != nil {
		return result, err
	}

	return result, nil
}
