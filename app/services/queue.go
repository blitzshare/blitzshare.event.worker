package services

import (
	"context"
	"fmt"
	"time"

	kubemq "github.com/kubemq-io/kubemq-go"
	log "github.com/sirupsen/logrus"
)

const clientId = "blitzshare-event-worker"

const (
	P2pNodeInstanceChannel = "p2p-node-instance-channel"
	P2pPeerRegistry        = "p2p-peer-registry"
)

func Subscribe(queueUrl string, topic string, callback func(i interface{})) {
	fmt.Println("Infinite Loop 2")
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
		time.Sleep(time.Second * 1)
	}
}
