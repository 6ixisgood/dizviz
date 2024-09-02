package data

import (
    "context"
    "log"
    "github.com/chromedp/chromedp"
    "errors"
    "time"
    "io/ioutil"
)

type Song struct {
	Name string
	Artist *Artist
	Album *Album
	Lengthm int
}

type Artist struct {
	Name string
}

type Album struct {
	Name string
	Year string
	Image interface{}
}

func NewSong() *Song {
    return &Song{
        Artist: &Artist{},
        Album: &Album{},
    }
}



type YoutubeMusicAPI struct {
    ctx    context.Context
    cancel context.CancelFunc
}

// func NewYoutubeMusicAPI() YoutubeMusicAPI {
// 	return &YoutubeMusicAPI{
// 		ctx: 
// 	}

// }

// Init initializes the browser context and sets up everything.
func (s *YoutubeMusicAPI) Init() error {
    log.Println("Init")

    // create chromedp context
    opts := append(chromedp.DefaultExecAllocatorOptions[:],
        chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36"),
    )
    allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
    defer cancel()
    s.ctx, s.cancel = chromedp.NewContext(allocCtx)
    defer s.cancel()

    var html string
    // Navigate to youtube music
    err := chromedp.Run(s.ctx, chromedp.Navigate("https://music.youtube.com/"))

    // click the "Sign in" link
    err = chromedp.Run(s.ctx,
        chromedp.WaitVisible(`a.sign-in-link`, chromedp.ByQuery),
        chromedp.Click(`a.sign-in-link`, chromedp.ByQuery),
    )
    if err != nil {
        log.Fatal(err)
    }


    err = chromedp.Run(s.ctx,
        chromedp.WaitVisible(`#identifierId`, chromedp.ByID),
        chromedp.SendKeys(`#identifierId`, ""),
        chromedp.Click(`#identifierNext`, chromedp.ByID),
        chromedp.Sleep(10*time.Second),  // wait for the password field to appear
        chromedp.OuterHTML("html", &html, chromedp.ByQuery),
        //chromedp.WaitVisible(`input[type="password"]`, chromedp.ByQuery),
        // chromedp.SendKeys(`input[type="password"]`, "your-password"),
        // chromedp.Click(`#passwordNext`, chromedp.ByID),
    )
    if err != nil {
        log.Fatal(err)
    }

    log.Println("Writing to file")
    err = ioutil.WriteFile("page.html", []byte(html), 0644)
    if err != nil {
        log.Fatal(err)
    }

    //log.Printf("Page HTML:\n%s", html)
    return err
}


// GetCurrentSong fetches the currently playing song.
func (s *YoutubeMusicAPI) GetCurrentSong() (*Song, error) {
    song := NewSong()

    log.Println("Making request")
    err := chromedp.Run(s.ctx,
        //chromedp.WaitVisible(`ytmusic-player-bar`, chromedp.ByQuery),
        chromedp.Text(`.title.ytmusic-player-bar`, &song.Name, chromedp.ByQuery),
        chromedp.Text(`.byline.ytmusic-player-bar`, &song.Artist.Name, chromedp.ByQuery),
        //chromedp.AttributeValue(`.image.ytmusic-player-bar img`, "src", & nil, chromedp.ByQuery)
    )

    if err != nil {
        return song, errors.New("Error getting current song")
    }

    return song, nil
}

// Close cleans up and closes the browser.
func (s *YoutubeMusicAPI) Close() {
    s.cancel()
}
