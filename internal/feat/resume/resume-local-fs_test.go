package resume_test

import (
	"path/filepath"

	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok

	"github.com/snivilised/li18ngo"
	nef "github.com/snivilised/nefilim"
	tv "github.com/snivilised/traverse"
	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/enums"
	lab "github.com/snivilised/traverse/internal/laboratory"
	"github.com/snivilised/traverse/internal/services"
	"github.com/snivilised/traverse/locale"
	"github.com/snivilised/traverse/pref"
	"github.com/snivilised/traverse/test/hydra"
)

var _ = Describe("Resume local-fs", Ordered, func() {
	var (
		from, tree, resumeAt string
		strategy             enums.ResumeStrategy
	)

	BeforeAll(func() {
		Expect(li18ngo.Use(
			func(o *li18ngo.UseOptions) {
				o.From.Sources = li18ngo.TranslationFiles{
					locale.SourceID: li18ngo.TranslationSource{Name: "traverse"},
				}
			},
		)).To(Succeed())

		// For these tests, the navigation tree is the 'test' directory
		// and the resume point is the 'data' directory.
		//
		from = lab.GetJSONPath()
		tree = hydra.Repo("test")
		resumeAt = filepath.Join(tree, "data")
	})

	BeforeEach(func() {
		services.Reset()
	})

	Context("fs:absolute", func() {
		Context("given: resume path exists", func() {
			It("🧪 should: resume traverse ok", func(ctx SpecContext) {
				strategy = enums.ResumeStrategyFastward
				_, err := tv.Walk().Configure().Extent(tv.Resume(
					&pref.Relic{
						Head: pref.Head{
							Subscription: enums.SubscribeUniversal,
							Handler: func(servant tv.Servant) error {
								node := servant.Node()
								depth := node.Extension.Depth
								GinkgoWriter.Printf(
									"---> (resume-abs-local-fs) 🐷 %v: (depth:%v) '%v'\n",
									strategy, depth, node.Path,
								)

								return nil
							},

							// Create an absolute fs because the default is relative.
							// Actually, the type of file system we use has to be inline
							// with the file system type that was used in the corresponding
							// primary run that we are resuming from, but how to enforce?
							// (see issue #301)
							//
							GetForest: func(_ string) *core.Forest {
								return &core.Forest{
									T: nef.NewTraverseABS(),
									R: nef.NewTraverseABS(),
								}
							},
						},
						From:     from,
						Strategy: enums.ResumeStrategyFastward,
						Restorer: func(_ *pref.Options, active *core.ActiveState) error {
							active.Tree = tree
							active.CurrentPath = resumeAt

							return nil
						},
					},
					tv.WithOnBegin(lab.Begin("🛡️")),
					tv.WithOnEnd(lab.End("🏁")),
				)).Navigate(ctx)

				Expect(err).To(Succeed())
			})
		})
	})
})