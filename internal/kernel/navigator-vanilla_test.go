package kernel_test

import (
	"fmt"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok

	"github.com/snivilised/li18ngo"
	nef "github.com/snivilised/nefilim"
	"github.com/snivilised/nefilim/test/luna"
	tv "github.com/snivilised/traverse"
	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/enums"
	lab "github.com/snivilised/traverse/internal/laboratory"
	"github.com/snivilised/traverse/internal/services"
	"github.com/snivilised/traverse/internal/third/lo"
	"github.com/snivilised/traverse/locale"
	"github.com/snivilised/traverse/test/hydra"
)

var _ = Describe("NavigatorUniversal", Ordered, func() {
	var (
		fS *luna.MemFS
	)

	BeforeAll(func() {
		const (
			verbose = false
		)

		fS = hydra.Nuxx(verbose,
			lab.Static.RetroWave,
			filepath.Join("rock", "metal"),
		)
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

	DescribeTable("Ensure Callback Invoked Once", Label("vanilla"),
		func(ctx SpecContext, entry *lab.NaviTE) {
			recall := make(lab.Recall)
			once := func(servant tv.Servant) error {
				node := servant.Node()
				_, found := recall[node.Path] // TODO: should this be name not path?
				Expect(found).To(BeFalse())
				recall[node.Path] = len(node.Children)

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
					GetForest: func(_ string) *core.Forest {
						return &core.Forest{
							T: fS,
							R: nef.NewTraverseABS(),
						}
					},
				},
				tv.WithOnBegin(lab.Begin("🛡️")),
				tv.WithOnEnd(lab.End("🏁")),

				tv.IfOption(entry.CaseSensitive, tv.WithHookCaseSensitiveSort()),
			)).Navigate(ctx)

			lab.AssertNavigation(entry, &lab.TestOptions{
				FS:        fS,
				Recording: recall,
				Path:      path,
				Result:    result,
				Err:       err,
				Every: func(p string) bool {
					_, found := recall[p]
					return found
				},
				ByPassMetrics: entry.ByPassMetrics,
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
				Files:       4,
				Directories: 1,
			},
		}),

		Entry(nil, Label(lab.Static.RetroWave), &lab.NaviTE{
			Given:        "universal: Path contains directories",
			Relative:     lab.Static.RetroWave,
			Subscription: enums.SubscribeUniversal,
			Callback:     lab.UniversalCallback("CONTAINS-DIRECTORIES"),
			ExpectedNoOf: lab.Quantities{
				Files:       14,
				Directories: 8,
			},
		}),

		Entry(nil, Label(lab.Static.RetroWave), &lab.NaviTE{
			Given:        "universal: Path contains directories (visit)",
			Relative:     lab.Static.RetroWave,
			Visit:        true,
			Subscription: enums.SubscribeUniversal,
			Callback:     lab.UniversalCallback("VISIT-CONTAINS-DIRECTORIES"),
			ExpectedNoOf: lab.Quantities{
				Files:       14,
				Directories: 8,
			},
		}),

		Entry(nil, Label(lab.Static.RetroWave), &lab.NaviTE{
			Given:         "universal: Path is Root",
			Relative:      ".",
			Subscription:  enums.SubscribeUniversal,
			Callback:      lab.UniversalCallback("ROOT-PATH"),
			ByPassMetrics: true,
		}),

		// === directories ===================================================

		Entry(nil, Label(lab.Static.RetroWave), &lab.NaviTE{
			Given:        "directories: Path is leaf",
			Relative:     "RETRO-WAVE/Chromatics/Night Drive",
			Subscription: enums.SubscribeDirectories,
			Callback:     lab.DirectoriesCallback("LEAF-PATH"),
			ExpectedNoOf: lab.Quantities{
				Directories: 1,
			},
		}),

		Entry(nil, Label(lab.Static.RetroWave), &lab.NaviTE{
			Given:        "directories: Path contains directories",
			Relative:     lab.Static.RetroWave,
			Subscription: enums.SubscribeDirectories,
			Callback:     lab.DirectoriesCallback("CONTAINS-DIRECTORIES"),
			ExpectedNoOf: lab.Quantities{
				Directories: 8,
			},
		}),

		Entry(nil, Label(lab.Static.RetroWave), &lab.NaviTE{
			Given:        "directories: Path contains directories (check all invoked)",
			Relative:     lab.Static.RetroWave,
			Visit:        true,
			Subscription: enums.SubscribeDirectories,
			Callback:     lab.DirectoriesCallback("CONTAINS-DIRECTORIES (check all invoked)"),
			ExpectedNoOf: lab.Quantities{
				Directories: 8,
			},
		}),

		Entry(nil, Label("metal"), &lab.NaviTE{
			Given:         "directories: case sensitive sort",
			Relative:      "rock/metal",
			Subscription:  enums.SubscribeDirectories,
			CaseSensitive: true,
			Callback: lab.DirectoriesCaseSensitiveCallback(
				"rock/metal/HARD-METAL", "rock/metal/dark",
			),
			ExpectedNoOf: lab.Quantities{
				Files:       0,
				Directories: 41,
			},
		}),

		// === files =========================================================

		Entry(nil, Label(lab.Static.RetroWave), &lab.NaviTE{
			Given:        "files: Path is leaf",
			Relative:     "RETRO-WAVE/Chromatics/Night Drive",
			Subscription: enums.SubscribeFiles,
			Callback:     lab.FilesCallback("LEAF-PATH"),
			ExpectedNoOf: lab.Quantities{
				Files:       4,
				Directories: 0,
			},
		}),

		Entry(nil, Label(lab.Static.RetroWave), &lab.NaviTE{
			Given:        "files: Path contains directories",
			Relative:     lab.Static.RetroWave,
			Subscription: enums.SubscribeFiles,
			Callback:     lab.FilesCallback("CONTAINS-DIRECTORIES"),
			ExpectedNoOf: lab.Quantities{
				Files:       14,
				Directories: 0,
			},
		}),

		Entry(nil, Label(lab.Static.RetroWave), &lab.NaviTE{
			Given:        "files: Path contains directories",
			Relative:     lab.Static.RetroWave,
			Visit:        true,
			Subscription: enums.SubscribeFiles,
			Callback:     lab.FilesCallback("VISIT-CONTAINS-DIRECTORIES"),
			ExpectedNoOf: lab.Quantities{
				Files:       14,
				Directories: 0,
			},
		}),
	)
})
