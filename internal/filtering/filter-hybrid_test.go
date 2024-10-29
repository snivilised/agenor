package filtering_test

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok

	"github.com/snivilised/li18ngo"
	nef "github.com/snivilised/nefilim"
	"github.com/snivilised/nefilim/luna"
	tv "github.com/snivilised/traverse"
	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/enums"
	lab "github.com/snivilised/traverse/internal/laboratory"
	"github.com/snivilised/traverse/internal/services"
	"github.com/snivilised/traverse/locale"
	"github.com/snivilised/traverse/pref"
	"github.com/snivilised/traverse/test/hydra"
)

var _ = Describe("feature", Ordered, func() {
	var (
		fS *luna.MemFS
	)

	BeforeAll(func() {
		const (
			verbose = false
		)

		fS = hydra.Nuxx(verbose, lab.Static.RetroWave)
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

	DescribeTable("directories with files filtered",
		func(ctx SpecContext, entry *lab.HybridFilterTE) {
			var (
				childFilter core.ChildTraverseFilter
			)

			recording := make(lab.RecordingMap)
			filterDefs := &pref.FilterOptions{
				Node:  &entry.NodeDef,
				Child: &entry.ChildDef,
				Sink: func(reply pref.FilterReply) {
					childFilter = reply.Child
				},
			}

			path := entry.Relative
			callback := func(servant tv.Servant) error {
				node := servant.Node()
				actualNoChildren := len(node.Children)
				GinkgoWriter.Printf(
					"===> 💠 Child Glob Filter(%v, children: %v)"+
						"source: '%v', node-name: '%v', node-scope: '%v', depth: '%v'\n",
					childFilter.Description(),
					actualNoChildren,
					childFilter.Source(),
					node.Extension.Name,
					node.Extension.Scope,
					node.Extension.Depth,
				)

				recording[node.Extension.Name] = len(node.Children)
				return nil
			}

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

				tv.WithFilter(filterDefs),
			)).Navigate(ctx)

			lab.AssertNavigation(&entry.NaviTE, &lab.TestOptions{
				FS:          fS,
				Recording:   recording,
				Path:        path,
				Result:      result,
				Err:         err,
				ExpectedErr: entry.ExpectedErr,
			})
		},
		func(entry *lab.HybridFilterTE) string {
			return fmt.Sprintf("🧪 ===> given: '%v'", entry.Given)
		},

		Entry(nil, &lab.HybridFilterTE{
			NaviTE: lab.NaviTE{
				Given:        "directory(with files): glob child filter",
				Relative:     lab.Static.RetroWave,
				Subscription: enums.SubscribeDirectoriesWithFiles,
				ExpectedNoOf: lab.Quantities{
					Directories: 6,
					Children: map[string]int{
						"Northern Council": 2,
						"Teenage Color":    2,
						"Innerworld":       2,
					},
				},
			},
			NodeDef: core.FilterDef{
				Type:        enums.FilterTypeGlob,
				Description: "directories contains o",
				Pattern:     "*o*",
				Scope:       enums.ScopeDirectory,
			},
			ChildDef: core.ChildFilterDef{
				Type:        enums.FilterTypeGlob,
				Description: "items with '.flac' suffix",
				Pattern:     "*.flac",
			},
		}),

		Entry(nil, &lab.HybridFilterTE{
			NaviTE: lab.NaviTE{
				Given:        "directory(with files): glob child filter (negate)",
				Relative:     lab.Static.RetroWave,
				Subscription: enums.SubscribeDirectoriesWithFiles,
				ExpectedNoOf: lab.Quantities{
					Directories: 2,
					Children: map[string]int{
						"Night Drive": 3,
					},
				},
			},
			NodeDef: core.FilterDef{
				Type:        enums.FilterTypeGlob,
				Description: "directories don't contain o",
				Pattern:     "*o*",
				Scope:       enums.ScopeDirectory,
				Negate:      true,
			},
			ChildDef: core.ChildFilterDef{
				Type:        enums.FilterTypeGlob,
				Description: "items without '.txt' suffix",
				Pattern:     "*.txt",
				Negate:      true,
			},
		}),

		Entry(nil, &lab.HybridFilterTE{
			NaviTE: lab.NaviTE{
				Given:        "directory(with files): regex child filter",
				Relative:     lab.Static.RetroWave,
				Subscription: enums.SubscribeDirectoriesWithFiles,
				ExpectedNoOf: lab.Quantities{
					Directories: 6,
					Children: map[string]int{
						"Northern Council": 2,
						"Teenage Color":    2,
						"Innerworld":       2,
					},
				},
			},
			NodeDef: core.FilterDef{
				Type:        enums.FilterTypeGlob,
				Description: "directories contains o",
				Pattern:     "*o*",
				Scope:       enums.ScopeDirectory,
			},
			ChildDef: core.ChildFilterDef{
				Type:        enums.FilterTypeRegex,
				Description: "items with '.flac' suffix",
				Pattern:     `\.flac`,
			},
		}),

		Entry(nil, &lab.HybridFilterTE{
			NaviTE: lab.NaviTE{
				Given:        "directory(with files): regex child filter (negate)",
				Relative:     lab.Static.RetroWave,
				Subscription: enums.SubscribeDirectoriesWithFiles,
				ExpectedNoOf: lab.Quantities{
					Directories: 2,
					Children: map[string]int{
						"Night Drive": 3,
					},
				},
			},
			NodeDef: core.FilterDef{
				Type:        enums.FilterTypeGlob,
				Description: "directories don't contain o",
				Pattern:     "*o*",
				Scope:       enums.ScopeDirectory,
				Negate:      true,
			},
			ChildDef: core.ChildFilterDef{
				Type:        enums.FilterTypeRegex,
				Description: "items without '.txt' suffix",
				Pattern:     `\.txt$`,
				Negate:      true,
			},
		}),

		Entry(nil, &lab.HybridFilterTE{
			NaviTE: lab.NaviTE{
				Given:        "directory(with files): glob child filter",
				Relative:     lab.Static.RetroWave,
				Subscription: enums.SubscribeDirectoriesWithFiles,
				ExpectedNoOf: lab.Quantities{
					Directories: 6,
					Children: map[string]int{
						"Northern Council": 2,
						"Teenage Color":    2,
						"Innerworld":       2,
					},
				},
			},
			NodeDef: core.FilterDef{
				Type:        enums.FilterTypeGlob,
				Description: "directories contains o",
				Pattern:     "*o*",
				Scope:       enums.ScopeDirectory,
			},
			ChildDef: core.ChildFilterDef{
				Type:        enums.FilterTypeGlobEx,
				Description: "items with '.flac' suffix",
				Pattern:     "*|flac",
			},
		}),

		Entry(nil, &lab.HybridFilterTE{
			NaviTE: lab.NaviTE{
				Given:        "directory(with files): glob child filter (negate)",
				Relative:     lab.Static.RetroWave,
				Subscription: enums.SubscribeDirectoriesWithFiles,
				ExpectedNoOf: lab.Quantities{
					Directories: 2,
					Children: map[string]int{
						"Night Drive": 3,
					},
				},
			},
			NodeDef: core.FilterDef{
				Type:        enums.FilterTypeGlob,
				Description: "directories don't contain o",
				Pattern:     "*o*",
				Scope:       enums.ScopeDirectory,
				Negate:      true,
			},
			ChildDef: core.ChildFilterDef{
				Type:        enums.FilterTypeGlobEx,
				Description: "items without '.txt' suffix",
				Pattern:     "*|txt",
				Negate:      true,
			},
		}),

		Entry(nil, &lab.HybridFilterTE{
			NaviTE: lab.NaviTE{
				Given:        "directory(with files): glob child filter",
				Relative:     lab.Static.RetroWave,
				Subscription: enums.SubscribeDirectoriesWithFiles,
				ExpectedNoOf: lab.Quantities{
					Directories: 6,
					Children: map[string]int{
						"Northern Council": 2,
						"Teenage Color":    2,
						"Innerworld":       2,
					},
				},
				ExpectedErr: locale.ErrFilterCustomNotSupported,
			},
			NodeDef: core.FilterDef{
				Type:        enums.FilterTypeGlob,
				Description: "directories contains o",
				Pattern:     "*o*",
				Scope:       enums.ScopeDirectory,
			},
			ChildDef: core.ChildFilterDef{
				Type:        enums.FilterTypeCustom,
				Description: "items with '.flac' suffix",
				Pattern:     "*|flac",
			},
		}),

		Entry(nil, &lab.HybridFilterTE{
			NaviTE: lab.NaviTE{
				Given:        "directory(with files): glob child filter (negate)",
				Relative:     lab.Static.RetroWave,
				Subscription: enums.SubscribeDirectoriesWithFiles,
				ExpectedNoOf: lab.Quantities{
					Directories: 2,
					Children: map[string]int{
						"Night Drive": 3,
					},
				},
			},
			NodeDef: core.FilterDef{
				Type:        enums.FilterTypeGlob,
				Description: "directories don't contain o",
				Pattern:     "*o*",
				Scope:       enums.ScopeDirectory,
				Negate:      true,
			},
			ChildDef: core.ChildFilterDef{
				Type:        enums.FilterTypeGlobEx,
				Description: "items without '.txt' suffix",
				Pattern:     "*|txt",
				Negate:      true,
			},
		}),

		// === error ==============================================================

		Entry(nil, &lab.HybridFilterTE{
			NaviTE: lab.NaviTE{
				Given:        "malformed glob ex filter (missing |)",
				Should:       "fail",
				Relative:     lab.Static.RetroWave,
				Subscription: enums.SubscribeDirectoriesWithFiles,
				ExpectedErr:  locale.NewInvalidIncaseFilterDefError("*.flac"),
			},
			NodeDef: core.FilterDef{
				Type:        enums.FilterTypeGlob,
				Description: "directories contains o",
				Pattern:     "*o*",
				Scope:       enums.ScopeDirectory,
			},
			ChildDef: core.ChildFilterDef{
				Type:        enums.FilterTypeGlobEx,
				Description: "items with '.flac' suffix",
				Pattern:     "*.flac",
			},
		}),

		Entry(nil, &lab.HybridFilterTE{
			NaviTE: lab.NaviTE{
				Given:        "malformed glob ex filter, missing type",
				Should:       "fail",
				Relative:     lab.Static.RetroWave,
				Subscription: enums.SubscribeDirectoriesWithFiles,
				ExpectedErr:  locale.ErrFilterUndefined,
			},
			NodeDef: core.FilterDef{
				Type:        enums.FilterTypeGlob,
				Description: "directories contains o",
				Pattern:     "*o*",
				Scope:       enums.ScopeDirectory,
			},
			ChildDef: core.ChildFilterDef{
				Description: "type missing",
				Pattern:     "*.flac",
			},
		}),
	)
})
