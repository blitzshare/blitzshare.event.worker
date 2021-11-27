package registry

import (
	"time"

	"blitzshare.event.worker/app/dependencies"
	"blitzshare.event.worker/app/domain"
	"blitzshare.event.worker/app/services/str"

	"github.com/go-redis/redis"

	log "github.com/sirupsen/logrus"
)

var client *redis.Client

func getClient(d *dependencies.Dependencies) *redis.Client {
	if client == nil {
		client = redis.NewClient(&redis.Options{
			Addr:     d.Config.RedisUrl,
			Password: "",
		})
		pong, _ := client.Ping().Result()
		log.Infoln("getClient pong", pong)
	}
	return client
}

func set(d *dependencies.Dependencies, key string, value string) (string, error) {
	client := getClient(d)
	return client.Set(key, value, time.Second*10000).Result()
}

func RegisterPeer(d *dependencies.Dependencies, peer *domain.P2pPeerRegistryCmd) (string, error) {
	return set(d, str.SanatizeStr(peer.OneTimePass), str.SanatizeStr(peer.MultiAddr))
}
