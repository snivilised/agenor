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
		vfs  fstest.MapFS
		root string
	)

	BeforeAll(func() {
		const (
			verbose = true
		)

		vfs, root = helpers.Musico(verbose,
			filepath.Join("MUSICO", "RETRO-WAVE"),
		)
		Expect(root).NotTo(BeEmpty())
	})

	BeforeEach(func() {
		services.Reset()
	})

	Context("glob", func() {
		DescribeTable("Filter Children (glob)",
			func(ctx SpecContext, entry *naviTE) {
				recording := make(recordingMap)
				once := func(node *tv.Node) error {
					_, found := recording[node.Extension.Name]
					Expect(found).To(BeFalse())
					recording[node.Extension.Name] = len(node.Children)

					return entry.callback(node)
				}
				path := helpers.Path(root, entry.relative)
				result, err := tv.Walk().Configure().Extent(tv.Prime(
					&tv.Using{
						Root:         path,
						Subscription: entry.subscription,
						Handler:      once,
						GetReadDirFS: func() fs.ReadDirFS {
							return vfs
						},
						GetQueryStatusFS: func(_ fs.FS) fs.StatFS {
							return vfs
						},
					},
					tv.If(entry.caseSensitive, tv.WithHookCaseSensitiveSort()),
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

				assertNavigation(entry, testOptions{
					recording: recording,
					path:      path,
					result:    result,
					err:       err,
				})
			},

			func(entry *naviTE) string {
				return fmt.Sprintf("🧪 ===> given: '%v'", entry.given)
			},

			// === folders (with files) ==========================================

			Entry(nil, &naviTE{
				given:        "folders(with files): Path is leaf",
				relative:     "RETRO-WAVE/Chromatics/Night Drive",
				subscription: enums.SubscribeFoldersWithFiles,
				callback:     foldersCallback("LEAF-PATH"),
				expectedNoOf: quantities{
					files:   0,
					folders: 1,
					children: map[string]int{
						"Night Drive": 4,
					},
				},
			}),

			Entry(nil, &naviTE{
				given:        "folders(with files): Path contains folders (check all invoked)",
				relative:     "RETRO-WAVE",
				visit:        true,
				subscription: enums.SubscribeFoldersWithFiles,
				expectedNoOf: quantities{
					files:   0,
					folders: 8,
					children: map[string]int{
						"Night Drive":      4,
						"Northern Council": 4,
						"Teenage Color":    3,
						"Innerworld":       3,
					},
				},
				callback: foldersCallback("CONTAINS-FOLDERS (check all invoked)"),
			}),
		)
	})
})
