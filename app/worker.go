package app

import (
	"blitzshare.event.worker/app/dependencies"
	"blitzshare.event.worker/app/services"
	log "github.com/sirupsen/logrus"
)

func onPspNodeInstanceChannelEvent(arg interface{})  {
	log.Infoln("onPspNodeInstanceChannelEvent", arg)
}
func Start(dep *dependencies.Dependencies) {
	log.Infoln("SubscribeToQueue", services.P2pNodeInstanceChannel)
	services.SubscribeToQueue(dep.Config.QueueUrl, services.P2pNodeInstanceChannel, onPspNodeInstanceChannelEvent)
}
