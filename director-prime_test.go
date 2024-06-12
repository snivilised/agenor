package tv_test

import (
	"context"
	"sync"

	"github.com/fortytw2/leaktest"
	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok
	tv "github.com/snivilised/traverse"
	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/internal/services"
	"github.com/snivilised/traverse/pref"
)

var _ = Describe("Director(Prime)", func() {
	BeforeEach(func() {
		services.Reset()
	})

	Context("simple", func() {
		Context("Walk", func() {
			When("Options", func() {
				It("🧪 should: walk primary navigation successfully", func(specCtx SpecContext) {
					defer leaktest.Check(GinkgoT())()

					ctx, cancel := context.WithCancel(specCtx)
					defer cancel()

					_, err := tv.Walk().Configure().Extent(tv.Prime(
						&tv.Using{
							Root:         RootPath,
							Subscription: tv.SubscribeFiles,
							Handler:      noOpHandler,
						},
						tv.WithOnAscend(func(_ *core.Node) {}),
					)).Navigate(ctx)

					Expect(err).To(Succeed())
				})
			})

			When("Pushed Options", func() {
				It("🧪 should: walk primary navigation successfully", func(specCtx SpecContext) {
					defer leaktest.Check(GinkgoT())()

					ctx, cancel := context.WithCancel(specCtx)
					defer cancel()

					o, _ := pref.Get()
					_, err := tv.Walk().Configure().Extent(tv.Prime(
						&tv.Using{
							Root:         RootPath,
							Subscription: tv.SubscribeFiles,
							Handler:      noOpHandler,
							O:            o,
						},
						tv.WithOnDescend(func(_ *core.Node) {}),
					)).Navigate(ctx)

					Expect(err).To(Succeed())
				})
			})
		})

		Context("Run", func() {
			When("Options", func() {
				It("🧪 should: perform run navigation successfully", func(specCtx SpecContext) {
					defer leaktest.Check(GinkgoT())()

					ctx, cancel := context.WithCancel(specCtx)
					defer cancel()

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

					_, err := tv.Run(&wg).Configure().Extent(tv.Prime(
						&tv.Using{
							Root:         RootPath,
							Subscription: tv.SubscribeFiles,
							Handler:      noOpHandler,
						},
						tv.WithOnBegin(func(_ string) {}),
					)).Navigate(ctx)

					wg.Wait()
					Expect(err).To(Succeed())
				})
			})

			When("Options Pushed", func() {
				It("🧪 should: run primary navigation successfully", func(specCtx SpecContext) {
					defer leaktest.Check(GinkgoT())()

					ctx, cancel := context.WithCancel(specCtx)
					defer cancel()

					var wg sync.WaitGroup

					o, _ := pref.Get()
					_, err := tv.Run(&wg).Configure().Extent(tv.Prime(
						&tv.Using{
							Root:         RootPath,
							Subscription: tv.SubscribeFiles,
							Handler:      noOpHandler,
							O:            o,
						},
						tv.WithOnEnd(func(_ core.TraverseResult) {}),
					)).Navigate(ctx)

					wg.Wait()
					Expect(err).To(Succeed())
				})
			})
		})
	})

	Context("features", func() {
		Context("Walk", func() {
			When("filter", func() {
				It("🧪 should: register ok", func(specCtx SpecContext) {
					defer leaktest.Check(GinkgoT())()

					ctx, cancel := context.WithCancel(specCtx)
					defer cancel()

					_, err := tv.Walk().Configure().Extent(tv.Prime(
						&tv.Using{
							Root:         RootPath,
							Subscription: tv.SubscribeFiles,
							Handler:      noOpHandler,
						},
						tv.WithFilter(&core.FilterDef{}),
						tv.WithOnStart(func(_ string) {}),
					)).Navigate(ctx)

					Expect(err).To(Succeed())
				})
			})

			When("hibernate", func() {
				It("🧪 should: register ok", func(specCtx SpecContext) {
					defer leaktest.Check(GinkgoT())()

					ctx, cancel := context.WithCancel(specCtx)
					defer cancel()

					_, err := tv.Walk().Configure().Extent(tv.Prime(
						&tv.Using{
							Root:         RootPath,
							Subscription: tv.SubscribeFiles,
							Handler:      noOpHandler,
						},
						tv.WithHibernationWake(&core.FilterDef{}),
						tv.WithOnStop(func(_ string) {}),
					)).Navigate(ctx)

					Expect(err).To(Succeed())
				})
			})

			When("sample", func() {
				It("🧪 should: register ok", func(specCtx SpecContext) {
					defer leaktest.Check(GinkgoT())()

					ctx, cancel := context.WithCancel(specCtx)
					defer cancel()

					_, err := tv.Walk().Configure().Extent(tv.Prime(
						&tv.Using{
							Root:         RootPath,
							Subscription: tv.SubscribeFiles,
							Handler:      noOpHandler,
						},
						tv.WithHibernationSleep(&core.FilterDef{}),
						tv.WithOnStop(func(_ string) {}),
					)).Navigate(ctx)

					Expect(err).To(Succeed())
				})
			})
		})
	})
})
