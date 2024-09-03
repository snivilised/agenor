package kernel_test

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"testing/fstest"

	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok

	tv "github.com/snivilised/traverse"
	"github.com/snivilised/traverse/enums"
	"github.com/snivilised/traverse/internal/helpers"
	"github.com/snivilised/traverse/internal/services"
)

var _ = Describe("NavigatorFoldersWithFiles", Ordered, func() {
	var (
		FS   fstest.MapFS
		root string
	)

	BeforeAll(func() {
		const (
			verbose = false
		)

		FS, root = helpers.Musico(verbose,
			filepath.Join("MUSICO", "RETRO-WAVE"),
		)
		Expect(root).NotTo(BeEmpty())
	})

	BeforeEach(func() {
		services.Reset()
	})

	Context("glob", func() {
		DescribeTable("Filter Children (glob)",
			func(ctx SpecContext, entry *helpers.NaviTE) {
				recording := make(helpers.RecordingMap)
				once := func(node *tv.Node) error {
					_, found := recording[node.Extension.Name]
					Expect(found).To(BeFalse())
					recording[node.Extension.Name] = len(node.Children)

					return entry.Callback(node)
				}
				path := helpers.Path(root, entry.Relative)
				result, err := tv.Walk().Configure().Extent(tv.Prime(
					&tv.Using{
						Root:         path,
						Subscription: entry.Subscription,
						Handler:      once,
						GetReadDirFS: func() fs.ReadDirFS {
							return FS
						},
						GetQueryStatusFS: func(_ fs.FS) fs.StatFS {
							return FS
						},
					},
					tv.IfOption(entry.CaseSensitive, tv.WithHookCaseSensitiveSort()),
					tv.WithHookQueryStatus(
						func(qsys fs.StatFS, path string) (fs.FileInfo, error) {
							return qsys.Stat(helpers.TrimRoot(path))
						},
					),
					tv.WithHookReadDirectory(
						func(rfs fs.ReadDirFS, dirname string) ([]fs.DirEntry, error) {
							return rfs.ReadDir(helpers.TrimRoot(dirname))
						},
					),
				)).Navigate(ctx)

				AssertNavigation(entry, &TestOptions{
					Recording: recording,
					Path:      path,
					Result:    result,
					Err:       err,
				})
			},

			func(entry *helpers.NaviTE) string {
				return fmt.Sprintf("🧪 ===> given: '%v'", entry.Given)
			},

			// === folders (with files) ==========================================

			Entry(nil, &helpers.NaviTE{
				Given:        "folders(with files): Path is leaf",
				Relative:     "RETRO-WAVE/Chromatics/Night Drive",
				Subscription: enums.SubscribeFoldersWithFiles,
				Callback:     FoldersCallback("LEAF-PATH"),
				ExpectedNoOf: helpers.Quantities{
					Files:   0,
					Folders: 1,
					Children: map[string]int{
						"Night Drive": 4,
					},
				},
			}),

			Entry(nil, &helpers.NaviTE{
				Given:        "folders(with files): Path contains folders (check all invoked)",
				Relative:     "RETRO-WAVE",
				Visit:        true,
				Subscription: enums.SubscribeFoldersWithFiles,
				ExpectedNoOf: helpers.Quantities{
					Files:   0,
					Folders: 8,
					Children: map[string]int{
						"Night Drive":      4,
						"Northern Council": 4,
						"Teenage Color":    3,
						"Innerworld":       3,
					},
				},
				Callback: FoldersCallback("CONTAINS-FOLDERS (check all invoked)"),
			}),
		)
	})
})
