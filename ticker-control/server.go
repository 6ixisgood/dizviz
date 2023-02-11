package main

import (
	"bufio"
	"os"
	"fmt"
)

func Serve() {

	for {
		// "listening"
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Println("Enter text: ")
		scanner.Scan()
		text := scanner.Text()
		var text_bytes []byte

		switch text {
			case "a": 
			text_bytes = []byte("{\"type\":\"textscroll\",\"config\":{\"size\":[64,32],\"textColor\":[129,17,198,255],\"bgColor\":[0,0,0,255],\"direction\":[1,0],\"font\":{\"size\":10,\"type\":\"normal\"},\"text\":\"This is my text for the scrolling thing\"}}")	
			case "b":
			text_bytes = []byte("{\"type\":\"textscroll\",\"config\":{\"size\":[64,32],\"textColor\":[58,147,13,255],\"bgColor\":[0,0,0,255],\"direction\":[1,0],\"font\":{\"size\":10,\"type\":\"normal\"},\"text\":\"This is my text for the scrolling thing\"}}")	
			case "c":
			text_bytes = []byte("{\"type\":\"textscroll\",\"config\":{\"size\":[64,32],\"textColor\":[221,162,53,255],\"bgColor\":[0,0,0,255],\"direction\":[1,0],\"font\":{\"size\":10,\"type\":\"normal\"},\"text\":\"This is my text for the scrolling thing\"}}")	
			default:
			text_bytes = []byte("{\"type\":\"textscroll\",\"config\":{\"size\":[64,32],\"textColor\":[255,255,255,255],\"bgColor\":[0,0,0,255],\"direction\":[1,0],\"font\":{\"size\":10,\"type\":\"normal\"},\"text\":\"This is my text for the scrolling thing\"}}")
		}

		go HandleRequest(text_bytes)
	}
}
