package event_test

import (
	"blitzshare.event.worker/app/config"
	"blitzshare.event.worker/app/dependencies"
	"blitzshare.event.worker/app/domain"
	"blitzshare.event.worker/app/services/event"
	"blitzshare.event.worker/app/services/queue"
	"blitzshare.event.worker/mocks"
	"blitzshare.event.worker/test"
	"encoding/json"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
)

var _ = Describe("Test str modue", func() {
	var deps *dependencies.Dependencies
	peerCmd := domain.P2pPeerRegistryCmd{
		Otp:       "test-otp",
		Token:     "test-token",
		Mode:      "chat",
		MultiAddr: "tcp/ip4/0.0.0.0",
	}
	Context("Given Mq sub functional", func() {
		It("expected to call RegisterPeer on queue sub PeerRegisterCmd", func() {
			registerPeerCallCh := make(chan bool)
			mq := &mocks.Mq{}
			reg := &mocks.Registry{}
			registerPeerCalled := false
			reg.On("RegisterPeer", mock.MatchedBy(func(cmd *domain.P2pPeerRegistryCmd) bool {
				registerPeerCalled = cmd.Token == peerCmd.Token && cmd.Otp == peerCmd.Otp
				if registerPeerCalled {
					registerPeerCallCh <- registerPeerCalled
				}
				return registerPeerCalled
			})).Return("ackId", nil)
			mq.On("Sub",
				mock.MatchedBy(func(input string) bool {
					return input == queue.PeerRegisterCmd
				}),
				mock.MatchedBy(func(cb func(*[]byte)) bool {
					bytes, _ := json.Marshal(peerCmd)
					cb(&bytes)
					return cb != nil
				}),
			).Return()
			mq.On("Sub",
				mock.MatchedBy(test.MatchAny),
				mock.MatchedBy(test.MatchAny),
			).Return()
			deps = &dependencies.Dependencies{
				Registry: reg,
				Mq:       mq,
				Config: config.Config{
					QueueUrl: "QueueUrl",
					RedisUrl: "RedisUrl",
				},
			}
			event.PeerRegistry(deps)
			called := <-registerPeerCallCh
			Expect(called).To(BeTrue())
		})
	})
})
