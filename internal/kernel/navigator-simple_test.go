package kernel_test

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"testing/fstest"

	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok

	"github.com/snivilised/li18ngo"
	tv "github.com/snivilised/traverse"
	"github.com/snivilised/traverse/enums"
	"github.com/snivilised/traverse/internal/helpers"
	"github.com/snivilised/traverse/internal/services"
	"github.com/snivilised/traverse/internal/third/lo"
	"github.com/snivilised/traverse/locale"
)

var _ = Describe("NavigatorUniversal", Ordered, func() {
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
			filepath.Join("MUSICO", "rock", "metal"),
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

	DescribeTable("Ensure Callback Invoked Once", Label("simple"),
		func(ctx SpecContext, entry *helpers.NaviTE) {
			recording := make(helpers.RecordingMap)
			once := func(node *tv.Node) error {
				_, found := recording[node.Path] // TODO: should this be name not path?
				Expect(found).To(BeFalse())
				recording[node.Path] = len(node.Children)

				return entry.Callback(node)
			}

			visitor := func(node *tv.Node) error {
				return once(node)
			}

			callback := lo.Ternary(entry.Once, once,
				lo.Ternary(entry.Visit, visitor, entry.Callback),
			)
			path := helpers.Path(root, entry.Relative)

			result, err := tv.Walk().Configure().Extent(tv.Prime(
				&tv.Using{
					Root:         path,
					Subscription: entry.Subscription,
					Handler:      callback,
					GetReadDirFS: func() fs.ReadDirFS {
						return FS
					},
					GetQueryStatusFS: func(_ fs.FS) fs.StatFS {
						return FS
					},
				},
				tv.WithOnBegin(helpers.Begin("🛡️")),
				tv.WithOnEnd(helpers.End("🏁")),
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

			helpers.AssertNavigation(entry, &helpers.TestOptions{
				FS:        FS,
				Recording: recording,
				Path:      path,
				Result:    result,
				Err:       err,
				Every: func(p string) bool {
					_, found := recording[p]
					return found
				},
			})
		},
		func(entry *helpers.NaviTE) string {
			return fmt.Sprintf("🧪 ===> given: '%v'", entry.Given)
		},

		// === universal =====================================================

		Entry(nil, Label("RETRO-WAVE"), &helpers.NaviTE{
			Given:        "universal: Path is leaf",
			Relative:     "RETRO-WAVE/Chromatics/Night Drive",
			Subscription: enums.SubscribeUniversal,
			Callback:     helpers.UniversalCallback("LEAF-PATH"),
			ExpectedNoOf: helpers.Quantities{
				Files:   4,
				Folders: 1,
			},
		}),

		Entry(nil, Label("RETRO-WAVE"), &helpers.NaviTE{
			Given:        "universal: Path contains folders",
			Relative:     "RETRO-WAVE",
			Subscription: enums.SubscribeUniversal,
			Callback:     helpers.UniversalCallback("CONTAINS-FOLDERS"),
			ExpectedNoOf: helpers.Quantities{
				Files:   14,
				Folders: 8,
			},
		}),

		Entry(nil, Label("RETRO-WAVE"), &helpers.NaviTE{
			Given:        "universal: Path contains folders (visit)",
			Relative:     "RETRO-WAVE",
			Visit:        true,
			Subscription: enums.SubscribeUniversal,
			Callback:     helpers.UniversalCallback("VISIT-CONTAINS-FOLDERS"),
			ExpectedNoOf: helpers.Quantities{
				Files:   14,
				Folders: 8,
			},
		}),

		// === folders =======================================================

		Entry(nil, Label("RETRO-WAVE"), &helpers.NaviTE{
			Given:        "folders: Path is leaf",
			Relative:     "RETRO-WAVE/Chromatics/Night Drive",
			Subscription: enums.SubscribeFolders,
			Callback:     helpers.FoldersCallback("LEAF-PATH"),
			ExpectedNoOf: helpers.Quantities{
				Folders: 1,
			},
		}),

		Entry(nil, Label("RETRO-WAVE"), &helpers.NaviTE{
			Given:        "folders: Path contains folders",
			Relative:     "RETRO-WAVE",
			Subscription: enums.SubscribeFolders,
			Callback:     helpers.FoldersCallback("CONTAINS-FOLDERS"),
			ExpectedNoOf: helpers.Quantities{
				Folders: 8,
			},
		}),

		Entry(nil, Label("RETRO-WAVE"), &helpers.NaviTE{
			Given:        "folders: Path contains folders (check all invoked)",
			Relative:     "RETRO-WAVE",
			Visit:        true,
			Subscription: enums.SubscribeFolders,
			Callback:     helpers.FoldersCallback("CONTAINS-FOLDERS (check all invoked)"),
			ExpectedNoOf: helpers.Quantities{
				Folders: 8,
			},
		}),

		Entry(nil, Label("metal"), &helpers.NaviTE{
			Given:         "folders: case sensitive sort",
			Relative:      "rock/metal",
			Subscription:  enums.SubscribeFolders,
			CaseSensitive: true,
			Callback: helpers.FoldersCaseSensitiveCallback(
				"rock/metal/HARD-METAL", "rock/metal/dark",
			),
			ExpectedNoOf: helpers.Quantities{
				Files:   0,
				Folders: 41,
			},
		}),

		// === files =========================================================

		Entry(nil, Label("RETRO-WAVE"), &helpers.NaviTE{
			Given:        "files: Path is leaf",
			Relative:     "RETRO-WAVE/Chromatics/Night Drive",
			Subscription: enums.SubscribeFiles,
			Callback:     helpers.FilesCallback("LEAF-PATH"),
			ExpectedNoOf: helpers.Quantities{
				Files:   4,
				Folders: 0,
			},
		}),

		Entry(nil, Label("RETRO-WAVE"), &helpers.NaviTE{
			Given:        "files: Path contains folders",
			Relative:     "RETRO-WAVE",
			Subscription: enums.SubscribeFiles,
			Callback:     helpers.FilesCallback("CONTAINS-FOLDERS"),
			ExpectedNoOf: helpers.Quantities{
				Files:   14,
				Folders: 0,
			},
		}),

		Entry(nil, Label("RETRO-WAVE"), &helpers.NaviTE{
			Given:        "files: Path contains folders",
			Relative:     "RETRO-WAVE",
			Visit:        true,
			Subscription: enums.SubscribeFiles,
			Callback:     helpers.FilesCallback("VISIT-CONTAINS-FOLDERS"),
			ExpectedNoOf: helpers.Quantities{
				Files:   14,
				Folders: 0,
			},
		}),
	)
})
