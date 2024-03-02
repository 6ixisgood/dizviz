package api

import (
	"encoding/json"
	"fmt"
	"github.com/6ixisgood/matrix-ticker/pkg/view"
	viewCommon "github.com/6ixisgood/matrix-ticker/pkg/view/common"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type AppServer struct {
	router *gin.Engine
}

type AppServerConfig struct {
	AllowedHost string
	Port        string
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
	Server.router.GET("/views/configs", getAllViewConfigs)
	Server.router.GET("/views/:id", getViewById)
	Server.router.POST("/display", displayView)
}

func getAllViewConfigs(c *gin.Context) {
	configs := make(map[string]interface{})
	for name, regView := range viewCommon.RegisteredViews {
		configSpec := viewCommon.GenerateViewConfigSpecJson(regView.NewConfig())
		configs[name] = configSpec
	}
	c.JSON(http.StatusOK, configs)
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
	var body viewCommon.ViewDefinitionRaw
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad request body"})
		return
	}

	regView, exists := viewCommon.RegisteredViews[body.Type]
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"message": "View type does not exist"})
		return
	}

	configInstance := regView.NewConfig()
	if err := json.Unmarshal(body.Config, &configInstance); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad view config passed"})
		return
	}

	newView, err := regView.NewView(configInstance)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to create view with given config", "error": err.Error()})
		return
	}

	log.Printf("Initializing the %s view", body.Type)
	animation := view.GetAnimation()
	animation.Init(newView)

	c.JSON(http.StatusOK, gin.H{"Status": "Created"})
}

func Run() {
	InitializeRoutes()
	Server.router.Run(fmt.Sprintf("%s:%s", Config.AllowedHost, Config.Port))
}
