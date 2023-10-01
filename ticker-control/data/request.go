package data

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"log"
	"net/http"
	"time"
	"crypto/sha1"
    "encoding/hex"
	"path/filepath"
	"strings"
	"github.com/nfnt/resize"
	"image"
	"image/color/palette"
	"image/draw"
	//"image/color"
	"image/gif"
	"image/png"
	"image/jpeg"
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
	ticker		 *time.Ticker	// ticker to reset rate limit
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

	client := &APIClient{
		options:     options,
		rateLimiter: rateLimiter,
		ticker:		time.NewTicker(options.RateLimit),
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

// Like FetchFile, but will resize images and redraw gifs
func FetchImage(file string, x int, y int) ([]byte, string, error) {
	// use FetchFile to get raw image
	var rawData []byte
	var rawFilePath string
	var err error

	rawData, rawFilePath, err = FetchFile(file)
	if (err != nil) {
		log.Fatal(err)
	}

    cacheDir := "/home/andrew/Lab/matrix-ticker/ticker-control/data/cache/trans/"  // Define your cache directory path here

	extension := strings.ToLower(filepath.Ext(rawFilePath))
	baseName := strings.TrimSuffix(filepath.Base(rawFilePath), extension)
	newPath := filepath.Join(cacheDir, fmt.Sprintf("%s_%dx%d%s", baseName, x, y, extension))

	// check cache for already transformed image
	_, err = os.Stat(newPath)
	if err == nil {
		log.Println("Transformed image already exists, fetch from cache")
	} else {
		log.Println("Transforming image...")

		// Open a file to write the new GIF into
		outFile, err := os.Create(newPath)
		if err != nil {
		    log.Fatal(err)
		}
		defer outFile.Close()

		// Decode based on extension
	    if extension == ".gif" {
		    // redraw in full image frames
	    	var redrawnGIF gif.GIF
		    gifData, err := gif.DecodeAll(bytes.NewReader(rawData))
	        if err != nil {
	            log.Fatal(err)
	        }
	        redrawnGIF.Delay = gifData.Delay

			// Background frame
			bg := image.NewPaletted(gifData.Image[0].Bounds(), gifData.Image[0].Palette)

			// Loop through each frame in the GIF
			for ix, frame := range gifData.Image {
				bounds := frame.Bounds()
				// Create a new frame that starts as a copy of the background
				newFrame := image.NewPaletted(gifData.Image[0].Bounds(), frame.Palette)
				draw.Draw(newFrame, newFrame.Bounds(), bg, image.Point{}, draw.Over)

				// Draw the new frame onto the background, respecting the bounds
				draw.Draw(newFrame, frame.Bounds(), frame, bounds.Min, draw.Over)

				// resize
				resizedImg := resize.Resize(uint(x), uint(y), newFrame, resize.Lanczos3)
				// Convert the resized image.Image to *image.Paletted
				resizedBounds := resizedImg.Bounds()
				palettedImage := image.NewPaletted(resizedBounds, frame.Palette)
				draw.FloydSteinberg.Draw(palettedImage, resizedBounds, resizedImg, image.Point{})
				// append to gif images
				redrawnGIF.Image = append(redrawnGIF.Image, palettedImage)

				// Update the background based on the disposal method
				switch gifData.Disposal[ix] {
				case gif.DisposalNone:
					bg = newFrame
				case gif.DisposalBackground:
					// Reset to original background (or however you want to handle it)
					bg = image.NewPaletted(gifData.Image[0].Bounds(), palette.Plan9)
				}
			}

			// Encode the new GIF and write to the file
			err = gif.EncodeAll(outFile, &redrawnGIF)
			if err != nil {
			    log.Fatal(err)
			}

	    } else if extension == ".png" {
	    	img, _, err := image.Decode(bytes.NewReader(rawData))
	        if err != nil {
	            log.Fatal(err)
	        }
	    	img = resize.Resize(uint(x), uint(y), img, resize.Lanczos3)
	    	err = png.Encode(outFile, img)
			if err != nil {
			    log.Fatalf("Failed to encode image: %s", err)
			}
	    } else if extension == ".jpg" || extension == "jpeg" {
	    	img, _, err := image.Decode(bytes.NewReader(rawData))
		    if err != nil {
		        log.Fatal(err)
		    }
		    // Resize the image
		    img = resize.Resize(uint(x), uint(y), img, resize.Lanczos3)
		    
		    // Encode the image to the outFile as a JPEG
		    err = jpeg.Encode(outFile, img, nil) // nil means use the default quality settings
		    if err != nil {
		        log.Fatalf("Failed to encode image: %s", err)
		    }
	    }
		log.Println("Image Transformed")
	}

	return FetchFile(newPath)
}


// first checks in cache for file. Will fetch from url if needed
func FetchFile(file string) ([]byte, string, error){
	var data []byte
    var err error

    cacheDir := "/home/andrew/Lab/matrix-ticker/ticker-control/data/cache/raw/"  // Define your cache directory path here
    var cachePath string

    if strings.HasPrefix(file, "http://") || strings.HasPrefix(file, "https://") {
    	// get set up to download the file
	    extension := strings.ToLower(filepath.Ext(file))
	    if extension == ".gifv" {
	    	extension = ".gif"
	    }

	    // Create a hash of the URL to use as the filename
	    h := sha1.New()
	    h.Write([]byte(file))
	    hash := hex.EncodeToString(h.Sum(nil))
    	cachePath = filepath.Join(cacheDir, hash + extension)
	} else {
		// go directly to file
		cachePath = file
	}


	// Check if the file already exists in the cache
    if _, err := os.Stat(cachePath); os.IsNotExist(err) {
	    // Download the file
	    log.Printf("Downloading file from %s", file)
        resp, err := http.Get(file)
        if err != nil {
            log.Fatal(err)
        }
        defer resp.Body.Close()

        data, err = ioutil.ReadAll(resp.Body)
        if err != nil {
            log.Fatal(err)
        }

	    // Save to cache
	    err = os.WriteFile(cachePath, data, 0644)
	    if err != nil {
	        log.Fatal(err)
	    } 
	    log.Printf("Saved %s to %s", file, cachePath)

    } else {
        // Load from cache
        data, err = ioutil.ReadFile(cachePath)
        if err != nil {
            log.Fatal(err)
        }
        log.Printf("Loaded %s from cache", file)
    }

    return data, cachePath, err
}

