package tv_test

import (
	"context"

	"github.com/fortytw2/leaktest"
	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok

	tv "github.com/snivilised/traverse"
	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/internal/services"
	"github.com/snivilised/traverse/pref"
)

const (
	RootPath    = "/traversal-root-path"
	RestorePath = "/from-restore-path"
	files       = 3
	folders     = 2
)

var _ = Describe("Traverse", Ordered, func() {
	var restore pref.Option

	BeforeAll(func() {
		restore = func(o *tv.Options) error {
			o.Events.Begin.On(func(_ string) {})

			return nil
		}
	})

	BeforeEach(func() {
		services.Reset()
	})

	Context("simple", func() {
		Context("Walk", func() {
			// We don't need to provide a context. For walk
			// cancellations, we use an internal context instead.
			//
			When("Prime", func() {
				It("🧪 should: walk primary navigation successfully", func() {
					defer leaktest.Check(GinkgoT())()

					_, err := tv.Walk().Configure().Extent(tv.Prime(
						tv.Using{
							Root:         RootPath,
							Subscription: tv.SubscribeFiles,
							Handler: func(_ *tv.Node) error {
								return nil
							},
						},
						tv.WithSubscription(tv.SubscribeFiles),
					)).Navigate()

					Expect(err).To(Succeed())
				})
			})

			When("Resume", func() {
				It("🧪 should: walk resume navigation successfully", func() {
					defer leaktest.Check(GinkgoT())()

					_, err := tv.Walk().Configure().Extent(tv.Resume(
						tv.As{
							Using: tv.Using{
								Root:         RootPath,
								Subscription: tv.SubscribeFiles,
								Handler: func(_ *tv.Node) error {
									return nil
								},
							},
							From:     RestorePath,
							Strategy: tv.ResumeStrategyFastward,
						},
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
					_, err := tv.Run().Configure().Extent(tv.Prime(
						tv.Using{
							Root:         RootPath,
							Subscription: tv.SubscribeFiles,
							Handler: func(_ *tv.Node) error {
								return nil
							},
						},
						tv.WithSubscription(tv.SubscribeFiles),
						tv.WithContext(ctx),
					)).Navigate()

					Expect(err).To(Succeed())
				})
			})

			When("Prime with cancel", func() {
				It("🧪 should: perform run navigation successfully", func() {
					defer leaktest.Check(GinkgoT())()

					ctx, cancel := context.WithCancel(context.Background())

					_, err := tv.Run().Configure().Extent(tv.Prime(
						tv.Using{
							Root:         RootPath,
							Subscription: tv.SubscribeFiles,
							Handler: func(_ *tv.Node) error {
								return nil
							},
						},
						tv.WithSubscription(tv.SubscribeFiles),
						tv.WithContext(ctx),
						tv.WithCancel(cancel),
					)).Navigate()

					Expect(err).To(Succeed())
				})
			})

			When("Resume", func() {
				It("🧪 should: perform run navigation successfully", func() {
					defer leaktest.Check(GinkgoT())()

					_, err := tv.Run().Configure().Extent(tv.Resume(
						tv.As{
							Using: tv.Using{
								Root:         RootPath,
								Subscription: tv.SubscribeFiles,
								Handler: func(_ *tv.Node) error {
									return nil
								},
							},
							From:     RestorePath,
							Strategy: tv.ResumeStrategySpawn,
						},
						restore,
					)).Navigate()

					Expect(err).To(Succeed())
				})
			})
		})
	})

	Context("features", func() {
		Context("Prime", func() {
			When("hibernate", func() {
				It("🧪 should: register ok", func() {
					defer leaktest.Check(GinkgoT())()

					_, err := tv.Walk().Configure().Extent(tv.Prime(
						tv.Using{
							Root:         RootPath,
							Subscription: tv.SubscribeFiles,
							Handler: func(_ *tv.Node) error {
								return nil
							},
						},
						tv.WithSubscription(tv.SubscribeFiles),
						tv.WithHibernation(&core.FilterDef{}),
					)).Navigate()

					Expect(err).To(Succeed())
				})
			})

			When("filter", func() {
				It("🧪 should: register ok", func() {
					defer leaktest.Check(GinkgoT())()

					_, err := tv.Walk().Configure().Extent(tv.Prime(
						tv.Using{
							Root:         RootPath,
							Subscription: tv.SubscribeFiles,
							Handler: func(_ *tv.Node) error {
								return nil
							},
						},
						tv.WithSubscription(tv.SubscribeFiles),
						tv.WithFilter(&core.FilterDef{}),
					)).Navigate()

					Expect(err).To(Succeed())
				})
			})

			When("sample", func() {
				It("🧪 should: register ok", func() {
					defer leaktest.Check(GinkgoT())()

					_, err := tv.Walk().Configure().Extent(tv.Prime(
						tv.Using{
							Root:         RootPath,
							Subscription: tv.SubscribeFiles,
							Handler: func(_ *tv.Node) error {
								return nil
							},
						},
						tv.WithSubscription(tv.SubscribeFiles),
						tv.WithSampling(files, folders),
					)).Navigate()

					Expect(err).To(Succeed())
				})
			})
		})

		Context("Resume", func() {
			When("hibernate", func() {
				It("🧪 should: register ok", func() {
					defer leaktest.Check(GinkgoT())()

					_, err := tv.Run().Configure().Extent(tv.Resume(
						tv.As{
							Using: tv.Using{
								Root:         RootPath,
								Subscription: tv.SubscribeFiles,
								Handler: func(_ *tv.Node) error {
									return nil
								},
							},
							From:     RestorePath,
							Strategy: tv.ResumeStrategySpawn,
						},
						tv.WithHibernation(&core.FilterDef{}),
					)).Navigate()

					Expect(err).To(Succeed())
				})
			})

			When("filter", func() {
				It("🧪 should: register ok", func() {
					defer leaktest.Check(GinkgoT())()

					_, err := tv.Run().Configure().Extent(tv.Resume(
						tv.As{
							Using: tv.Using{
								Root:         RootPath,
								Subscription: tv.SubscribeFiles,
								Handler: func(_ *tv.Node) error {
									return nil
								},
							},
							From:     RestorePath,
							Strategy: tv.ResumeStrategySpawn,
						},
						tv.WithFilter(&core.FilterDef{}),
					)).Navigate()

					Expect(err).To(Succeed())
				})
			})

			When("sample", func() {
				It("🧪 should: register ok", func() {
					defer leaktest.Check(GinkgoT())()

					_, err := tv.Run().Configure().Extent(tv.Resume(
						tv.As{
							Using: tv.Using{
								Root:         RootPath,
								Subscription: tv.SubscribeFiles,
								Handler: func(_ *tv.Node) error {
									return nil
								},
							},
							From:     RestorePath,
							Strategy: tv.ResumeStrategySpawn,
						},
						tv.WithSampling(files, folders),
					)).Navigate()

					Expect(err).To(Succeed())
				})
			})
		})
	})
})
