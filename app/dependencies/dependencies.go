package dependencies

import (
	"blitzshare.event.worker/app/config"
	"blitzshare.event.worker/app/services/queue"
	"blitzshare.event.worker/app/services/registry"
)

type Dependencies struct {
	Config   config.Config
	Registry registry.Registry
	Mq       queue.Mq
}

func NewDependencies(config config.Config) (*Dependencies, error) {
	return &Dependencies{
		Config:   config,
		Registry: registry.NewRegistry(config.RedisUrl),
		Mq:       queue.NewMq(config.QueueUrl),
	}, nil
}
