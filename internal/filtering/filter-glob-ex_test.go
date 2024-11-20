package filtering_test

import (
	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok

	age "github.com/snivilised/agenor"
	"github.com/snivilised/agenor/core"
	"github.com/snivilised/agenor/enums"
	lab "github.com/snivilised/agenor/internal/laboratory"
	"github.com/snivilised/agenor/internal/services"
	"github.com/snivilised/agenor/internal/third/lo"
	"github.com/snivilised/agenor/pref"
	"github.com/snivilised/agenor/test/hanno"
	"github.com/snivilised/agenor/tfs"
	"github.com/snivilised/li18ngo"
	"github.com/snivilised/nefilim/test/luna"
)

var _ = Describe("filtering", Ordered, func() {
	var (
		fS *luna.MemFS
	)

	BeforeAll(func() {
		const (
			verbose = false
		)

		fS = hanno.Nuxx(verbose, "rock")
		Expect(li18ngo.Use()).To(Succeed())
	})

	BeforeEach(func() {
		services.Reset()
	})

	Context("comprehension", func() {
		When("universal: filtering with glob ex", func() {
			It("should: invoke for filtered nodes only", Label("example"),
				func(ctx SpecContext) {
					path := lab.Static.RetroWave
					filterDefs := &pref.FilterOptions{
						Node: &core.FilterDef{
							Type:        enums.FilterTypeGlobEx,
							Description: "nodes with 'flac' suffix",
							Pattern:     "*|flac",
							Scope:       enums.ScopeAll,
						},
					}
					result, _ := age.Walk().Configure().Extent(age.Prime(
						&pref.Using{
							Subscription: enums.SubscribeUniversal,
							Head: pref.Head{
								Handler: func(servant age.Servant) error {
									node := servant.Node()
									GinkgoWriter.Printf(
										"---> 🍯 EXAMPLE-EXTENDED-GLOB-FILTER-CALLBACK: '%v'\n", node.Path,
									)
									return nil
								},
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

					GinkgoWriter.Printf("===> 🍭 invoked '%v' directories, '%v' files.\n",
						result.Metrics().Count(enums.MetricNoDirectoriesInvoked),
						result.Metrics().Count(enums.MetricNoFilesInvoked),
					)
				},
			)
		})
	})

	DescribeTable("directories with files",
		func(ctx SpecContext, entry *lab.FilterTE) {
			var (
				traverseFilter core.TraverseFilter
			)

			recall := make(lab.Recall)
			filterDefs := &pref.FilterOptions{
				Node: &core.FilterDef{
					Type:            enums.FilterTypeGlobEx,
					Description:     entry.Description,
					Pattern:         entry.Pattern,
					Scope:           entry.Scope,
					Negate:          entry.Negate,
					IfNotApplicable: entry.IfNotApplicable,
				},
				Sink: func(reply pref.FilterReply) {
					traverseFilter = reply.Node
				},
			}

			path := entry.Relative
			callback := func(servant age.Servant) error {
				node := servant.Node()
				indicator := lo.Ternary(node.IsDirectory(), "📁", "💠")
				GinkgoWriter.Printf(
					"===> %v Glob Filter(%v) source: '%v', node-name: '%v', node-scope(fs): '%v(%v)'\n",
					indicator,
					traverseFilter.Description(),
					traverseFilter.Source(),
					node.Extension.Name,
					node.Extension.Scope,
					traverseFilter.Scope(),
				)
				if lo.Contains(entry.Mandatory, node.Extension.Name) {
					Expect(node).Should(MatchCurrentExtendedFilter(traverseFilter))
				}

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
				FS:        fS,
				Recording: recall,
				Path:      path,
				Result:    result,
				Err:       err,
			})
		},
		lab.FormatFilterTestDescription,

		// === universal =====================================================

		Entry(nil, &lab.FilterTE{
			DescribedTE: lab.DescribedTE{
				Given: "universal(any scope): glob ex filter",
			},
			NaviTE: lab.NaviTE{
				Relative:     "rock/PROGRESSIVE-ROCK/Marillion",
				Subscription: enums.SubscribeUniversal,
				ExpectedNoOf: lab.Quantities{
					Files:       16,
					Directories: 5,
				},
				Prohibited: []string{"cover-clutching-at-straws-jpg"},
			},
			Description: "nodes with 'flac' suffix",
			Pattern:     "*|flac",
			Scope:       enums.ScopeAll,
		}),

		Entry(nil, &lab.FilterTE{
			DescribedTE: lab.DescribedTE{
				Given: "universal(any scope): glob ex filter, with dot extension",
			},
			NaviTE: lab.NaviTE{
				Relative:     "rock/PROGRESSIVE-ROCK/Marillion",
				Subscription: enums.SubscribeUniversal,
				ExpectedNoOf: lab.Quantities{
					Files:       16,
					Directories: 5,
				},
				Prohibited: []string{"cover-clutching-at-straws-jpg"},
			},
			Description: "items with 'flac' suffix",
			Pattern:     "*|.flac",
			Scope:       enums.ScopeAll,
		}),

		Entry(nil, &lab.FilterTE{
			DescribedTE: lab.DescribedTE{
				Given: "universal(any scope): glob ex filter, with multiple extensions",
			},
			NaviTE: lab.NaviTE{
				Relative:     "rock/PROGRESSIVE-ROCK/Marillion",
				Subscription: enums.SubscribeUniversal,
				ExpectedNoOf: lab.Quantities{
					Files:       19,
					Directories: 5,
				},
				Mandatory:  []string{"front.jpg"},
				Prohibited: []string{"cover-clutching-at-straws-jpg"},
			},
			Description: "items with 'flac' suffix",
			Pattern:     "*|flac,jpg",
			Scope:       enums.ScopeAll,
		}),

		Entry(nil, &lab.FilterTE{
			DescribedTE: lab.DescribedTE{
				Given: "universal(any scope): glob ex filter, without extension",
			},
			NaviTE: lab.NaviTE{
				Relative:     "rock/PROGRESSIVE-ROCK/Marillion",
				Subscription: enums.SubscribeUniversal,
				ExpectedNoOf: lab.Quantities{
					Files:       3,
					Directories: 5,
				},
				Mandatory:  []string{"cover-clutching-at-straws-jpg"},
				Prohibited: []string{"01 - Hotel Hobbies.flac"},
			},
			Description: "items with 'flac' suffix",
			Pattern:     "*|",
			Scope:       enums.ScopeAll,
		}),

		Entry(nil, &lab.FilterTE{
			DescribedTE: lab.DescribedTE{
				Given: "universal(file scope): glob ex filter (negate)",
			},
			NaviTE: lab.NaviTE{
				Relative:     "rock/PROGRESSIVE-ROCK/Marillion",
				Subscription: enums.SubscribeUniversal,
				ExpectedNoOf: lab.Quantities{
					Files:       7,
					Directories: 5,
				},
				Prohibited: []string{"01 - Hotel Hobbies.flac"},
			},
			Description: "files without .flac suffix",
			Pattern:     "*|flac",
			Scope:       enums.ScopeFile,
			Negate:      true,
		}),

		Entry(nil, &lab.FilterTE{
			DescribedTE: lab.DescribedTE{
				Given: "universal(undefined scope): glob ex filter",
			},
			NaviTE: lab.NaviTE{
				Relative:     "rock/PROGRESSIVE-ROCK/Marillion",
				Subscription: enums.SubscribeUniversal,
				ExpectedNoOf: lab.Quantities{
					Files:       16,
					Directories: 5,
				},
				Prohibited: []string{"cover-clutching-at-straws-jpg"},
			},
			Description: "items with '.flac' suffix",
			Pattern:     "*|flac",
		}),

		Entry(nil, &lab.FilterTE{
			DescribedTE: lab.DescribedTE{
				Given: "universal(any scope): glob ex filter, any extension",
			},
			NaviTE: lab.NaviTE{
				Relative:     "rock/PROGRESSIVE-ROCK/Marillion",
				Subscription: enums.SubscribeUniversal,
				ExpectedNoOf: lab.Quantities{
					Files:       4,
					Directories: 1,
				},
				Mandatory:  []string{"cover-clutching-at-straws-jpg"},
				Prohibited: []string{"01 - Hotel Hobbies.flac"},
			},
			Description: "starts with c, any extension",
			Pattern:     "c*|*",
			Scope:       enums.ScopeAll,
		}),

		// === directories ===================================================

		Entry(nil, &lab.FilterTE{
			DescribedTE: lab.DescribedTE{
				Given: "directories(any scope): glob ex filter",
			},
			NaviTE: lab.NaviTE{
				Relative:     "rock/PROGRESSIVE-ROCK/Marillion",
				Subscription: enums.SubscribeDirectories,
				ExpectedNoOf: lab.Quantities{
					Files:       0,
					Directories: 2,
				},
				Mandatory:  []string{"Marillion"},
				Prohibited: []string{"Fugazi"},
			},
			Description: "directories starting with M",
			Pattern:     "M*|",
			Scope:       enums.ScopeDirectory,
		}),

		Entry(nil, &lab.FilterTE{
			DescribedTE: lab.DescribedTE{
				Given: "directories(directory scope): glob ex filter (negate)",
			},
			NaviTE: lab.NaviTE{
				Relative:     "rock/PROGRESSIVE-ROCK/Marillion",
				Subscription: enums.SubscribeDirectories,
				ExpectedNoOf: lab.Quantities{
					Files:       0,
					Directories: 3,
				},
				Mandatory:  []string{"Fugazi"},
				Prohibited: []string{"Marillion"},
			},
			Description: "directories NOT starting with M",
			Pattern:     "M*|",
			Scope:       enums.ScopeDirectory,
			Negate:      true,
		}),

		Entry(nil, &lab.FilterTE{
			DescribedTE: lab.DescribedTE{
				Given: "universal(undefined scope): glob ex filter",
			},
			NaviTE: lab.NaviTE{
				Relative:     "rock/PROGRESSIVE-ROCK/Marillion",
				Subscription: enums.SubscribeDirectories,
				ExpectedNoOf: lab.Quantities{
					Files:       0,
					Directories: 2,
				},
				Mandatory:  []string{"Marillion"},
				Prohibited: []string{"Fugazi"},
			},
			Description: "directories starting with M",
			Pattern:     "M*|",
		}),

		// === files =========================================================

		Entry(nil, &lab.FilterTE{
			DescribedTE: lab.DescribedTE{
				Given: "files(file scope): glob ex filter",
			},
			NaviTE: lab.NaviTE{
				Relative:     "rock/PROGRESSIVE-ROCK/Marillion",
				Subscription: enums.SubscribeFiles,
				ExpectedNoOf: lab.Quantities{
					Files:       16,
					Directories: 0,
				},
				Mandatory:  []string{"01 - Hotel Hobbies.flac"},
				Prohibited: []string{"cover-clutching-at-straws-jpg"},
			},
			Description: "items with 'flac' suffix",
			Pattern:     "*|flac",
			Scope:       enums.ScopeFile,
		}),

		Entry(nil, &lab.FilterTE{
			DescribedTE: lab.DescribedTE{
				Given: "files(any scope): glob ex filter, with dot extension",
			},
			NaviTE: lab.NaviTE{
				Relative:     "rock/PROGRESSIVE-ROCK/Marillion",
				Subscription: enums.SubscribeFiles,
				ExpectedNoOf: lab.Quantities{
					Files:       16,
					Directories: 0,
				},
				Mandatory:  []string{"01 - Hotel Hobbies.flac"},
				Prohibited: []string{"cover-clutching-at-straws-jpg"},
			},
			Description: "items with 'flac' suffix",
			Pattern:     "*|.flac",
			Scope:       enums.ScopeFile,
		}),

		Entry(nil, &lab.FilterTE{
			DescribedTE: lab.DescribedTE{
				Given: "files(file scope): glob ex filter, with multiple extensions",
			},
			NaviTE: lab.NaviTE{
				Relative:     "rock/PROGRESSIVE-ROCK/Marillion",
				Subscription: enums.SubscribeFiles,
				ExpectedNoOf: lab.Quantities{
					Files:       19,
					Directories: 0,
				},
				Mandatory:  []string{"front.jpg"},
				Prohibited: []string{"cover-clutching-at-straws-jpg"},
			},
			Description: "items with 'flac' suffix",
			Pattern:     "*|flac,jpg",
			Scope:       enums.ScopeFile,
		}),

		Entry(nil, &lab.FilterTE{
			DescribedTE: lab.DescribedTE{
				Given: "file(file scope): glob ex filter, without extension",
			},
			NaviTE: lab.NaviTE{
				Relative:     "rock/PROGRESSIVE-ROCK/Marillion",
				Subscription: enums.SubscribeFiles,
				ExpectedNoOf: lab.Quantities{
					Files:       3,
					Directories: 0,
				},
				Mandatory:  []string{"cover-clutching-at-straws-jpg"},
				Prohibited: []string{"01 - Hotel Hobbies.flac"},
			},
			Description: "items with 'flac' suffix",
			Pattern:     "*|",
			Scope:       enums.ScopeFile,
		}),

		Entry(nil, &lab.FilterTE{
			DescribedTE: lab.DescribedTE{
				Given: "file(file scope): glob ex filter (negate)",
			},
			NaviTE: lab.NaviTE{
				Relative:     "rock/PROGRESSIVE-ROCK/Marillion",
				Subscription: enums.SubscribeFiles,
				ExpectedNoOf: lab.Quantities{
					Files:       7,
					Directories: 0,
				},
				Mandatory:  []string{"cover-clutching-at-straws-jpg"},
				Prohibited: []string{"01 - Hotel Hobbies.flac"},
			},
			Description: "files without .flac suffix",
			Pattern:     "*|flac",
			Scope:       enums.ScopeFile,
			Negate:      true,
		}),

		Entry(nil, &lab.FilterTE{
			DescribedTE: lab.DescribedTE{
				Given: "file(undefined scope): glob ex filter",
			},
			NaviTE: lab.NaviTE{
				Relative:     "rock/PROGRESSIVE-ROCK/Marillion",
				Subscription: enums.SubscribeFiles,
				ExpectedNoOf: lab.Quantities{
					Files:       16,
					Directories: 0,
				},
				Mandatory:  []string{"01 - Hotel Hobbies.flac"},
				Prohibited: []string{"cover-clutching-at-straws-jpg"},
			},
			Description: "items with '.flac' suffix",
			Pattern:     "*|flac",
		}),

		Entry(nil, &lab.FilterTE{
			DescribedTE: lab.DescribedTE{
				Given: "file(any scope): glob ex filter, any extension",
			},
			NaviTE: lab.NaviTE{
				Relative:     "rock/PROGRESSIVE-ROCK/Marillion",
				Subscription: enums.SubscribeFiles,
				ExpectedNoOf: lab.Quantities{
					Files:       4,
					Directories: 0,
				},
				Mandatory:  []string{"cover-clutching-at-straws-jpg"},
				Prohibited: []string{"01 - Hotel Hobbies.flac"},
			},
			Description: "starts with c, any extension",
			Pattern:     "c*|*",
			Scope:       enums.ScopeAll,
		}),

		// === ifNotApplicable ===============================================

		Entry(nil, &lab.FilterTE{
			DescribedTE: lab.DescribedTE{
				Given: "universal(leaf scope): glob ex filter (ifNotApplicable=true)",
			},
			NaviTE: lab.NaviTE{
				Relative:     "rock/PROGRESSIVE-ROCK/Marillion",
				Subscription: enums.SubscribeUniversal,
				ExpectedNoOf: lab.Quantities{
					Files:       16,
					Directories: 5,
				},
				Mandatory:  []string{"Marillion"},
				Prohibited: []string{"cover-clutching-at-straws-jpg"},
			},
			Description:     "leaf items with 'flac' suffix",
			Pattern:         "*|flac",
			Scope:           enums.ScopeLeaf,
			IfNotApplicable: enums.TriStateBoolTrue,
		}),

		Entry(nil, &lab.FilterTE{
			DescribedTE: lab.DescribedTE{
				Given: "universal(leaf scope): glob ex filter (ifNotApplicable=false)",
			},
			NaviTE: lab.NaviTE{
				Relative:     "rock/PROGRESSIVE-ROCK/Marillion",
				Subscription: enums.SubscribeUniversal,
				ExpectedNoOf: lab.Quantities{
					Files:       16,
					Directories: 4,
				},
				Prohibited: []string{"Marillion"},
			},
			Description:     "items with '.flac' suffix",
			Pattern:         "*|flac",
			Scope:           enums.ScopeLeaf,
			IfNotApplicable: enums.TriStateBoolFalse,
		}),

		// === with-exclusion ================================================

		Entry(nil, &lab.FilterTE{
			DescribedTE: lab.DescribedTE{
				Given: "universal(any scope): glob ex filter with exclusion",
			},
			NaviTE: lab.NaviTE{
				Relative:     "rock/PROGRESSIVE-ROCK/Marillion",
				Subscription: enums.SubscribeFiles,
				ExpectedNoOf: lab.Quantities{
					Files:       12,
					Directories: 0,
				},
				Prohibited: []string{"01 - Hotel Hobbies.flac"},
			},
			Description: "files starting with 0, except 01 items and flac suffix",
			Pattern:     "0*/*01*|flac",
			Scope:       enums.ScopeFile,
		}),
	)
})
