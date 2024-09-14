package api

import (
	"encoding/json"
	"fmt"
	"github.com/6ixisgood/matrix-ticker/pkg/view"
	viewCommon "github.com/6ixisgood/matrix-ticker/pkg/view/common"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
	Server.router.GET("/views/configSpecs", getAllViewConfigSpecs)
	Server.router.GET("/views/definitions", getAllViewDefinitions)
	Server.router.POST("/views/definitions", saveViewDefinition)
	Server.router.GET("/views/definitions/:id", getViewDefinition)
	Server.router.DELETE("/views/definitions/:id", deleteViewDefinition)
	Server.router.GET("/views/:id", getViewById)
	Server.router.POST("/display/:id", displayViewById)
}

func getAllViewDefinitions(c *gin.Context) {
	definitions, err := viewCommon.GetAllViewDefinitions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error retrieving view definitions"})
		return
	}
	c.JSON(http.StatusOK, definitions)
}

func getViewDefinition(c *gin.Context) {
	id := c.Param("id")
	definition, err := viewCommon.GetViewDefinition(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "View definition not found"})
		return
	}
	c.JSON(http.StatusOK, definition)
}

func saveViewDefinition(c *gin.Context) {
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
	log.Printf(string(body.Config))
	if err := json.Unmarshal(body.Config, &configInstance); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad view config passed"})
		return
	}

	// Generate a UUID if the definition doesn't have an ID
	if body.Id == "" {
		body.Id = uuid.New().String()
	}

	definition := viewCommon.ViewDefinition{
		Id:     body.Id,
		Name:   body.Name,
		Type:   body.Type,
		Config: configInstance,
	}

	err := viewCommon.SaveViewDefinition(definition)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error saving view definition"})
		return
	}
	// Return the ID of the saved definition to the client
	c.JSON(http.StatusOK, gin.H{"message": "View definition saved successfully", "id": definition.Id})
}

// deleteViewDefinition handler function
func deleteViewDefinition(c *gin.Context) {
	id := c.Param("id") // Extract the ID from the URL parameter

	err := viewCommon.DeleteViewDefinition(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("Unable to delete ID: %s", id)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "View definition deleted successfully"})
}

func getAllViewConfigSpecs(c *gin.Context) {
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

func displayViewById(c *gin.Context) {
	viewDefinitionId := c.Param("id")

	// was ID given?
	if viewDefinitionId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "No valid ID provided"})
		return
	}

	// fetch by ID
	viewDefinition, err := viewCommon.GetViewDefinition(viewDefinitionId)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			gin.H{"message": fmt.Sprintf("Error fetching View Defintion with ID: %s", viewDefinitionId)})
		return
	}

	regView, exists := viewCommon.RegisteredViews[viewDefinition.Type]
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"message": "View type does not exist"})
		return
	}

	// create new view and trigger
	newView, err := regView.NewView(viewDefinition.Config)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to create view with saved config", "error": err.Error()})
		return
	}

	log.Printf("Initializing the %s view", viewDefinition.Id)
	animation := view.GetAnimation()
	animation.Init(newView)

	c.JSON(http.StatusOK, gin.H{"Status": "Created"})
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
