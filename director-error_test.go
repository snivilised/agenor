package tv_test

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/fortytw2/leaktest"
	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok

	"github.com/snivilised/li18ngo"
	tv "github.com/snivilised/traverse"
	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/internal/services"
	"github.com/snivilised/traverse/locale"
	"github.com/snivilised/traverse/pref"
)

type traverseErrorTE struct {
	given string
	using *tv.Using
	relic *tv.Relic
}

var _ = Describe("director error", Ordered, func() {
	var handler tv.Client

	BeforeAll(func() {
		handler = func(_ tv.Servant) error {
			return nil
		}
		Expect(li18ngo.Use(
			func(o *li18ngo.UseOptions) {
				o.From.Sources = li18ngo.TranslationFiles{
					locale.SourceID: li18ngo.TranslationSource{Name: "traverse"},
				}
			},
		)).To(Succeed())
	})

	BeforeEach(func() {
		services.Reset()
	})

	DescribeTable("Validate",
		func(entry *traverseErrorTE) {
			if entry.using != nil {
				Expect(entry.using.Validate()).NotTo(Succeed())

				return
			}

			if entry.relic != nil {
				Expect(entry.relic.Validate()).NotTo(Succeed())

				return
			}
		},
		func(entry *traverseErrorTE) string {
			return fmt.Sprintf("given: %v, 🧪 should fail", entry.given)
		},

		Entry(nil, &traverseErrorTE{
			given: "using missing tree path",
			using: &tv.Using{
				Head: tv.Head{
					Subscription: tv.SubscribeFiles,
					Handler:      handler,
				},
			},
		}),

		Entry(nil, &traverseErrorTE{
			given: "using missing subscription",
			using: &tv.Using{
				Head: tv.Head{
					Handler: handler,
				},
				Tree: "/tree-traverse-path",
			},
		}),

		Entry(nil, &traverseErrorTE{
			given: "using missing handler",
			using: &tv.Using{
				Head: tv.Head{
					Subscription: tv.SubscribeFiles,
				},
				Tree: "/tree-traverse-path",
			},
		}),

		Entry(nil, &traverseErrorTE{
			given: "as missing restore from path",
			relic: &tv.Relic{
				Head: tv.Head{
					Subscription: tv.SubscribeFiles,
					Handler:      handler,
				},
				From:     "/resume-from-path",
				Strategy: tv.ResumeStrategySpawn,
			},
		}),

		Entry(nil, &traverseErrorTE{
			given: "as missing resume strategy",
			relic: &tv.Relic{
				Head: tv.Head{
					Subscription: tv.SubscribeFiles,
					Handler:      handler,
				},
				From: "/resume-from-path",
			},
		}),
	)

	When("Prime with subscription error", func() {
		It("🧪 should: fail", func(specCtx SpecContext) {
			defer leaktest.Check(GinkgoT())()

			ctx, cancel := context.WithCancel(specCtx)
			defer cancel()

			_, err := tv.Walk().Configure().Extent(tv.Prime(
				&pref.Using{
					Head: pref.Head{
						Handler: noOpHandler,
					},
					Tree: TreePath,
				},
			)).Navigate(ctx)

			Expect(err).To(MatchError(locale.ErrUsageMissingSubscription))
		})
	})

	When("Prime with options build error", func() {
		It("🧪 should: fail", func(specCtx SpecContext) {
			defer leaktest.Check(GinkgoT())()

			ctx, cancel := context.WithCancel(specCtx)
			defer cancel()

			_, err := tv.Walk().Configure().Extent(tv.Prime(
				&pref.Using{
					Head: pref.Head{
						Subscription: tv.SubscribeFiles,
						Handler:      noOpHandler,
					},
					Tree: TreePath,
				},
				func(_ *pref.Options) error {
					return errBuildOptions
				},
			)).Navigate(ctx)

			Expect(err).To(MatchError(errBuildOptions))
		})
	})

	When("Prime with subscription error", func() {
		It("🧪 should: fail", func(specCtx SpecContext) {
			defer leaktest.Check(GinkgoT())()

			ctx, cancel := context.WithCancel(specCtx)
			defer cancel()

			var wg sync.WaitGroup

			_, err := tv.Run(&wg).Configure().Extent(tv.Prime(
				&pref.Using{
					Head: pref.Head{
						Handler: noOpHandler,
					},
					Tree: TreePath,
				},
			)).Navigate(ctx)

			wg.Wait()
			Expect(err).NotTo(Succeed())
		})

		It("🧪 should: log error", func(specCtx SpecContext) {
			defer leaktest.Check(GinkgoT())()

			ctx, cancel := context.WithCancel(specCtx)
			defer cancel()

			invoked := false
			_, _ = tv.Walk().Configure().Extent(tv.Prime(
				&pref.Using{
					Head: pref.Head{
						Handler: noOpHandler,
					},
					Tree: TreePath,
				},
				tv.WithLogger(
					slog.New(slog.NewTextHandler(&TestWriter{
						assertFn: func() {
							invoked = true
						},
					}, nil)),
				),
			)).Navigate(ctx)

			Expect(invoked).To(BeTrue(), "validation error not logged")
		})
	})

	When("incorrect facade", func() {
		Context("primary (expected using)", func() {
			It("🧪 should: return error", func(specCtx SpecContext) {
				defer leaktest.Check(GinkgoT())()

				ctx, cancel := context.WithCancel(specCtx)
				defer cancel()

				_, err := tv.Walk().Configure().Extent(tv.Prime(
					&pref.Relic{
						Head: pref.Head{
							Subscription: tv.SubscribeFiles,
							Handler:      noOpHandler,
						},
						From: "/from-path/wrong-facade/primary/relic",
					},
				)).Navigate(ctx)

				Expect(err).To(MatchError(core.ErrWrongPrimaryFacade))
			})
		})

		Context("resume (expected relic)", func() {
			It("🧪 should: return error", func(specCtx SpecContext) {
				defer leaktest.Check(GinkgoT())()

				ctx, cancel := context.WithCancel(specCtx)
				defer cancel()

				_, err := tv.Walk().Configure().Extent(tv.Resume(
					&pref.Using{
						Head: pref.Head{
							Subscription: tv.SubscribeFiles,
							Handler:      noOpHandler,
						},
						Tree: "/tree-path/wrong-facade/resume/using",
					},
				)).Navigate(ctx)

				Expect(err).To(MatchError(core.ErrWrongResumeFacade))
			})
		})
	})
})
