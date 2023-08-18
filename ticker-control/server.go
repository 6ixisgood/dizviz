package main

import (
	"log"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	comp "github.com/sixisgoood/matrix-ticker/components"


)

type AppServer struct {
	router *gin.Engine
}

func NewAppServer() *AppServer {
	return &AppServer{
		router: gin.Default(),
	}
}

func (s *AppServer) InitializeRoutes() {
	s.router.GET("/views", s.getAllViews)
	s.router.GET("/views/:id", s.getViewById)
	s.router.POST("/display/:id", s.displayView)
}

func (s *AppServer) getAllViews(c *gin.Context) {
	var views []string
	for view, _ := range comp.RegisteredViews {
		views = append(views, view)
	}
	c.JSON(http.StatusOK, views)
}

func (s *AppServer) getViewById(c *gin.Context) {
	id := c.Param("id")
	if id != "" {
		c.JSON(http.StatusOK, "")
		return
	}
	c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
}

func (s *AppServer) displayView(c *gin.Context) {
	id := c.Param("id")

	var config map[string]string
	if err := c.BindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad configuration passed"})
		return
	}

	log.Printf("Initializing the %s view", id)
	animation := GetAnimation()
	animation.Init(id, config)

	c.JSON(http.StatusOK, gin.H{"Status": "Created"})
}

func (s *AppServer) Run(port string) {
	s.router.Run(fmt.Sprintf("127.0.0.1:%s", port))
}