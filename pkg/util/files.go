package util

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"github.com/nfnt/resize"
	"image"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// Like FetchFile, but will resize images and redraw gifs
func FetchImage(file string, x int, y int) ([]byte, string, error) {
	// use FetchFile to get raw image
	var rawData []byte
	var rawFilePath string
	var err error

	rawData, rawFilePath, err = FetchFile(file)
	if err != nil {
		log.Fatal(err)
	}

	extension := strings.ToLower(filepath.Ext(rawFilePath))
	baseName := strings.TrimSuffix(filepath.Base(rawFilePath), extension)
	newPath := filepath.Join(Config.CacheDir, fmt.Sprintf("%s_%dx%d%s", baseName, x, y, extension))

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
func FetchFile(file string) ([]byte, string, error) {
	var data []byte
	var err error

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
		cachePath = filepath.Join(Config.CacheDir, hash+extension)
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
