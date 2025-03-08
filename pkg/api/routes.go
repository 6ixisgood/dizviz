package api

import (
	"github.com/6ixisgood/matrix-ticker/pkg/api/handlers"
	"github.com/gin-gonic/gin"
)


// RegisterRoutes sets up all the API routes
func Router() *gin.Engine {
	engine := gin.Default()
	r := engine.Group("/")
	//r.Use(middleware.AuthMiddleware())

	registerViewRoutes(r)
	registerDisplayRoutes(r)

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
