package agent

import (
	"time"
)

// Config holds agent configuration
type Config struct {
	// Agent identification
	AgentName string
	Version   string

	// Control plane connection
	ControlPlaneAddr  string        // "localhost:50051"
	ReconnectDelay    time.Duration // Delay before reconnecting
	HeartbeatInterval time.Duration // How often to send heartbeats

	// Display configuration
	DisplayType       string
	DisplayConfigPath string

	// Runtime settings
	FPS        int
	BufferSize int

	// Features
	EnableAutoReconnect bool
	EnableHealthChecks  bool
}

// DefaultConfig returns a configuration with sensible defaults
func DefaultConfig() Config {
	return Config{
		AgentName:           "disco-agent",
		Version:             "1.0.0",
		ControlPlaneAddr:    "localhost:50051",
		ReconnectDelay:      5 * time.Second,
		HeartbeatInterval:   10 * time.Second,
		DisplayType:         "led-matrix",
		FPS:                 30,
		BufferSize:          20,
		EnableAutoReconnect: true,
		EnableHealthChecks:  true,
	}
}
