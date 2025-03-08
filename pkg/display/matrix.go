package display

import (
	"os"
	"time"
	"github.com/sixisgoood/go-rpi-rgb-led-matrix"
)

type Display interface {
	StartDisplay(view view.View)
	StopDisplay()
	Next() (image.Image, <-chan time.Time)
}

type MatrixDisplay struct {
	config *rgbmatrix.DefaultConfig
	animation *Animation
	stopChan  chan struct{}
	imageBuffer chan image.Image
}

func (d *MatrixDisplay) StartDisplay(view viewCommon.View, config *rgbmatrix.DefaultConfig)  {
	fmt.Println("Starting Matrix\n")

	// create animation

	a := newMatrixAnimation()

	// setup matrix
	d.config, err := rgbmatrix.NewRGBLedMatrix(config)
	fatal(err)

	// create and init the controller
	d.controller = viewCommon.NewController()
	d.controller.Init(view)

	// create toolkit and start playing
	tk := rgbmatrix.NewToolKit(d.config)
	go tk.PlayAnimation(d)
}


func (d *MatrixDisplay) StopDisplay() {
	close(d.stopChan)
}


type MatrixAnimation struct {
	buffer		chan image.Image
	stopChan  chan struct{}

}

func newMatrixAnimation() MatrixAnimation{
	return MatrixAnimation{
		buffer: make(chan image.Image, 10),
		stopChan: make(chan stuct{}),
	}
}


func (a *MatrixAnimation) Next() (image.Image, <-chan time.Time, error) {
	var im image.Image
	for {
		select {
		case im = <-a.Buffer:
			return im, time.After(time.Millisecond * 10), nil
		case <-d.stopChan:
			return nil, io.EOF
		default:
			time.Sleep(100 * time.Millisecond)
		}
	}
}
