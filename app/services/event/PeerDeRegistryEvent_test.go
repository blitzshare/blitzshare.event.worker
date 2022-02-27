package event_test

import (
	"blitzshare.event.worker/app/config"
	"blitzshare.event.worker/app/dependencies"
	"blitzshare.event.worker/app/domain"
	"blitzshare.event.worker/app/services/event"
	"blitzshare.event.worker/mocks"
	"blitzshare.event.worker/test"
	"encoding/json"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
)

var _ = Describe("Test str modue", func() {
	var deps *dependencies.Dependencies
	Context("Given Mq sub functional", func() {
		It("expected to call DeregisterPeer on queue sub MqPeerDeregisterCmd", func() {
			deRegistryCmd := domain.P2pPeerDeregisterCmd{
				Otp:   "test-otp",
				Token: "deregister-token",
			}
			registerPeerCallCh := make(chan bool)
			mq := &mocks.Mq{}
			reg := &mocks.Registry{}
			mq.On("Sub",
				mock.MatchedBy(func(input string) bool {
					return input == config.MqPeerDeregisterCmd
				}),
				mock.MatchedBy(func(cb func(*[]byte)) bool {
					bytes, _ := json.Marshal(deRegistryCmd)
					cb(&bytes)
					return cb != nil
				}),
			).Return()
			mq.On("Sub",
				mock.MatchedBy(test.MatchAny),
				mock.MatchedBy(test.MatchAny),
			).Return()
			reg.On("DeregisterPeer", mock.MatchedBy(func(cmd *domain.P2pPeerDeregisterCmd) bool {
				check := cmd.Otp == deRegistryCmd.Otp && cmd.Token == deRegistryCmd.Token
				if check {
					registerPeerCallCh <- check
				}
				return check
			})).Return(nil)

			deps = &dependencies.Dependencies{
				Registry: reg,
				Mq:       mq,
				Config: config.Config{
					QueueUrl: "QueueUrl",
					RedisUrl: "RedisUrl",
				},
			}
			event.PeerDeRegistry(deps)
			called := <-registerPeerCallCh
			Expect(called).To(BeTrue())
		})
	})
})
