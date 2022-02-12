package app

import (
	"blitzshare.event.worker/app/dependencies"
	"blitzshare.event.worker/app/domain"
	event2 "blitzshare.event.worker/app/services/event"
	"blitzshare.event.worker/app/services/queue"
	"encoding/json"
	log "github.com/sirupsen/logrus"
)

func NodeRegistry(dep *dependencies.Dependencies) {
	go dep.Mq.Sub(queue.P2pBootstrapNodeRegistryCmd, func(msg *[]byte) {
		var node domain.P2pBootstrapNodeRegistryCmd
		json.Unmarshal(*msg, &node)
		res, err := dep.Registry.RegisterNode(&node)
		if err == nil {
			log.Errorln("SUCCESS Node Registry", res, err)
		} else {
			log.Infoln("FAILED Node Registry", res, err)
		}
	})
}

func Start(dep *dependencies.Dependencies) {
	log.Infoln("worker Subscribed To Queue", dep.Config.QueueUrl)
	// local test: kubemqctl queues send p2p-peer-register-cmd '{"multiAddr": "multiAddr", "otp":"otp", "mode": "mode", "token":"token"}'
	event2.PeerRegistry(dep)
	// local test: kubemqctl queues send p2p-bootstrap-node-registry-cmd '{"nodeId":"nodeId", "port": 123}'
	NodeRegistry(dep)
	//local test: kubemqctl queues send p2p-peer-deregister-cmd  '{"token":"token", "otp": "otp"}'
	event2.PeerDeRegistry(dep)
}
