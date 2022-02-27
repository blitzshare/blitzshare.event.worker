package event

import (
	"blitzshare.event.worker/app/config"
	"blitzshare.event.worker/app/dependencies"
	"blitzshare.event.worker/app/domain"
	"encoding/json"
	log "github.com/sirupsen/logrus"
)

func PeerDeRegistry(dep *dependencies.Dependencies) {
	go dep.Mq.Sub(config.MqPeerDeregisterCmd, func(msg *[]byte) {
		var peer domain.P2pPeerDeregisterCmd
		json.Unmarshal(*msg, &peer)
		err := dep.Registry.DeregisterPeer(&peer)
		if err == nil {
			log.Infoln("SUCCESS PeerDeRegistry", err)
		} else {
			log.Infoln("FAILED PeerDeRegistry", err)
		}
	})
}
