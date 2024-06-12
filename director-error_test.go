package tv_test

import (
	"context"
	"fmt"
	"sync"

	"github.com/fortytw2/leaktest"
	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok

	tv "github.com/snivilised/traverse"
	"github.com/snivilised/traverse/internal/services"
	"github.com/snivilised/traverse/pref"
)

type traverseErrorTE struct {
	given string
	using *tv.Using
	was   *tv.Was
}

var _ = Describe("director error", Ordered, func() {
	var handler tv.Client

	BeforeAll(func() {
		handler = func(_ *tv.Node) error {
			return nil
		}
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

			if entry.was != nil {
				Expect(entry.was.Validate()).NotTo(Succeed())

				return
			}
		},
		func(entry *traverseErrorTE) string {
			return fmt.Sprintf("given: %v, 🧪 should fail", entry.given)
		},

		Entry(nil, &traverseErrorTE{
			given: "using missing root path",
			using: &tv.Using{
				Subscription: tv.SubscribeFiles,
				Handler:      handler,
			},
		}),

		Entry(nil, &traverseErrorTE{
			given: "using missing subscription",
			using: &tv.Using{
				Root:    "/root-traverse-path",
				Handler: handler,
			},
		}),

		Entry(nil, &traverseErrorTE{
			given: "using missing handler",
			using: &tv.Using{
				Root:         "/root-traverse-path",
				Subscription: tv.SubscribeFiles,
			},
		}),

		Entry(nil, &traverseErrorTE{
			given: "as missing restore from path",
			was: &tv.Was{
				Using: tv.Using{
					Root:         "/root-traverse-path",
					Subscription: tv.SubscribeFiles,
					Handler:      handler,
				},
				Strategy: tv.ResumeStrategySpawn,
			},
		}),

		Entry(nil, &traverseErrorTE{
			given: "as missing resume strategy",
			was: &tv.Was{
				Using: tv.Using{
					Root:         "/root-traverse-path",
					Subscription: tv.SubscribeFiles,
					Handler:      handler,
				},
				From: "/restore-from-path",
			},
		}),
	)

	When("Prime with subscription error", func() {
		It("🧪 should: fail", func(specCtx SpecContext) {
			defer leaktest.Check(GinkgoT())()

			ctx, cancel := context.WithCancel(specCtx)
			defer cancel()

			_, err := tv.Walk().Configure().Extent(tv.Prime(
				&tv.Using{
					Root: RootPath,
					Handler: func(_ *tv.Node) error {
						return nil
					},
				},
			)).Navigate(ctx)

			Expect(err).NotTo(Succeed())
		})
	})

	When("Prime with options build error", func() {
		It("🧪 should: fail", func(specCtx SpecContext) {
			defer leaktest.Check(GinkgoT())()

			ctx, cancel := context.WithCancel(specCtx)
			defer cancel()

			_, err := tv.Walk().Configure().Extent(tv.Prime(
				&tv.Using{
					Root:         RootPath,
					Subscription: tv.SubscribeFiles,
					Handler: func(_ *tv.Node) error {
						return nil
					},
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
				&tv.Using{
					Root: RootPath,
					Handler: func(_ *tv.Node) error {
						return nil
					},
				},
			)).Navigate(ctx)

			wg.Wait()
			Expect(err).NotTo(Succeed())
		})
	})
})
