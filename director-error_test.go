package age_test

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/fortytw2/leaktest"
	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok

	age "github.com/snivilised/agenor"
	"github.com/snivilised/agenor/core"
	"github.com/snivilised/agenor/enums"
	"github.com/snivilised/agenor/internal/enclave"
	lab "github.com/snivilised/agenor/internal/laboratory"
	"github.com/snivilised/agenor/internal/services"
	"github.com/snivilised/agenor/locale"
	"github.com/snivilised/agenor/pref"
	"github.com/snivilised/agenor/test/hydra"
	"github.com/snivilised/li18ngo"
)

type traverseErrorTE struct {
	given string
	using *age.Using
	relic *age.Relic
}

var _ = Describe("director error", Ordered, func() {
	var handler age.Client

	BeforeAll(func() {
		handler = func(_ age.Servant) error {
			return nil
		}
		Expect(li18ngo.Use(
			func(o *li18ngo.UseOptions) {
				o.From.Sources = li18ngo.TranslationFiles{
					locale.SourceID: li18ngo.TranslationSource{Name: "agenor"},
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
			using: &age.Using{
				Subscription: age.SubscribeFiles,
				Head: age.Head{
					Handler: handler,
				},
			},
		}),

		Entry(nil, &traverseErrorTE{
			given: "using missing subscription",
			using: &age.Using{
				Head: age.Head{
					Handler: handler,
				},
				Tree: "/tree-traverse-path",
			},
		}),

		Entry(nil, &traverseErrorTE{
			given: "using missing handler",
			using: &age.Using{
				Subscription: age.SubscribeFiles,
				Head:         age.Head{},
				Tree:         "/tree-traverse-path",
			},
		}),

		Entry(nil, &traverseErrorTE{
			given: "as missing restore from path",
			relic: &age.Relic{
				Head: age.Head{
					Handler: handler,
				},
				From:     "/resume-from-path",
				Strategy: age.ResumeStrategySpawn,
			},
		}),

		Entry(nil, &traverseErrorTE{
			given: "as missing resume strategy",
			relic: &age.Relic{
				Head: age.Head{
					Handler: handler,
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

			_, err := age.Walk().Configure().Extent(age.Prime(
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

			_, err := age.Walk().Configure().Extent(age.Prime(
				&pref.Using{
					Subscription: age.SubscribeFiles,
					Head: pref.Head{
						Handler: noOpHandler,
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

	When("Resume with subscription error", func() {
		It("🧪 should: fail", func(specCtx SpecContext) {
			// In case user has tampered with the json file changing
			// the subscription to an inappropriate value
			defer leaktest.Check(GinkgoT())()

			ctx, cancel := context.WithCancel(specCtx)
			defer cancel()

			_, err := age.Walk().Configure(enclave.Loader(func(active *core.ActiveState) {
				active.Tree = hydra.Repo("test")
				active.TraverseDescription.IsRelative = false
				active.ResumeDescription.IsRelative = false
				active.Subscription = enums.SubscribeUndefined
			})).Extent(age.Resume(
				&pref.Relic{
					Head: pref.Head{
						Handler: noOpHandler,
					},
					From:     lab.GetJSONPath(),
					Strategy: age.ResumeStrategyFastward,
				},
				age.WithFaultHandler(age.Accepter(lab.IgnoreFault)),
			)).Navigate(ctx)

			Expect(err).To(MatchError(locale.ErrUsageMissingSubscription))
		})
	})

	When("Prime with subscription error", func() {
		It("🧪 should: fail", func(specCtx SpecContext) {
			defer leaktest.Check(GinkgoT())()

			ctx, cancel := context.WithCancel(specCtx)
			defer cancel()

			var wg sync.WaitGroup

			_, err := age.Run(&wg).Configure().Extent(age.Prime(
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
			_, _ = age.Walk().Configure().Extent(age.Prime(
				&pref.Using{
					Head: pref.Head{
						Handler: noOpHandler,
					},
					Tree: TreePath,
				},
				age.WithLogger(
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

				_, err := age.Walk().Configure().Extent(age.Prime(
					&pref.Relic{
						Head: pref.Head{
							Handler: noOpHandler,
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

				_, err := age.Walk().Configure().Extent(age.Resume(
					&pref.Using{
						Subscription: age.SubscribeFiles,
						Head: pref.Head{
							Handler: noOpHandler,
						},
						Tree: "/tree-path/wrong-facade/resume/using",
					},
				)).Navigate(ctx)

				Expect(err).To(MatchError(core.ErrWrongResumeFacade))
			})
		})
	})
})
