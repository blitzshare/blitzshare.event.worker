package event

import (
	"blitzshare.event.worker/app/dependencies"
	"blitzshare.event.worker/app/domain"
	"blitzshare.event.worker/app/services/queue"
	"encoding/json"
	log "github.com/sirupsen/logrus"
)

func PeerRegistry(dep *dependencies.Dependencies) {
	go dep.Mq.Sub(queue.PeerRegisterCmd, func(msg *[]byte) {
		var peer domain.P2pPeerRegistryCmd
		json.Unmarshal(*msg, &peer)
		res, err := dep.Registry.RegisterPeer(&peer)
		if err == nil {
			log.Errorln("SUCCESS PeerRegistry", res, err)
		} else {
			log.Infoln("FAILED PeerRegistry", res, err)
		}
	})
}
