package test

import (
	"blitzshare.event.worker/app/config"
	"blitzshare.event.worker/app/domain"
	"context"
	"encoding/json"
	"fmt"
	"github.com/cucumber/godog"
	"github.com/go-redis/redis"
	"github.com/kubemq-io/kubemq-go"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"time"
)

const (
	ClientId           = "acceptance-tests"
	KubeMqUrl          = "kubemq-cluster-grpc.kubemq.svc.cluster.local"
	REDIS_URL          = "redis-svc.blitzshare-ns.svc.cluster.local:6379"
	P2pPeersDb         = 0
	P2pBootstrapNodeDb = 1
	BootstrapNodeKey   = "BootstrapNode"
)

var BootstrapNodeRegistryCmdMessage = domain.P2pBootstrapNodeRegistryCmd{
	NodeId: "acceptance-test-node-id",
	Port:   8000,
}

func getRnd() int {
	max := 9999999999
	min := 1000000000
	return rand.Intn(max-min) + min
}

var PeerRegistryCmd = domain.P2pPeerRegistryCmd{
	MultiAddr: fmt.Sprintf("acceptance-test/10.101.18.26/tcp/63785/p2p/12D3KooWPGR-%d", getRnd()),
	Otp:       fmt.Sprintf("acceptance-test-gelandelaufer-astromancer-scurvyweed-%d", getRnd()),
	Token:     fmt.Sprintf("deregistration-token-%d", getRnd()),
	Mode:      "file",
}

var PeerDeregistryCmd = domain.P2pPeerDeregisterCmd{
	Otp:   PeerRegistryCmd.Otp,
	Token: PeerRegistryCmd.Token,
}

type P2pBootstrapNodeRegistryCmdKey struct{}
type PeerRegistryCmdKey struct{}
type BootstrapInitConfig struct{}

var p2pPeersDbClient *redis.Client
var p2pBootstraoNodeDbClient *redis.Client

func getPeersClient(redisUrl string) *redis.Client {
	if p2pPeersDbClient == nil {
		p2pPeersDbClient = redis.NewClient(&redis.Options{
			Addr:     redisUrl,
			Password: "",
			DB:       P2pPeersDb,
		})
	}
	return p2pPeersDbClient
}

func getBootstrapNodeDbClient(redisUrl string) *redis.Client {
	if p2pBootstraoNodeDbClient == nil {
		p2pBootstraoNodeDbClient = redis.NewClient(&redis.Options{
			Addr:     redisUrl,
			Password: "",
			DB:       P2pBootstrapNodeDb,
		})
	}
	result, err := p2pBootstraoNodeDbClient.Ping().Result()
	Expect(err).To(BeNil())
	Expect(result, "OK")
	return p2pBootstraoNodeDbClient
}

func emitEvent(event []byte, topic string) string {
	ctx, _ := context.WithCancel(context.Background())
	client, err := kubemq.NewClient(ctx,
		kubemq.WithAddress(KubeMqUrl, 50000),
		kubemq.WithClientId(ClientId),
		kubemq.WithTransportType(kubemq.TransportTypeGRPC))
	if err != nil {
		log.Fatalln(err)
	}
	defer client.Close()
	sendResult, err := client.NewQueueMessage().
		SetChannel(topic).
		SetBody(event).
		Send(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	return sendResult.MessageID
}

func dispatchPeerRegistryEvent(ctx context.Context) context.Context {
	bEvent, err := json.Marshal(PeerRegistryCmd)
	if err != nil {
		log.Fatalln(err)
	}
	mId := emitEvent(bEvent, config.MqPeerRegisterCmd)
	return context.WithValue(ctx, PeerRegistryCmdKey{}, mId)
}

func dispatchPeerDeregistryEvent(ctx context.Context) context.Context {
	bEvent, err := json.Marshal(PeerDeregistryCmd)
	if err != nil {
		log.Fatalln(err)
	}
	mId := emitEvent(bEvent, config.MqPeerDeregisterCmd)
	return context.WithValue(ctx, PeerRegistryCmdKey{}, mId)
}
func dispatchNodeRegistryEvent(ctx context.Context) context.Context {
	bEvent, err := json.Marshal(BootstrapNodeRegistryCmdMessage)
	if err != nil {
		log.Fatalln(err)
	}
	mId := emitEvent(bEvent, config.MqP2pBootstrapNodeRegistryCmd)
	return context.WithValue(ctx, P2pBootstrapNodeRegistryCmdKey{}, mId)
}

func getNodeDbRecord() domain.P2pBootstrapNodeRegistryCmd {
	client := getBootstrapNodeDbClient(REDIS_URL)
	dbRecord, err := client.Get(BootstrapNodeKey).Result()
	var nodeConfig domain.P2pBootstrapNodeRegistryCmd
	err = json.Unmarshal([]byte(dbRecord), &nodeConfig)
	if err != nil {
		log.Fatalln(err)
	}
	return nodeConfig
}

func getPeerDbRecord(otp string) (domain.P2pPeerRegistryCmd, error) {
	client := getPeersClient(REDIS_URL)
	dbRecord, err := client.Get(otp).Result()
	if err != nil {
		return domain.P2pPeerRegistryCmd{}, err
	}
	var registry domain.P2pPeerRegistryCmd
	err = json.Unmarshal([]byte(dbRecord), &registry)
	if err != nil {
		log.Fatalln(err)
	}
	return registry, nil
}

func assertNodeRegistryRecordsExists(ctx context.Context) context.Context {
	nodeConfig := getNodeDbRecord()
	Expect(nodeConfig.NodeId).To(Equal(BootstrapNodeRegistryCmdMessage.NodeId))
	Expect(nodeConfig.Port).To(Equal(BootstrapNodeRegistryCmdMessage.Port))
	return ctx
}

func assertPeerRegistryExists(ctx context.Context) context.Context {
	dbRecord, err := getPeerDbRecord(PeerRegistryCmd.Otp)
	if err != nil {
		log.Fatalln(err)
	}
	Expect(dbRecord.Otp).To(Equal(PeerRegistryCmd.Otp))
	Expect(dbRecord.Mode).To(Equal(PeerRegistryCmd.Mode))
	Expect(dbRecord.MultiAddr).To(Equal(PeerRegistryCmd.MultiAddr))
	Expect(dbRecord.Token).To(Equal(PeerRegistryCmd.Token))
	return ctx
}
func assertPeerRegistryNotExists(ctx context.Context) context.Context {
	_, err := getPeerDbRecord(PeerRegistryCmd.Otp)
	Expect(err).To(Not(BeNil()))
	return ctx
}

func cleanup(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
	nodeConfig := ctx.Value(BootstrapInitConfig{}).(domain.P2pBootstrapNodeRegistryCmd)
	if nodeConfig.NodeId != "" && nodeConfig.NodeId != BootstrapNodeRegistryCmdMessage.NodeId {
		client := getBootstrapNodeDbClient(REDIS_URL)
		bEvent, err := json.Marshal(nodeConfig)
		if err != nil {
			return ctx, err
		}
		log.Infoln("Restoring bootstrap record to", string(bEvent))
		_, err = client.Set(BootstrapNodeKey, string(bEvent), 0).Result()
		if err != nil {
			return ctx, err
		}
		_, err = client.Del(PeerRegistryCmd.Otp).Result()
		log.Infoln("Deleted peer registry", PeerRegistryCmd.Otp)

	}
	return ctx, nil
}
func InitializeScenario(ctx *godog.ScenarioContext) {
	RegisterFailHandler(Fail)
	PanicWith(Panic)
	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		nodeConfig := getNodeDbRecord()
		return context.WithValue(ctx, BootstrapInitConfig{}, nodeConfig), nil
	})
	ctx.Step(`^Node registry event is dispatched$`, dispatchNodeRegistryEvent)
	ctx.Step(`^Node registry record is created$`, assertNodeRegistryRecordsExists)
	ctx.Step(`^Peer registry event is dispatched$`, dispatchPeerRegistryEvent)
	ctx.Step(`^Peer registry record is created$`, assertPeerRegistryExists)
	ctx.Step(`^Peer deregistry event is dispatched$`, dispatchPeerDeregistryEvent)
	ctx.Step(`^Peer registry record is deleted$`, assertPeerRegistryNotExists)
	ctx.Step(`^Test Wait for 1 seconds$`, func(ctx context.Context) context.Context {
		time.Sleep(time.Second * 1)
		return ctx
	})
	ctx.After(cleanup)
}
