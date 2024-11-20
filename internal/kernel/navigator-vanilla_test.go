package kernel_test

import (
	"path/filepath"

	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok

	age "github.com/snivilised/agenor"
	"github.com/snivilised/agenor/core"
	"github.com/snivilised/agenor/enums"
	lab "github.com/snivilised/agenor/internal/laboratory"
	"github.com/snivilised/agenor/internal/services"
	"github.com/snivilised/agenor/internal/third/lo"
	"github.com/snivilised/agenor/locale"
	"github.com/snivilised/agenor/pref"
	"github.com/snivilised/agenor/test/hanno"
	"github.com/snivilised/agenor/tfs"
	"github.com/snivilised/li18ngo"
	"github.com/snivilised/nefilim/test/luna"
)

var _ = Describe("NavigatorUniversal", Ordered, func() {
	var (
		fS *luna.MemFS
	)

	BeforeAll(func() {
		const (
			verbose = false
		)

		fS = hanno.Nuxx(verbose,
			lab.Static.RetroWave,
			filepath.Join("rock", "metal"),
		)
		Expect(li18ngo.Use(
			func(o *li18ngo.UseOptions) {
				o.From.Sources = li18ngo.TranslationFiles{
					locale.SourceID: li18ngo.TranslationSource{Name: "agenor"},
				}
			},
		)).To(Succeed())
	})

	BeforeEach(func() {
		services.Reset()
	})

	DescribeTable("Ensure Callback Invoked Once", Label("vanilla"),
		func(ctx SpecContext, entry *lab.GeneralTE) {
			recall := make(lab.Recall)
			once := func(servant age.Servant) error {
				node := servant.Node()
				_, found := recall[node.Path] // TODO: should this be name not path?
				Expect(found).To(BeFalse())
				recall[node.Path] = len(node.Children)

				return entry.Callback(servant)
			}

			visitor := func(servant age.Servant) error {
				return once(servant)
			}

			callback := lo.Ternary(entry.Once, once,
				lo.Ternary(entry.Visit, visitor, entry.Callback),
			)
			path := entry.Relative

			result, err := age.Walk().Configure().Extent(age.Prime(
				&pref.Using{
					Subscription: entry.Subscription,
					Head: pref.Head{
						Handler: callback,
						GetForest: func(_ string) *core.Forest {
							return &core.Forest{
								T: fS,
								R: tfs.New(),
							}
						},
					},
					Tree: path,
				},
				age.WithOnBegin(lab.Begin("🛡️")),
				age.WithOnEnd(lab.End("🏁")),

				age.IfOption(entry.CaseSensitive, age.WithHookCaseSensitiveSort()),
			)).Navigate(ctx)

			lab.AssertNavigation(&entry.NaviTE, &lab.TestOptions{
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
		lab.FormatGeneralTestDescription,

		// === universal =====================================================

		Entry(nil, Label(lab.Static.RetroWave), &lab.GeneralTE{
			DescribedTE: lab.DescribedTE{
				Given: "universal: Path is leaf",
			},
			NaviTE: lab.NaviTE{
				Relative:     "RETRO-WAVE/Chromatics/Night Drive",
				Subscription: enums.SubscribeUniversal,
				Callback:     lab.UniversalCallback("LEAF-PATH"),
				ExpectedNoOf: lab.Quantities{
					Files:       4,
					Directories: 1,
				},
			},
		}),

		Entry(nil, Label(lab.Static.RetroWave), &lab.GeneralTE{
			DescribedTE: lab.DescribedTE{
				Given: "universal: Path contains directories",
			},
			NaviTE: lab.NaviTE{
				Relative:     lab.Static.RetroWave,
				Subscription: enums.SubscribeUniversal,
				Callback:     lab.UniversalCallback("CONTAINS-DIRECTORIES"),
				ExpectedNoOf: lab.Quantities{
					Files:       14,
					Directories: 8,
				},
			},
		}),

		Entry(nil, Label(lab.Static.RetroWave), &lab.GeneralTE{
			DescribedTE: lab.DescribedTE{
				Given: "universal: Path contains directories (visit)",
			},
			NaviTE: lab.NaviTE{
				Relative:     lab.Static.RetroWave,
				Visit:        true,
				Subscription: enums.SubscribeUniversal,
				Callback:     lab.UniversalCallback("VISIT-CONTAINS-DIRECTORIES"),
				ExpectedNoOf: lab.Quantities{
					Files:       14,
					Directories: 8,
				},
			},
		}),

		Entry(nil, Label(lab.Static.RetroWave), &lab.GeneralTE{
			DescribedTE: lab.DescribedTE{
				Given: "universal: Path is Root",
			},
			NaviTE: lab.NaviTE{
				Relative:      ".",
				Subscription:  enums.SubscribeUniversal,
				Callback:      lab.UniversalCallback("ROOT-PATH"),
				ByPassMetrics: true,
			},
		}),

		// === directories ===================================================

		Entry(nil, Label(lab.Static.RetroWave), &lab.GeneralTE{
			DescribedTE: lab.DescribedTE{
				Given: "directories: Path is leaf",
			},
			NaviTE: lab.NaviTE{
				Relative:     "RETRO-WAVE/Chromatics/Night Drive",
				Subscription: enums.SubscribeDirectories,
				Callback:     lab.DirectoriesCallback("LEAF-PATH"),
				ExpectedNoOf: lab.Quantities{
					Directories: 1,
				},
			},
		}),

		Entry(nil, Label(lab.Static.RetroWave), &lab.GeneralTE{
			DescribedTE: lab.DescribedTE{
				Given: "directories: Path contains directories",
			},
			NaviTE: lab.NaviTE{
				Relative:     lab.Static.RetroWave,
				Subscription: enums.SubscribeDirectories,
				Callback:     lab.DirectoriesCallback("CONTAINS-DIRECTORIES"),
				ExpectedNoOf: lab.Quantities{
					Directories: 8,
				},
			},
		}),

		Entry(nil, Label(lab.Static.RetroWave), &lab.GeneralTE{
			DescribedTE: lab.DescribedTE{
				Given: "directories: Path contains directories (check all invoked)",
			},
			NaviTE: lab.NaviTE{
				Relative:     lab.Static.RetroWave,
				Visit:        true,
				Subscription: enums.SubscribeDirectories,
				Callback:     lab.DirectoriesCallback("CONTAINS-DIRECTORIES (check all invoked)"),
				ExpectedNoOf: lab.Quantities{
					Directories: 8,
				},
			},
		}),

		Entry(nil, Label("metal"), &lab.GeneralTE{
			DescribedTE: lab.DescribedTE{
				Given: "directories: case sensitive sort",
			},
			NaviTE: lab.NaviTE{
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
			},
		}),

		// === files =========================================================

		Entry(nil, Label(lab.Static.RetroWave), &lab.GeneralTE{
			DescribedTE: lab.DescribedTE{
				Given: "files: Path is leaf",
			},
			NaviTE: lab.NaviTE{
				Relative:     "RETRO-WAVE/Chromatics/Night Drive",
				Subscription: enums.SubscribeFiles,
				Callback:     lab.FilesCallback("LEAF-PATH"),
				ExpectedNoOf: lab.Quantities{
					Files:       4,
					Directories: 0,
				},
			},
		}),

		Entry(nil, Label(lab.Static.RetroWave), &lab.GeneralTE{
			DescribedTE: lab.DescribedTE{
				Given: "files: Path contains directories",
			},
			NaviTE: lab.NaviTE{
				Relative:     lab.Static.RetroWave,
				Subscription: enums.SubscribeFiles,
				Callback:     lab.FilesCallback("CONTAINS-DIRECTORIES"),
				ExpectedNoOf: lab.Quantities{
					Files:       14,
					Directories: 0,
				},
			},
		}),

		Entry(nil, Label(lab.Static.RetroWave), &lab.GeneralTE{
			DescribedTE: lab.DescribedTE{
				Given: "files: Path contains directories",
			},
			NaviTE: lab.NaviTE{
				Relative:     lab.Static.RetroWave,
				Visit:        true,
				Subscription: enums.SubscribeFiles,
				Callback:     lab.FilesCallback("VISIT-CONTAINS-DIRECTORIES"),
				ExpectedNoOf: lab.Quantities{
					Files:       14,
					Directories: 0,
				},
			},
		}),
	)
})
