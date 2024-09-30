package kernel_test

import (
	"fmt"
	"io/fs"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok

	"github.com/snivilised/li18ngo"
	tv "github.com/snivilised/traverse"
	"github.com/snivilised/traverse/enums"
	lab "github.com/snivilised/traverse/internal/laboratory"
	"github.com/snivilised/traverse/internal/services"
	"github.com/snivilised/traverse/lfs"
	"github.com/snivilised/traverse/locale"
)

var _ = Describe("NavigatorFoldersWithFiles", Ordered, func() {
	var (
		FS   *lab.TestTraverseFS
		root string
	)

	BeforeAll(func() {
		const (
			verbose = false
		)

		FS, root = lab.Musico(verbose,
			filepath.Join("MUSICO", "RETRO-WAVE"),
		)
		Expect(root).NotTo(BeEmpty())
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
				once := func(node *tv.Node) error {
					_, found := recording[node.Extension.Name]
					Expect(found).To(BeFalse())
					recording[node.Extension.Name] = len(node.Children)

					return entry.Callback(node)
				}
				path := lab.Path(root, entry.Relative)
				result, err := tv.Walk().Configure().Extent(tv.Prime(
					&tv.Using{
						Root:         path,
						Subscription: entry.Subscription,
						Handler:      once,
						GetTraverseFS: func(_ string) lfs.TraverseFS {
							return FS
						},
					},
					tv.IfOption(entry.CaseSensitive, tv.WithHookCaseSensitiveSort()),
					tv.WithHookQueryStatus(
						func(qsys fs.StatFS, path string) (fs.FileInfo, error) {
							return qsys.Stat(lab.TrimRoot(path))
						},
					),
					tv.WithHookReadDirectory(
						func(rfs fs.ReadDirFS, dirname string) ([]fs.DirEntry, error) {
							return rfs.ReadDir(lab.TrimRoot(dirname))
						},
					),
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
					Relative:     "RETRO-WAVE",
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
