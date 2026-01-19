package api

import (
	"github.com/6ixisgood/matrix-ticker/pkg/api/handlers"
	"github.com/6ixisgood/matrix-ticker/pkg/controlplane"
	"github.com/gin-gonic/gin"
)

// SetRegistry sets the agent registry instance for the API handlers
func SetRegistry(registry *controlplane.AgentRegistry) {
	handlers.SetRegistry(registry)
}

// SetServer sets the control plane server instance for the API handlers
func SetServer(server *controlplane.Server) {
	handlers.SetServer(server)
}

// SetStoreService sets the store service instance for the API handlers
func SetStoreService(storeService *controlplane.StoreService) {
	handlers.SetStoreService(storeService)
}

// Router sets up all the API routes for the control plane
func Router() *gin.Engine {
	engine := gin.Default()
	r := engine.Group("/")

	registerViewRoutes(r)
	registerControlPlaneRoutes(r)

	return engine
}

// registerViewRoutes sets up the /view routes for view definition management
func registerViewRoutes(r *gin.RouterGroup) {
	view := r.Group("/view")
	view.GET("/configSpecs", handlers.GetAllViewConfigSpecs)
	view.GET("/definitions", handlers.GetAllViewDefinitions)
	view.POST("/definitions", handlers.SaveViewDefinition)
	view.GET("/definitions/:id", handlers.GetViewDefinition)
	view.DELETE("/definitions/:id", handlers.DeleteViewDefinition)
	view.GET("/:id", handlers.GetViewById)
}

// registerControlPlaneRoutes sets up the /controlplane routes
func registerControlPlaneRoutes(r *gin.RouterGroup) {
	cp := r.Group("/controlplane")

	// Control plane stats and health (no params)
	cp.GET("/stats", handlers.GetControlPlaneStats)
	cp.GET("/health", handlers.GetControlPlaneHealth)

	// Agent management
	cp.GET("/agents", handlers.ListAgents)
	cp.GET("/agents/:agent_id", handlers.GetAgent)

	// Display management for specific agent
	cp.GET("/agents/:agent_id/displays", handlers.GetAgentDisplays)
	cp.GET("/agents/:agent_id/displays/:display_id", handlers.GetAgentDisplay)

	// View assignment to specific display
	cp.POST("/agents/:agent_id/displays/:display_id/view", handlers.AssignViewToDisplay)
	cp.GET("/agents/:agent_id/displays/:display_id/view", handlers.GetDisplayCurrentView)
}
