package kernel_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok

	age "github.com/snivilised/agenor"
	lab "github.com/snivilised/agenor/internal/laboratory"
	"github.com/snivilised/agenor/internal/services"
	"github.com/snivilised/agenor/locale"
	"github.com/snivilised/agenor/pref"
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
					_, err := age.Walk().Configure().Extent(age.Prime(
						&pref.Using{
							Subscription: age.SubscribeFiles,
							Head: pref.Head{
								Handler: func(_ age.Servant) error {
									return nil
								},
							},
							Tree: RootPath,
						},
						age.WithFaultHandler(age.Accepter(lab.IgnoreFault)),
					)).Navigate(ctx)

					Expect(err).To(Succeed())

				})
			})
		})
	})
})
