package kernel_test

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok

	"github.com/snivilised/li18ngo"
	"github.com/snivilised/nefilim/luna"
	tv "github.com/snivilised/traverse"
	"github.com/snivilised/traverse/enums"
	"github.com/snivilised/traverse/hydra"
	lab "github.com/snivilised/traverse/internal/laboratory"
	"github.com/snivilised/traverse/internal/services"
	"github.com/snivilised/traverse/locale"
)

var _ = Describe("NavigatorFoldersWithFiles", Ordered, func() {
	var (
		fS *luna.MemFS
	)

	BeforeAll(func() {
		const (
			verbose = false
		)

		fS = hydra.Nuxx(verbose, lab.Static.RetroWave)
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

	Context("glob", func() {
		DescribeTable("Filter Children (glob)",
			func(ctx SpecContext, entry *lab.FilterTE) {
				recording := make(lab.RecordingMap)
				once := func(servant tv.Servant) error {
					node := servant.Node()
					_, found := recording[node.Extension.Name]
					Expect(found).To(BeFalse())
					recording[node.Extension.Name] = len(node.Children)

					return entry.Callback(servant)
				}
				path := entry.Relative
				result, err := tv.Walk().Configure().Extent(tv.Prime(
					&tv.Using{
						Tree:         path,
						Subscription: entry.Subscription,
						Handler:      once,
						GetTraverseFS: func(_ string) tv.TraverseFS {
							return fS
						},
					},
					tv.WithOnBegin(lab.Begin("🛡️")),
					tv.WithOnEnd(lab.End("🏁")),

					tv.IfOption(entry.CaseSensitive, tv.WithHookCaseSensitiveSort()),
				)).Navigate(ctx)

				lab.AssertNavigation(&entry.NaviTE, &lab.TestOptions{
					Recording: recording,
					Path:      path,
					Result:    result,
					Err:       err,
				})
			},

			func(entry *lab.FilterTE) string {
				return fmt.Sprintf("🧪 ===> given: '%v'", entry.Given)
			},

			// === folders (with files) ==========================================

			Entry(nil, &lab.FilterTE{
				NaviTE: lab.NaviTE{
					Given:        "folders(with files): Path is leaf",
					Relative:     "RETRO-WAVE/Chromatics/Night Drive",
					Subscription: enums.SubscribeFoldersWithFiles,
					Callback:     lab.FoldersCallback("LEAF-PATH"),
					ExpectedNoOf: lab.Quantities{
						Files:   0,
						Folders: 1,
						Children: map[string]int{
							"Night Drive": 4,
						},
					},
				},
			}),

			Entry(nil, &lab.FilterTE{
				NaviTE: lab.NaviTE{
					Given:        "folders(with files): Path contains folders (check all invoked)",
					Relative:     lab.Static.RetroWave,
					Visit:        true,
					Subscription: enums.SubscribeFoldersWithFiles,
					ExpectedNoOf: lab.Quantities{
						Files:   0,
						Folders: 8,
						Children: map[string]int{
							"Night Drive":      4,
							"Northern Council": 4,
							"Teenage Color":    3,
							"Innerworld":       3,
						},
					},
					Callback: lab.FoldersCallback("CONTAINS-FOLDERS (check all invoked)"),
				},
			}),
		)
	})
})
