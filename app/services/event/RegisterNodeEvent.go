package event

import (
	"blitzshare.event.worker/app/config"
	"blitzshare.event.worker/app/dependencies"
	"blitzshare.event.worker/app/domain"
	"encoding/json"
	log "github.com/sirupsen/logrus"
)

func NodeRegistry(dep *dependencies.Dependencies) {
	go dep.Mq.Sub(config.MqP2pBootstrapNodeRegistryCmd, func(msg *[]byte) {
		var node domain.P2pBootstrapNodeRegistryCmd
		json.Unmarshal(*msg, &node)
		res, err := dep.Registry.RegisterNode(&node)
		if err == nil {
			log.Errorln("SUCCESS NodeRegistry", res, err)
		} else {
			log.Infoln("FAILED NodeRegistry", res, err)
		}
	})
}
