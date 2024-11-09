package resume_test

import (
	"io/fs"

	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok

	"github.com/snivilised/li18ngo"
	nef "github.com/snivilised/nefilim"
	"github.com/snivilised/nefilim/test/luna"
	tv "github.com/snivilised/traverse"
	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/enums"
	lab "github.com/snivilised/traverse/internal/laboratory"
	"github.com/snivilised/traverse/internal/services"
	"github.com/snivilised/traverse/locale"
	"github.com/snivilised/traverse/pref"
	"github.com/snivilised/traverse/test/hydra"
)

var _ = Describe("Resume Error", Ordered, func() {
	var (
		from string
		fS   *luna.MemFS
	)

	BeforeAll(func() {
		Expect(li18ngo.Use(
			func(o *li18ngo.UseOptions) {
				o.From.Sources = li18ngo.TranslationFiles{
					locale.SourceID: li18ngo.TranslationSource{Name: "traverse"},
				}
			},
		)).To(Succeed())

		fS = hydra.Nuxx(verbose, lab.Static.RetroWave)
		from = lab.GetJSONPath()
	})

	BeforeEach(func() {
		services.Reset()
	})

	Context("given: resume path does not exist", func() {
		It("🧪 should: return error", func(ctx SpecContext) {
			from = "/invalid-path"
			_, err := tv.Walk().Configure().Extent(tv.Resume(
				&pref.Relic{
					Head: pref.Head{
						Handler: func(_ tv.Servant) error {
							return nil
						},
						GetForest: func(_ string) *core.Forest {
							return &core.Forest{
								T: fS,
								R: nef.NewTraverseABS(),
							}
						},
					},
					From:     from,
					Strategy: enums.ResumeStrategyFastward,
				},
				tv.WithOnBegin(lab.Begin("🛡️")),
				tv.WithOnEnd(lab.End("🏁")),
			)).Navigate(ctx)

			Expect(err).To(MatchError(fs.ErrNotExist))
		})
	})
})
