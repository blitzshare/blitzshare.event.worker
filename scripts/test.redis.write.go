package main

import (
	"blitzshare.event.worker/app/domain"
	"blitzshare.event.worker/app/services/registry"
	log "github.com/sirupsen/logrus"
)

func redisP2pRegistry(redisUrl string) {
	res, err := registry.NewRegistry(redisUrl).RegisterPeer(&domain.P2pPeerRegistryCmd{
		MultiAddr: "test",
		Otp:       "test",
	})
	log.Infoln(res, err)
	if err == nil {
		log.Infoln("SUCCESS Peer Registry", res, err)
	} else {
		log.Errorln("FAILED Peer Registry", res, err)
	}
}

func nodeP2pRegistry(redisUrl string) {
	res, err := registry.NewRegistry(redisUrl).RegisterNode(&domain.P2pBootstrapNodeRegistryCmd{
		NodeId: "sdfsdfsdfsdfsdfsdfsdfsdf",
		Port:   63785,
	})
	log.Infoln(res, err)
	if err == nil {
		log.Infoln("SUCCESS Peer Registry", res, err)
	} else {
		log.Errorln("FAILED Peer Registry", res, err)
	}
}

func main() {
	const RedisUrl = "redis-svc.blitzshare-ns.svc.cluster.local:6379"
	nodeP2pRegistry(RedisUrl)
	redisP2pRegistry(RedisUrl)
}
