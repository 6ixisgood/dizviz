package controlplane

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	pb "github.com/6ixisgood/disco/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// DataSourceConfig holds configuration for external data sources
type DataSourceConfig struct {
	Sleeper struct {
		BaseUrl string
	}
	SportsFeed struct {
		BaseUrl  string
		Username string
		Password string
	}
	Weather struct {
		BaseUrl string
		Key     string
	}
}

// Server is the control plane gRPC server
type Server struct {
	pb.UnimplementedAgentServiceServer

	addr           string
	registry       *AgentRegistry
	storeService   *StoreService
	dataSourceCfg  *DataSourceConfig
	grpcSrv        *grpc.Server
	commandStreams map[string]pb.AgentService_StreamCommandsServer // agentID -> stream
	mu             sync.RWMutex
	running        bool
}

// NewServer creates a new control plane server
func NewServer(addr string, storeService *StoreService, dataSourceCfg *DataSourceConfig) *Server {
	return &Server{
		addr:           addr,
		registry:       NewAgentRegistry(),
		storeService:   storeService,
		dataSourceCfg:  dataSourceCfg,
		commandStreams: make(map[string]pb.AgentService_StreamCommandsServer),
	}
}

// Start starts the gRPC server
func (s *Server) Start() error {
	s.mu.Lock()
	if s.running {
		s.mu.Unlock()
		return fmt.Errorf("server already running")
	}
	s.running = true
	s.mu.Unlock()

	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	s.grpcSrv = grpc.NewServer()
	pb.RegisterAgentServiceServer(s.grpcSrv, s)

	log.Printf("[ControlPlane] Server starting on %s", s.addr)

	// Start cleanup routine
	go s.cleanupRoutine()

	if err := s.grpcSrv.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %w", err)
	}

	return nil
}

// Stop gracefully stops the server
func (s *Server) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.running {
		return
	}

	log.Printf("[ControlPlane] Server stopping...")
	if s.grpcSrv != nil {
		s.grpcSrv.GracefulStop()
	}
	s.running = false
	log.Printf("[ControlPlane] Server stopped")
}

// RegisterAgent handles agent registration
func (s *Server) RegisterAgent(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	log.Printf("[ControlPlane] Registration request from: %s", req.AgentName)

	if req.AgentName == "" {
		return &pb.RegisterResponse{
			Success: false,
			Message: "agent name is required",
		}, status.Error(codes.InvalidArgument, "agent name is required")
	}

	if req.Capabilities == nil {
		return &pb.RegisterResponse{
			Success: false,
			Message: "capabilities are required",
		}, status.Error(codes.InvalidArgument, "capabilities are required")
	}

	// Register the agent
	agentID, err := s.registry.Register(req.AgentName, req.Capabilities, req.Version)
	if err != nil {
		return &pb.RegisterResponse{
			Success: false,
			Message: err.Error(),
		}, status.Error(codes.Internal, err.Error())
	}

	log.Printf("[ControlPlane] Agent registered: %s (ID: %s)", req.AgentName, agentID)

	// Build runtime config from data source configuration
	runtimeConfig := s.buildRuntimeConfig()

	return &pb.RegisterResponse{
		AgentId:       agentID,
		Success:       true,
		Message:       "registration successful",
		RuntimeConfig: runtimeConfig,
	}, nil
}

// buildRuntimeConfig creates AgentRuntimeConfig from the control plane's data source configuration
func (s *Server) buildRuntimeConfig() *pb.AgentRuntimeConfig {
	if s.dataSourceCfg == nil {
		return nil
	}

	dataSources := make(map[string]*pb.DataSourceConfig)

	// Sleeper configuration
	if s.dataSourceCfg.Sleeper.BaseUrl != "" {
		dataSources["sleeper"] = &pb.DataSourceConfig{
			Type: "sleeper",
			Config: map[string]string{
				"base_url": s.dataSourceCfg.Sleeper.BaseUrl,
			},
		}
	}

	// SportsFeed configuration
	if s.dataSourceCfg.SportsFeed.BaseUrl != "" {
		dataSources["sportsfeed"] = &pb.DataSourceConfig{
			Type: "sportsfeed",
			Config: map[string]string{
				"base_url": s.dataSourceCfg.SportsFeed.BaseUrl,
				"username": s.dataSourceCfg.SportsFeed.Username,
				"password": s.dataSourceCfg.SportsFeed.Password,
			},
		}
	}

	// Weather configuration
	if s.dataSourceCfg.Weather.BaseUrl != "" {
		dataSources["weather"] = &pb.DataSourceConfig{
			Type: "weather",
			Config: map[string]string{
				"base_url": s.dataSourceCfg.Weather.BaseUrl,
				"key":      s.dataSourceCfg.Weather.Key,
			},
		}
	}

	return &pb.AgentRuntimeConfig{
		DataSources: dataSources,
	}
}

// StreamCommands establishes a bidirectional stream for commands and status
func (s *Server) StreamCommands(stream pb.AgentService_StreamCommandsServer) error {
	log.Printf("[ControlPlane] New command stream established")

	var agentID string
	defer func() {
		if agentID != "" {
			s.mu.Lock()
			delete(s.commandStreams, agentID)
			s.mu.Unlock()
			log.Printf("[ControlPlane] Command stream closed for agent %s", agentID)
		}
	}()

	// Wait for messages from agent
	for {
		msg, err := stream.Recv()
		if err != nil {
			log.Printf("[ControlPlane] Stream error: %v", err)
			return err
		}

		// Handle agent messages
		if msg.AgentId == "" {
			log.Printf("[ControlPlane] Received message without agent ID")
			continue
		}

		// Register the stream for this agent on first message
		if agentID == "" {
			agentID = msg.AgentId
			s.mu.Lock()
			s.commandStreams[agentID] = stream
			s.mu.Unlock()
			log.Printf("[ControlPlane] Registered command stream for agent %s", agentID)
		}

		switch payload := msg.Payload.(type) {
		case *pb.AgentMessage_Status:
			s.handleStatusUpdate(msg.AgentId, payload.Status)
		case *pb.AgentMessage_Response:
			s.handleCommandResponse(msg.AgentId, payload.Response)
		default:
			log.Printf("[ControlPlane] Unknown message type from agent %s", msg.AgentId)
		}
	}
}

// Heartbeat handles periodic heartbeats from agents
func (s *Server) Heartbeat(ctx context.Context, req *pb.HeartbeatRequest) (*pb.HeartbeatResponse, error) {
	if req.AgentId == "" {
		return &pb.HeartbeatResponse{Success: false}, status.Error(codes.InvalidArgument, "agent ID is required")
	}

	// Update last heartbeat time
	if err := s.registry.UpdateHeartbeat(req.AgentId); err != nil {
		log.Printf("[ControlPlane] Heartbeat from unknown agent: %s", req.AgentId)
		return &pb.HeartbeatResponse{Success: false}, status.Error(codes.NotFound, "agent not found")
	}

	return &pb.HeartbeatResponse{Success: true}, nil
}

// UpdateStatus handles status updates from agents
func (s *Server) UpdateStatus(ctx context.Context, req *pb.StatusUpdate) (*pb.StatusResponse, error) {
	if req.AgentId == "" {
		return &pb.StatusResponse{Success: false}, status.Error(codes.InvalidArgument, "agent ID is required")
	}

	if req.Status == nil {
		return &pb.StatusResponse{Success: false}, status.Error(codes.InvalidArgument, "status is required")
	}

	// Update agent status
	if err := s.registry.UpdateStatus(req.AgentId, req.Status); err != nil {
		log.Printf("[ControlPlane] Status update from unknown agent: %s", req.AgentId)
		return &pb.StatusResponse{Success: false}, status.Error(codes.NotFound, "agent not found")
	}

	// Log status update with display count
	displayCount := len(req.Status.Displays)
	log.Printf("[ControlPlane] Status update from %s: health=%s, displays=%d", req.AgentId, req.Status.Health, displayCount)

	return &pb.StatusResponse{Success: true}, nil
}

// handleStatusUpdate processes status updates from agents
func (s *Server) handleStatusUpdate(agentID string, update *pb.StatusUpdate) {
	if update == nil || update.Status == nil {
		return
	}

	if err := s.registry.UpdateStatus(agentID, update.Status); err != nil {
		log.Printf("[ControlPlane] Failed to update status for agent %s: %v", agentID, err)
		return
	}

	log.Printf("[ControlPlane] Status update from %s: health=%s", agentID, update.Status.Health)
}

// handleCommandResponse processes command responses from agents
func (s *Server) handleCommandResponse(agentID string, response *pb.CommandResponse) {
	if response == nil {
		return
	}

	log.Printf("[ControlPlane] Command response from %s: type=%s, success=%v, message=%s",
		agentID, response.CommandType, response.Success, response.Message)
}

// SendCommandToAgent sends a command to a specific agent via its active stream
func (s *Server) SendCommandToAgent(agentID string, cmd *pb.ControlPlaneMessage) error {
	s.mu.RLock()
	stream, exists := s.commandStreams[agentID]
	s.mu.RUnlock()

	if !exists {
		return fmt.Errorf("no active command stream for agent %s", agentID)
	}

	log.Printf("[ControlPlane] Sending command to agent %s", agentID)

	if err := stream.Send(cmd); err != nil {
		// Remove the stream if send fails
		s.mu.Lock()
		delete(s.commandStreams, agentID)
		s.mu.Unlock()
		return fmt.Errorf("failed to send command to agent %s: %w", agentID, err)
	}

	return nil
}

// cleanupRoutine removes stale agents
func (s *Server) cleanupRoutine() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		s.mu.RLock()
		running := s.running
		s.mu.RUnlock()

		if !running {
			return
		}

		<-ticker.C
		s.registry.CleanupStaleAgents(60 * time.Second)
	}
}

// GetRegistry returns the agent registry
func (s *Server) GetRegistry() *AgentRegistry {
	return s.registry
}

// GetStoreService returns the store service
func (s *Server) GetStoreService() *StoreService {
	return s.storeService
}
