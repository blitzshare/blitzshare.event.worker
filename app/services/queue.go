package services

import (
	"context"
	"encoding/json"
	"time"

	"blitzshare.event.worker/app/domain"
	kubemq "github.com/kubemq-io/kubemq-go"
	log "github.com/sirupsen/logrus"
)

const clientId = "blitzshare-event-worker"

const (
	P2pNodeInstanceChannel = "p2p-node-instance-channel"
	P2pPeerRegistryCmd     = "p2p-peer-registry-cmd"
)

type Message struct {
	MessageID string
	Body      []byte
}

func SubscribeP2pJoinedEvent(queueUrl string, callback func(p2p *domain.P2pPeerRegistryCmd)) {
	Subscribe(queueUrl, P2pPeerRegistryCmd, func(message *kubemq.QueueMessage) {
		var registry domain.P2pPeerRegistryCmd
		json.Unmarshal(message.Body, &registry)
		callback(&registry)
	})
}

func Subscribe(queueUrl string, topic string, callback func(message *kubemq.QueueMessage)) {
	for {
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
			log.Infoln("Received Messages", receiveResult.MessagesReceived)
			for _, msg := range receiveResult.Messages {
				// log.Printf("MessageID: %s, Body: %s", msg.MessageID, string(msg.Body))
				callback(msg)
			}
		}
		time.Sleep(time.Second * 1)
	}
}
