package traverse_test

import (
	"context"

	"github.com/fortytw2/leaktest"
	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok

	"github.com/snivilised/traverse"
	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/enums"
	"github.com/snivilised/traverse/pref"
)

var _ = Describe("Traverse", func() {
	Context("speculation", func() {
		Context("Walk", func() {
			// We don't need to provide a context. For walk
			// cancellations, we use an internal context instead.
			//
			When("Primary", func() {
				It("should: walk primary navigation successfully", func() {
					defer leaktest.Check(GinkgoT())()

					_, err := traverse.Walk().Primary(
						"/root-path",
						func(node *core.Node) error {
							_ = node

							return nil
						},
						pref.WithSubscription(enums.SubscribeFiles),
					).Navigate()
					Expect(err).To(Succeed())
				})
			})

			When("Resume", func() {
				It("should: walk resume navigation successfully", func() {
					defer leaktest.Check(GinkgoT())()

					_, err := traverse.Walk().Resume(
						"/from-restore-path",
					).Navigate()
					Expect(err).To(Succeed())
				})
			})
		})

		Context("Run", func() {
			When("Primary without cancel", func() {
				It("should: perform run navigation successfully", func() {
					// need to make sure that when a ctrl-c occurs, who is
					// responsible for handling the cancellation; ie if a
					// ctrl-c occurs should the client handle it or do we?
					//
					// Internally, we could create our own child context
					// from this parent content which contains a cancelFunc.
					// This way, when ctrl-c occurs, we can trap that,
					// and perform a save. If we don't do this, then how
					// can we tap into cancellation?
					//
					// The context has a lifetime. The kernel will know when
					// it has become invalidated, at which point a message
					// is sent on the message bus, on a topic called
					// something like "context.expired"
					//
					ctx := context.Background()
					_, err := traverse.Run(
						pref.WithContext(ctx),
					).Primary(
						"/root-path",
						func(node *core.Node) error {
							_ = node

							return nil
						},
						pref.WithSubscription(enums.SubscribeFiles),
					).Navigate()
					Expect(err).To(Succeed())
				})
			})

			When("Primary with cancel", func() {
				It("should: run primary navigation successfully", func() {
					defer leaktest.Check(GinkgoT())()

					ctx, cancel := context.WithCancel(context.Background())
					_, err := traverse.Run(
						pref.WithContext(ctx),
						pref.WithCancel(cancel),
					).Primary(
						"/root-path",
						func(node *core.Node) error {
							_ = node

							return nil
						},
						pref.WithSubscription(enums.SubscribeFiles),
					).Navigate()
					Expect(err).To(Succeed())
				})
			})

			When("Resume", func() {
				It("should: run primary navigation successfully", func() {
					defer leaktest.Check(GinkgoT())()

					ctx := context.Background()
					_, err := traverse.Run(
						pref.WithContext(ctx),
					).Resume(
						"/from-restore-path",
					).Navigate()
					Expect(err).To(Succeed())
				})
			})
		})
	})
})
