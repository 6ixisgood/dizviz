package app

import (
	"fmt"
	"log"

	"github.com/6ixisgood/matrix-ticker/pkg/display"
	"github.com/6ixisgood/matrix-ticker/pkg/view/common"
)

// Application manages the main application state and display system
type Application struct {
	displayManager *display.Manager
	currentView    common.View
	running        bool
}

// New creates a new application instance
func New() *Application {
	return &Application{}
}

// Start initializes and starts the display system with the given display and config
func (app *Application) Start(disp display.Display, displayConfig interface{}) error {
	if app.running {
		return fmt.Errorf("application already running")
	}

	// Create display manager with default compositor
	app.displayManager = display.NewManager(disp)

	// If we have a current view, start with it
	if app.currentView != nil {
		if err := app.displayManager.Start(app.currentView, displayConfig); err != nil {
			return fmt.Errorf("failed to start display manager: %w", err)
		}
	}

	app.running = true
	log.Printf("Application started successfully")
	return nil
}

// StartWithCustomSettings starts the application with custom display manager settings
func (app *Application) StartWithCustomSettings(disp display.Display, fps, bufferSize int, displayConfig interface{}) error {
	if app.running {
		return fmt.Errorf("application already running")
	}

	// Create display manager with custom settings
	app.displayManager = display.NewManagerWithSettings(disp, fps, bufferSize)

	// If we have a current view, start with it
	if app.currentView != nil {
		if err := app.displayManager.Start(app.currentView, displayConfig); err != nil {
			return fmt.Errorf("failed to start display manager: %w", err)
		}
	}

	app.running = true
	log.Printf("Application started with custom settings")
	return nil
}

// Stop shuts down the application
func (app *Application) Stop() {
	if !app.running {
		return
	}

	if app.displayManager != nil {
		app.displayManager.Stop()
	}

	app.running = false
	log.Printf("Application stopped")
}

// ChangeView switches to a new view
func (app *Application) ChangeView(view common.View) error {
	if !app.running {
		return fmt.Errorf("application not running")
	}

	if app.displayManager == nil {
		return fmt.Errorf("display manager not initialized")
	}

	app.currentView = view
	return app.displayManager.ChangeView(view)
}

// SetInitialView sets the view to use when the application starts
func (app *Application) SetInitialView(view common.View) {
	app.currentView = view
}

// GetCurrentView returns the current view
func (app *Application) GetCurrentView() common.View {
	return app.currentView
}

// IsRunning returns whether the application is currently running
func (app *Application) IsRunning() bool {
	return app.running
}

// GetDisplayManager returns the display manager (for advanced use cases)
func (app *Application) GetDisplayManager() *display.Manager {
	return app.displayManager
}
