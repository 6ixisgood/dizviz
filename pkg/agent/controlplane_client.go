package agent

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// ControlPlaneClient handles communication with the control plane
type ControlPlaneClient struct {
	addr  string
	agent *Agent
	conn  *grpc.ClientConn
	// client     pb.AgentServiceClient  // Will be added when we generate protobuf code
	connected bool
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
	c.connected = true

	// TODO: Create gRPC client when proto is compiled
	// c.client = pb.NewAgentServiceClient(conn)

	log.Printf("[ControlPlaneClient] Connected successfully")
	return nil
}

// Disconnect closes the connection to the control plane
func (c *ControlPlaneClient) Disconnect() {
	if c.conn != nil {
		c.conn.Close()
		c.connected = false
		log.Printf("[ControlPlaneClient] Disconnected")
	}
}

// Register registers the agent with the control plane
func (c *ControlPlaneClient) Register(ctx context.Context, name string, caps Capabilities, version string) (string, error) {
	if !c.connected {
		return "", fmt.Errorf("not connected to control plane")
	}

	log.Printf("[ControlPlaneClient] Registering agent: %s", name)

	// TODO: Implement actual gRPC registration call
	// For now, return a mock ID
	mockID := fmt.Sprintf("agent-%d", time.Now().Unix())

	log.Printf("[ControlPlaneClient] Registration successful, ID: %s", mockID)
	return mockID, nil
}

// SendHeartbeat sends a heartbeat to the control plane
func (c *ControlPlaneClient) SendHeartbeat(ctx context.Context) error {
	if !c.connected {
		return fmt.Errorf("not connected to control plane")
	}

	// TODO: Implement actual gRPC heartbeat call
	// For now, just log
	// log.Printf("[ControlPlaneClient] Heartbeat sent")

	return nil
}

// UpdateStatus sends a status update to the control plane
func (c *ControlPlaneClient) UpdateStatus(ctx context.Context, status Status) error {
	if !c.connected {
		return fmt.Errorf("not connected to control plane")
	}

	// TODO: Implement actual gRPC status update call
	log.Printf("[ControlPlaneClient] Status update sent")

	return nil
}

// IsConnected returns whether the client is connected
func (c *ControlPlaneClient) IsConnected() bool {
	return c.connected
}
