package services

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	kubemq "github.com/kubemq-io/kubemq-go"
	"time"
)
const clientId = "blitzshare-event-worker"



const (
	P2pNodeInstanceChannel = "p2p-node-instance-channel"
)

func SubscribeToQueue(queueUrl string, topic string, callback func(i interface{})) {
	fmt.Println("Infinite Loop 2")
	log.Info("SubmitEvent:SubscribeToQueue")
	for true {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		client, err := kubemq.NewClient(ctx,
			kubemq.WithAddress(queueUrl, 50000),
			kubemq.WithClientId(clientId),
			kubemq.WithTransportType(kubemq.TransportTypeGRPC))
		if err != nil {
			log.Fatal(err)
		}
		defer client.Close()

		receiveResult, err := client.NewReceiveQueueMessagesRequest().
			SetChannel(topic).
			SetMaxNumberOfMessages(1).
			SetWaitTimeSeconds(5).
			Send(ctx)
		if err != nil {
			log.Fatal(err)
		}
		if receiveResult.MessagesReceived > 0 {
			log.Printf("Received %d Messages:\n", receiveResult.MessagesReceived)
			for _, msg := range receiveResult.Messages {
				log.Printf("MessageID: %s, Body: %s", msg.MessageID, string(msg.Body))
			}
		}
		time.Sleep(time.Second * 10)
	}
}