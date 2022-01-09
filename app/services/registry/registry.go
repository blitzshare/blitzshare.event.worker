package registry

import (
	"encoding/json"
	"time"

	"blitzshare.event.worker/app/domain"
	"blitzshare.event.worker/app/services/str"

	"github.com/go-redis/redis"

	log "github.com/sirupsen/logrus"
)

type Registry interface {
	RegisterPeer(peer *domain.P2pPeerRegistryCmd) (string, error)
	RegisterNode(node *domain.P2pBootstrapNodeRegistryCmd) (string, error)
}

type RegistryIml struct {
	RedisUrl string
}

func NewRegistry(RedisUrl string) Registry {
	return &RegistryIml{
		RedisUrl: RedisUrl,
	}
}

var client *redis.Client

const (
	P2pPeersDb          = 0
	P2pBootstraoNodeDb  = 1
	DefaultKeyTimeout   = time.Second * 10000
	NoExpirationTimeout = 0
	BootstrapNodeId     = "BootstrapNodeId"
)

func (r *RegistryIml) getClient(db int) *redis.Client {
	if client == nil {
		client = redis.NewClient(&redis.Options{
			Addr:     r.RedisUrl,
			Password: "",
			DB:       db,
		})
		pong, _ := client.Ping().Result()
		log.Infoln("getClient", pong)
	}
	return client
}

func (r *RegistryIml) RegisterPeer(peer *domain.P2pPeerRegistryCmd) (string, error) {
	client := r.getClient(P2pPeersDb)
	return client.Set(str.SanatizeStr(peer.Otp), str.SanatizeStr(peer.MultiAddr), DefaultKeyTimeout).Result()
}

func (r *RegistryIml) RegisterNode(node *domain.P2pBootstrapNodeRegistryCmd) (string, error) {
	bEvent, err := json.Marshal(node)
	if err != nil {
		log.Fatal(err)
	}
	client := r.getClient(P2pBootstraoNodeDb)
	return client.Set(BootstrapNodeId, string(bEvent), NoExpirationTimeout).Result()
}
