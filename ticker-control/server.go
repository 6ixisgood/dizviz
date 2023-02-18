package main

import (
	"bufio"
	"os"
	"fmt"
	"github.com/sixisgoood/matrix-ticker/content"
)

func Serve() {

	for {
		// "listening"
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Println("Enter text: ")
		scanner.Scan()
		text := scanner.Text()
		var text_bytes []byte

		white := "[255,255,255,255]"
		black := "[0,0,0,255]"
		str := "{\"type\":\"textscroll\",\"config\":{\"size\":[64,32],\"textColor\":%s,\"bgColor\":%s,\"direction\":[1,0],\"font\":{\"size\":10,\"type\":\"normal\"},\"text\":\"%s\"}}"
		content := getNHLContent()

		switch text {
			default:
			text_bytes = []byte(fmt.Sprintf(str, white, black, content))
		}

		go HandleRequest(text_bytes)
	}
}

func getNHLContent() string {
	return content.NHLDailyGamesTicker("2022-2023-regular", "20230216")	
}
