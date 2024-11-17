package kernel_test

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok
	age "github.com/snivilised/agenor"
	"github.com/snivilised/agenor/core"
	"github.com/snivilised/agenor/enums"
	lab "github.com/snivilised/agenor/internal/laboratory"
	"github.com/snivilised/agenor/internal/services"
	"github.com/snivilised/agenor/locale"
	"github.com/snivilised/agenor/pref"
	"github.com/snivilised/agenor/test/hydra"
	"github.com/snivilised/agenor/tfs"
	"github.com/snivilised/li18ngo"
	"github.com/snivilised/nefilim/test/luna"
)

const (
	verbose = false
)

var _ = Describe("navigator", Ordered, func() {
	var (
		fS *luna.MemFS
	)

	BeforeAll(func() {
		Expect(li18ngo.Use(
			func(o *li18ngo.UseOptions) {
				o.From.Sources = li18ngo.TranslationFiles{
					locale.SourceID: li18ngo.TranslationSource{Name: "agenor"},
				}
			},
		)).To(Succeed())

		fS = hydra.Nuxx(verbose, lab.Static.RetroWave)
	})

	BeforeEach(func() {
		services.Reset()
	})

	DescribeTable("cascade",
		func(ctx SpecContext, entry *lab.CascadeTE) {
			path := entry.Relative
			result, err := age.Walk().Configure().Extent(age.Prime(
				&pref.Using{
					Tree:         path,
					Subscription: entry.Subscription,
					Head: pref.Head{
						Handler: entry.Callback,
						GetForest: func(_ string) *core.Forest {
							return &core.Forest{
								T: fS,
								R: tfs.New(),
							}
						},
					},
				},
				age.WithOnBegin(lab.Begin("🛡️")),
				age.WithOnEnd(lab.End("🏁")),
				age.IfOptionF(entry.Depth > 0, func() pref.Option {
					return age.WithDepth(entry.Depth)
				}),
				age.IfOptionF(entry.NoRecurse, func() pref.Option {
					return age.WithNoRecurse()
				}),
			)).Navigate(ctx)

			lab.AssertNavigation(&entry.NaviTE, &lab.TestOptions{
				FS:     fS,
				Path:   path,
				Result: result,
				Err:    err,
			})
		},
		func(entry *lab.CascadeTE) string {
			return fmt.Sprintf("🧪 ===> given: '%v'", entry.Given)
		},

		// === universal =====================================================

		Entry(nil, &lab.CascadeTE{
			NaviTE: lab.NaviTE{
				Given:        "universal: Path contains folders only, no-recurse",
				Should:       "traverse single level",
				Relative:     "RETRO-WAVE",
				Subscription: enums.SubscribeUniversal,
				Callback:     lab.UniversalCallback("CONTAINS-FOLDERS"),
				ExpectedNoOf: lab.Quantities{
					Files:       0,
					Directories: 4,
				},
			},
			NoRecurse: true,
		}),

		Entry(nil, &lab.CascadeTE{
			NaviTE: lab.NaviTE{
				Given:        "universal: Path contains files only, no-recurse",
				Should:       "traverse single level (containing files)",
				Relative:     "RETRO-WAVE/Chromatics/Night Drive",
				Subscription: enums.SubscribeUniversal,
				Callback:     lab.UniversalCallback("CONTAINS-FILES"),
				ExpectedNoOf: lab.Quantities{
					Files:       4,
					Directories: 1,
				},
			},
			NoRecurse: true,
		}),

		Entry(nil, &lab.CascadeTE{
			NaviTE: lab.NaviTE{
				Given:        "universal: Path contains folders only, depth=1",
				Should:       "traverse single level",
				Relative:     "RETRO-WAVE",
				Subscription: enums.SubscribeUniversal,
				Callback:     lab.UniversalCallback("CONTAINS-FOLDERS"),
				ExpectedNoOf: lab.Quantities{
					Files:       0,
					Directories: 4,
				},
			},
			Depth: 1,
		}),

		Entry(nil, &lab.CascadeTE{
			NaviTE: lab.NaviTE{
				Given:        "universal: Path contains folders only, depth=2",
				Should:       "traverse 2 levels",
				Relative:     "RETRO-WAVE",
				Subscription: enums.SubscribeUniversal,
				Callback:     lab.UniversalCallback("CONTAINS-FOLDERS"),
				ExpectedNoOf: lab.Quantities{
					Files:       0,
					Directories: 8,
				},
			},
			Depth: 2,
		}),

		Entry(nil, &lab.CascadeTE{
			NaviTE: lab.NaviTE{
				Given:        "universal: Path contains folders only, depth=3",
				Should:       "traverse 3 levels (containing files)",
				Relative:     "RETRO-WAVE",
				Subscription: enums.SubscribeUniversal,
				Callback:     lab.UniversalCallback("CONTAINS-FOLDERS"),
				ExpectedNoOf: lab.Quantities{
					Files:       14,
					Directories: 8,
				},
			},
			Depth: 3,
		}),

		// === folders =======================================================

		Entry(nil, &lab.CascadeTE{
			NaviTE: lab.NaviTE{
				Given:        "universal: Path contains folders only, no-recurse",
				Should:       "traverse single level",
				Relative:     "RETRO-WAVE",
				Subscription: enums.SubscribeDirectories,
				Callback:     lab.DirectoriesCallback("CONTAINS-FILES"),
				ExpectedNoOf: lab.Quantities{
					Files:       0,
					Directories: 4,
				},
			},
			NoRecurse: true,
		}),

		Entry(nil, &lab.CascadeTE{
			NaviTE: lab.NaviTE{
				Given:        "universal: Path contains files only, no-recurse",
				Should:       "traverse single level (containing files)",
				Relative:     "RETRO-WAVE/Chromatics/Night Drive",
				Subscription: enums.SubscribeDirectories,
				Callback:     lab.UniversalCallback("LEAF-PATH"),
				ExpectedNoOf: lab.Quantities{
					Files:       0,
					Directories: 1,
				},
			},
			NoRecurse: true,
		}),

		Entry(nil, &lab.CascadeTE{
			NaviTE: lab.NaviTE{
				Given:        "universal: Path contains folders only, depth=1",
				Should:       "traverse single level",
				Relative:     "RETRO-WAVE",
				Subscription: enums.SubscribeDirectories,
				Callback:     lab.UniversalCallback("CONTAINS-FOLDERS"),
				ExpectedNoOf: lab.Quantities{
					Files:       0,
					Directories: 4,
				},
			},
			Depth: 1,
		}),

		Entry(nil, &lab.CascadeTE{
			NaviTE: lab.NaviTE{
				Given:        "universal: Path contains folders only, depth=2",
				Should:       "traverse 2 levels",
				Relative:     "RETRO-WAVE",
				Subscription: enums.SubscribeDirectories,
				Callback:     lab.UniversalCallback("CONTAINS-FOLDERS"),
				ExpectedNoOf: lab.Quantities{
					Files:       0,
					Directories: 8,
				},
			},
			Depth: 2,
		}),

		Entry(nil, &lab.CascadeTE{
			NaviTE: lab.NaviTE{
				Given:        "universal: Path contains folders only, depth=3",
				Should:       "traverse 3 levels (containing files)",
				Relative:     "RETRO-WAVE",
				Subscription: enums.SubscribeDirectories,
				Callback:     lab.UniversalCallback("CONTAINS-FOLDERS"),
				ExpectedNoOf: lab.Quantities{
					Files:       0,
					Directories: 8,
				},
			},
			Depth: 3,
		}),

		// === files =========================================================

		Entry(nil, &lab.CascadeTE{
			NaviTE: lab.NaviTE{
				Given:        "file: Path contains folders only, no-recurse",
				Should:       "traverse single level",
				Relative:     "RETRO-WAVE/Chromatics/Night Drive",
				Subscription: enums.SubscribeFiles,
				Callback:     lab.FilesCallback("FILE"),
				ExpectedNoOf: lab.Quantities{
					Files:       4,
					Directories: 0,
				},
			},
			NoRecurse: true,
		}),

		Entry(nil, &lab.CascadeTE{
			NaviTE: lab.NaviTE{
				Given:        "file: Path contains folders only, depth=1",
				Should:       "traverse single level",
				Relative:     "RETRO-WAVE/Chromatics/Night Drive",
				Subscription: enums.SubscribeFiles,
				Callback:     lab.FilesCallback("FILE"),
				ExpectedNoOf: lab.Quantities{
					Files:       4,
					Directories: 0,
				},
			},
			Depth: 1,
		}),
	)
})
