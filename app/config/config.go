package config

import "github.com/kelseyhightower/envconfig"

type Config struct {
	QueueUrl string `envconfig:"QUEUE_URL"`
	RedisUrl string `envconfig:"REDIS_URL"`
}

const (
	MqP2pBootstrapNodeRegistryCmd = "p2p-bootstrap-node-registry-cmd"
	MqPeerRegisterCmd             = "p2p-peer-register-cmd"
	MqPeerDeregisterCmd           = "p2p-peer-deregister-cmd"
	MqPort                        = 50000
)

func Load() (Config, error) {
	err := LoadEnvironment()
	cfg := Config{}
	err = envconfig.Process("settings", &cfg)
	return cfg, err
}
