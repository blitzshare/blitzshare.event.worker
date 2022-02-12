package app

import (
	"blitzshare.event.worker/app/dependencies"
	"blitzshare.event.worker/app/services/event"
	log "github.com/sirupsen/logrus"
)

func Start(dep *dependencies.Dependencies) {
	log.Infoln("worker Subscribed To Queue", dep.Config.QueueUrl)
	// local test: kubemqctl queues send p2p-peer-register-cmd '{"multiAddr": "multiAddr", "otp":"otp", "mode": "mode", "token":"token"}'
	event.PeerRegistry(dep)
	// local test: kubemqctl queues send p2p-bootstrap-node-registry-cmd '{"nodeId":"nodeId", "port": 123}'
	event.NodeRegistry(dep)
	//local test: kubemqctl queues send p2p-peer-deregister-cmd  '{"token":"token", "otp": "otp"}'
	event.PeerDeRegistry(dep)
}
