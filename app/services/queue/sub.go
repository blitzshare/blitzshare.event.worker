package queue

import (
	"blitzshare.event.worker/app/config"
	"context"
	"time"

	mq "github.com/kubemq-io/kubemq-go"
	log "github.com/sirupsen/logrus"
)

const clientId = "blitzshare-event-worker"

type Message struct {
	MessageID string
	Body      []byte
}
type Mq interface {
	Sub(topic string, cb func(message *[]byte))
}
type MqImpl struct {
	QueueUrl string
}

func NewMq(queueUrl string) Mq {
	return &MqImpl{
		QueueUrl: queueUrl,
	}
}

func subscribe(queueUrl string, topic string, cb func(result *mq.ReceiveQueueMessagesResponse)) {
	for {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		client, err := mq.NewClient(ctx,
			mq.WithAddress(queueUrl, config.MqPort),
			mq.WithClientId(clientId),
			mq.WithTransportType(mq.TransportTypeGRPC))
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
		cb(receiveResult)
		time.Sleep(time.Second * 1)
	}
}

func (impl *MqImpl) Sub(topic string, cb func(message *[]byte)) {
	log.Infoln("subscribed to", impl.QueueUrl, topic)
	subscribe(impl.QueueUrl, topic, func(resp *mq.ReceiveQueueMessagesResponse) {
		if resp.MessagesReceived > 0 {
			for _, msg := range resp.Messages {
				cb(&msg.Body)
			}
		}
	})
}
