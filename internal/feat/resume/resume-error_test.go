package resume_test

import (
	"context"
	"fmt"
	"io/fs"

	"github.com/fortytw2/leaktest"
	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok

	"github.com/snivilised/li18ngo"
	"github.com/snivilised/nefilim/test/luna"
	tv "github.com/snivilised/traverse"
	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/enums"
	"github.com/snivilised/traverse/internal/enclave"
	lab "github.com/snivilised/traverse/internal/laboratory"
	"github.com/snivilised/traverse/internal/services"
	"github.com/snivilised/traverse/locale"
	"github.com/snivilised/traverse/pref"
	"github.com/snivilised/traverse/test/hydra"
	"github.com/snivilised/traverse/tfs"
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
	})

	BeforeEach(func() {
		services.Reset()
		fS = hydra.Nuxx(verbose, lab.Static.RetroWave)
		from = lab.GetJSONPath()
	})

	Context("given: resume path does not exist", func() {
		It("🧪 should: return error", func(specCtx SpecContext) {
			defer leaktest.Check(GinkgoT())()

			ctx, cancel := context.WithCancel(specCtx)
			defer cancel()

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
								R: tfs.New(),
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

	Context("forest creation failure", func() {
		DescribeTable("fs type mismatch",
			func(specCtx SpecContext, _ string, travIsRelative, resIsRelative bool) {
				defer leaktest.Check(GinkgoT())()

				ctx, cancel := context.WithCancel(specCtx)
				defer cancel()

				_, err := tv.Walk().Configure(enclave.Loader(func(active *core.ActiveState) {
					active.Tree = lab.Static.RetroWave
					active.CurrentPath = ResumeAtTeenageColor
					active.Subscription = enums.SubscribeUniversal
					active.TraverseDescription.IsRelative = travIsRelative
					active.ResumeDescription.IsRelative = resIsRelative
				})).Extent(tv.Resume(
					&pref.Relic{
						Head: pref.Head{
							Handler: func(_ tv.Servant) error {
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
					tv.WithOnBegin(lab.Begin("🛡️")),
					tv.WithOnEnd(lab.End("🏁")),
				)).Navigate(ctx)

				Expect(err).To(MatchError(locale.ErrCoreResumeFsMismatch))
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
			defer leaktest.Check(GinkgoT())()

			ctx, cancel := context.WithCancel(specCtx)
			defer cancel()

			_, err := tv.Walk().Configure(enclave.Loader(func(active *core.ActiveState) {
				active.Tree = lab.Static.RetroWave
				active.CurrentPath = ResumeAtTeenageColor
				active.Subscription = enums.SubscribeUniversal
				active.TraverseDescription.IsRelative = true
				active.ResumeDescription.IsRelative = false
			})).Extent(tv.Resume(
				&pref.Relic{
					Head: pref.Head{
						Handler: func(_ tv.Servant) error {
							return nil
						},
						GetForest: func(_ string) *core.Forest {
							return nil
						},
					},
					From:     from,
					Strategy: enums.ResumeStrategyFastward,
				},
				tv.WithOnBegin(lab.Begin("🛡️")),
				tv.WithOnEnd(lab.End("🏁")),
			)).Navigate(ctx)

			Expect(err).NotTo(Succeed())
			Expect(err).To(MatchError(core.ErrNilForest))
		})
	})
})
