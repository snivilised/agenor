package traverse_test

import (
	"context"

	"github.com/fortytw2/leaktest"
	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok

	"github.com/snivilised/traverse"
	"github.com/snivilised/traverse/enums"
	"github.com/snivilised/traverse/internal/services"
	"github.com/snivilised/traverse/pref"
)

var _ = Describe("Traverse", Ordered, func() {
	BeforeEach(func() {
		services.Reset()
	})

	Context("speculation", func() {
		Context("Walk", func() {
			// We don't need to provide a context. For walk
			// cancellations, we use an internal context instead.
			//
			When("Prime", func() {
				It("🧪 should: walk primary navigation successfully", func() {
					defer leaktest.Check(GinkgoT())()

					_, err := traverse.Walk().Configure().Extent(traverse.Prime(
						pref.WithSubscription(enums.SubscribeFiles),
					)).Navigate()

					Expect(err).To(Succeed())
				})
			})

			When("Resume", func() {
				It("🧪 should: walk resume navigation successfully", func() {
					defer leaktest.Check(GinkgoT())()

					restore := func(_ *pref.Options) error {
						return nil
					}

					_, err := traverse.Walk().Configure().Extent(traverse.Resume(
						"/from-restore-path",
						enums.ResumeStrategyFastward,
						restore,
					)).Navigate()

					Expect(err).To(Succeed())
				})
			})
		})

		Context("Run", func() {
			When("Prime without cancel", func() {
				It("🧪 should: perform run navigation successfully", func() {
					defer leaktest.Check(GinkgoT())()

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
					_, err := traverse.Run().Configure().Extent(traverse.Prime(
						pref.WithSubscription(enums.SubscribeFiles),
						pref.WithContext(ctx),
					)).Navigate()

					Expect(err).To(Succeed())
				})
			})

			When("Prime with cancel", func() {
				It("🧪 should: perform run navigation successfully", func() {
					defer leaktest.Check(GinkgoT())()

					ctx, cancel := context.WithCancel(context.Background())

					_, err := traverse.Run().Configure().Extent(traverse.Prime(
						pref.WithSubscription(enums.SubscribeFiles),
						pref.WithContext(ctx),
						pref.WithCancel(cancel),
					)).Navigate()

					Expect(err).To(Succeed())
				})
			})

			When("Resume", func() {
				It("🧪 should: perform run navigation successfully", func() {
					defer leaktest.Check(GinkgoT())()

					restore := func(_ *pref.Options) error {
						return nil
					}

					_, err := traverse.Run().Configure().Extent(traverse.Resume(
						"/from-restore-path",
						enums.ResumeStrategySpawn,
						restore,
					)).Navigate()

					Expect(err).To(Succeed())
				})
			})
		})
	})
})
