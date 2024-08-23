package view

import (
	compCommon "github.com/6ixisgood/matrix-ticker/pkg/component/common"
	viewCommon "github.com/6ixisgood/matrix-ticker/pkg/view/common"
	_ "github.com/6ixisgood/matrix-ticker/pkg/view/types"
	"image"
	"image/draw"
	"log"
	"time"
	"context"
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

func cloneImage(img image.Image) image.Image {
    bounds := img.Bounds()
    dst := image.NewRGBA(bounds)
    draw.Draw(dst, bounds, img, bounds.Min, draw.Src)
    return dst
}


func (a *Animation) startRendering(ctx context.Context) {
    for {
        select {
        case <-ctx.Done():
            // Context was cancelled, exit the goroutine
            return
        default:
			if len(a.buffer) < cap(a.buffer) {
				im := cloneImage(a.view.Template().Render())
				a.buffer <- im
			} else {
				time.Sleep(100 * time.Millisecond)
			}
        }
    }
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
