package data

import (
    "testing"
    //"github.com/stretchr/testify/assert"
    "log"
)

func TestYoutubeMusicAPI(t *testing.T) {
	log.Println("Testing Youtube Music API")

	var api YoutubeMusicAPI
	err := api.Init()
    if err != nil {
        log.Fatalf("Failed to initialize scraper: %v", err)
    }
    defer api.Close()


	song, err := api.GetCurrentSong()
    if err != nil {
    	log.Fatalf("Failed to get current song: %v", err)
    }

    log.Println(song)

}

