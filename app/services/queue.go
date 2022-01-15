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
	P2pBootstrapNodeRegistryCmd = "p2p-bootstrap-node-registry-cmd"
	PeerRegisterCmd             = "p2p-peer-register-cmd"
	PeerDeregisterCmd           = "p2p-peer-deregister-cmd"
)

type Message struct {
	MessageID string
	Body      []byte
}

func SubscribePeerDeregisterCmd(queueUrl string, callback func(p2p *domain.P2pPeerDeregisterCmd)) {
	log.Infoln("subscribed to", PeerDeregisterCmd)
	Subscribe(queueUrl, PeerDeregisterCmd, func(message *kubemq.QueueMessage) {
		var registry domain.P2pPeerDeregisterCmd
		json.Unmarshal(message.Body, &registry)
		callback(&registry)
	})
}

func SubscribeBoostrapNodeJoinedCmd(queueUrl string, callback func(p2p *domain.P2pBootstrapNodeRegistryCmd)) {
	log.Infoln("subscribed to", P2pBootstrapNodeRegistryCmd)
	Subscribe(queueUrl, P2pBootstrapNodeRegistryCmd, func(message *kubemq.QueueMessage) {
		var registry domain.P2pBootstrapNodeRegistryCmd
		json.Unmarshal(message.Body, &registry)
		callback(&registry)
	})
}

func SubscribeP2pJoinedCmd(queueUrl string, callback func(p2p *domain.P2pPeerRegistryCmd)) {
	log.Infoln("subscribed to", PeerRegisterCmd)
	Subscribe(queueUrl, PeerRegisterCmd, func(message *kubemq.QueueMessage) {
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
