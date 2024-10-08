package data

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
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
	BaseURL   string                // The base URL for API requests.
	BasicAuth *BasicAuthCredentials // Optional basic auth credentials.
	Headers   map[string]string     // Headers to be added to each request.
	Timeout   time.Duration         // Max time to wait for a response.
	RateLimit time.Duration         // Duration to wait between API calls.
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
	options     APIClientOptions // Configuration options.
	rateLimiter chan struct{}    // Rate limiter channel.
	ticker      *time.Ticker     // ticker to reset rate limit
	httpClient  *http.Client     // Underlying HTTP client.
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

	client := &APIClient{
		options:     options,
		rateLimiter: rateLimiter,
		ticker:      time.NewTicker(options.RateLimit),
		httpClient:  &http.Client{Timeout: options.Timeout},
	}

	go client.refillTokens()

	return client
}

func (c *APIClient) refillTokens() {
	for range c.ticker.C {
		select {
		case c.rateLimiter <- struct{}{}:
		default:
		}
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

// DoAndUnmarshal sends a request and unmarshals the response into a provided struct
func (c *APIClient) DoAndUnmarshal(req *APIRequest, v interface{}) (int, error) {
	// Send the request
	resp, err := c.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, err
	}

	// Handle status codes
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		// Return error for non-2xx responses
		return resp.StatusCode, errors.New(fmt.Sprintf("unexpected status code: %d", resp.StatusCode))
	}

	// Unmarshal the response body if the status code is 2xx and the body is not empty
	if len(body) > 0 {
		if err := json.Unmarshal(body, v); err != nil {
			return resp.StatusCode, err
		}
	}

	// Return the status code and no error
	return resp.StatusCode, nil
}
