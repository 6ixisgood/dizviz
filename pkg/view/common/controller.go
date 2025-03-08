type Controller struct {
	view     View
	template compCommon.Template
	ctx      context.Context
	cancel   context.CancelFunc
	Buffer   chan image.Image
}

func NewController(bufferSize) {
	return &Controller{
		buffer: make(chan image.Image, bufferSize),
	}
}

func (c *Controller) Init(newView viewCommon.View) {
	log.Printf("Initializing view in controller")

	// init new view in background
	newView.Init()
	viewCommon.TemplateRefresh(newView)

	// stop the old view and switch to new view
	if c.view != nil {
		c.view.Stop()
	}
	c.view = newView

	// close the running rendering task
	if c.cancel != nil {
		c.cancel()
	}

	

	// Start the new rendering task
	go c.startRendering(c.ctx)
}


func cloneImage(img image.Image) image.Image {
	bounds := img.Bounds()
	dst := image.NewRGBA(bounds)
	draw.Draw(dst, bounds, img, bounds.Min, draw.Src)
	return dst
}

func (c *Controller) RenderToBuffer(buffer chan image.Image) {
	// Create a new context for the new rendering task
	c.ctx, c.cancel = context.WithCancel(context.Background())

	go func() {
		for {
			select {
			case <-ctx.Done():
				// Context was cancelled, exit the goroutine
				return
			default:
				if len(buffer) < cap(buffer) {
					im := cloneImage(c.view.Template().Render())
					buffer <- im
				} else {
					time.Sleep(100 * time.Millisecond)
				}
			}
		}
	}
}

func (c *Controller) startRendering(ctx context.Context) {
	
}