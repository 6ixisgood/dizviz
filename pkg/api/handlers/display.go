package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	viewCommon "github.com/6ixisgood/matrix-ticker/pkg/view/common"
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
	
	// Use the application to change views
	app := GetApplication()
	if app == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Application not initialized"})
		return
	}

	if err := app.ChangeView(newView); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to set view", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Status": "Created"})
}

func PostView(c *gin.Context) {
	log.Printf("POST /view")

	var body viewCommon.ViewDefinition
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad request body"})
		return
	}

	regView, exists := viewCommon.RegisteredViews[body.Type]
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"message": "View type does not exist"})
		return
	}

	newView, err := regView.NewView(body.Config)
	if err != nil {
		log.Printf("Failed to create view of type %s with given config. Error: %s", body.Type, err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad view config passed"})
		return
	}

	// Use the application to change views
	app := GetApplication()
	if app == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Application not initialized"})
		return
	}

	if err := app.ChangeView(newView); err != nil {
		log.Printf("Failed to change view: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to change view"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "View updated successfully"})
}