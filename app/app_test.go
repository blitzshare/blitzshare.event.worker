package app_test

import (
	"testing"

	"blitzshare.event.worker/app"
	"blitzshare.event.worker/app/config"
	"blitzshare.event.worker/app/dependencies"
	"blitzshare.event.worker/mocks"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestApp(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Worker test")
}

var _ = Describe("Test str modue", func() {
	var deps *dependencies.Dependencies
	BeforeEach(func() {
		deps = &dependencies.Dependencies{
			Registry: &mocks.Registry{},
			Config: config.Config{
				QueueUrl: "QueueUrl",
				RedisUrl: "RedisUrl",
			},
		}
	})
	Context("?", func() {
		It("?", func() {
			// TODO: Simulate queue sub functionality
			app.Start(deps)
			Expect(2).To(Equal(1 + 1))
		})
	})
})
