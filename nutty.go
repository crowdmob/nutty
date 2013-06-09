// Framework v 0.0.0.0.0.0.1b

package nutty

import (
  "os"
  "log"
  "fmt"
  configfile "github.com/crowdmob/goconfig"
  "github.com/crowdmob/goamz/aws"
  "bitbucket.org/kardianos/osext"
)

type App struct {
  configFileName  string
  Name            string
  Env             string
  Port            int64
  ExePath         string
  Logfile         string
  SnsArn          string
  PaypalUsername  string
  PaypalPassword  string
  PaypalSignature string
  AwsRegion       aws.Region
  AwsAuth         aws.Auth
  KafkaHost       string
  KafkaPort       int64
  KafkaHostname   string
  KafkaPartition  int64
  Routes          Router

  Globals         map[string]interface{}
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
  returnedApp.Logfile, err = config.GetString("default", "logfile")
  if err != nil { log.Fatalf("Error reading Nuts config: [default].logfile %#v\n", err) }
  returnedApp.Port, err = config.GetInt64("default", "port")
  if err != nil { log.Fatalf("Error reading Nuts config: [default].port %#v\n", err) }
  
  // Paypal
  returnedApp.PaypalUsername, err = config.GetString("paypal", "username")
  if err != nil { log.Fatalf("Error reading Nuts config: [paypal].username %#v\n", err) }
  returnedApp.PaypalPassword, err = config.GetString("paypal", "password")
  if err != nil { log.Fatalf("Error reading Nuts config: [paypal].password %#v\n", err) }
  returnedApp.PaypalSignature, err = config.GetString("paypal", "signature")
  if err != nil { log.Fatalf("Error reading Nuts config: [paypal].signature %#v\n", err) }
  
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
  
  // Kafka
  returnedApp.KafkaHost, err = config.GetString("kafka", "host")
  if err != nil { log.Fatalf("Error reading Nuts config: [kafka].host %#v\n", err) }
  returnedApp.KafkaPort, err = config.GetInt64("kafka", "port")
  if err != nil { log.Fatalf("Error reading Nuts config: [kafka].port %#v\n", err) }
  returnedApp.KafkaHostname = fmt.Sprintf("%s:%d", returnedApp.KafkaHost, returnedApp.KafkaPort)
  returnedApp.KafkaPartition, err = config.GetInt64("kafka", "partition")
  if err != nil { log.Fatalf("Error reading Nuts config: [kafka].partition %#v\n", err) }
  
  // ExePath
  returnedApp.ExePath, err = osext.ExecutableFolder()
  if err != nil { log.Fatalf("Error setting Nuts Config: osext.ExecutableFolder() %#v\n", err) }
  
  returnedApp.Routes.handlers = make(map[string](map[string]interface{}))
  returnedApp.Routes.initializations = make(map[string]bool)

  returnedApp.Globals = make(map[string]interface{})

  if returnedApp.Logfile != "stdout" {
    f, err := os.Open(returnedApp.Logfile)
    if err != nil {
      log.Fatalf("Couldn't open logfile %s because: %#v\n", returnedApp.Logfile, err)
    } else {
      log.SetOutput(f)
    }
  }

  return returnedApp
}
