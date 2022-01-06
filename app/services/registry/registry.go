package registry

import (
	"encoding/json"
	"time"

	"blitzshare.event.worker/app/dependencies"
	"blitzshare.event.worker/app/domain"
	"blitzshare.event.worker/app/services/str"

	"github.com/go-redis/redis"

	log "github.com/sirupsen/logrus"
)

var client *redis.Client

func getClient(d *dependencies.Dependencies, db int) *redis.Client {
	if client == nil {
		client = redis.NewClient(&redis.Options{
			Addr:     d.Config.RedisUrl,
			Password: "",
			DB:       db,
		})
		pong, _ := client.Ping().Result()
		log.Infoln("getClient pong", pong)
	}
	return client
}

const (
	P2pPeersDb         = 0
	P2pBootstraoNodeDb = 1
	NodeList           = "p2p-bootstrap-nodes"
	DefaultKeyTimeout  = time.Second * 10000
)

func RegisterPeer(d *dependencies.Dependencies, peer *domain.P2pPeerRegistryCmd) (string, error) {
	client := getClient(d, P2pPeersDb)
	return client.Set(str.SanatizeStr(peer.Otp), str.SanatizeStr(peer.MultiAddr), DefaultKeyTimeout).Result()
}

func RegisterNode(d *dependencies.Dependencies, node *domain.P2pBootstrapNodeRegistryCmd) (string, error) {
	bEvent, err := json.Marshal(node)
	if err != nil {
		log.Fatal(err)
	}
	client := getClient(d, P2pBootstraoNodeDb)
	return client.Set(str.SanatizeStr(node.NodeId), string(bEvent), DefaultKeyTimeout).Result()
}
