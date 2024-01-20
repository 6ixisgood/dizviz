package api

import (
	"log"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/6ixisgood/matrix-ticker/pkg/view"
	viewCommon "github.com/6ixisgood/matrix-ticker/pkg/view/common"
)

type AppServer struct {
	router *gin.Engine
}

type AppServerConfig struct {
	AllowedHost	string
	Port		string
}

var (
	Server = &AppServer{
		router: gin.Default(),
	}
	Config = &AppServerConfig{}
)

func SetAppServerConfig(config *AppServerConfig) {
	Config = config
}

func InitializeRoutes() {
	Server.router.GET("/views", getAllViews)
	Server.router.GET("/views/:id", getViewById)
	Server.router.POST("/display/:id", displayView)
}

func getAllViews(c *gin.Context) {
	var views []string
	for view, _ := range viewCommon.RegisteredViews {
		views = append(views, view)
	}
	c.JSON(http.StatusOK, views)
}

func getViewById(c *gin.Context) {
	id := c.Param("id")
	if id != "" {
		c.JSON(http.StatusOK, "")
		return
	}
	c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
}

func displayView(c *gin.Context) {
	id := c.Param("id")

	var config map[string]string
	if err := c.BindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad configuration passed"})
		return
	}

	log.Printf("Initializing the %s view", id)
	animation := view.GetAnimation()
	animation.Init(id, config)

	c.JSON(http.StatusOK, gin.H{"Status": "Created"})
}

func Run() {
	InitializeRoutes()
	Server.router.Run(fmt.Sprintf("%s:%s", Config.AllowedHost, Config.Port))
}