package app

import (
	"blitzshare.event.worker/app/dependencies"
	"blitzshare.event.worker/app/services/event"
	log "github.com/sirupsen/logrus"
)

func Start(dep *dependencies.Dependencies) {
	log.Infoln("Worker Subscribed To Queue", dep.Config.QueueUrl)
	event.PeerRegistry(dep)
	event.NodeRegistry(dep)
	event.PeerDeRegistry(dep)
}
