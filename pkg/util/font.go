package util

import (
	"fmt"
	"github.com/golang/freetype/truetype"
	"io/ioutil"
	"log"
	"path/filepath"
)

func LoadFont(fontName string) *truetype.Font {
	// Read font file from disk
	filePath := filepath.Join(Config.FontDir, fmt.Sprintf("%s.ttf", fontName))
	fontBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Failed to read font file: %v", err)
	}
	// Parse font file into a truetype.Font
	font, err := truetype.Parse(fontBytes)
	if err != nil {
		log.Fatalf("Failed to parse font file: %v", err)
	}
	return font
}
