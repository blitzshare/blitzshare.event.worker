package event

import (
	"blitzshare.event.worker/app/config"
	"blitzshare.event.worker/app/dependencies"
	"blitzshare.event.worker/app/domain"
	"encoding/json"
	log "github.com/sirupsen/logrus"
)

func PeerRegistry(dep *dependencies.Dependencies) {
	go dep.Mq.Sub(config.MqPeerRegisterCmd, func(msg *[]byte) {
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
