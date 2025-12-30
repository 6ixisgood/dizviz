package api

import (
	"github.com/6ixisgood/matrix-ticker/pkg/api/handlers"
	"github.com/6ixisgood/matrix-ticker/pkg/app"
	"github.com/6ixisgood/matrix-ticker/pkg/controlplane"
	"github.com/gin-gonic/gin"
)

// SetApplication sets the application instance for the API handlers
func SetApplication(application *app.Application) {
	handlers.SetApplication(application)
}

// SetRegistry sets the agent registry instance for the API handlers
func SetRegistry(registry *controlplane.AgentRegistry) {
	handlers.SetRegistry(registry)
}

// RegisterRoutes sets up all the API routes
func Router() *gin.Engine {
	engine := gin.Default()
	r := engine.Group("/")
	//r.Use(middleware.AuthMiddleware())

	registerViewRoutes(r)
	registerDisplayRoutes(r)
	registerControlPlaneRoutes(r)

	return engine
}

// registerViewRoutes sets up the /view routes
func registerViewRoutes(r *gin.RouterGroup) {
	view := r.Group("/view")
	view.GET("/configSpecs", handlers.GetAllViewConfigSpecs)
	view.GET("/definitions", handlers.GetAllViewDefinitions)
	view.POST("/definitions", handlers.SaveViewDefinition)
	view.GET("/definitions/:id", handlers.GetViewDefinition)
	view.DELETE("/definitions/:id", handlers.DeleteViewDefinition)
	view.GET("/:id", handlers.GetViewById)
}

// registerDisplayRoutes sets up the /display routes
func registerDisplayRoutes(r *gin.RouterGroup) {
	display := r.Group("/display")
	display.POST("/:id", handlers.DisplayViewById)
}

// registerControlPlaneRoutes sets up the /controlplane routes
func registerControlPlaneRoutes(r *gin.RouterGroup) {
	cp := r.Group("/controlplane")
	cp.GET("/agents", handlers.ListAgents)
	cp.GET("/agents/:id", handlers.GetAgent)
	cp.GET("/stats", handlers.GetControlPlaneStats)
	cp.GET("/health", handlers.GetControlPlaneHealth)
}
