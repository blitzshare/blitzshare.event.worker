package dependencies

import "blitzshare.event.worker/app/config"

type Dependencies struct {
	Config config.Config
}

func NewDependencies(config config.Config) (*Dependencies, error) {
	return &Dependencies{Config: config}, nil
}
