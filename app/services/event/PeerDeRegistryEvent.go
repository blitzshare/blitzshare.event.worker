package event

import (
	"blitzshare.event.worker/app/dependencies"
	"blitzshare.event.worker/app/domain"
	"blitzshare.event.worker/app/services/queue"
	"encoding/json"
	log "github.com/sirupsen/logrus"
)

func PeerDeRegistry(dep *dependencies.Dependencies) {
	go dep.Mq.Sub(queue.PeerDeregisterCmd, func(msg *[]byte) {
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
