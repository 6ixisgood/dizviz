package store

import (
	"encoding/json"
	"errors"
	"fmt"
	view "github.com/6ixisgood/matrix-ticker/pkg/view/common"
	"github.com/linxGnu/grocksdb"
)

// Store represents the RocksDB store for view configurations
type Store struct {
	db *grocksdb.DB
}

const (
	viewDefinitionPrefix = "VD"
)


// NewStore initializes and returns a new Store
func NewStore(dbPath string) (*Store, error) {
	opts := grocksdb.NewDefaultOptions()
	opts.SetCreateIfMissing(true)
	db, err := grocksdb.OpenDb(opts, dbPath)
	if err != nil {
		return nil, err
	}
	return &Store{db: db}, nil
}

// SaveViewDefinition saves a view definition to the store
func (s *Store) SaveViewDefinition(definition view.ViewDefinition) error {
	data, err := json.Marshal(definition)
	if err != nil {
		return err
	}
	wo := grocksdb.NewDefaultWriteOptions()
	defer wo.Destroy()
	fmt.Printf("\n\n%s\n\n", data)

	return s.db.Put(wo, []byte(viewDefinitionPrefix + "-" + definition.Id), data)
}

// GetViewConfig retrieves a view definition from the store by name
func (s *Store) GetViewDefinition(id string) (view.ViewDefinition, error) {
    ro := grocksdb.NewDefaultReadOptions()
    defer ro.Destroy()
    data, err := s.db.Get(ro, []byte(viewDefinitionPrefix + "-" + id))
    if err != nil {
        return view.ViewDefinition{}, err
    }
    defer data.Free()
	
	vd, err :=unmarshalViewDefinition(data.Data())
	

    return vd, err
}

// GetAllViewDefinitions retrieves all saved view definitions from the store
func (s *Store) GetAllViewDefinitions() ([]view.ViewDefinition, error) {
    var definitions []view.ViewDefinition
    ro := grocksdb.NewDefaultReadOptions()
    defer ro.Destroy()
    it := s.db.NewIterator(ro)
    defer it.Close()

    prefix := []byte(viewDefinitionPrefix + "-")
    for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
        definition, err := unmarshalViewDefinition(it.Value().Data())
        if err != nil {
            return nil, err
        }
        definitions = append(definitions, definition)
    }

    if err := it.Err(); err != nil {
        return nil, err
    }

    return definitions, nil
}

// DeleteViewDefinition removes a view definition from the store by its ID
func (s *Store) DeleteViewDefinition(id string) error {
    key := []byte(viewDefinitionPrefix + "-" + id)

    // Create a write options object
    wo := grocksdb.NewDefaultWriteOptions()
    defer wo.Destroy()

    // Perform the deletion
    err := s.db.Delete(wo, key)
    if err != nil {
        return fmt.Errorf("failed to delete view definition with ID %s: %w", id, err)
    }

    return nil
}

func unmarshalViewDefinition(data []byte) (view.ViewDefinition, error) {
    var definition view.ViewDefinition

    // Unmarshal the whole object
    err := json.Unmarshal(data, &definition)
    if err != nil {
        return definition, err
    }

    // Unmarshal into a map to get the raw config
    var rawDef map[string]json.RawMessage
    err = json.Unmarshal(data, &rawDef)
    if err != nil {
        return definition, err
    }

    // Get the RegisteredView for this type
    regView, ok := view.RegisteredViews[definition.Type]
    if !ok {
        return definition, errors.New(fmt.Sprintf("No registered view of type %s", definition.Type))
    }

    // Use the RegisteredView's NewConfig function to unmarshal the config
    definition.Config = regView.NewConfig()
    err = json.Unmarshal(rawDef["config"], &definition.Config)
    if err != nil {
        return definition, err
    }

    return definition, nil
}

// Close closes the store and releases the database resources
func (s *Store) Close() {
	s.db.Close()
}

// Helper function to handle common logic for reading data from RocksDB
func (s *Store) readData(name string) ([]byte, error) {
	ro := grocksdb.NewDefaultReadOptions()
	defer ro.Destroy()
	data, err := s.db.Get(ro, []byte(name))
	if err != nil {
		return nil, err
	}
	defer data.Free()
	return data.Data(), nil
}

