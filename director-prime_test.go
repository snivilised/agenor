package age_test

import (
	"context"
	"sync"

	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok

	age "github.com/snivilised/agenor"
	"github.com/snivilised/agenor/core"
	"github.com/snivilised/agenor/enums"
	lab "github.com/snivilised/agenor/internal/laboratory"
	"github.com/snivilised/agenor/internal/opts"
	"github.com/snivilised/agenor/internal/services"
	"github.com/snivilised/agenor/life"
	"github.com/snivilised/agenor/locale"
	"github.com/snivilised/agenor/pref"
	"github.com/snivilised/agenor/test/hanno"
	"github.com/snivilised/li18ngo"
)

var _ = Describe("Director(Prime)", Ordered, func() {
	var (
		tree string
	)

	// 👽 These tests are not using Nuxx therefore they are traversing the
	// local test directory.
	BeforeAll(func() {
		Expect(li18ngo.Use(
			func(o *li18ngo.UseOptions) {
				o.From.Sources = li18ngo.TranslationFiles{
					locale.SourceID: li18ngo.TranslationSource{Name: "agenor"},
				}
			},
		)).To(Succeed())

		tree = hanno.Repo("test")
	})

	BeforeEach(func() {
		services.Reset()
	})

	Context("simple", func() {
		Context("Walk", func() {
			When("Options", func() {
				It("🧪 should: walk primary navigation successfully", func(specCtx SpecContext) {
					lab.WithTestContext(specCtx, func(ctx context.Context, _ context.CancelFunc) {
						_, err := age.Walk().Configure().Extent(age.Prime(
							&pref.Using{
								Subscription: age.SubscribeFiles,
								Head: pref.Head{
									Handler: noOpHandler,
								},
								Tree: tree,
							},
							age.WithOnAscend(func(_ *core.Node) {}),
							age.WithNoRecurse(),
							age.WithFaultHandler(age.Accepter(lab.IgnoreFault)),
						)).Navigate(ctx)

						Expect(err).To(Succeed())
					})
				})
			})

			When("Push Options", func() {
				It("🧪 should: walk primary navigation successfully", func(specCtx SpecContext) {
					lab.WithTestContext(specCtx, func(ctx context.Context, _ context.CancelFunc) {
						o, _, _ := opts.Get()
						o.Defects.Fault = age.Accepter(lab.IgnoreFault)

						_, err := age.Walk().Configure().Extent(age.Prime(
							&pref.Using{
								Subscription: age.SubscribeFiles,
								Head: pref.Head{
									Handler: noOpHandler,
								},
								Tree: TreePath,
								O:    o,
							},
							age.WithOnDescend(func(_ *core.Node) {}),
						)).Navigate(ctx)

						Expect(err).To(Succeed())
					})
				})
			})
		})

		Context("Run", func() {
			When("Options", func() {
				It("🧪 should: perform run navigation successfully", func(specCtx SpecContext) {
					lab.WithTestContext(specCtx, func(ctx context.Context, _ context.CancelFunc) {
						var wg sync.WaitGroup

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

						_, err := age.Run(&wg).Configure().Extent(age.Prime(
							&pref.Using{
								Subscription: age.SubscribeFiles,
								Head: pref.Head{
									Handler: noOpHandler,
								},
								Tree: TreePath,
							},
							age.WithOnBegin(func(_ *life.BeginState) {}),
							age.WithCPU(),
							age.WithFaultHandler(age.Accepter(lab.IgnoreFault)),
						)).Navigate(ctx)

						wg.Wait()
						Expect(err).To(Succeed())
					})
				})
			})

			When("Push Options", func() {
				It("🧪 should: run primary navigation successfully", func(specCtx SpecContext) {
					lab.WithTestContext(specCtx, func(ctx context.Context, _ context.CancelFunc) {
						var wg sync.WaitGroup

						o, _, _ := opts.Get()
						o.Defects.Fault = age.Accepter(lab.IgnoreFault)
						_, err := age.Run(&wg).Configure().Extent(age.Prime(
							&pref.Using{
								Subscription: age.SubscribeFiles,
								Head: pref.Head{
									Handler: noOpHandler,
								},
								Tree: TreePath,
								O:    o,
							},
							age.WithOnEnd(func(_ core.TraverseResult) {}),
						)).Navigate(ctx)

						wg.Wait()
						Expect(err).To(Succeed())

					})
				})
			})
		})
	})

	Context("features", func() {
		Context("Walk", func() {
			When("filter", func() {
				It("🧪 should: register ok", func(specCtx SpecContext) {
					lab.WithTestContext(specCtx, func(ctx context.Context, _ context.CancelFunc) {
						_, err := age.Walk().Configure().Extent(age.Prime(
							&pref.Using{
								Subscription: age.SubscribeFiles,
								Head: pref.Head{
									Handler: noOpHandler,
								},
								Tree: TreePath,
							},
							age.WithFilter(&pref.FilterOptions{}),
							age.WithOnWake(func(_ string) {}),
							age.WithFaultHandler(age.Accepter(lab.IgnoreFault)),
						)).Navigate(ctx)

						Expect(err).To(Succeed())
					})
				})
			})

			When("hibernate", func() {
				It("🧪 should: register ok", func(specCtx SpecContext) {
					lab.WithTestContext(specCtx, func(ctx context.Context, _ context.CancelFunc) {
						_, err := age.Walk().Configure().Extent(age.Prime(
							&pref.Using{
								Subscription: age.SubscribeFiles,
								Head: pref.Head{
									Handler: noOpHandler,
								},
								Tree: TreePath,
							},
							age.WithHibernationFilterWake(&core.FilterDef{
								Description: "nonsense",
								Type:        enums.FilterTypeGlob,
								Pattern:     "*",
							}),
							age.WithOnSleep(func(_ string) {}),
							age.WithFaultHandler(age.Accepter(lab.IgnoreFault)),
						)).Navigate(ctx)

						Expect(err).To(Succeed())
					})
				})
			})

			When("sample", func() {
				It("🧪 should: register ok", func(specCtx SpecContext) {
					lab.WithTestContext(specCtx, func(ctx context.Context, _ context.CancelFunc) {
						_, err := age.Walk().Configure().Extent(age.Prime(
							&pref.Using{
								Subscription: age.SubscribeFiles,
								Head: pref.Head{
									Handler: noOpHandler,
								},
								Tree: TreePath,
							},
							age.WithHibernationFilterSleep(&core.FilterDef{
								Description: "nonsense",
								Type:        enums.FilterTypeGlob,
								Pattern:     "*",
							}),
							age.WithOnSleep(func(_ string) {}),
							age.WithFaultHandler(age.Accepter(lab.IgnoreFault)),
						)).Navigate(ctx)

						Expect(err).To(Succeed())
					})
				})
			})
		})
	})
})
