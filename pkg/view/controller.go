package view

import (
	compCommon "github.com/6ixisgood/matrix-ticker/pkg/component/common"
	viewCommon "github.com/6ixisgood/matrix-ticker/pkg/view/common"
	_ "github.com/6ixisgood/matrix-ticker/pkg/view/types"
	"image"
	"log"
	"time"
	"context"

	"image/color"
	"image/draw"
	"math/rand"
)

var (
	animation = &Animation{}
)

type Animation struct {
	view     viewCommon.View
	template compCommon.Template
	ctx      context.Context
    cancel   context.CancelFunc
	buffer  chan image.Image
}

func (a *Animation) Init(newView viewCommon.View) {
	log.Printf("Initializing view in controller")

	// init new view in background
	newView.Init()
	viewCommon.TemplateRefresh(newView)

	// stop the old view and switch to new view
	if a.view != nil {
		a.view.Stop()
	}
	a.view = newView

	// close the running rendering task
	if a.cancel != nil {
		a.cancel()
	}

	// Create a new context for the new rendering task
    a.ctx, a.cancel = context.WithCancel(context.Background())

	a.buffer = make(chan image.Image, 10)

    // Start the new rendering task
    go a.startRendering(a.ctx)
}

func (a *Animation) startRendering(ctx context.Context) {
    for {
        select {
        case <-ctx.Done():
            // Context was cancelled, exit the goroutine
            return
        default:
			if len(a.buffer) < cap(a.buffer) {
				//im := a.view.Template().Render()
				im := generateRandomColorImage(192, 192)
				a.buffer <- im
			} else {
				time.Sleep(100 * time.Millisecond)
			}
        }
    }
}

var (
	randomColor = color.RGBA{0,0,0,255}
	count = 0
)

func generateRandomColorImage(width, height int) image.Image {
	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	// Generate random color
	if (count > 15) {
		r := uint8(rand.Intn(256))
		g := uint8(rand.Intn(256))
		b := uint8(rand.Intn(256))
		randomColor = color.RGBA{r, g, b, 255}
		count = 0
	} else {
		count = count + 1
	}

	// Create a new image with the specified size
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Fill the image with the random color
	draw.Draw(img, img.Bounds(), &image.Uniform{randomColor}, image.Point{}, draw.Src)

	return img
}

func (a *Animation) Next() (image.Image, <-chan time.Time, error) {
	var im image.Image
	for {
        select {
        case im = <-a.buffer:
            return im, time.After(time.Millisecond * 10), nil
        default:
            time.Sleep(100 * time.Millisecond)
        }
    }

}

func GetAnimation() *Animation {
	return animation
}
