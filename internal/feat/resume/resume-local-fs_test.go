package resume_test

import (
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	age "github.com/snivilised/agenor"
	"github.com/snivilised/agenor/core"
	"github.com/snivilised/agenor/enums"
	"github.com/snivilised/agenor/internal/enclave"
	lab "github.com/snivilised/agenor/internal/laboratory"
	"github.com/snivilised/agenor/internal/services"
	"github.com/snivilised/agenor/locale"
	"github.com/snivilised/agenor/pref"
	"github.com/snivilised/agenor/test/hanno"
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
				_, err := age.Walk().Configure(enclave.Loader(func(active *core.ActiveState) {
					active.Tree = tree
					active.CurrentPath = resumeAt
					active.Subscription = enums.SubscribeUniversal
				})).Extent(age.Resume(
					&pref.Relic{
						Head: pref.Head{
							Handler: func(servant age.Servant) error {
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
					age.WithOnBegin(lab.Begin("🛡️")),
					age.WithOnEnd(lab.End("🏁")),
				)).Navigate(ctx)

				Expect(err).To(Succeed())
			})
		})
	})
})
