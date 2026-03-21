package resume_test

import (
	"path/filepath"

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

var _ = Describe("Resume local-fs", Ordered, func() {
	var (
		from, tree, resumeAt string
		strategy             enums.ResumeStrategy
	)

	BeforeAll(func() {
		Expect(li18ngo.Use(
			func(o *li18ngo.UseOptions) {
				o.From.Sources = li18ngo.TranslationFiles{
					locale.SourceID: li18ngo.TranslationSource{Name: "agenor"},
				}
			},
		)).To(Succeed())

		// For these tests, the navigation tree is the 'test' directory
		// and the resume point is the 'data' directory.
		//
		from = lab.GetJSONPath()
		tree = hanno.Repo("test")
		resumeAt = filepath.Join(tree, "data")
	})

	BeforeEach(func() {
		services.Reset()
	})

	Context("fs:absolute", func() {
		Context("given: resume path exists", func() {
			It("🧪 should: resume traverse ok", func(ctx SpecContext) {
				strategy = enums.ResumeStrategyFastward
				_, err := agenor.Walk().Configure(enclave.Loader(func(active *core.ActiveState) {
					active.Tree = tree
					active.CurrentPath = resumeAt
					active.Subscription = enums.SubscribeUniversal
				})).Extent(agenor.Resume(
					&pref.Relic{
						Head: pref.Head{
							Handler: func(servant agenor.Servant) error {
								node := servant.Node()
								depth := node.Extension.Depth
								GinkgoWriter.Printf(
									"---> (resume-abs-local-fs) 🐷 %v: (depth:%v) '%v'\n",
									strategy, depth, node.Path,
								)

								return nil
							},
						},
						From:     from,
						Strategy: enums.ResumeStrategyFastward,
					},
					agenor.WithOnBegin(lab.Begin("🛡️")),
					agenor.WithOnEnd(lab.End("🏁")),
				)).Navigate(ctx)

				Expect(err).To(Succeed())
			})
		})
	})
})
