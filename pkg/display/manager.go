package display

import (
	"context"
	"fmt"
	"image"
	"image/draw"
	"log"
	"time"

	"github.com/6ixisgood/matrix-ticker/pkg/view/common"
)

// Manager coordinates view rendering and display hardware
type Manager struct {
	display Display
	view    common.View
	ctx     context.Context
	cancel  context.CancelFunc
	running bool

	// Compositor functionality built-in
	fps    int
	frames chan image.Image
}

// NewManager creates a new display manager with default settings (30 FPS, buffer size 10)
func NewManager(display Display) *Manager {
	return &Manager{
		display: display,
		fps:     30,
		frames:  make(chan image.Image, 20),
	}
}

// NewManagerWithSettings creates a new display manager with custom settings
func NewManagerWithSettings(display Display, fps, bufferSize int) *Manager {
	return &Manager{
		display: display,
		fps:     fps,
		frames:  make(chan image.Image, bufferSize),
	}
}

// Start begins the display pipeline with the given view and display config
func (m *Manager) Start(view common.View, displayConfig interface{}) error {
	if m.running {
		return fmt.Errorf("manager already running")
	}

	// Start the display hardware
	if err := m.display.Start(displayConfig); err != nil {
		return fmt.Errorf("failed to start display: %w", err)
	}

	// view.Init()
	common.TemplateRefresh(view)

	// Stop the old view if it exists
	if m.view != nil {
		m.view.Stop()
	}
	m.view = view

	// Start the rendering pipeline
	m.ctx, m.cancel = context.WithCancel(context.Background())
	go m.compositionLoop()
	go m.renderLoop()

	m.running = true
	log.Printf("Display manager started successfully")
	return nil
}

// Stop shuts down the display pipeline
func (m *Manager) Stop() {
	if !m.running {
		return
	}

	log.Printf("Stopping display manager")

	if m.cancel != nil {
		m.cancel()
	}

	if m.view != nil {
		m.view.Stop()
	}

	m.display.Stop()

	m.running = false
	log.Printf("Display manager stopped")
}

// ChangeView switches to a new view
func (m *Manager) ChangeView(view common.View) error {
	if !m.running {
		return fmt.Errorf("manager not running")
	}

	// Note: View is already initialized by the agent before being passed here
	// Refresh the template to ensure it's up to date
	common.TemplateRefresh(view)

	// Stop the old view
	if m.view != nil {
		m.view.Stop()
	}

	// Switch to the new view
	m.view = view

	return nil
}

// SetFPS changes the frames per second
func (m *Manager) SetFPS(fps int) {
	m.fps = fps
}

// compositionLoop generates frames on demand from the current view
func (m *Manager) compositionLoop() {
	for {
		select {
		case <-m.ctx.Done():
			return
		default:
			// Only compose if buffer has space and we have a view
			if len(m.frames) < cap(m.frames) && m.view != nil {
				frame := m.composeFrame()
				if frame != nil {
					select {
					case m.frames <- frame:
						// Frame sent successfully
					case <-m.ctx.Done():
						return
					}
				}
			} else {
				// Buffer is full or no view, wait a bit to avoid spinning
				time.Sleep(time.Millisecond)
			}
		}
	}
}

// composeFrame creates a single frame from the current view
func (m *Manager) composeFrame() image.Image {
	if m.view == nil {
		return nil
	}

	// Render the view template
	source := m.view.RenderTemplate()
	if source == nil {
		return nil
	}

	// Clone the image to avoid memory issues
	bounds := source.Bounds()
	frame := image.NewRGBA(bounds)
	draw.Draw(frame, bounds, source, bounds.Min, draw.Src)

	return frame
}

// renderLoop continuously renders frames at the configured FPS
func (m *Manager) renderLoop() {
	ticker := time.NewTicker(time.Second / time.Duration(m.fps))
	defer ticker.Stop()

	for {
		select {
		case <-m.ctx.Done():
			return
		case <-ticker.C:
			// Try to get a frame from the buffer
			select {
			case frame := <-m.frames:
				if err := m.display.Render(frame); err != nil {
					log.Printf("Error rendering frame: %v", err)
				}
			default:
				// No frame available, skip this tick
				// This is normal and prevents blocking
			}
		}
	}
}
