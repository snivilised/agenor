package kernel_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/snivilised/jaywalk/src/agenor"
	lab "github.com/snivilised/jaywalk/src/agenor/internal/laboratory"
	"github.com/snivilised/jaywalk/src/internal/services"
	"github.com/snivilised/jaywalk/locale"
	"github.com/snivilised/jaywalk/src/agenor/pref"
	"github.com/snivilised/li18ngo"
)

var _ = Describe("NavigatorFiles", Ordered, func() {
	BeforeAll(func() {
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

	Context("nav", func() {
		When("foo", func() {
			It("🧪 should: not fail", func(specCtx SpecContext) {
				lab.WithTestContext(specCtx, func(ctx context.Context, _ context.CancelFunc) {
					_, err := agenor.Walk().Configure().Extent(agenor.Prime(
						&pref.Using{
							Subscription: agenor.SubscribeFiles,
							Head: pref.Head{
								Handler: func(_ agenor.Servant) error {
									return nil
								},
							},
							Tree: RootPath,
						},
						agenor.WithFaultHandler(agenor.Accepter(lab.IgnoreFault)),
					)).Navigate(ctx)

					Expect(err).To(Succeed())
				})
			})
		})
	})
})
