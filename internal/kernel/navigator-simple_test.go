package kernel_test

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
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
			visited := []string{}

			once := func(node *tv.Node) error {
				_, found := recording[node.Path]
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
					GetFS: func() fs.FS {
						return vfs
					},
				},
				tv.WithOnBegin(begin("🛡️")),
				tv.If(entry.caseSensitive, tv.WithHookCaseSensitiveSort()),
				tv.WithHookQueryStatus(func(path string) (fs.FileInfo, error) {
					return vfs.Stat(helpers.TrimRoot(path))
				}),
				tv.WithHookReadDirectory(func(_ fs.FS, dirname string) ([]fs.DirEntry, error) {
					return vfs.ReadDir(helpers.TrimRoot(dirname))
				}),
			),
			).Navigate(ctx)

			_ = result.Session().StartedAt()
			_ = result.Session().Elapsed()

			if entry.visit {
				_ = fs.WalkDir(vfs, path, func(path string, de fs.DirEntry, _ error) error {
					if strings.HasSuffix(path, ".DS_Store") {
						return nil
					}

					if subscribes(entry.subscription, de) {
						visited = append(visited, path)
					}
					return nil
				})
			}

			if entry.visit {
				every := lo.EveryBy(visited, func(p string) bool {
					_, found := recording[p]
					return found
				})
				Expect(every).To(BeTrue())
			}

			Expect(err).To(Succeed())
			Expect(result.Metrics().Count(enums.MetricNoFilesInvoked)).To(
				Equal(entry.expectedNoOf.files), "Incorrect no of files",
			)
			Expect(result.Metrics().Count(enums.MetricNoFoldersInvoked)).To(
				Equal(entry.expectedNoOf.folders), "Incorrect no of folders",
			)
		},
		func(entry *naviTE) string {
			return fmt.Sprintf("🧪 ===> given: '%v'", entry.message)
		},

		// === universal =====================================================

		Entry(nil, Label("RETRO-WAVE"), &naviTE{
			message:      "universal: Path is leaf",
			relative:     "RETRO-WAVE/Chromatics/Night Drive",
			subscription: enums.SubscribeUniversal,
			callback:     universalCallback("LEAF-PATH"),
			expectedNoOf: quantities{
				files:   4,
				folders: 1,
			},
		}),

		Entry(nil, Label("RETRO-WAVE"), &naviTE{
			message:      "universal: Path contains folders",
			relative:     "RETRO-WAVE",
			subscription: enums.SubscribeUniversal,
			callback:     universalCallback("CONTAINS-FOLDERS"),
			expectedNoOf: quantities{
				files:   14,
				folders: 8,
			},
		}),

		Entry(nil, Label("RETRO-WAVE"), &naviTE{
			message:      "universal: Path contains folders (visit)",
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
			message:      "folders: Path is leaf",
			relative:     "RETRO-WAVE/Chromatics/Night Drive",
			subscription: enums.SubscribeFolders,
			callback:     foldersCallback("LEAF-PATH"),
			expectedNoOf: quantities{
				folders: 1,
			},
		}),

		Entry(nil, Label("RETRO-WAVE"), &naviTE{
			message:      "folders: Path contains folders",
			relative:     "RETRO-WAVE",
			subscription: enums.SubscribeFolders,
			callback:     foldersCallback("CONTAINS-FOLDERS"),
			expectedNoOf: quantities{
				folders: 8,
			},
		}),

		Entry(nil, Label("RETRO-WAVE"), &naviTE{
			message:      "folders: Path contains folders (check all invoked)",
			relative:     "RETRO-WAVE",
			visit:        true,
			subscription: enums.SubscribeFolders,
			callback:     foldersCallback("CONTAINS-FOLDERS (check all invoked)"),
			expectedNoOf: quantities{
				folders: 8,
			},
		}),

		Entry(nil, Label("metal"), &naviTE{
			message:       "folders: case sensitive sort",
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
			message:      "files: Path is leaf",
			relative:     "RETRO-WAVE/Chromatics/Night Drive",
			subscription: enums.SubscribeFiles,
			callback:     filesCallback("LEAF-PATH"),
			expectedNoOf: quantities{
				files:   4,
				folders: 0,
			},
		}),

		Entry(nil, Label("RETRO-WAVE"), &naviTE{
			message:      "files: Path contains folders",
			relative:     "RETRO-WAVE",
			subscription: enums.SubscribeFiles,
			callback:     filesCallback("CONTAINS-FOLDERS"),
			expectedNoOf: quantities{
				files:   14,
				folders: 0,
			},
		}),

		Entry(nil, Label("RETRO-WAVE"), &naviTE{
			message:      "files: Path contains folders",
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
