package checkurl

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

// Attributes is the message data structure
type Attributes struct {
	BrainURL string `json:"brainurl"`
}

// PubSubMessage is the payload of a Pub/Sub event. Please refer to the docs for
// additional information regarding Pub/Sub events.
type PubSubMessage struct {
	Data       []byte     `json:"data"`
	Attributes Attributes `json:"attributes"`
}

// CheckURL contains the url and last http response
type CheckURL struct {
	ID             string `json:"id"`
	URL            string `json:"url"`
	ActualResponse int16  `json:"actual_response"`
}

// Response is a object related with checks
type Response struct {
	ID           string  `json:"id"`
	URL          string  `json:"url"`
	Response     int16   `json:"response"`
	LastResponse int16   `json:"last_response"`
	RequestTime  float64 `json:"request_time"`
}

// responsesJSON transform responses to a json array
func responsesJSON(responses []Response) []byte {
	output, err := json.Marshal(responses)
	if err != nil {
		return []byte("[]")
	}
	return output
}

// CheckUrl sends a new email with the given info
func CheckUrl(ctx context.Context, m PubSubMessage) error {
	log.Printf("CheckURL: %s!", m.Data)
	var checkURLS []CheckURL
	err := json.Unmarshal(m.Data, &checkURLS)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("BrainURL: %s", m.Attributes.BrainURL)
	jwtToken := os.Getenv("TOKEN")
	brainURL := m.Attributes.BrainURL
	var responses []Response
	var status int16

	var wg sync.WaitGroup
	wg.Add(len(checkURLS))

	for _, monitor := range checkURLS {
		go func(id string, r string, a int16) {
			defer wg.Done()
			startTime := time.Now()
			// Request
			resp, err := http.Get(r)
			d := time.Since(startTime)
			// Calc request's duration
			duration := float64(d) / float64(time.Second)
			if err != nil {
				status = 600
			} else {
				status = int16(resp.StatusCode)
			}
			// Build responses' array
			responses = append(responses, Response{ID: id, URL: r, Response: status, LastResponse: a, RequestTime: duration})
		}(monitor.ID, monitor.URL, monitor.ActualResponse)
	}

	wg.Wait()

	// Make json string
	result := responsesJSON(responses)
	log.Println("JSON: ", string(result))
	// Send to API
	if brainURL != "" {
		// Create request object.
		req, err := http.NewRequest("POST", brainURL, bytes.NewBuffer(result))

		// Set the header in the request.
		req.Header.Set("Authorization", "Bearer "+jwtToken)

		// Execute the request.
		resp, err := http.DefaultClient.Do(req)

		if err != nil {
			log.Fatal("Error connection to Brain API", err)
		}
		defer resp.Body.Close()
	}

	return nil
}
