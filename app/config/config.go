package config


import "github.com/kelseyhightower/envconfig"

type Config struct {
	QueueUrl string `envconfig:"QUEUE_URL"`
	RedisUrl string `envconfig:"REDIS_URL"`
}

func Load() (Config, error) {
	err := LoadEnvironment()
	cfg := Config{}
	err = envconfig.Process("settings", &cfg)
	return cfg, err
}
