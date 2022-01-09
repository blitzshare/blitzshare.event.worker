package app

import (
	"blitzshare.event.worker/app/dependencies"
	"blitzshare.event.worker/app/domain"
	"blitzshare.event.worker/app/services"
	log "github.com/sirupsen/logrus"
)

func Start(dep *dependencies.Dependencies) {
	log.Infoln("worker Subscribed To Queue", dep.Config.QueueUrl)
	go services.SubscribeP2pJoinedCmd(dep.Config.QueueUrl, func(peer *domain.P2pPeerRegistryCmd) {
		log.Printf("Peer Registry [%s], [%s]`n", peer.MultiAddr, peer.Otp)
		res, err := dep.Registry.RegisterPeer(peer)
		if err == nil {
			log.Errorln("SUCCESS Peer Registry", res, err)
		} else {
			log.Infoln("FAILED Peer Registry", res, err)
		}
	})
	go services.SubscribeBoostrapNodeJoinedCmd(dep.Config.QueueUrl, func(node *domain.P2pBootstrapNodeRegistryCmd) {
		log.Infoln("Node Registry", node)
		res, err := dep.Registry.RegisterNode(node)
		if err == nil {
			log.Errorln("SUCCESS Node Registry", res, err)
		} else {
			log.Infoln("FAILED Node Registry", res, err)
		}
	})
}
