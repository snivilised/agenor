package agenor_test

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/snivilised/jaywalk/src/agenor"
	"github.com/snivilised/jaywalk/src/agenor/core"
	"github.com/snivilised/jaywalk/src/agenor/enums"
	"github.com/snivilised/jaywalk/src/agenor/internal/enclave"
	lab "github.com/snivilised/jaywalk/src/agenor/internal/laboratory"
	"github.com/snivilised/jaywalk/src/internal/services"
	"github.com/snivilised/jaywalk/locale"
	"github.com/snivilised/jaywalk/src/agenor/pref"
	"github.com/snivilised/jaywalk/src/agenor/test/hanno"
	"github.com/snivilised/li18ngo"
)

type traverseErrorTE struct {
	given string
	using *agenor.Using
	relic *agenor.Relic
}

var _ = Describe("director error", Ordered, func() {
	var handler agenor.Client

	BeforeAll(func() {
		handler = func(_ agenor.Servant) error {
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
		func(entry *traverseErrorTE) string { // !!!
			return fmt.Sprintf("given: %v, 🧪 should fail", entry.given)
		},

		Entry(nil, &traverseErrorTE{
			given: "using missing tree path",
			using: &agenor.Using{
				Subscription: agenor.SubscribeFiles,
				Head: agenor.Head{
					Handler: handler,
				},
			},
		}),

		Entry(nil, &traverseErrorTE{
			given: "using missing subscription",
			using: &agenor.Using{
				Head: agenor.Head{
					Handler: handler,
				},
				Tree: "/tree-traverse-path",
			},
		}),

		Entry(nil, &traverseErrorTE{
			given: "using missing handler",
			using: &agenor.Using{
				Subscription: agenor.SubscribeFiles,
				Head:         agenor.Head{},
				Tree:         "/tree-traverse-path",
			},
		}),

		Entry(nil, &traverseErrorTE{
			given: "as missing restore from path",
			relic: &agenor.Relic{
				Head: agenor.Head{
					Handler: handler,
				},
				From:     "/resume-from-path",
				Strategy: agenor.ResumeStrategySpawn,
			},
		}),

		Entry(nil, &traverseErrorTE{
			given: "as missing resume strategy",
			relic: &agenor.Relic{
				Head: agenor.Head{
					Handler: handler,
				},
				From: "/resume-from-path",
			},
		}),
	)

	When("Prime with subscription error", func() {
		It("🧪 should: fail", func(specCtx SpecContext) {
			lab.WithTestContext(specCtx, func(ctx context.Context, _ context.CancelFunc) {
				_, err := agenor.Walk().Configure().Extent(agenor.Prime(
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
	})

	When("Prime with options build error", func() {
		It("🧪 should: fail", func(specCtx SpecContext) {
			lab.WithTestContext(specCtx, func(ctx context.Context, _ context.CancelFunc) {
				_, err := agenor.Walk().Configure().Extent(agenor.Prime(
					&pref.Using{
						Subscription: agenor.SubscribeFiles,
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
	})

	When("Resume with subscription error", func() {
		It("🧪 should: fail", func(specCtx SpecContext) {
			lab.WithTestContext(specCtx, func(ctx context.Context, _ context.CancelFunc) {
				// In case user has tampered with the json file changing
				// the subscription to an inappropriate value
				_, err := agenor.Walk().Configure(enclave.Loader(func(active *core.ActiveState) {
					active.Tree = hanno.Repo("test")
					active.TraverseDescription.IsRelative = false
					active.ResumeDescription.IsRelative = false
					active.Subscription = enums.SubscribeUndefined
				})).Extent(agenor.Resume(
					&pref.Relic{
						Head: pref.Head{
							Handler: noOpHandler,
						},
						From:     lab.GetJSONPath(),
						Strategy: agenor.ResumeStrategyFastward,
					},
					agenor.WithFaultHandler(agenor.Accepter(lab.IgnoreFault)),
				)).Navigate(ctx)

				Expect(err).To(MatchError(locale.ErrUsageMissingSubscription))
			})
		})
	})

	When("Prime with subscription error", func() {
		It("🧪 should: fail", func(specCtx SpecContext) {
			lab.WithTestContext(specCtx, func(ctx context.Context, _ context.CancelFunc) {
				var wg sync.WaitGroup

				_, err := agenor.Run(&wg).Configure().Extent(agenor.Prime(
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
		})

		It("🧪 should: log error", func(specCtx SpecContext) {
			lab.WithTestContext(specCtx, func(ctx context.Context, _ context.CancelFunc) {
				invoked := false
				_, _ = agenor.Walk().Configure().Extent(agenor.Prime(
					&pref.Using{
						Head: pref.Head{
							Handler: noOpHandler,
						},
						Tree: TreePath,
					},
					agenor.WithLogger(
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
	})

	When("incorrect facade", func() {
		Context("primary (expected using)", func() {
			It("🧪 should: return error", func(specCtx SpecContext) {
				lab.WithTestContext(specCtx, func(ctx context.Context, _ context.CancelFunc) {
					_, err := agenor.Walk().Configure().Extent(agenor.Prime(
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
		})

		Context("resume (expected relic)", func() {
			It("🧪 should: return error", func(specCtx SpecContext) {
				lab.WithTestContext(specCtx, func(ctx context.Context, _ context.CancelFunc) {
					_, err := agenor.Walk().Configure().Extent(agenor.Resume(
						&pref.Using{
							Subscription: agenor.SubscribeFiles,
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
})
