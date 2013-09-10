package nutty

import (
	"fmt"
	"github.com/crowdmob/goamz/exp/sns"
	"log"
)

// Writes to SNS on the topic specified in nuttyApp.SnsArn if on production.  Otherwise, outputs with log.Println
func (nuttyApp *App) SNSPublish(subject string, message string) {
	go func() {
		_, snsErr := sns.New(nuttyApp.AwsAuth, nuttyApp.AwsRegion).Publish(&sns.PublishOpt{message, "", fmt.Sprintf("[%s] %s", nuttyApp.Name, subject), nuttyApp.SnsArn})
		if snsErr != nil {
			log.Println(fmt.Sprintf("SNS error: %#v: %s", snsErr, message))
		}
	}()
}
