package app

import (
	"blitzshare.event.worker/app/dependencies"
	"blitzshare.event.worker/app/domain"
	"blitzshare.event.worker/app/services"
	"blitzshare.event.worker/app/services/registry"
	log "github.com/sirupsen/logrus"
)

func Start(dep *dependencies.Dependencies) {
	log.Infoln("SubscribeToQueue")
	// go services.Subscribe(dep.Config.QueueUrl, services.P2pPeerRegistry, onP2pPeerJoined)

	go services.SubscribeP2pJoinedEvent(dep.Config.QueueUrl, func(peer *domain.P2pPeerRegistryCmd) {
		log.Infoln("onP2pPeerJoined", peer)
		res, err := registry.RegisterPeer(dep, peer)
		if err != nil {
			log.Errorln("onP2pPeerJoined", res, err)
		} else {
			log.Infoln("onP2pPeerJoined", res, err)
		}
	})

}
