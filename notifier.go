package nutty

import (
  "log"
  "fmt"
  "github.com/crowdmob/goamz/exp/sns"
)

// Writes to SNS on the topic specified in nuttyApp.SnsArn if on production.  Otherwise, outputs with log.Println
func (nuttyApp *App) SNSPublish(subject string, message string) {
  go func() {
  	_, snsErr := sns.New(nuttyApp.AwsAuth, nuttyApp.AwsRegion).Publish(&sns.PublishOpt{message, "", fmt.Sprintf("[%s] %s", nuttyApp.Name), nuttyApp.SnsArn})
    if snsErr != nil {
      log.Println(fmt.Sprintf("SNS error: %#v during report of error writing to kafka: %s", snsErr, message))
    }
  }()
}