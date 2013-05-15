nutty
=====

An opinionated web framework in go (golang) that is built to work with Kafka, DynamoDB, and SNS.


Example `server.go`
-------------------


```go
package main

import (
  "flag"
  "fmt"
  "log"
  "net/http"
  "os"

  "github.com/crowdmob/nutty"
  "./routes"
)

var ConfigFilename string

func init() {
  flag.StringVar(&ConfigFilename, "c", "config/nutty.properties", "path to config file")
}

func main() {
  flag.Parse()
  log.Printf("Loading config file %s", ConfigFilename)
  
  nuttyApp := nutty.New(&ConfigFilename)
  routes.Init(nuttyApp)
  nuttyApp.Globals["my_global"] = someGlobalVariableYouWantInAllControllers

  log.Printf("HTTP server listening on %d...\n", nuttyApp.Port)
  err = http.ListenAndServe(fmt.Sprintf(":%d",nuttyApp.Port), nil)
  if err != nil {
    log.Fatal("ListenAndServe: ", err)
  }
}
```