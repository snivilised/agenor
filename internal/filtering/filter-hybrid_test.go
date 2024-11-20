package filtering_test

import (
	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok

	age "github.com/snivilised/agenor"
	"github.com/snivilised/agenor/core"
	"github.com/snivilised/agenor/enums"
	lab "github.com/snivilised/agenor/internal/laboratory"
	"github.com/snivilised/agenor/internal/services"
	"github.com/snivilised/agenor/locale"
	"github.com/snivilised/agenor/pref"
	"github.com/snivilised/agenor/test/hanno"
	"github.com/snivilised/agenor/tfs"
	"github.com/snivilised/li18ngo"
	"github.com/snivilised/nefilim/test/luna"
)

var _ = Describe("feature", Ordered, func() {
	var (
		fS *luna.MemFS
	)

	BeforeAll(func() {
		const (
			verbose = false
		)

		fS = hanno.Nuxx(verbose, lab.Static.RetroWave)
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

	DescribeTable("directories with files filtered",
		func(ctx SpecContext, entry *lab.HybridFilterTE) {
			var (
				childFilter core.ChildTraverseFilter
			)

			recall := make(lab.Recall)
			filterDefs := &pref.FilterOptions{
				Node:  &entry.NodeDef,
				Child: &entry.ChildDef,
				Sink: func(reply pref.FilterReply) {
					childFilter = reply.Child
				},
			}

			path := entry.Relative
			callback := func(servant age.Servant) error {
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

				recall[node.Extension.Name] = len(node.Children)
				return nil
			}

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

				age.WithFilter(filterDefs),
			)).Navigate(ctx)

			lab.AssertNavigation(&entry.NaviTE, &lab.TestOptions{
				FS:          fS,
				Recording:   recall,
				Path:        path,
				Result:      result,
				Err:         err,
				ExpectedErr: entry.ExpectedErr,
			})
		},
		lab.FormatHybridFilterTestDescription,

		Entry(nil, &lab.HybridFilterTE{
			DescribedTE: lab.DescribedTE{
				Given: "directory(with files): glob child filter",
			},
			NaviTE: lab.NaviTE{
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
			DescribedTE: lab.DescribedTE{
				Given: "directory(with files): glob child filter (negate)",
			},
			NaviTE: lab.NaviTE{
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
			DescribedTE: lab.DescribedTE{
				Given: "directory(with files): regex child filter",
			},
			NaviTE: lab.NaviTE{
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
			DescribedTE: lab.DescribedTE{
				Given: "directory(with files): regex child filter (negate)",
			},
			NaviTE: lab.NaviTE{
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
			DescribedTE: lab.DescribedTE{
				Given: "directory(with files): glob child filter",
			},
			NaviTE: lab.NaviTE{
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
			DescribedTE: lab.DescribedTE{
				Given: "directory(with files): glob child filter (negate)",
			},
			NaviTE: lab.NaviTE{
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
			DescribedTE: lab.DescribedTE{
				Given: "directory(with files): glob child filter",
			},
			NaviTE: lab.NaviTE{
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
			DescribedTE: lab.DescribedTE{
				Given: "directory(with files): glob child filter (negate)",
			},
			NaviTE: lab.NaviTE{
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
			DescribedTE: lab.DescribedTE{
				Given:  "malformed glob ex filter (missing |)",
				Should: "fail",
			},
			NaviTE: lab.NaviTE{
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
			DescribedTE: lab.DescribedTE{
				Given:  "malformed glob ex filter, missing type",
				Should: "fail",
			},
			NaviTE: lab.NaviTE{
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
