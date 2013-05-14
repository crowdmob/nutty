// Framework v 0.0.0.0.0.0.1b

package nutty

import (
  "log"
  "fmt"
  "os/exec"
  "strings"
  "net/http"
  configfile "github.com/crowdmob/goconfig"
  "github.com/crowdmob/goamz/aws"
  // "github.com/crowdmob/goamz/dynamodb"
  // "github.com/crowdmob/kafka"
)

type ControllerWithIndex interface {
  Index(*App, http.ResponseWriter, *http.Request)
}

type ControllerWithCreate interface {
  Create(*App, http.ResponseWriter, *http.Request)
}


type Router struct {
  handlers          map[string](map[string]interface{})
  initializations   map[string]bool
}

type App struct {
  configFileName  string
  Name            string
  Env             string
  Port            int64
  SnsArn          string
  AwsRegion       aws.Region
  AwsAuth         aws.Auth
  KafkaHost       string
  KafkaPort       int64
  KafkaHostname   string
  KafkaPartition  int64
  Routes          Router

  Globals         map[string]interface{}
}


// TODO implement resources controller that should be GET/POST/PUT/DELETE
// func (routes *Router) Resources(resourceName *string, controller interface{}) {
//   
// }

func (routes *Router) Map(uri string, controller interface{}, httpMethods []string, nuttyApp *App) {
  if !routes.initializations[uri] {
    routes.initializations[uri] = true
    routes.handlers[uri] = make(map[string]interface{})
    http.HandleFunc(uri, func(resp http.ResponseWriter, req *http.Request) {
      if routes.handlers[uri][req.Method] == nil {
        http.NotFound(resp, req)
      } else {
        if req.Method == "POST" {
          (routes.handlers[uri][req.Method]).(ControllerWithCreate).Create(nuttyApp, resp, req)
        } else {
          (routes.handlers[uri][req.Method]).(ControllerWithIndex).Index(nuttyApp, resp, req)
        }
      }
    })
  }
  
  for _, method := range httpMethods {
    routes.handlers[uri][strings.ToUpper(method)] = controller
  }
}

// Defaults to GET if no http methods sent
func (routes *Router) Root(ctrl ControllerWithIndex, nuttyApp *App) {
  if !routes.initializations["/"] {
    routes.initializations["/"] = true
    routes.handlers["/"] = make(map[string]interface{})
  }
  
  handler := func(w http.ResponseWriter, r *http.Request) { ctrl.Index(nuttyApp,w,r) }
  routes.handlers["/"]["GET"] = handler
  http.HandleFunc("/", handler)
  http.HandleFunc("/index", handler)
  http.HandleFunc("/index.html", handler)
}

func New(configFileName *string) *App {
  returnedApp := &App{}
  returnedApp.configFileName = *configFileName
  
  config, err := configfile.ReadConfigFile(returnedApp.configFileName)
  if err != nil {
    log.Fatalf("Couldn't read config file %s because: %#v\n", returnedApp.configFileName, err)
  }
  
  // Defaults
  returnedApp.Name, err = config.GetString("default", "appname")
  if err != nil { log.Fatalf("Error reading Nuts config: [default].appname %#v\n", err) }
  returnedApp.Env, err = config.GetString("default", "env")
  if err != nil { log.Fatalf("Error reading Nuts config: [default].env %#v\n", err) }
  if len(returnedApp.Env) == 0 { returnedApp.Env = "development" }
  returnedApp.Port, err = config.GetInt64("default", "port")
  if err != nil { log.Fatalf("Error reading Nuts config: [default].port %#v\n", err) }
  
  // AWS
  awsRegion, err := config.GetString("aws", "region")
  if err != nil { log.Fatalf("Error reading Nuts config: [aws].region %#v\n", err) }
  if len(awsRegion) == 0 { awsRegion = "us-east-1" }
  returnedApp.AwsRegion = aws.Regions[awsRegion]
  awsKey, err := config.GetString("aws", "accesskey")
  if err != nil { log.Fatalf("Error reading Nuts config: [aws].accesskey %#v\n", err) }
  awsSecret, err := config.GetString("aws", "secretkey")
  if err != nil { log.Fatalf("Error reading Nuts config: [aws].secretkey %#v\n", err) }
  returnedApp.AwsAuth = aws.Auth{awsKey, awsSecret}
  
  // SNS
  returnedApp.SnsArn, err = config.GetString("sns", "arn")
  if err != nil { log.Fatalf("Error reading Nuts config: [sns].arn %#v\n", err) }
  
  // returnedApp.DynamoDbServer = dynamodb.Server{returnedApp.AwsAuth, returnedApp.AwsRegion}
  
  // Kafka
  returnedApp.KafkaHost, err = config.GetString("kafka", "host")
  if err != nil { log.Fatalf("Error reading Nuts config: [kafka].host %#v\n", err) }
  returnedApp.KafkaPort, err = config.GetInt64("kafka", "port")
  if err != nil { log.Fatalf("Error reading Nuts config: [kafka].port %#v\n", err) }
  returnedApp.KafkaHostname = fmt.Sprintf("%s:%d", returnedApp.KafkaHost, returnedApp.KafkaPort)
  returnedApp.KafkaPartition, err = config.GetInt64("kafka", "partition")
  if err != nil { log.Fatalf("Error reading Nuts config: [kafka].partition %#v\n", err) }
  
  returnedApp.Routes.handlers = make(map[string](map[string]interface{}))
  returnedApp.Routes.initializations = make(map[string]bool)

  returnedApp.Globals = make(map[string]interface{})

  return returnedApp
}

func GenerateUUID() (string, error) {
  b, err := exec.Command("uuidgen").Output()
  if err != nil {
    return string(b), err
  }
  return strings.TrimSpace(string(b)), nil
}
