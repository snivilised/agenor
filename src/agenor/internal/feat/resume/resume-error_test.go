package resume_test

import (
	"context"
	"fmt"
	"io/fs"

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
	"github.com/snivilised/jaywalk/src/agenor/tfs"
	"github.com/snivilised/li18ngo"
	"github.com/snivilised/nefilim/test/luna"
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
					locale.SourceID: li18ngo.TranslationSource{Name: "agenor"},
				}
			},
		)).To(Succeed())
	})

	BeforeEach(func() {
		services.Reset()

		fS = hanno.Nuxx(verbose, lab.Static.RetroWave)
		from = lab.GetJSONPath()
	})

	Context("given: resume path does not exist", func() {
		It("🧪 should: return error", func(specCtx SpecContext) {
			lab.WithTestContext(specCtx, func(ctx context.Context, _ context.CancelFunc) {
				from = "/invalid-path"
				_, err := agenor.Walk().Configure().Extent(agenor.Resume(
					&pref.Relic{
						Head: pref.Head{
							Handler: func(_ agenor.Servant) error {
								return nil
							},
							GetForest: func(_ string) *core.Forest {
								return &core.Forest{
									T: fS,
									R: tfs.New(),
								}
							},
						},
						From:     from,
						Strategy: enums.ResumeStrategyFastward,
					},
					agenor.WithOnBegin(lab.Begin("🛡️")),
					agenor.WithOnEnd(lab.End("🏁")),
				)).Navigate(ctx)

				Expect(err).To(MatchError(fs.ErrNotExist))
			})
		})
	})

	Context("forest inception failure", func() {
		DescribeTable("fs type mismatch",
			func(specCtx SpecContext, _ string, travIsRelative, resIsRelative bool) {
				lab.WithTestContext(specCtx, func(ctx context.Context, _ context.CancelFunc) {
					_, err := agenor.Walk().Configure(enclave.Loader(func(active *core.ActiveState) {
						active.Tree = lab.Static.RetroWave
						active.CurrentPath = ResumeAtTeenageColor
						active.Subscription = enums.SubscribeUniversal
						active.TraverseDescription.IsRelative = travIsRelative
						active.ResumeDescription.IsRelative = resIsRelative
					})).Extent(agenor.Resume(
						&pref.Relic{
							Head: pref.Head{
								Handler: func(_ agenor.Servant) error {
									return nil
								},
								GetForest: func(_ string) *core.Forest {
									return &core.Forest{
										T: fS,
										R: tfs.New(),
									}
								},
							},
							From:     from,
							Strategy: enums.ResumeStrategyFastward,
						},
						agenor.WithOnBegin(lab.Begin("🛡️")),
						agenor.WithOnEnd(lab.End("🏁")),
					)).Navigate(ctx)

					Expect(err).To(MatchError(locale.ErrCoreResumeFsMismatch))
				})
			},
			func(given string, _, _ bool) string {
				return fmt.Sprintf("🧪 ===> given: '%v'", given)
			},
			Entry(nil, "traverse-fs does not match", false, false),
			Entry(nil, "resume-fs does not match", false, true),
		)
	})

	When("custom forest creator returns nil", func() {
		It("should: fail", func(specCtx SpecContext) {
			lab.WithTestContext(specCtx, func(ctx context.Context, _ context.CancelFunc) {
				_, err := agenor.Walk().Configure(enclave.Loader(func(active *core.ActiveState) {
					active.Tree = lab.Static.RetroWave
					active.CurrentPath = ResumeAtTeenageColor
					active.Subscription = enums.SubscribeUniversal
					active.TraverseDescription.IsRelative = true
					active.ResumeDescription.IsRelative = false
				})).Extent(agenor.Resume(
					&pref.Relic{
						Head: pref.Head{
							Handler: func(_ agenor.Servant) error {
								return nil
							},
							GetForest: func(_ string) *core.Forest {
								return nil
							},
						},
						From:     from,
						Strategy: enums.ResumeStrategyFastward,
					},
					agenor.WithOnBegin(lab.Begin("🛡️")),
					agenor.WithOnEnd(lab.End("🏁")),
				)).Navigate(ctx)

				Expect(err).NotTo(Succeed())
				Expect(err).To(MatchError(core.ErrNilForest))
			})
		})
	})
})
