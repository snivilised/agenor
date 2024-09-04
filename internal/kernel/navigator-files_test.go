package kernel_test

import (
	"context"

	"github.com/fortytw2/leaktest"
	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok

	"github.com/snivilised/li18ngo"
	tv "github.com/snivilised/traverse"
	"github.com/snivilised/traverse/internal/services"
	"github.com/snivilised/traverse/locale"
)

var _ = Describe("NavigatorFiles", Ordered, func() {
	BeforeAll(func() {
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

	Context("nav", func() {
		When("foo", func() {
			It("🧪 should: not fail", func(specCtx SpecContext) {
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
				)).Navigate(ctx)

				Expect(err).To(Succeed())
			})
		})
	})
})
