package agent

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/6ixisgood/matrix-ticker/pkg/display"
	viewCommon "github.com/6ixisgood/matrix-ticker/pkg/view/common"
	pb "github.com/6ixisgood/matrix-ticker/proto"
)

// DisplayInfo holds information about a single display managed by the agent
type DisplayInfo struct {
	manager     *display.Manager
	display     display.Display
	currentView viewCommon.View
	currentFPS  int
	active      bool
	mu          sync.RWMutex
}

// Agent represents a display agent that runs views and connects to control plane
type Agent struct {
	// Identification
	id   string
	name string

	// Configuration
	config       Config
	capabilities Capabilities

	// Multi-display management
	displays map[string]*DisplayInfo

	// Control plane connection
	controlPlaneClient *ControlPlaneClient

	// Runtime state
	ctx       context.Context
	cancel    context.CancelFunc
	running   bool
	startTime time.Time
	mu        sync.RWMutex

	// Status tracking
	health   string
	errorMsg string
}

// New creates a new agent instance
func New(config Config, capabilities Capabilities) *Agent {
	return &Agent{
		name:         config.AgentName,
		config:       config,
		capabilities: capabilities,
		displays:     make(map[string]*DisplayInfo),
		health:       "initializing",
	}
}

// AddDisplay adds a display to the agent
func (a *Agent) AddDisplay(displayID string, disp display.Display, fps, bufferSize int) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if _, exists := a.displays[displayID]; exists {
		return fmt.Errorf("display %s already exists", displayID)
	}

	a.displays[displayID] = &DisplayInfo{
		display:    disp,
		manager:    display.NewManagerWithSettings(disp, fps, bufferSize),
		currentFPS: fps,
		active:     false,
	}

	log.Printf("[Agent:%s] Added display: %s", a.name, displayID)
	return nil
}

// Start initializes and starts the agent
func (a *Agent) Start(displayConfig interface{}) error {
	a.mu.Lock()
	if a.running {
		a.mu.Unlock()
		return fmt.Errorf("agent already running")
	}

	if len(a.displays) == 0 {
		a.mu.Unlock()
		return fmt.Errorf("no displays configured")
	}

	a.ctx, a.cancel = context.WithCancel(context.Background())
	a.startTime = time.Now()
	a.running = true
	a.mu.Unlock()

	// Start all display managers
	for displayID, displayInfo := range a.displays {
		displayInfo.mu.Lock()
		// If there's an initial view set, start with it
		if displayInfo.currentView != nil {
			if err := displayInfo.manager.Start(displayInfo.currentView, displayConfig); err != nil {
				displayInfo.mu.Unlock()
				return fmt.Errorf("failed to start display %s: %w", displayID, err)
			}
			displayInfo.active = true
		}
		displayInfo.mu.Unlock()
		log.Printf("[Agent:%s] Display %s started", a.name, displayID)
	}

	log.Printf("[Agent:%s] Agent started with %d display(s)", a.name, len(a.displays))
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

	// Stop all display managers
	for displayID, displayInfo := range a.displays {
		displayInfo.mu.Lock()
		if displayInfo.manager != nil {
			displayInfo.manager.Stop()
		}
		displayInfo.active = false
		displayInfo.mu.Unlock()
		log.Printf("[Agent:%s] Display %s stopped", a.name, displayID)
	}

	log.Printf("[Agent:%s] Stopped", a.name)
}

// ChangeView switches to a new view on a specific display
func (a *Agent) ChangeView(displayID string, view viewCommon.View) error {
	a.mu.RLock()
	displayInfo, exists := a.displays[displayID]
	a.mu.RUnlock()

	if !exists {
		return fmt.Errorf("display %s not found", displayID)
	}

	displayInfo.mu.Lock()
	defer displayInfo.mu.Unlock()

	if err := displayInfo.manager.ChangeView(view); err != nil {
		a.setHealth("degraded", fmt.Sprintf("failed to change view on display %s: %v", displayID, err))
		return err
	}

	displayInfo.currentView = view

	log.Printf("[Agent:%s] Changed display %s to view: %T", a.name, displayID, view)
	return nil
}

// SetInitialView sets the view to display when a specific display starts
func (a *Agent) SetInitialView(displayID string, view viewCommon.View) error {
	a.mu.RLock()
	displayInfo, exists := a.displays[displayID]
	a.mu.RUnlock()

	if !exists {
		return fmt.Errorf("display %s not found", displayID)
	}

	displayInfo.mu.Lock()
	displayInfo.currentView = view
	displayInfo.mu.Unlock()

	log.Printf("[Agent:%s] Set initial view for display %s", a.name, displayID)
	return nil
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

// GetDisplayIDs returns all display IDs managed by this agent
func (a *Agent) GetDisplayIDs() []string {
	a.mu.RLock()
	defer a.mu.RUnlock()

	ids := make([]string, 0, len(a.displays))
	for id := range a.displays {
		ids = append(ids, id)
	}
	return ids
}

// GetDisplayInfo returns information about a specific display
func (a *Agent) GetDisplayInfo(displayID string) (*DisplayInfo, bool) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	info, exists := a.displays[displayID]
	return info, exists
}

// GetStatus returns the current agent status
func (a *Agent) GetStatus() Status {
	a.mu.RLock()
	defer a.mu.RUnlock()

	uptime := time.Duration(0)
	if a.running {
		uptime = time.Since(a.startTime)
	}

	// Build display statuses from all managed displays
	displays := make([]DisplayStatus, 0, len(a.displays))
	for displayID, displayInfo := range a.displays {
		displayInfo.mu.RLock()
		currentViewType := ""
		if displayInfo.currentView != nil {
			currentViewType = fmt.Sprintf("%T", displayInfo.currentView)
		}
		displays = append(displays, DisplayStatus{
			DisplayID:   displayID,
			CurrentView: currentViewType,
			CurrentFPS:  displayInfo.currentFPS,
			Active:      displayInfo.active,
		})
		displayInfo.mu.RUnlock()
	}

	return Status{
		Health:        a.health,
		Displays:      displays,
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

	consecutiveFailures := 0
	maxFailures := 3

	for {
		select {
		case <-a.ctx.Done():
			return
		case <-ticker.C:
			if a.controlPlaneClient != nil {
				if err := a.controlPlaneClient.SendHeartbeat(a.ctx); err != nil {
					consecutiveFailures++
					log.Printf("[Agent:%s] Heartbeat failed (%d/%d): %v", a.name, consecutiveFailures, maxFailures, err)

					// If too many consecutive failures and auto-reconnect is enabled, attempt reconnection
					if consecutiveFailures >= maxFailures && a.config.EnableAutoReconnect {
						log.Printf("[Agent:%s] Too many heartbeat failures, attempting reconnection", a.name)
						go a.attemptReconnection()
					}
				} else {
					// Reset failure counter on success
					if consecutiveFailures > 0 {
						consecutiveFailures = 0
						log.Printf("[Agent:%s] Heartbeat recovered", a.name)
					}
				}
			}
		}
	}
}

// attemptReconnection tries to reconnect to the control plane
func (a *Agent) attemptReconnection() {
	log.Printf("[Agent:%s] Starting reconnection process", a.name)

	// Disconnect current client
	if a.controlPlaneClient != nil {
		a.controlPlaneClient.Disconnect()
	}

	// Attempt to reconnect
	if err := a.connectToControlPlane(); err != nil {
		log.Printf("[Agent:%s] Reconnection failed: %v", a.name, err)
		a.setHealth("degraded", fmt.Sprintf("control plane disconnected: %v", err))

		// Schedule another retry if auto-reconnect is enabled
		if a.config.EnableAutoReconnect {
			time.Sleep(a.config.ReconnectDelay)
			go a.attemptReconnection()
		}
		return
	}

	log.Printf("[Agent:%s] Reconnection successful", a.name)
	a.setHealth("healthy", "")
}

// commandLoop listens for commands from control plane
func (a *Agent) commandLoop() {
	if a.controlPlaneClient == nil {
		log.Printf("[Agent:%s] No control plane client, command loop exiting", a.name)
		return
	}

	// Start command stream
	commandChan, err := a.controlPlaneClient.StartCommandStream(a.ctx)
	if err != nil {
		log.Printf("[Agent:%s] Failed to start command stream: %v", a.name, err)
		return
	}

	log.Printf("[Agent:%s] Command loop started", a.name)

	for {
		select {
		case <-a.ctx.Done():
			log.Printf("[Agent:%s] Command loop shutting down", a.name)
			return

		case cmd, ok := <-commandChan:
			if !ok {
				log.Printf("[Agent:%s] Command channel closed", a.name)
				return
			}

			// Handle command based on type
			a.handleCommand(cmd)
		}
	}
}

// handleCommand processes a command from the control plane
func (a *Agent) handleCommand(msg *pb.ControlPlaneMessage) {
	if msg == nil {
		return
	}

	switch payload := msg.Payload.(type) {
	case *pb.ControlPlaneMessage_AssignView:
		a.handleAssignView(payload.AssignView)
	case *pb.ControlPlaneMessage_UpdateConfig:
		a.handleUpdateConfig(payload.UpdateConfig)
	case *pb.ControlPlaneMessage_Stop:
		a.handleStopCommand(payload.Stop)
	case *pb.ControlPlaneMessage_Ping:
		a.handlePing(payload.Ping)
	default:
		log.Printf("[Agent:%s] Unknown command type", a.name)
	}
}

// handleAssignView handles view assignment command
func (a *Agent) handleAssignView(cmd *pb.AssignViewCommand) {
	log.Printf("[Agent:%s] Received AssignView command: type=%s, display=%s",
		a.name, cmd.ViewType, cmd.DisplayId)

	// Validate display ID (for now, only "primary" is supported)
	displayID := cmd.DisplayId
	if displayID == "" {
		displayID = "primary"
	}

	if displayID != "primary" {
		a.sendCommandResponse(false,
			fmt.Sprintf("display %s not found", displayID),
			"AssignView")
		return
	}

	// Fetch view definition by ID (following same pattern as DisplayViewById handler)
	viewDefinition, err := viewCommon.GetViewDefinition(cmd.ViewType)
	if err != nil {
		log.Printf("[Agent:%s] View definition not found: %s", a.name, cmd.ViewType)
		a.sendCommandResponse(false,
			fmt.Sprintf("view definition %s not found: %v", cmd.ViewType, err),
			"AssignView")
		return
	}

	// Get the registered view type factory
	regView, exists := viewCommon.RegisteredViews[viewDefinition.Type]
	if !exists {
		log.Printf("[Agent:%s] View type not registered: %s", a.name, viewDefinition.Type)
		a.sendCommandResponse(false,
			fmt.Sprintf("view type %s not registered", viewDefinition.Type),
			"AssignView")
		return
	}

	// Use config from command if provided, otherwise use stored config
	config := viewDefinition.Config
	if cmd.ViewConfigJson != "" && cmd.ViewConfigJson != "{}" {
		config = cmd.ViewConfigJson
	}

	// Create view instance from config
	newView, err := regView.NewView(config)
	if err != nil {
		log.Printf("[Agent:%s] Failed to create view: %v", a.name, err)
		a.sendCommandResponse(false,
			fmt.Sprintf("failed to create view: %v", err),
			"AssignView")
		return
	}

	// Change to the new view on the specified display
	if err := a.ChangeView(displayID, newView); err != nil {
		log.Printf("[Agent:%s] Failed to change view on display %s: %v", a.name, displayID, err)
		a.sendCommandResponse(false,
			fmt.Sprintf("failed to change view: %v", err),
			"AssignView")
		return
	}

	// Success!
	log.Printf("[Agent:%s] Successfully changed display %s to view: %s (type: %s)",
		a.name, displayID, viewDefinition.Id, viewDefinition.Type)
	a.sendCommandResponse(true,
		fmt.Sprintf("view %s assigned successfully to display %s", cmd.ViewType, displayID),
		"AssignView")
}

// handleUpdateConfig handles configuration update command
func (a *Agent) handleUpdateConfig(cmd *pb.UpdateConfigCommand) {
	log.Printf("[Agent:%s] Received UpdateConfig command", a.name)

	success := true
	message := "Configuration updated"

	// Update FPS if provided
	if cmd.Fps != nil {
		newFPS := int(*cmd.Fps)
		a.mu.Lock()
		a.config.FPS = newFPS
		a.mu.Unlock()
		log.Printf("[Agent:%s] Updated FPS to %d", a.name, newFPS)
		// Note: FPS is applied per display manager when started
	}

	// Update buffer size if provided
	if cmd.BufferSize != nil {
		newBufferSize := int(*cmd.BufferSize)
		a.mu.Lock()
		a.config.BufferSize = newBufferSize
		a.mu.Unlock()
		log.Printf("[Agent:%s] Updated buffer size to %d", a.name, newBufferSize)
		// Note: Buffer size is applied per display manager when started
	}

	if err := a.controlPlaneClient.SendCommandResponse(a.ctx, success, message, "UpdateConfig"); err != nil {
		log.Printf("[Agent:%s] Failed to send command response: %v", a.name, err)
	}
}

// handleStopCommand handles stop command from control plane
func (a *Agent) handleStopCommand(cmd *pb.StopCommand) {
	log.Printf("[Agent:%s] Received Stop command: %s (graceful=%v)", a.name, cmd.Reason, cmd.Graceful)

	// Send response before stopping
	if err := a.controlPlaneClient.SendCommandResponse(a.ctx, true, "Shutting down", "Stop"); err != nil {
		log.Printf("[Agent:%s] Failed to send command response: %v", a.name, err)
	}

	// Stop the agent
	if cmd.Graceful {
		go func() {
			time.Sleep(100 * time.Millisecond) // Brief delay to send response
			a.Stop()
		}()
	} else {
		a.Stop()
	}
}

// handlePing handles ping command from control plane
func (a *Agent) handlePing(cmd *pb.PingCommand) {
	log.Printf("[Agent:%s] Received Ping command", a.name)

	message := fmt.Sprintf("Pong at %d", time.Now().Unix())
	if err := a.controlPlaneClient.SendCommandResponse(a.ctx, true, message, "Ping"); err != nil {
		log.Printf("[Agent:%s] Failed to send command response: %v", a.name, err)
	}
}

// sendCommandResponse is a helper method for sending command responses to control plane
func (a *Agent) sendCommandResponse(success bool, message, commandType string) {
	if a.controlPlaneClient != nil {
		if err := a.controlPlaneClient.SendCommandResponse(a.ctx, success, message, commandType); err != nil {
			log.Printf("[Agent:%s] Failed to send command response: %v", a.name, err)
		}
	}
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
	Displays      []DisplayStatus
	UptimeSeconds int64
	MemoryUsage   int64
	ErrorMessage  string
}

// DisplayStatus represents the status of a single display
type DisplayStatus struct {
	DisplayID   string
	CurrentView string
	CurrentFPS  int
	Active      bool
}
