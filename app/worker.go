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

	go services.SubscribeP2pJoinedCmd(dep.Config.QueueUrl, func(peer *domain.P2pPeerRegistryCmd) {
		log.Infoln("Peer Registry", peer)
		res, err := registry.RegisterPeer(dep, peer)
		if err != nil {
			log.Errorln("SUCCESS Peer Registry", res, err)
		} else {
			log.Infoln("FAILED Peer Registry", res, err)
		}
	})
	// TODO
	go services.SubscribeBoostrapNodeJoinedCmd(dep.Config.QueueUrl, func(node *domain.P2pBootstrapNodeRegistryCmd) {
		log.Infoln("Node Registry", node)
		// TODO: change the way we store node info
		// res, err := registry.RegisterNode(dep, node)
		// if err != nil {
		// 	log.Errorln("Node Registry", res, err)
		// } else {
		// 	log.Infoln("Node Registry", res, err)
		// }
	})

}
