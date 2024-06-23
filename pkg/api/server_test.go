package api

import (
    "net/http"
    "net/http/httptest"
    "testing"
    "os"
    "io/ioutil"
	"strings"
    "github.com/stretchr/testify/assert"
	"github.com/6ixisgood/matrix-ticker/pkg/store"
)

var testStore *store.Store

func TestMain(m *testing.M) {
    // Setup
    dbPath, err := ioutil.TempDir("", "rocksdb")
    if err != nil {
        os.Exit(1)
    }
    defer os.RemoveAll(dbPath)

    testStore, err = store.NewStore(dbPath)
    if err != nil {
        os.Exit(1)
    }

    SetAppServerConfig(&AppServerConfig{
        AllowedHost: "127.0.0.1",
        Port:        "8081",
        Store:       testStore,
    })
	InitializeRoutes()

    // Run the tests
    code := m.Run()

    // Teardown
    // Add any necessary teardown steps

    os.Exit(code)
}

func TestSaveViewDefinition(t *testing.T) {
    definition := `{
        "Id": "123456",
        "Name": "myTestViewDefinition",
        "Type": "text",
        "Config": {
            "Text": "My Text"
        }
    }` // replace with your actual view definition structure

    req, _ := http.NewRequest("POST", "/views/definitions", strings.NewReader(definition))
    req.Header.Set("Content-Type", "application/json")
    resp := httptest.NewRecorder()
    Server.router.ServeHTTP(resp, req)

    assert.Equal(t, http.StatusOK, resp.Code)
}

func TestGetAllViewDefinitions(t *testing.T) {
    req, _ := http.NewRequest("GET", "/views/definitions", nil)
    resp := httptest.NewRecorder()
    Server.router.ServeHTTP(resp, req)

    assert.Equal(t, http.StatusOK, resp.Code)
}

func TestGetViewDefinition(t *testing.T) {
    req, _ := http.NewRequest("GET", "/views/definitions/123456", nil) // replace '1' with a valid ID
    resp := httptest.NewRecorder()
    Server.router.ServeHTTP(resp, req)

    assert.Equal(t, http.StatusOK, resp.Code)
}

func TestGetAllViewConfigs(t *testing.T) {
    req, _ := http.NewRequest("GET", "/views/configs", nil)
    resp := httptest.NewRecorder()
    Server.router.ServeHTTP(resp, req)

    assert.Equal(t, http.StatusOK, resp.Code)
}

func TestGetViewById(t *testing.T) {
    req, _ := http.NewRequest("GET", "/views/1", nil) // replace '1' with a valid ID
    resp := httptest.NewRecorder()
    Server.router.ServeHTTP(resp, req)

    assert.Equal(t, http.StatusOK, resp.Code)
}