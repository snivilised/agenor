package filtering_test

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok
	"github.com/snivilised/li18ngo"
	"github.com/snivilised/nefilim/luna"
	tv "github.com/snivilised/traverse"
	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/enums"
	lab "github.com/snivilised/traverse/internal/laboratory"
	"github.com/snivilised/traverse/internal/services"
	"github.com/snivilised/traverse/internal/third/lo"
	"github.com/snivilised/traverse/pref"
	"github.com/snivilised/traverse/test/hydra"
)

var _ = Describe("filtering", Ordered, func() {
	var (
		fS *luna.MemFS
	)

	BeforeAll(func() {
		const (
			verbose = false
		)

		fS = hydra.Nuxx(verbose, "rock")
		Expect(li18ngo.Use()).To(Succeed())
	})

	BeforeEach(func() {
		services.Reset()
	})

	Context("comprehension", func() {
		When("universal: filtering with extended glob", func() {
			It("should: invoke for filtered nodes only", Label("example"),
				func(ctx SpecContext) {
					path := lab.Static.RetroWave
					filterDefs := &pref.FilterOptions{
						Node: &core.FilterDef{
							Type:        enums.FilterTypeExtendedGlob,
							Description: "nodes with 'flac' suffix",
							Pattern:     "*|flac",
							Scope:       enums.ScopeAll,
						},
					}
					result, _ := tv.Walk().Configure().Extent(tv.Prime(
						&tv.Using{
							Tree:         path,
							Subscription: enums.SubscribeUniversal,
							Handler: func(servant tv.Servant) error {
								node := servant.Node()
								GinkgoWriter.Printf(
									"---> 🍯 EXAMPLE-EXTENDED-GLOB-FILTER-CALLBACK: '%v'\n", node.Path,
								)
								return nil
							},
							GetTraverseFS: func(_ string) tv.TraverseFS {
								return fS
							},
						},
						tv.WithOnBegin(lab.Begin("🛡️")),
						tv.WithOnEnd(lab.End("🏁")),

						tv.WithFilter(filterDefs),
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

			recording := make(lab.RecordingMap)
			filterDefs := &pref.FilterOptions{
				Node: &core.FilterDef{
					Type:            enums.FilterTypeExtendedGlob,
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
			callback := func(servant tv.Servant) error {
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

				recording[node.Extension.Name] = len(node.Children)
				return nil
			}
			result, err := tv.Walk().Configure().Extent(tv.Prime(
				&tv.Using{
					Tree:         path,
					Subscription: entry.Subscription,
					Handler:      callback,
					GetTraverseFS: func(_ string) tv.TraverseFS {
						return fS
					},
				},
				tv.WithOnBegin(lab.Begin("🛡️")),
				tv.WithOnEnd(lab.End("🏁")),

				tv.WithFilter(filterDefs),
			)).Navigate(ctx)

			lab.AssertNavigation(&entry.NaviTE, &lab.TestOptions{
				FS:        fS,
				Recording: recording,
				Path:      path,
				Result:    result,
				Err:       err,
			})
		},

		func(entry *lab.FilterTE) string {
			return fmt.Sprintf("🧪 ===> given: '%v'", entry.Given)
		},

		// === universal =====================================================

		Entry(nil, &lab.FilterTE{
			NaviTE: lab.NaviTE{
				Given:        "universal(any scope): extended glob filter",
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
			NaviTE: lab.NaviTE{
				Given:        "universal(any scope): extended glob filter, with dot extension",
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
			NaviTE: lab.NaviTE{
				Given:        "universal(any scope): extended glob filter, with multiple extensions",
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
			NaviTE: lab.NaviTE{
				Given:        "universal(any scope): extended glob filter, without extension",
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
			NaviTE: lab.NaviTE{
				Given:        "universal(file scope): extended glob filter (negate)",
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
			NaviTE: lab.NaviTE{
				Given:        "universal(undefined scope): extended glob filter",
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
			NaviTE: lab.NaviTE{
				Given:        "universal(any scope): extended glob filter, any extension",
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
			NaviTE: lab.NaviTE{
				Given:        "directories(any scope): extended glob filter",
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
			NaviTE: lab.NaviTE{
				Given:        "directories(directory scope): extended glob filter (negate)",
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
			NaviTE: lab.NaviTE{
				Given:        "universal(undefined scope): extended glob filter",
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
			NaviTE: lab.NaviTE{
				Given:        "files(file scope): extended glob filter",
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
			NaviTE: lab.NaviTE{
				Given:        "files(any scope): extended glob filter, with dot extension",
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
			NaviTE: lab.NaviTE{
				Given:        "files(file scope): extended glob filter, with multiple extensions",
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
			NaviTE: lab.NaviTE{
				Given:        "file(file scope): extended glob filter, without extension",
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
			NaviTE: lab.NaviTE{
				Given:        "file(file scope): extended glob filter (negate)",
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
			NaviTE: lab.NaviTE{
				Given:        "file(undefined scope): extended glob filter",
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
			NaviTE: lab.NaviTE{
				Given:        "file(any scope): extended glob filter, any extension",
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
			NaviTE: lab.NaviTE{
				Given:        "universal(leaf scope): extended glob filter (ifNotApplicable=true)",
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
			NaviTE: lab.NaviTE{
				Given:        "universal(leaf scope): extended glob filter (ifNotApplicable=false)",
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
			NaviTE: lab.NaviTE{
				Given:        "universal(any scope): extended glob filter with exclusion",
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
