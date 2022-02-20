package registry

import (
	"encoding/json"
	"errors"
	"time"

	"blitzshare.event.worker/app/domain"
	"blitzshare.event.worker/app/services/str"

	"github.com/go-redis/redis"

	log "github.com/sirupsen/logrus"
)

type Registry interface {
	RegisterPeer(peer *domain.P2pPeerRegistryCmd) (string, error)
	DeregisterPeer(peer *domain.P2pPeerDeregisterCmd) error
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

var p2pPeersDbClient *redis.Client
var p2pBootstraoNodeDbClient *redis.Client

const (
	P2pPeersDb          = 0
	P2pBootstrapNodeDb  = 1
	DefaultKeyTimeout   = time.Second * 900
	NoExpirationTimeout = 0
	BootstrapNode       = "BootstrapNode"
)

func (r *RegistryIml) getPeersClient() *redis.Client {
	if p2pPeersDbClient == nil {
		p2pPeersDbClient = redis.NewClient(&redis.Options{
			Addr:     r.RedisUrl,
			Password: "",
			DB:       P2pPeersDb,
		})
	}
	return p2pPeersDbClient
}
func (r *RegistryIml) getBootstrapNodeDbClient() *redis.Client {
	if p2pBootstraoNodeDbClient == nil {
		p2pBootstraoNodeDbClient = redis.NewClient(&redis.Options{
			Addr:     r.RedisUrl,
			Password: "",
			DB:       P2pBootstrapNodeDb,
		})
	}
	return p2pBootstraoNodeDbClient
}

func (r *RegistryIml) RegisterPeer(peer *domain.P2pPeerRegistryCmd) (string, error) {
	bEvent, err := json.Marshal(peer)
	if err != nil {
		log.Fatal(err)
	}
	client := r.getPeersClient()
	return client.Set(str.SanatizeStr(peer.Otp), string(bEvent), DefaultKeyTimeout).Result()
}

func (r *RegistryIml) RegisterNode(node *domain.P2pBootstrapNodeRegistryCmd) (string, error) {
	bEvent, err := json.Marshal(node)
	if err != nil {
		log.Fatal(err)
	}
	client := r.getBootstrapNodeDbClient()
	return client.Set(BootstrapNode, string(bEvent), NoExpirationTimeout).Result()
}

func (r *RegistryIml) DeregisterPeer(cmd *domain.P2pPeerDeregisterCmd) error {
	client := r.getPeersClient()
	value, err := client.Get(cmd.Otp).Result()
	if err != nil {
		log.Fatal(err)
	}
	var peer domain.P2pPeerRegistryCmd
	err = json.Unmarshal([]byte(value), &peer)
	if err != nil {
		log.Fatal(err)
	}
	if peer.Token == cmd.Token {
		_, err := client.Del(cmd.Otp).Result()
		if err == nil {
			return nil
		}
		return err
	}
	return errors.New("cannot deregistre, token missmatch")
}
