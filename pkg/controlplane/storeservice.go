package controlplane

import (
	"encoding/json"
	"fmt"

	"github.com/6ixisgood/matrix-ticker/pkg/store"
	viewCommon "github.com/6ixisgood/matrix-ticker/pkg/view/common"
)

const (
	viewDefinitionPrefix = "VD"
)

// StoreService provides persistence operations for the control plane
type StoreService struct {
	store *store.Store
}

// NewStoreService creates a new store service
func NewStoreService(store *store.Store) *StoreService {
	return &StoreService{
		store: store,
	}
}

// Close closes the underlying store
func (s *StoreService) Close() {
	if s.store != nil {
		s.store.Close()
	}
}

// GetStore returns the underlying store (use sparingly, prefer specific methods)
func (s *StoreService) GetStore() *store.Store {
	return s.store
}

// === View Definition Operations ===

// SaveViewDefinition saves a view definition to the store
func (s *StoreService) SaveViewDefinition(definition viewCommon.ViewDefinition) error {
	data, err := json.Marshal(definition)
	if err != nil {
		return fmt.Errorf("failed to marshal view definition: %w", err)
	}
	_, err = s.store.SaveItem(viewDefinitionPrefix, definition.Id, data)
	if err != nil {
		return fmt.Errorf("failed to save view definition: %w", err)
	}
	return nil
}

// GetViewDefinition retrieves a view definition from the store by ID
func (s *StoreService) GetViewDefinition(id string) (viewCommon.ViewDefinition, error) {
	data, err := s.store.GetItem(viewDefinitionPrefix + "-" + id)
	if err != nil {
		return viewCommon.ViewDefinition{}, fmt.Errorf("failed to get view definition: %w", err)
	}

	definition, err := s.unmarshalViewDefinition(data)
	if err != nil {
		return viewCommon.ViewDefinition{}, fmt.Errorf("failed to unmarshal view definition: %w", err)
	}

	return definition, nil
}

// GetAllViewDefinitions retrieves all saved view definitions from the store
func (s *StoreService) GetAllViewDefinitions() ([]viewCommon.ViewDefinition, error) {
	var definitions []viewCommon.ViewDefinition
	viewDatas, err := s.store.GetPrefix(viewDefinitionPrefix + "-")
	if err != nil {
		return nil, fmt.Errorf("failed to get view definitions: %w", err)
	}

	for _, viewData := range viewDatas {
		definition, err := s.unmarshalViewDefinition(viewData)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal view definition: %w", err)
		}
		definitions = append(definitions, definition)
	}

	return definitions, nil
}

// DeleteViewDefinition deletes a view definition from the store by ID
func (s *StoreService) DeleteViewDefinition(id string) error {
	err := s.store.DeleteItem(viewDefinitionPrefix + "-" + id)
	if err != nil {
		return fmt.Errorf("failed to delete view definition: %w", err)
	}
	return nil
}

// unmarshalViewDefinition helper function to get a ViewDefinition from []byte
func (s *StoreService) unmarshalViewDefinition(data []byte) (viewCommon.ViewDefinition, error) {
	// First unmarshal into ViewDefinitionRaw to get the config as json.RawMessage
	var rawDef viewCommon.ViewDefinitionRaw
	err := json.Unmarshal(data, &rawDef)
	if err != nil {
		return viewCommon.ViewDefinition{}, err
	}

	// Verify the view type is registered
	if _, ok := viewCommon.RegisteredViews[rawDef.Type]; !ok {
		return viewCommon.ViewDefinition{}, fmt.Errorf("no registered view of type %s", rawDef.Type)
	}

	// Create the ViewDefinition with Config as json.RawMessage ([]byte)
	// The config will be unmarshaled by the view's Init() method when the view is created
	definition := viewCommon.ViewDefinition{
		Id:     rawDef.Id,
		Name:   rawDef.Name,
		Type:   rawDef.Type,
		Config: rawDef.Config, // Keep as json.RawMessage
	}

	return definition, nil
}

// === Future: Add more store operations here ===
// Example:
// - SaveAgentConfig(agentID, config)
// - GetAgentConfig(agentID)
// - SaveDisplayLayout(layoutID, layout)
// - GetDisplayLayout(layoutID)
