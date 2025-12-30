package agent

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/6ixisgood/matrix-ticker/pkg/app"
	"github.com/6ixisgood/matrix-ticker/pkg/display"
	viewCommon "github.com/6ixisgood/matrix-ticker/pkg/view/common"
)

// Agent represents a display agent that runs views and connects to control plane
type Agent struct {
	// Identification
	id   string
	name string

	// Configuration
	config       Config
	capabilities Capabilities

	// Core application
	application *app.Application
	display     display.Display

	// Control plane connection
	controlPlaneClient *ControlPlaneClient

	// Runtime state
	ctx       context.Context
	cancel    context.CancelFunc
	running   bool
	startTime time.Time
	mu        sync.RWMutex

	// Status tracking
	currentView string
	currentFPS  int
	health      string
	errorMsg    string
}

// New creates a new agent instance
func New(config Config, capabilities Capabilities) *Agent {
	return &Agent{
		name:         config.AgentName,
		config:       config,
		capabilities: capabilities,
		application:  app.New(),
		currentFPS:   config.FPS,
		health:       "initializing",
	}
}

// SetDisplay sets the display device for this agent
func (a *Agent) SetDisplay(disp display.Display) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.display = disp
}

// Start initializes and starts the agent
func (a *Agent) Start(displayConfig interface{}) error {
	a.mu.Lock()
	if a.running {
		a.mu.Unlock()
		return fmt.Errorf("agent already running")
	}

	if a.display == nil {
		a.mu.Unlock()
		return fmt.Errorf("display not set")
	}

	a.ctx, a.cancel = context.WithCancel(context.Background())
	a.startTime = time.Now()
	a.running = true
	a.mu.Unlock()

	// Start the display application with custom settings
	if err := a.application.StartWithCustomSettings(
		a.display,
		a.config.FPS,
		a.config.BufferSize,
		displayConfig,
	); err != nil {
		return fmt.Errorf("failed to start application: %w", err)
	}

	log.Printf("[Agent:%s] Application started", a.name)
	a.setHealth("healthy", "")

	// Connect to control plane if address is provided
	if a.config.ControlPlaneAddr != "" {
		if err := a.connectToControlPlane(); err != nil {
			log.Printf("[Agent:%s] Failed to connect to control plane: %v", a.name, err)
			if !a.config.EnableAutoReconnect {
				return err
			}
			// Continue running without control plane, will retry in background
		}
	}

	log.Printf("[Agent:%s] Started successfully", a.name)
	return nil
}

// Stop gracefully shuts down the agent
func (a *Agent) Stop() {
	a.mu.Lock()
	if !a.running {
		a.mu.Unlock()
		return
	}
	a.running = false
	a.mu.Unlock()

	log.Printf("[Agent:%s] Stopping...", a.name)

	// Cancel context
	if a.cancel != nil {
		a.cancel()
	}

	// Disconnect from control plane
	if a.controlPlaneClient != nil {
		a.controlPlaneClient.Disconnect()
	}

	// Stop application
	a.application.Stop()

	log.Printf("[Agent:%s] Stopped", a.name)
}

// ChangeView switches to a new view
func (a *Agent) ChangeView(view viewCommon.View) error {
	if err := a.application.ChangeView(view); err != nil {
		a.setHealth("degraded", fmt.Sprintf("failed to change view: %v", err))
		return err
	}

	a.mu.Lock()
	a.currentView = fmt.Sprintf("%T", view)
	a.mu.Unlock()

	log.Printf("[Agent:%s] Changed to view: %s", a.name, a.currentView)
	return nil
}

// SetInitialView sets the view to display when agent starts
func (a *Agent) SetInitialView(view viewCommon.View) {
	a.application.SetInitialView(view)
	a.mu.Lock()
	a.currentView = fmt.Sprintf("%T", view)
	a.mu.Unlock()
}

// GetID returns the agent's ID (assigned by control plane)
func (a *Agent) GetID() string {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.id
}

// GetName returns the agent's name
func (a *Agent) GetName() string {
	return a.name
}

// GetCapabilities returns the agent's capabilities
func (a *Agent) GetCapabilities() Capabilities {
	return a.capabilities
}

// GetStatus returns the current agent status
func (a *Agent) GetStatus() Status {
	a.mu.RLock()
	defer a.mu.RUnlock()

	uptime := time.Duration(0)
	if a.running {
		uptime = time.Since(a.startTime)
	}

	return Status{
		Health:        a.health,
		CurrentView:   a.currentView,
		CurrentFPS:    a.currentFPS,
		UptimeSeconds: int64(uptime.Seconds()),
		ErrorMessage:  a.errorMsg,
	}
}

// IsRunning returns whether the agent is currently running
func (a *Agent) IsRunning() bool {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.running
}

// connectToControlPlane establishes connection to control plane
func (a *Agent) connectToControlPlane() error {
	log.Printf("[Agent:%s] Connecting to control plane at %s", a.name, a.config.ControlPlaneAddr)

	client := NewControlPlaneClient(a.config.ControlPlaneAddr, a)
	if err := client.Connect(a.ctx); err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}

	a.controlPlaneClient = client

	// Register with control plane
	agentID, err := client.Register(a.ctx, a.name, a.capabilities, a.config.Version)
	if err != nil {
		return fmt.Errorf("failed to register: %w", err)
	}

	a.mu.Lock()
	a.id = agentID
	a.mu.Unlock()

	log.Printf("[Agent:%s] Registered with control plane, assigned ID: %s", a.name, agentID)

	// Start heartbeat routine
	go a.heartbeatLoop()

	// Start listening for commands
	go a.commandLoop()

	return nil
}

// heartbeatLoop sends periodic heartbeats to control plane
func (a *Agent) heartbeatLoop() {
	ticker := time.NewTicker(a.config.HeartbeatInterval)
	defer ticker.Stop()

	for {
		select {
		case <-a.ctx.Done():
			return
		case <-ticker.C:
			if a.controlPlaneClient != nil {
				if err := a.controlPlaneClient.SendHeartbeat(a.ctx); err != nil {
					log.Printf("[Agent:%s] Heartbeat failed: %v", a.name, err)
				}
			}
		}
	}
}

// commandLoop listens for commands from control plane
func (a *Agent) commandLoop() {
	// This will be implemented in Phase 2
	// For now, just a placeholder
	<-a.ctx.Done()
}

// setHealth updates the agent's health status
func (a *Agent) setHealth(health, errorMsg string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.health = health
	a.errorMsg = errorMsg
}

// Status represents the current status of an agent
type Status struct {
	Health        string
	CurrentView   string
	CurrentFPS    int
	UptimeSeconds int64
	MemoryUsage   int64
	ErrorMessage  string
}
