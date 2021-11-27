package app

import (
	"blitzshare.event.worker/app/dependencies"
	"blitzshare.event.worker/app/services"
	log "github.com/sirupsen/logrus"
)

func onP2pNodeInstanceChannelEvent(arg interface{}) {
	log.Infoln("onP2pNodeInstanceChannelEvent", arg)
}
func onP2pPeerJoined(arg interface{}) {
	// RegisterPeer

	// registry.RegisterPeer(registry.Peer{
	// 	MultiAddr:   "?",
	// 	OneTimePass: "?",
	// })
	log.Infoln("onP2pPeerJoined", arg)
}
func Start(dep *dependencies.Dependencies) {
	log.Infoln("SubscribeToQueue", services.P2pPeerRegistry)
	go services.Subscribe(dep.Config.QueueUrl, services.P2pPeerRegistry, onP2pPeerJoined)
	log.Infoln("SubscribeToQueue", services.P2pNodeInstanceChannel)
	go services.Subscribe(dep.Config.QueueUrl, services.P2pNodeInstanceChannel, onP2pNodeInstanceChannelEvent)

}
