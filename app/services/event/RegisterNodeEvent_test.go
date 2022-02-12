package event_test

import (
	"blitzshare.event.worker/app"
	"blitzshare.event.worker/app/config"
	"blitzshare.event.worker/app/dependencies"
	"blitzshare.event.worker/app/domain"
	"blitzshare.event.worker/app/services/queue"
	"blitzshare.event.worker/mocks"
	"blitzshare.event.worker/test"
	"encoding/json"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	"reflect"
	"runtime"
	"strings"
)

var _ = Describe("Node registry event testss", func() {
	var deps *dependencies.Dependencies
	Context("Given Mq sub functional", func() {
		It("expected to call RegisterNode on queue sub P2pBootstrapNodeRegistryCmd", func() {
			nodeRegistryCmd := domain.P2pBootstrapNodeRegistryCmd{
				Port:   6543,
				NodeId: "jsdfklsjdlkfjsdkfjklsdfjk",
			}
			registerPeerCallCh := make(chan bool)
			mq := &mocks.Mq{}
			reg := &mocks.Registry{}

			reg.On("RegisterNode", mock.MatchedBy(func(cmd *domain.P2pBootstrapNodeRegistryCmd) bool {
				check := cmd.NodeId == nodeRegistryCmd.NodeId && cmd.Port == nodeRegistryCmd.Port
				if check {
					registerPeerCallCh <- check
				}
				return check
			})).Return("ackId", nil)

			mq.On("Sub",
				mock.MatchedBy(func(input string) bool {
					return input == queue.P2pBootstrapNodeRegistryCmd
				}),
				mock.MatchedBy(func(cb func(*[]byte)) bool {
					// TODO: change  from reflection function name check to type check
					name := runtime.FuncForPC(reflect.ValueOf(cb).Pointer()).Name()
					check := strings.Contains(name, "NodeRegistry")
					if check {
						bytes, _ := json.Marshal(nodeRegistryCmd)
						cb(&bytes)
					}
					return check
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
			app.Start(deps)
			called := <-registerPeerCallCh
			Expect(called).To(BeTrue())
		})

	})
})
