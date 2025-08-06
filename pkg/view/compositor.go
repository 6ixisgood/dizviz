package view

import (
	"log"
	"sync"

	"github.com/6ixisgood/matrix-ticker/pkg/compositor"
	"github.com/6ixisgood/matrix-ticker/pkg/view/common"
)

// Global compositor instance for the application
var (
	globalCompositor *compositor.Compositor
	compositorMutex  sync.RWMutex
)

// InitGlobalCompositor initializes the global compositor
func InitGlobalCompositor() {
	compositorMutex.Lock()
	defer compositorMutex.Unlock()
	
	if globalCompositor == nil {
		globalCompositor = compositor.NewDefault()
		log.Printf("Global compositor initialized")
	}
}

// GetGlobalCompositor returns the global compositor instance
func GetGlobalCompositor() *compositor.Compositor {
	compositorMutex.RLock()
	defer compositorMutex.RUnlock()
	
	if globalCompositor == nil {
		log.Printf("Warning: Global compositor not initialized, creating default")
		return compositor.NewDefault()
	}
	
	return globalCompositor
}

// SetView changes the current view in the global compositor
func SetView(view common.View) error {
	comp := GetGlobalCompositor()
	return comp.Compose(view)
}

// Stop stops the global compositor
func StopGlobalCompositor() {
	compositorMutex.Lock()
	defer compositorMutex.Unlock()
	
	if globalCompositor != nil {
		globalCompositor.Stop()
		globalCompositor = nil
		log.Printf("Global compositor stopped")
	}
}
