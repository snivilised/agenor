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
	"github.com/snivilised/traverse/pref"
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
		var portion = filepath.Join("MUSICO", "RETRO-WAVE")
		vfs, root = helpers.Musico(portion, verbose)
		Expect(root).NotTo(BeEmpty())
	})

	BeforeEach(func() {
		services.Reset()
	})

	DescribeTable("Ensure Callback Invoked Once", Label("simple", "RETRO-WAVE"),
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
				tv.WithSortBehaviour(&pref.SortBehaviour{
					IsCaseSensitive: entry.caseSensitive,
				}),
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
				_ = filepath.WalkDir(path, func(path string, de fs.DirEntry, _ error) error {
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

		Entry(nil, &naviTE{
			message:      "universal: Path is leaf",
			relative:     "RETRO-WAVE/Chromatics/Night Drive",
			subscription: enums.SubscribeUniversal,
			callback:     universalCallback("LEAF-PATH"),
			expectedNoOf: directoryQuantities{
				files:   4,
				folders: 1,
			},
		}),

		Entry(nil, &naviTE{
			message:      "universal: Path contains folders",
			relative:     "RETRO-WAVE",
			subscription: enums.SubscribeUniversal,
			callback:     universalCallback("CONTAINS-FOLDERS"),
			expectedNoOf: directoryQuantities{
				files:   14,
				folders: 8,
			},
		}),
		Entry(nil, &naviTE{
			message:      "universal: Path contains folders (visit)",
			relative:     "RETRO-WAVE",
			visit:        true,
			subscription: enums.SubscribeUniversal,
			callback:     universalCallback("VISIT-CONTAINS-FOLDERS"),
			expectedNoOf: directoryQuantities{
				files:   14,
				folders: 8,
			},
		}),
	)
})
