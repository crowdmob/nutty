package nutty

import (
	"fmt"
	"github.com/crowdmob/kafka"
	"log"
)

// Writes to Kafka on the sent topic if on production.  Otherwise, outputs with log.Println
func (nuttyApp *App) KafkaPublish(topicName *string, message *string, completedNotice *chan bool) {
	go func() {
		if nuttyApp.Env != "production" {
			log.Println(*message)
		} else {
			broker := kafka.NewBrokerPublisher(nuttyApp.KafkaHostname, *topicName, int(nuttyApp.KafkaPartition))
			_, err := broker.Publish(kafka.NewMessage([]byte(*message)))
			if err != nil {
				nuttyApp.SNSPublish("ERROR Writing To Kafka", fmt.Sprintf("An error occurred when writing to kafka: %#v", err))
			}
		}

		if completedNotice != nil {
			*completedNotice <- true
		}
	}()
}
