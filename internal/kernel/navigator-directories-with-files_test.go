package kernel_test

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok

	age "github.com/snivilised/agenor"
	"github.com/snivilised/agenor/core"
	"github.com/snivilised/agenor/enums"
	lab "github.com/snivilised/agenor/internal/laboratory"
	"github.com/snivilised/agenor/internal/services"
	"github.com/snivilised/agenor/locale"
	"github.com/snivilised/agenor/pref"
	"github.com/snivilised/agenor/test/hanno"
	"github.com/snivilised/agenor/tfs"
	"github.com/snivilised/li18ngo"
	"github.com/snivilised/nefilim/test/luna"
)

var _ = Describe("NavigatorDirectoriesWithFiles", Ordered, func() {
	var (
		fS *luna.MemFS
	)

	BeforeAll(func() {
		const (
			verbose = false
		)

		fS = hanno.Nuxx(verbose, lab.Static.RetroWave)
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

	Context("glob", func() {
		DescribeTable("Filter Children (glob)",
			func(ctx SpecContext, entry *lab.FilterTE) {
				recall := make(lab.Recall)
				once := func(servant age.Servant) error {
					node := servant.Node()
					_, found := recall[node.Extension.Name]
					Expect(found).To(BeFalse())
					recall[node.Extension.Name] = len(node.Children)

					return entry.Callback(servant)
				}
				path := entry.Relative
				result, err := age.Walk().Configure().Extent(age.Prime(
					&pref.Using{
						Subscription: entry.Subscription,
						Head: pref.Head{
							Handler: once,
							GetForest: func(_ string) *core.Forest {
								return &core.Forest{
									T: fS,
									R: tfs.New(),
								}
							},
						},
						Tree: path,
					},
					age.WithOnBegin(lab.Begin("🛡️")),
					age.WithOnEnd(lab.End("🏁")),

					age.IfOption(entry.CaseSensitive, age.WithHookCaseSensitiveSort()),
				)).Navigate(ctx)

				lab.AssertNavigation(&entry.NaviTE, &lab.TestOptions{
					Recording: recall,
					Path:      path,
					Result:    result,
					Err:       err,
				})
			},

			func(entry *lab.FilterTE) string {
				return fmt.Sprintf("🧪 ===> given: '%v'", entry.Given)
			},

			// === directories (with files) ======================================

			Entry(nil, &lab.FilterTE{
				NaviTE: lab.NaviTE{
					Given:        "directories(with files): Path is leaf",
					Relative:     "RETRO-WAVE/Chromatics/Night Drive",
					Subscription: enums.SubscribeDirectoriesWithFiles,
					Callback:     lab.DirectoriesCallback("LEAF-PATH"),
					ExpectedNoOf: lab.Quantities{
						Files:       0,
						Directories: 1,
						Children: map[string]int{
							"Night Drive": 4,
						},
					},
				},
			}),

			Entry(nil, &lab.FilterTE{
				NaviTE: lab.NaviTE{
					Given:        "directories(with files): Path contains directories (check all invoked)",
					Relative:     lab.Static.RetroWave,
					Visit:        true,
					Subscription: enums.SubscribeDirectoriesWithFiles,
					ExpectedNoOf: lab.Quantities{
						Files:       0,
						Directories: 8,
						Children: map[string]int{
							"Night Drive":      4,
							"Northern Council": 4,
							"Teenage Color":    3,
							"Innerworld":       3,
						},
					},
					Callback: lab.DirectoriesCallback("CONTAINS-DIRECTORIES (check all invoked)"),
				},
			}),
		)
	})
})
