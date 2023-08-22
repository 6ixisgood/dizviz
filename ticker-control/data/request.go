package data

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"log"
	"net/http"
	"time"
)

const (
	defaultRateLimitDuration = time.Second
	maxRetries               = 3
	initialBackoffDuration   = time.Second
)

// APIClientOptions holds configuration options for the APIClient.
type APIClientOptions struct {
	BaseURL    string             // The base URL for API requests.
	BasicAuth  *BasicAuthCredentials // Optional basic auth credentials.
	Headers    map[string]string  // Headers to be added to each request.
	Timeout    time.Duration      // Max time to wait for a response.
	RateLimit  time.Duration      // Duration to wait between API calls.
}

// BasicAuthCredentials holds basic authentication username and password.
type BasicAuthCredentials struct {
	Username string
	Password string
}

// APIRequest holds information for making API requests.
type APIRequest struct {
	Method   string            // HTTP method (GET, POST, etc.)
	Endpoint string            // API endpoint, to be appended to the BaseURL.
	Headers  map[string]string // Additional headers for this specific request.
	Body     []byte            // Request body, if applicable.
}

// APIClient is the main API client type.
type APIClient struct {
	options      APIClientOptions  // Configuration options.
	rateLimiter  chan struct{}  // Rate limiter channel.
	httpClient   *http.Client   // Underlying HTTP client.
}

// TimeoutError is an error indicating a request timeout.
type TimeoutError struct{}

// JSONUnmarshalError is an error indicating a failure to unmarshal JSON.
type JSONUnmarshalError struct{ Err error }

func (t TimeoutError) Error() string {
	return "Request timed out"
}

func (j JSONUnmarshalError) Error() string {
	return fmt.Sprintf("Failed to unmarshal JSON: %v", j.Err)
}

// NewAPIClient initializes a new API client with the provided options.
func NewAPIClient(options APIClientOptions) *APIClient {
	if options.RateLimit == 0 {
		options.RateLimit = defaultRateLimitDuration
	}
	if options.Timeout == 0 {
		options.Timeout = 10 * time.Second
	}

	rateLimiter := make(chan struct{}, 1)
	rateLimiter <- struct{}{}

	return &APIClient{
		options:     options,
		rateLimiter: rateLimiter,
		httpClient:  &http.Client{Timeout: options.Timeout},
	}
}

// Do sends the API request and returns the response.
// It uses exponential backoff for retries in case of transient errors.
func (c *APIClient) Do(req *APIRequest) (*http.Response, error) {
	var lastError error
	backoff := initialBackoffDuration

	for attempt := 0; attempt < maxRetries; attempt++ {
		<-c.rateLimiter

		fullURL := c.options.BaseURL + req.Endpoint
		log.Printf("Making %s request to %s", req.Method, fullURL)
		httpReq, err := http.NewRequest(req.Method, fullURL, bytes.NewBuffer(req.Body))
		if err != nil {
			lastError = err
			time.Sleep(backoff)
			backoff *= 2
			continue
		}

		for key, val := range c.options.Headers {
			httpReq.Header.Set(key, val)
		}

		for key, val := range req.Headers {
			httpReq.Header.Set(key, val)
		}

		if c.options.BasicAuth != nil {
			httpReq.SetBasicAuth(c.options.BasicAuth.Username, c.options.BasicAuth.Password)
		}

		resp, err := c.httpClient.Do(httpReq)
		if err != nil {
			if errors.Is(err, http.ErrHandlerTimeout) {
				lastError = TimeoutError{}
			} else {
				lastError = err
			}
			time.Sleep(backoff)
			backoff *= 2
			continue
		}

		return resp, nil
	}

	return nil, lastError
}

// DoAndUnmarshal sends the API request, reads the response,
// and unmarshals it into the provided output structure.
func (c *APIClient) DoAndUnmarshal(req *APIRequest, out interface{}) error {
	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(body, out); err != nil {
		return JSONUnmarshalError{Err: err}
	}
	return nil
}


func ReadFileAndUnmarshal(path string, out interface{}) error {
	// Open the file
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	// Read the file content
	data, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("failed to read file: %v", err)
	}

	// Unmarshal the JSON data
	if err := json.Unmarshal(data, out); err != nil {
		return fmt.Errorf("failed to unmarshal json: %v", err)
	}

	return nil
}
