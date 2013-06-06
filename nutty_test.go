package nutty_test

import (
  "net/http"
  "../nutty"
	"testing"
)

var TesterController = TesterControllerDefinition{}
type TesterControllerDefinition struct {}

func (TesterControllerDefinition) Index(nuttyApp *nutty.App, w http.ResponseWriter, r *http.Request) {}
func (TesterControllerDefinition) Create(nuttyApp *nutty.App, w http.ResponseWriter, r *http.Request) {}
func (TesterControllerDefinition) Update(nuttyApp *nutty.App, w http.ResponseWriter, r *http.Request) {}
func (TesterControllerDefinition) Destroy(nuttyApp *nutty.App, w http.ResponseWriter, r *http.Request) {}
  
func TestInitRoutes(t *testing.T) {
  configFilename := "nutty.properties.example"
	nuttyApp := nutty.New(&configFilename)
  
  // Test Root
  nuttyApp.Routes.Root(TesterController, nuttyApp)
  
  // Test Individual Function
  nuttyApp.Routes.Map("/indexer",    TesterController, []string{"GET"}, nuttyApp)
  nuttyApp.Routes.Map("/creater",    TesterController, []string{"POST"}, nuttyApp)
  nuttyApp.Routes.Map("/updater",    TesterController, []string{"PUT"}, nuttyApp)
  nuttyApp.Routes.Map("/destroyer",  TesterController, []string{"DELETE"}, nuttyApp)

  // Test Combination Routes
  nuttyApp.Routes.Map("/teams.json",  TesterController, []string{"GET", "POST", "PUT", "DELETE"}, nuttyApp)
}

