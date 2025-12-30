package handlers

import (
	"net/http"
	"time"

	"github.com/6ixisgood/matrix-ticker/pkg/controlplane"
	pb "github.com/6ixisgood/matrix-ticker/proto"
	"github.com/gin-gonic/gin"
)

// Global registry instance for control plane handlers
var registryInstance *controlplane.AgentRegistry

// SetRegistry sets the agent registry instance for handlers to use
func SetRegistry(registry *controlplane.AgentRegistry) {
	registryInstance = registry
}

// GetRegistry returns the current registry instance
func GetRegistry() *controlplane.AgentRegistry {
	return registryInstance
}

// ListAgents returns all registered agents
func ListAgents(c *gin.Context) {
	if registryInstance == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "registry not available"})
		return
	}

	agents := registryInstance.List()

	response := make([]gin.H, 0, len(agents))
	for _, agent := range agents {
		response = append(response, gin.H{
			"id":   agent.ID,
			"name": agent.Name,
			"capabilities": gin.H{
				"display_type": agent.Capabilities.DisplayType,
				"width":        agent.Capabilities.Width,
				"height":       agent.Capabilities.Height,
				"max_fps":      agent.Capabilities.MaxFps,
			},
			"version":        agent.Version,
			"status":         formatStatus(agent.Status),
			"registered_at":  agent.RegisteredAt.Format(time.RFC3339),
			"last_heartbeat": agent.LastHeartbeat.Format(time.RFC3339),
		})
	}

	c.JSON(http.StatusOK, response)
}

// GetAgent returns information about a specific agent
func GetAgent(c *gin.Context) {
	if registryInstance == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "registry not available"})
		return
	}

	agentID := c.Param("id")
	if agentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "agent ID required"})
		return
	}

	agent, err := registryInstance.Get(agentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "agent not found"})
		return
	}

	response := gin.H{
		"id":   agent.ID,
		"name": agent.Name,
		"capabilities": gin.H{
			"display_type":    agent.Capabilities.DisplayType,
			"width":           agent.Capabilities.Width,
			"height":          agent.Capabilities.Height,
			"max_fps":         agent.Capabilities.MaxFps,
			"color_depth":     agent.Capabilities.ColorDepth,
			"supported_views": agent.Capabilities.SupportedViews,
		},
		"version":        agent.Version,
		"status":         formatStatus(agent.Status),
		"registered_at":  agent.RegisteredAt.Format(time.RFC3339),
		"last_heartbeat": agent.LastHeartbeat.Format(time.RFC3339),
	}

	c.JSON(http.StatusOK, response)
}

// GetControlPlaneStats returns control plane statistics
func GetControlPlaneStats(c *gin.Context) {
	if registryInstance == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "registry not available"})
		return
	}

	response := gin.H{
		"total_agents":   registryInstance.GetAgentCount(),
		"healthy_agents": registryInstance.GetHealthyAgentCount(),
		"timestamp":      time.Now().Format(time.RFC3339),
	}

	c.JSON(http.StatusOK, response)
}

// GetControlPlaneHealth returns server health status
func GetControlPlaneHealth(c *gin.Context) {
	response := gin.H{
		"status": "healthy",
		"time":   time.Now().Format(time.RFC3339),
	}
	c.JSON(http.StatusOK, response)
}

// formatStatus converts protobuf status to JSON-friendly format
func formatStatus(status *pb.AgentStatus) gin.H {
	if status == nil {
		return gin.H{
			"health": "unknown",
		}
	}

	return gin.H{
		"health":             status.Health,
		"current_view":       status.CurrentView,
		"current_fps":        status.CurrentFps,
		"uptime_seconds":     status.UptimeSeconds,
		"memory_usage_bytes": status.MemoryUsageBytes,
		"error_message":      status.ErrorMessage,
	}
}
