package store

import (
    "io/ioutil"
    "os"
    "testing"
    view "github.com/6ixisgood/matrix-ticker/pkg/view/common"
    viewT "github.com/6ixisgood/matrix-ticker/pkg/view/types"
    "github.com/stretchr/testify/assert"
)

func TestSaveAndGetViewFromStore(t *testing.T) {
    // Create a temporary directory for the RocksDB store
    dbPath, err := ioutil.TempDir("", "rocksdb")
    if err != nil {
        t.Fatal(err)
    }
    defer os.RemoveAll(dbPath)

    // Test creating a new store
    store, err := NewStore(dbPath)
    assert.NoError(t, err)

    // Create a view definition
    definition := view.ViewDefinition{
        Id: "123456",
        Name: "myTestViewDefinition",
		Type: "text",
		Config: &viewT.TextViewConfig{
            Text: "My Text",
        },
    }
    
    // create a view with the definition to compare
    ogView, err := viewT.TextViewCreate(definition.Config)
    assert.NoError(t, err)

    // Test saving the view definition
    err = store.SaveViewDefinition(definition)
    assert.NoError(t, err)

    // Test fetching the view definition back
    retrievedDefinition, err := store.GetViewDefinition(definition.Id)
    assert.NoError(t, err)
    
	// Create a new view with with the retrieved config
    newView, err := viewT.TextViewCreate(retrievedDefinition.Config)
    assert.NoError(t, err)
    assert.Equal(t, ogView, newView) 
    
    // Get all views
    allViews, err := store.GetAllViewDefinitions()
    assert.NoError(t, err)
    assert.IsType(t, []view.ViewDefinition{}, allViews)
    assert.Equal(t, len(allViews), 1)
    assert.IsType(t, view.ViewDefinition{}, allViews[0])
}

