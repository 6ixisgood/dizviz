package agent

import (
	"context"
	"fmt"
	"io"
	"log"
	"sync"
	"time"

	pb "github.com/6ixisgood/matrix-ticker/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// ControlPlaneClient handles communication with the control plane
type ControlPlaneClient struct {
	addr      string
	agent     *Agent
	conn      *grpc.ClientConn
	client    pb.AgentServiceClient
	connected bool
	mu        sync.RWMutex

	// Command stream
	commandStream pb.AgentService_StreamCommandsClient
	streamCtx     context.Context
	streamCancel  context.CancelFunc
	streamMu      sync.Mutex
}

// NewControlPlaneClient creates a new control plane client
func NewControlPlaneClient(addr string, agent *Agent) *ControlPlaneClient {
	return &ControlPlaneClient{
		addr:  addr,
		agent: agent,
	}
}

// Connect establishes a connection to the control plane
func (c *ControlPlaneClient) Connect(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	log.Printf("[ControlPlaneClient] Connecting to %s", c.addr)

	// Set up connection with timeout
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(
		ctx,
		c.addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return fmt.Errorf("failed to dial: %w", err)
	}

	c.conn = conn
	c.client = pb.NewAgentServiceClient(conn)
	c.connected = true

	log.Printf("[ControlPlaneClient] Connected successfully")
	return nil
}

// Disconnect closes the connection to the control plane
func (c *ControlPlaneClient) Disconnect() {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Close command stream
	if c.streamCancel != nil {
		c.streamCancel()
	}

	if c.conn != nil {
		c.conn.Close()
		c.connected = false
		log.Printf("[ControlPlaneClient] Disconnected")
	}
}

// Register registers the agent with the control plane
func (c *ControlPlaneClient) Register(ctx context.Context, name string, caps Capabilities, version string) (string, error) {
	c.mu.RLock()
	connected := c.connected
	client := c.client
	c.mu.RUnlock()

	if !connected {
		return "", fmt.Errorf("not connected to control plane")
	}

	log.Printf("[ControlPlaneClient] Registering agent: %s", name)

	// Convert agent capabilities to protobuf format
	pbDisplays := make([]*pb.DisplayCapability, 0, len(caps.Displays))
	for _, display := range caps.Displays {
		pbDisplays = append(pbDisplays, &pb.DisplayCapability{
			DisplayId:   display.DisplayID,
			DisplayType: display.DisplayType,
			Width:       int32(display.Width),
			Height:      int32(display.Height),
			MaxFps:      int32(display.MaxFPS),
			ColorDepth:  int32(display.ColorDepth),
		})
	}

	req := &pb.RegisterRequest{
		AgentName: name,
		Version:   version,
		Capabilities: &pb.Capabilities{
			Displays:       pbDisplays,
			SupportedViews: caps.SupportedViews,
		},
	}

	resp, err := client.RegisterAgent(ctx, req)
	if err != nil {
		return "", fmt.Errorf("registration failed: %w", err)
	}

	if !resp.Success {
		return "", fmt.Errorf("registration rejected: %s", resp.Message)
	}

	log.Printf("[ControlPlaneClient] Registration successful, ID: %s", resp.AgentId)
	return resp.AgentId, nil
}

// SendHeartbeat sends a heartbeat to the control plane
func (c *ControlPlaneClient) SendHeartbeat(ctx context.Context) error {
	c.mu.RLock()
	connected := c.connected
	client := c.client
	agentID := c.agent.GetID()
	c.mu.RUnlock()

	if !connected {
		return fmt.Errorf("not connected to control plane")
	}

	req := &pb.HeartbeatRequest{
		AgentId:   agentID,
		Timestamp: time.Now().Unix(),
	}

	resp, err := client.Heartbeat(ctx, req)
	if err != nil {
		return fmt.Errorf("heartbeat failed: %w", err)
	}

	if !resp.Success {
		return fmt.Errorf("heartbeat rejected")
	}

	return nil
}

// UpdateStatus sends a status update to the control plane
func (c *ControlPlaneClient) UpdateStatus(ctx context.Context, status Status) error {
	c.mu.RLock()
	connected := c.connected
	client := c.client
	agentID := c.agent.GetID()
	c.mu.RUnlock()

	if !connected {
		return fmt.Errorf("not connected to control plane")
	}

	// Convert display statuses to protobuf format
	pbDisplays := make([]*pb.DisplayStatus, 0, len(status.Displays))
	for _, display := range status.Displays {
		pbDisplays = append(pbDisplays, &pb.DisplayStatus{
			DisplayId:   display.DisplayID,
			CurrentView: display.CurrentView,
			CurrentFps:  int32(display.CurrentFPS),
			Active:      display.Active,
		})
	}

	req := &pb.StatusUpdate{
		AgentId:   agentID,
		Timestamp: time.Now().Unix(),
		Status: &pb.AgentStatus{
			Health:           status.Health,
			Displays:         pbDisplays,
			UptimeSeconds:    status.UptimeSeconds,
			MemoryUsageBytes: status.MemoryUsage,
			ErrorMessage:     status.ErrorMessage,
		},
	}

	resp, err := client.UpdateStatus(ctx, req)
	if err != nil {
		return fmt.Errorf("status update failed: %w", err)
	}

	if !resp.Success {
		return fmt.Errorf("status update rejected")
	}

	log.Printf("[ControlPlaneClient] Status update sent")
	return nil
}

// IsConnected returns whether the client is connected
func (c *ControlPlaneClient) IsConnected() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.connected
}

// StartCommandStream begins listening for commands from control plane
func (c *ControlPlaneClient) StartCommandStream(ctx context.Context) (<-chan *pb.ControlPlaneMessage, error) {
	c.streamMu.Lock()
	defer c.streamMu.Unlock()

	c.mu.RLock()
	connected := c.connected
	client := c.client
	c.mu.RUnlock()

	if !connected {
		return nil, fmt.Errorf("not connected to control plane")
	}

	// Create stream context
	c.streamCtx, c.streamCancel = context.WithCancel(ctx)

	// Establish bidirectional stream
	stream, err := client.StreamCommands(c.streamCtx)
	if err != nil {
		c.streamCancel()
		return nil, fmt.Errorf("failed to create command stream: %w", err)
	}

	c.commandStream = stream

	// Create channel for commands
	commandChan := make(chan *pb.ControlPlaneMessage, 10)

	// Start goroutine to receive commands
	go func() {
		defer close(commandChan)
		for {
			msg, err := stream.Recv()
			if err != nil {
				if err == io.EOF {
					log.Printf("[ControlPlaneClient] Command stream closed by server")
					return
				}
				if c.streamCtx.Err() != nil {
					// Context cancelled, normal shutdown
					return
				}
				log.Printf("[ControlPlaneClient] Error receiving command: %v", err)
				return
			}

			select {
			case commandChan <- msg:
			case <-c.streamCtx.Done():
				return
			}
		}
	}()

	log.Printf("[ControlPlaneClient] Command stream started")
	return commandChan, nil
}

// SendCommandResponse sends a response for an executed command
func (c *ControlPlaneClient) SendCommandResponse(ctx context.Context, success bool, message, commandType string) error {
	c.streamMu.Lock()
	stream := c.commandStream
	c.streamMu.Unlock()

	if stream == nil {
		return fmt.Errorf("command stream not active")
	}

	c.mu.RLock()
	agentID := c.agent.GetID()
	c.mu.RUnlock()

	msg := &pb.AgentMessage{
		AgentId: agentID,
		Payload: &pb.AgentMessage_Response{
			Response: &pb.CommandResponse{
				Success:     success,
				Message:     message,
				CommandType: commandType,
			},
		},
	}

	if err := stream.Send(msg); err != nil {
		return fmt.Errorf("failed to send command response: %w", err)
	}

	return nil
}

// Reconnect attempts to reconnect to the control plane with exponential backoff
func (c *ControlPlaneClient) Reconnect(ctx context.Context, maxRetries int) error {
	delay := 1 * time.Second
	maxDelay := 30 * time.Second

	for attempt := 1; attempt <= maxRetries; attempt++ {
		log.Printf("[ControlPlaneClient] Reconnection attempt %d/%d", attempt, maxRetries)

		// Try to connect
		if err := c.Connect(ctx); err != nil {
			log.Printf("[ControlPlaneClient] Reconnection attempt %d failed: %v", attempt, err)

			// Check if context is cancelled
			if ctx.Err() != nil {
				return fmt.Errorf("reconnection cancelled: %w", ctx.Err())
			}

			// Wait before retry with exponential backoff
			select {
			case <-time.After(delay):
				// Increase delay for next attempt (exponential backoff)
				delay *= 2
				if delay > maxDelay {
					delay = maxDelay
				}
			case <-ctx.Done():
				return fmt.Errorf("reconnection cancelled: %w", ctx.Err())
			}

			continue
		}

		// Reconnection successful
		log.Printf("[ControlPlaneClient] Reconnected successfully on attempt %d", attempt)
		return nil
	}

	return fmt.Errorf("failed to reconnect after %d attempts", maxRetries)
}
