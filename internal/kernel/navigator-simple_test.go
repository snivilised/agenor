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
	"github.com/snivilised/traverse/internal/lo"
	"github.com/snivilised/traverse/internal/services"
)

var _ = Describe("NavigatorUniversal", Ordered, func() {
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
			filepath.Join("MUSICO", "rock", "metal"),
		)
		Expect(root).NotTo(BeEmpty())
	})

	BeforeEach(func() {
		services.Reset()
	})

	DescribeTable("Ensure Callback Invoked Once", Label("simple"),
		func(ctx SpecContext, entry *naviTE) {
			recording := make(recordingMap)
			once := func(node *tv.Node) error {
				_, found := recording[node.Path] // TODO: should this be name not path?
				Expect(found).To(BeFalse())
				recording[node.Path] = len(node.Children)

				return entry.callback(node)
			}

			visitor := func(node *tv.Node) error {
				return once(node)
			}

			callback := lo.Ternary(entry.once, once,
				lo.Ternary(entry.visit, visitor, entry.callback),
			)
			path := helpers.Path(root, entry.relative)

			result, err := tv.Walk().Configure().Extent(tv.Prime(
				&tv.Using{
					Root:         path,
					Subscription: entry.subscription,
					Handler:      callback,
					GetReadDirFS: func() fs.ReadDirFS {
						return vfs
					},
					GetQueryStatusFS: func(_ fs.FS) fs.StatFS {
						return vfs
					},
				},
				tv.WithOnBegin(begin("🛡️")),
				tv.WithOnEnd(end("🏁")),
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
				vfs:       vfs,
				recording: recording,
				path:      path,
				result:    result,
				err:       err,
				every: func(p string) bool {
					_, found := recording[p]
					return found
				},
			})
		},
		func(entry *naviTE) string {
			return fmt.Sprintf("🧪 ===> given: '%v'", entry.given)
		},

		// === universal =====================================================

		Entry(nil, Label("RETRO-WAVE"), &naviTE{
			given:        "universal: Path is leaf",
			relative:     "RETRO-WAVE/Chromatics/Night Drive",
			subscription: enums.SubscribeUniversal,
			callback:     universalCallback("LEAF-PATH"),
			expectedNoOf: quantities{
				files:   4,
				folders: 1,
			},
		}),

		Entry(nil, Label("RETRO-WAVE"), &naviTE{
			given:        "universal: Path contains folders",
			relative:     "RETRO-WAVE",
			subscription: enums.SubscribeUniversal,
			callback:     universalCallback("CONTAINS-FOLDERS"),
			expectedNoOf: quantities{
				files:   14,
				folders: 8,
			},
		}),

		Entry(nil, Label("RETRO-WAVE"), &naviTE{
			given:        "universal: Path contains folders (visit)",
			relative:     "RETRO-WAVE",
			visit:        true,
			subscription: enums.SubscribeUniversal,
			callback:     universalCallback("VISIT-CONTAINS-FOLDERS"),
			expectedNoOf: quantities{
				files:   14,
				folders: 8,
			},
		}),

		// === folders =======================================================

		Entry(nil, Label("RETRO-WAVE"), &naviTE{
			given:        "folders: Path is leaf",
			relative:     "RETRO-WAVE/Chromatics/Night Drive",
			subscription: enums.SubscribeFolders,
			callback:     foldersCallback("LEAF-PATH"),
			expectedNoOf: quantities{
				folders: 1,
			},
		}),

		Entry(nil, Label("RETRO-WAVE"), &naviTE{
			given:        "folders: Path contains folders",
			relative:     "RETRO-WAVE",
			subscription: enums.SubscribeFolders,
			callback:     foldersCallback("CONTAINS-FOLDERS"),
			expectedNoOf: quantities{
				folders: 8,
			},
		}),

		Entry(nil, Label("RETRO-WAVE"), &naviTE{
			given:        "folders: Path contains folders (check all invoked)",
			relative:     "RETRO-WAVE",
			visit:        true,
			subscription: enums.SubscribeFolders,
			callback:     foldersCallback("CONTAINS-FOLDERS (check all invoked)"),
			expectedNoOf: quantities{
				folders: 8,
			},
		}),

		Entry(nil, Label("metal"), &naviTE{
			given:         "folders: case sensitive sort",
			relative:      "rock/metal",
			subscription:  enums.SubscribeFolders,
			caseSensitive: true,
			callback: foldersCaseSensitiveCallback(
				"rock/metal/HARD-METAL", "rock/metal/dark",
			),
			expectedNoOf: quantities{
				files:   0,
				folders: 41,
			},
		}),

		// === files =========================================================

		Entry(nil, Label("RETRO-WAVE"), &naviTE{
			given:        "files: Path is leaf",
			relative:     "RETRO-WAVE/Chromatics/Night Drive",
			subscription: enums.SubscribeFiles,
			callback:     filesCallback("LEAF-PATH"),
			expectedNoOf: quantities{
				files:   4,
				folders: 0,
			},
		}),

		Entry(nil, Label("RETRO-WAVE"), &naviTE{
			given:        "files: Path contains folders",
			relative:     "RETRO-WAVE",
			subscription: enums.SubscribeFiles,
			callback:     filesCallback("CONTAINS-FOLDERS"),
			expectedNoOf: quantities{
				files:   14,
				folders: 0,
			},
		}),

		Entry(nil, Label("RETRO-WAVE"), &naviTE{
			given:        "files: Path contains folders",
			relative:     "RETRO-WAVE",
			visit:        true,
			subscription: enums.SubscribeFiles,
			callback:     filesCallback("VISIT-CONTAINS-FOLDERS"),
			expectedNoOf: quantities{
				files:   14,
				folders: 0,
			},
		}),
	)
})
