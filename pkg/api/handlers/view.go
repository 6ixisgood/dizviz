package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	viewCommon "github.com/6ixisgood/matrix-ticker/pkg/view/common"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetAllViewDefinitions(c *gin.Context) {
	if storeServiceInstance == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"message": "Store service not available"})
		return
	}

	definitions, err := storeServiceInstance.GetAllViewDefinitions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error retrieving view definitions"})
		log.Println(err)
		return
	}
	c.JSON(http.StatusOK, definitions)
}

func GetViewDefinition(c *gin.Context) {
	if storeServiceInstance == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"message": "Store service not available"})
		return
	}

	id := c.Param("id")
	definition, err := storeServiceInstance.GetViewDefinition(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "View definition not found"})
		return
	}
	c.JSON(http.StatusOK, definition)
}

func SaveViewDefinition(c *gin.Context) {
	if storeServiceInstance == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"message": "Store service not available"})
		return
	}

	var body viewCommon.ViewDefinitionRaw
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad request body"})
		return
	}

	// Verify the view type exists
	_, exists := viewCommon.RegisteredViews[body.Type]
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"message": "View type does not exist"})
		return
	}

	// Validate the config JSON is valid JSON
	var configTest interface{}
	if err := json.Unmarshal(body.Config, &configTest); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid JSON in view config"})
		return
	}

	// Generate a UUID if the definition doesn't have an ID
	if body.Id == "" {
		body.Id = uuid.New().String()
	}

	// Store the definition with Config as json.RawMessage
	// The config will be validated when the view is actually created
	definition := viewCommon.ViewDefinition{
		Id:     body.Id,
		Name:   body.Name,
		Type:   body.Type,
		Config: body.Config, // Store as json.RawMessage
	}

	err := storeServiceInstance.SaveViewDefinition(definition)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error saving view definition"})
		return
	}
	// Return the ID of the saved definition to the client
	c.JSON(http.StatusOK, gin.H{"message": "View definition saved successfully", "id": definition.Id})
}

// DeleteViewDefinition handler function
func DeleteViewDefinition(c *gin.Context) {
	if storeServiceInstance == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"message": "Store service not available"})
		return
	}

	id := c.Param("id") // Extract the ID from the URL parameter

	err := storeServiceInstance.DeleteViewDefinition(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("Unable to delete ID: %s", id)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "View definition deleted successfully"})
}

func GetAllViewConfigSpecs(c *gin.Context) {
	configs := make(map[string]interface{})
	for name, factory := range viewCommon.RegisteredViews {
		// Create a view instance to introspect its config fields
		view := factory()

		// Generate spec from the view's structure (excluding BaseView context fields)
		configSpec := viewCommon.GenerateViewConfigSpecJson(view)
		configs[name] = configSpec
	}
	c.JSON(http.StatusOK, configs)
}

func GetViewById(c *gin.Context) {
	id := c.Param("id")
	if id != "" {
		c.JSON(http.StatusOK, "")
		return
	}
	c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
}
