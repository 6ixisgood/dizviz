package handlers

import (
	"encoding/json"
	"fmt"
	viewCommon "github.com/6ixisgood/matrix-ticker/pkg/view/common"
	"github.com/google/uuid"
	"log"
	"net/http"
)

func GetAllViewDefinitions(c *gin.Context) {
	definitions, err := viewCommon.GetAllViewDefinitions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error retrieving view definitions"})
		return
	}
	c.JSON(http.StatusOK, definitions)
}

func GetViewDefinition(c *gin.Context) {
	id := c.Param("id")
	definition, err := viewCommon.GetViewDefinition(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "View definition not found"})
		return
	}
	c.JSON(http.StatusOK, definition)
}

func SaveViewDefinition(c *gin.Context) {
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

// DeleteViewDefinition handler function
func DeleteViewDefinition(c *gin.Context) {
	id := c.Param("id") // Extract the ID from the URL parameter

	err := viewCommon.DeleteViewDefinition(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("Unable to delete ID: %s", id)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "View definition deleted successfully"})
}

func GetAllViewConfigSpecs(c *gin.Context) {
	configs := make(map[string]interface{})
	for name, regView := range viewCommon.RegisteredViews {
		configSpec := viewCommon.GenerateViewConfigSpecJson(regView.NewConfig())
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