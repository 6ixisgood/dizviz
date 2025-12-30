package controlplane

import (
	"fmt"
	"log"
	"sync"
	"time"

	pb "github.com/6ixisgood/matrix-ticker/proto"
	"github.com/google/uuid"
)

// AgentInfo holds information about a registered agent
type AgentInfo struct {
	ID            string
	Name          string
	Capabilities  *pb.Capabilities
	Version       string
	RegisteredAt  time.Time
	LastHeartbeat time.Time
	Status        *pb.AgentStatus
	CommandStream pb.AgentService_StreamCommandsServer
	mu            sync.RWMutex
}

// AgentRegistry manages registered agents
type AgentRegistry struct {
	agents map[string]*AgentInfo
	mu     sync.RWMutex
}

// NewAgentRegistry creates a new agent registry
func NewAgentRegistry() *AgentRegistry {
	return &AgentRegistry{
		agents: make(map[string]*AgentInfo),
	}
}

// Register registers a new agent and returns its assigned ID
func (r *AgentRegistry) Register(name string, caps *pb.Capabilities, version string) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Generate unique agent ID
	agentID := uuid.New().String()

	now := time.Now()
	info := &AgentInfo{
		ID:            agentID,
		Name:          name,
		Capabilities:  caps,
		Version:       version,
		RegisteredAt:  now,
		LastHeartbeat: now,
		Status: &pb.AgentStatus{
			Health:        "healthy",
			CurrentView:   "",
			CurrentFps:    0,
			UptimeSeconds: 0,
		},
	}

	r.agents[agentID] = info
	log.Printf("[Registry] Agent registered: %s (%s) - Display: %s (%dx%d)",
		name, agentID, caps.DisplayType, caps.Width, caps.Height)

	return agentID, nil
}

// Unregister removes an agent from the registry
func (r *AgentRegistry) Unregister(agentID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.agents[agentID]; !exists {
		return fmt.Errorf("agent not found: %s", agentID)
	}

	delete(r.agents, agentID)
	log.Printf("[Registry] Agent unregistered: %s", agentID)
	return nil
}

// Get retrieves agent information
func (r *AgentRegistry) Get(agentID string) (*AgentInfo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	info, exists := r.agents[agentID]
	if !exists {
		return nil, fmt.Errorf("agent not found: %s", agentID)
	}

	return info, nil
}

// List returns all registered agents
func (r *AgentRegistry) List() []*AgentInfo {
	r.mu.RLock()
	defer r.mu.RUnlock()

	agents := make([]*AgentInfo, 0, len(r.agents))
	for _, info := range r.agents {
		agents = append(agents, info)
	}

	return agents
}

// UpdateHeartbeat updates the last heartbeat time for an agent
func (r *AgentRegistry) UpdateHeartbeat(agentID string) error {
	r.mu.RLock()
	info, exists := r.agents[agentID]
	r.mu.RUnlock()

	if !exists {
		return fmt.Errorf("agent not found: %s", agentID)
	}

	info.mu.Lock()
	info.LastHeartbeat = time.Now()
	info.mu.Unlock()

	return nil
}

// UpdateStatus updates the status of an agent
func (r *AgentRegistry) UpdateStatus(agentID string, status *pb.AgentStatus) error {
	r.mu.RLock()
	info, exists := r.agents[agentID]
	r.mu.RUnlock()

	if !exists {
		return fmt.Errorf("agent not found: %s", agentID)
	}

	info.mu.Lock()
	info.Status = status
	info.mu.Unlock()

	return nil
}

// SetCommandStream sets the command stream for an agent
func (r *AgentRegistry) SetCommandStream(agentID string, stream pb.AgentService_StreamCommandsServer) error {
	r.mu.RLock()
	info, exists := r.agents[agentID]
	r.mu.RUnlock()

	if !exists {
		return fmt.Errorf("agent not found: %s", agentID)
	}

	info.mu.Lock()
	info.CommandStream = stream
	info.mu.Unlock()

	return nil
}

// CleanupStaleAgents removes agents that haven't sent heartbeats
func (r *AgentRegistry) CleanupStaleAgents(timeout time.Duration) {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	for agentID, info := range r.agents {
		info.mu.RLock()
		lastHeartbeat := info.LastHeartbeat
		info.mu.RUnlock()

		if now.Sub(lastHeartbeat) > timeout {
			log.Printf("[Registry] Removing stale agent: %s (last heartbeat: %v ago)",
				agentID, now.Sub(lastHeartbeat))
			delete(r.agents, agentID)
		}
	}
}

// Count returns the number of registered agents
func (r *AgentRegistry) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.agents)
}

// GetAgentsByDisplayType returns agents with a specific display type
func (r *AgentRegistry) GetAgentsByDisplayType(displayType string) []*AgentInfo {
	r.mu.RLock()
	defer r.mu.RUnlock()

	agents := make([]*AgentInfo, 0)
	for _, info := range r.agents {
		if info.Capabilities != nil && info.Capabilities.DisplayType == displayType {
			agents = append(agents, info)
		}
	}

	return agents
}

// GetHealthyAgents returns agents with healthy status
func (r *AgentRegistry) GetHealthyAgents() []*AgentInfo {
	r.mu.RLock()
	defer r.mu.RUnlock()

	agents := make([]*AgentInfo, 0)
	for _, info := range r.agents {
		info.mu.RLock()
		health := info.Status.Health
		info.mu.RUnlock()

		if health == "healthy" {
			agents = append(agents, info)
		}
	}

	return agents
}

// GetAgentCount returns the total number of registered agents
func (r *AgentRegistry) GetAgentCount() int {
	return r.Count()
}

// GetHealthyAgentCount returns the number of healthy agents
func (r *AgentRegistry) GetHealthyAgentCount() int {
	r.mu.RLock()
	defer r.mu.RUnlock()

	count := 0
	for _, info := range r.agents {
		info.mu.RLock()
		health := info.Status.Health
		info.mu.RUnlock()

		if health == "healthy" {
			count++
		}
	}

	return count
}
