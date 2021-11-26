package config_test

import (
	"os"
	"testing"

	"blitzshare.event.worker/app/config"
	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	setUp()
	cfg, err := config.Load()

	assert.Nil(t, err, "Unable to log the config")
	assert.Equal(t, cfg.QueueUrl, "kubemq-cluster-grpc.kubemq.svc.cluster.local")
	assert.Equal(t, cfg.RedisUrl, "redis-svc.blitzshare-api-ns.svc.cluster.local:6379")
	tearDown()
}

func setUp() {
	_ = os.Setenv("QUEUE_URL", "kubemq-cluster-grpc.kubemq.svc.cluster.local")
	_ = os.Setenv("REDIS_URL", "")
}

func tearDown() {
	_ = os.Unsetenv("ENV")
}
