package handlers

import (
	"encoding/json"
	"fmt"
	viewCommon "github.com/6ixisgood/matrix-ticker/pkg/view/common"
	"log"
	"net/http"
)

func DisplayViewById(c *gin.Context) {
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

func DisplayView(c *gin.Context) {
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