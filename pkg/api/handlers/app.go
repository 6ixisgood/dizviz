package handlers

import (
	"github.com/6ixisgood/matrix-ticker/pkg/app"
)

// Global application instance for handlers
var appInstance *app.Application

// SetApplication sets the application instance for handlers to use
func SetApplication(app *app.Application) {
	appInstance = app
}

// GetApplication returns the current application instance
func GetApplication() *app.Application {
	return appInstance
}
