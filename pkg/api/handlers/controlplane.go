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
			"id":             agent.ID,
			"name":           agent.Name,
			"capabilities":   formatCapabilities(agent.Capabilities),
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

	agentID := c.Param("agent_id")
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
		"id":             agent.ID,
		"name":           agent.Name,
		"capabilities":   formatCapabilities(agent.Capabilities),
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

// formatCapabilities converts protobuf capabilities to JSON-friendly format
func formatCapabilities(caps *pb.Capabilities) gin.H {
	if caps == nil {
		return gin.H{
			"displays":        []gin.H{},
			"supported_views": []string{},
		}
	}

	displays := make([]gin.H, 0, len(caps.Displays))
	for _, display := range caps.Displays {
		displays = append(displays, gin.H{
			"display_id":   display.DisplayId,
			"display_type": display.DisplayType,
			"width":        display.Width,
			"height":       display.Height,
			"max_fps":      display.MaxFps,
			"color_depth":  display.ColorDepth,
		})
	}

	return gin.H{
		"displays":        displays,
		"supported_views": caps.SupportedViews,
	}
}

// formatStatus converts protobuf status to JSON-friendly format
func formatStatus(status *pb.AgentStatus) gin.H {
	if status == nil {
		return gin.H{
			"health":   "unknown",
			"displays": []gin.H{},
		}
	}

	displays := make([]gin.H, 0, len(status.Displays))
	for _, display := range status.Displays {
		displays = append(displays, gin.H{
			"display_id":   display.DisplayId,
			"current_view": display.CurrentView,
			"current_fps":  display.CurrentFps,
			"active":       display.Active,
		})
	}

	return gin.H{
		"health":             status.Health,
		"displays":           displays,
		"uptime_seconds":     status.UptimeSeconds,
		"memory_usage_bytes": status.MemoryUsageBytes,
		"error_message":      status.ErrorMessage,
	}
}

// GetAgentDisplays returns all displays for a specific agent
func GetAgentDisplays(c *gin.Context) {
	if registryInstance == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "registry not available"})
		return
	}

	agentID := c.Param("agent_id")
	agent, err := registryInstance.Get(agentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "agent not found"})
		return
	}

	if agent.Capabilities == nil {
		c.JSON(http.StatusOK, gin.H{"displays": []gin.H{}})
		return
	}

	displays := make([]gin.H, 0, len(agent.Capabilities.Displays))
	for _, display := range agent.Capabilities.Displays {
		// Find corresponding status for this display
		var displayStatus *pb.DisplayStatus
		if agent.Status != nil {
			for _, status := range agent.Status.Displays {
				if status.DisplayId == display.DisplayId {
					displayStatus = status
					break
				}
			}
		}

		displayInfo := gin.H{
			"display_id":   display.DisplayId,
			"display_type": display.DisplayType,
			"width":        display.Width,
			"height":       display.Height,
			"max_fps":      display.MaxFps,
			"color_depth":  display.ColorDepth,
		}

		if displayStatus != nil {
			displayInfo["current_view"] = displayStatus.CurrentView
			displayInfo["current_fps"] = displayStatus.CurrentFps
			displayInfo["active"] = displayStatus.Active
		}

		displays = append(displays, displayInfo)
	}

	c.JSON(http.StatusOK, gin.H{"displays": displays})
}

// GetAgentDisplay returns information about a specific display on an agent
func GetAgentDisplay(c *gin.Context) {
	if registryInstance == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "registry not available"})
		return
	}

	agentID := c.Param("agent_id")
	displayID := c.Param("display_id")

	agent, err := registryInstance.Get(agentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "agent not found"})
		return
	}

	if agent.Capabilities == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "display not found"})
		return
	}

	// Find the display in capabilities
	var displayCap *pb.DisplayCapability
	for _, display := range agent.Capabilities.Displays {
		if display.DisplayId == displayID {
			displayCap = display
			break
		}
	}

	if displayCap == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "display not found"})
		return
	}

	// Find corresponding status
	var displayStatus *pb.DisplayStatus
	if agent.Status != nil {
		for _, status := range agent.Status.Displays {
			if status.DisplayId == displayID {
				displayStatus = status
				break
			}
		}
	}

	response := gin.H{
		"display_id":   displayCap.DisplayId,
		"display_type": displayCap.DisplayType,
		"width":        displayCap.Width,
		"height":       displayCap.Height,
		"max_fps":      displayCap.MaxFps,
		"color_depth":  displayCap.ColorDepth,
	}

	if displayStatus != nil {
		response["current_view"] = displayStatus.CurrentView
		response["current_fps"] = displayStatus.CurrentFps
		response["active"] = displayStatus.Active
	}

	c.JSON(http.StatusOK, response)
}

// AssignViewToDisplay assigns a view to a specific display on an agent
func AssignViewToDisplay(c *gin.Context) {
	if registryInstance == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "registry not available"})
		return
	}

	agentID := c.Param("agent_id")
	displayID := c.Param("display_id")

	var request struct {
		ViewID string `json:"view_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify agent exists
	agent, err := registryInstance.Get(agentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "agent not found"})
		return
	}

	// Verify display exists
	displayFound := false
	if agent.Capabilities != nil {
		for _, display := range agent.Capabilities.Displays {
			if display.DisplayId == displayID {
				displayFound = true
				break
			}
		}
	}

	if !displayFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "display not found on agent"})
		return
	}

	// TODO: Send command to agent via gRPC command stream
	// For now, return accepted - actual implementation needs command streaming
	c.JSON(http.StatusAccepted, gin.H{
		"message":    "view assignment command queued",
		"agent_id":   agentID,
		"display_id": displayID,
		"view_id":    request.ViewID,
	})
}

// GetDisplayCurrentView returns the current view assigned to a specific display
func GetDisplayCurrentView(c *gin.Context) {
	if registryInstance == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "registry not available"})
		return
	}

	agentID := c.Param("agent_id")
	displayID := c.Param("display_id")

	agent, err := registryInstance.Get(agentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "agent not found"})
		return
	}

	// Find the display status
	var displayStatus *pb.DisplayStatus
	if agent.Status != nil {
		for _, status := range agent.Status.Displays {
			if status.DisplayId == displayID {
				displayStatus = status
				break
			}
		}
	}

	if displayStatus == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "display not found or no status available"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"display_id":   displayID,
		"current_view": displayStatus.CurrentView,
		"current_fps":  displayStatus.CurrentFps,
		"active":       displayStatus.Active,
	})
}
