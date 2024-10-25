package kernel_test

import (
	"fmt"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok

	"github.com/snivilised/li18ngo"
	"github.com/snivilised/nefilim/luna"
	tv "github.com/snivilised/traverse"
	"github.com/snivilised/traverse/enums"
	lab "github.com/snivilised/traverse/internal/laboratory"
	"github.com/snivilised/traverse/internal/services"
	"github.com/snivilised/traverse/internal/third/lo"
	"github.com/snivilised/traverse/locale"
)

var _ = Describe("NavigatorUniversal", Ordered, func() {
	var (
		FS   *luna.MemFS
		root string
	)

	BeforeAll(func() {
		const (
			verbose = false
		)

		FS, root = lab.Musico(verbose,
			lab.Static.RetroWave,
			filepath.Join("rock", "metal"),
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
		func(ctx SpecContext, entry *lab.NaviTE) {
			recording := make(lab.RecordingMap)
			once := func(servant tv.Servant) error {
				node := servant.Node()
				_, found := recording[node.Path] // TODO: should this be name not path?
				Expect(found).To(BeFalse())
				recording[node.Path] = len(node.Children)

				return entry.Callback(servant)
			}

			visitor := func(servant tv.Servant) error {
				return once(servant)
			}

			callback := lo.Ternary(entry.Once, once,
				lo.Ternary(entry.Visit, visitor, entry.Callback),
			)
			path := entry.Relative

			result, err := tv.Walk().Configure().Extent(tv.Prime(
				&tv.Using{
					Tree:         path,
					Subscription: entry.Subscription,
					Handler:      callback,
					GetTraverseFS: func(_ string) tv.TraverseFS {
						return FS
					},
				},
				tv.WithOnBegin(lab.Begin("🛡️")),
				tv.WithOnEnd(lab.End("🏁")),

				tv.IfOption(entry.CaseSensitive, tv.WithHookCaseSensitiveSort()),
			)).Navigate(ctx)

			lab.AssertNavigation(entry, &lab.TestOptions{
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
		func(entry *lab.NaviTE) string {
			return fmt.Sprintf("🧪 ===> given: '%v'", entry.Given)
		},

		// === universal =====================================================

		Entry(nil, Label(lab.Static.RetroWave), &lab.NaviTE{
			Given:        "universal: Path is leaf",
			Relative:     "RETRO-WAVE/Chromatics/Night Drive",
			Subscription: enums.SubscribeUniversal,
			Callback:     lab.UniversalCallback("LEAF-PATH"),
			ExpectedNoOf: lab.Quantities{
				Files:   4,
				Folders: 1,
			},
		}),

		Entry(nil, Label(lab.Static.RetroWave), &lab.NaviTE{
			Given:        "universal: Path contains folders",
			Relative:     lab.Static.RetroWave,
			Subscription: enums.SubscribeUniversal,
			Callback:     lab.UniversalCallback("CONTAINS-FOLDERS"),
			ExpectedNoOf: lab.Quantities{
				Files:   14,
				Folders: 8,
			},
		}),

		Entry(nil, Label(lab.Static.RetroWave), &lab.NaviTE{
			Given:        "universal: Path contains folders (visit)",
			Relative:     lab.Static.RetroWave,
			Visit:        true,
			Subscription: enums.SubscribeUniversal,
			Callback:     lab.UniversalCallback("VISIT-CONTAINS-FOLDERS"),
			ExpectedNoOf: lab.Quantities{
				Files:   14,
				Folders: 8,
			},
		}),

		// === folders =======================================================

		Entry(nil, Label(lab.Static.RetroWave), &lab.NaviTE{
			Given:        "folders: Path is leaf",
			Relative:     "RETRO-WAVE/Chromatics/Night Drive",
			Subscription: enums.SubscribeFolders,
			Callback:     lab.FoldersCallback("LEAF-PATH"),
			ExpectedNoOf: lab.Quantities{
				Folders: 1,
			},
		}),

		Entry(nil, Label(lab.Static.RetroWave), &lab.NaviTE{
			Given:        "folders: Path contains folders",
			Relative:     lab.Static.RetroWave,
			Subscription: enums.SubscribeFolders,
			Callback:     lab.FoldersCallback("CONTAINS-FOLDERS"),
			ExpectedNoOf: lab.Quantities{
				Folders: 8,
			},
		}),

		Entry(nil, Label(lab.Static.RetroWave), &lab.NaviTE{
			Given:        "folders: Path contains folders (check all invoked)",
			Relative:     lab.Static.RetroWave,
			Visit:        true,
			Subscription: enums.SubscribeFolders,
			Callback:     lab.FoldersCallback("CONTAINS-FOLDERS (check all invoked)"),
			ExpectedNoOf: lab.Quantities{
				Folders: 8,
			},
		}),

		Entry(nil, Label("metal"), &lab.NaviTE{
			Given:         "folders: case sensitive sort",
			Relative:      "rock/metal",
			Subscription:  enums.SubscribeFolders,
			CaseSensitive: true,
			Callback: lab.FoldersCaseSensitiveCallback(
				"rock/metal/HARD-METAL", "rock/metal/dark",
			),
			ExpectedNoOf: lab.Quantities{
				Files:   0,
				Folders: 41,
			},
		}),

		// === files =========================================================

		Entry(nil, Label(lab.Static.RetroWave), &lab.NaviTE{
			Given:        "files: Path is leaf",
			Relative:     "RETRO-WAVE/Chromatics/Night Drive",
			Subscription: enums.SubscribeFiles,
			Callback:     lab.FilesCallback("LEAF-PATH"),
			ExpectedNoOf: lab.Quantities{
				Files:   4,
				Folders: 0,
			},
		}),

		Entry(nil, Label(lab.Static.RetroWave), &lab.NaviTE{
			Given:        "files: Path contains folders",
			Relative:     lab.Static.RetroWave,
			Subscription: enums.SubscribeFiles,
			Callback:     lab.FilesCallback("CONTAINS-FOLDERS"),
			ExpectedNoOf: lab.Quantities{
				Files:   14,
				Folders: 0,
			},
		}),

		Entry(nil, Label(lab.Static.RetroWave), &lab.NaviTE{
			Given:        "files: Path contains folders",
			Relative:     lab.Static.RetroWave,
			Visit:        true,
			Subscription: enums.SubscribeFiles,
			Callback:     lab.FilesCallback("VISIT-CONTAINS-FOLDERS"),
			ExpectedNoOf: lab.Quantities{
				Files:   14,
				Folders: 0,
			},
		}),
	)
})
