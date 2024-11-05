package filtering_test

import (
	"fmt"
	"regexp/syntax"

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

		fS = hydra.Nuxx(verbose,
			lab.Static.RetroWave, "PROGRESSIVE-HOUSE",
		)
		Expect(li18ngo.Use()).To(Succeed())
	})

	BeforeEach(func() {
		services.Reset()
	})

	Context("comprehension", func() {
		When("files: filtering with regex", func() {
			It("should: invoke for filtered nodes only", Label("example"),
				func(ctx SpecContext) {
					path := lab.Static.RetroWave
					filterDefs := &pref.FilterOptions{
						Node: &core.FilterDef{
							Type:        enums.FilterTypeRegex,
							Description: "items that start with 'vinyl'",
							Pattern:     "^vinyl",
							Scope:       enums.ScopeAll,
						},
					}
					result, _ := tv.Walk().Configure().Extent(tv.Prime(
						&pref.Using{
							Head: pref.Head{
								Subscription: enums.SubscribeUniversal,
								Handler: func(servant tv.Servant) error {
									node := servant.Node()
									GinkgoWriter.Printf(
										"---> 🍯 EXAMPLE-REGEX-FILTER-CALLBACK: '%v'\n", node.Path,
									)
									return nil
								},
								GetForest: func(_ string) *core.Forest {
									return &core.Forest{
										T: fS,
										R: nef.NewTraverseABS(),
									}
								},
							},
							Tree: path,
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

	DescribeTable("regex",
		func(ctx SpecContext, entry *lab.FilterTE) {
			var (
				traverseFilter core.TraverseFilter
			)

			recall := make(lab.Recall)
			filterDefs := &pref.FilterOptions{
				Node: &core.FilterDef{
					Type:            enums.FilterTypeRegex,
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
					Expect(node).Should(MatchCurrentRegexFilter(traverseFilter))
				}

				recall[node.Extension.Name] = len(node.Children)
				return nil
			}
			result, err := tv.Walk().Configure().Extent(tv.Prime(
				&pref.Using{
					Head: pref.Head{
						Subscription: entry.Subscription,
						Handler:      callback,
						GetForest: func(_ string) *core.Forest {
							return &core.Forest{
								T: fS,
								R: nef.NewTraverseABS(),
							}
						},
					},
					Tree: path,
				},
				tv.WithOnBegin(lab.Begin("🛡️")),
				tv.WithOnEnd(lab.End("🏁")),

				tv.WithFilter(filterDefs),
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
		func(entry *lab.FilterTE) string {
			return fmt.Sprintf("🧪 ===> given: '%v'", entry.Given)
		},

		// === files =========================================================

		Entry(nil, &lab.FilterTE{
			NaviTE: lab.NaviTE{
				Given:        "files(any scope): regex filter",
				Relative:     lab.Static.RetroWave,
				Subscription: enums.SubscribeFiles,
				ExpectedNoOf: lab.Quantities{
					Files:       4,
					Directories: 0,
				},
			},
			Description: "items that start with 'vinyl'",
			Pattern:     "^vinyl",
			Scope:       enums.ScopeAll,
		}),

		Entry(nil, &lab.FilterTE{
			NaviTE: lab.NaviTE{
				Given:        "files(any scope): regex filter (negate)",
				Relative:     lab.Static.RetroWave,
				Subscription: enums.SubscribeFiles,
				ExpectedNoOf: lab.Quantities{
					Files:       10,
					Directories: 0,
				},
			},
			Description: "items that don't start with 'vinyl'",
			Pattern:     "^vinyl",
			Scope:       enums.ScopeAll,
			Negate:      true,
		}),

		Entry(nil, &lab.FilterTE{
			NaviTE: lab.NaviTE{
				Given:        "files(default to any scope): regex filter",
				Relative:     lab.Static.RetroWave,
				Subscription: enums.SubscribeFiles,
				ExpectedNoOf: lab.Quantities{
					Files:       4,
					Directories: 0,
				},
			},
			Description: "items that start with 'vinyl'",
			Pattern:     "^vinyl",
		}),

		// === directories ===================================================

		Entry(nil, &lab.FilterTE{
			NaviTE: lab.NaviTE{
				Given:        "directories(any scope): regex filter",
				Relative:     lab.Static.RetroWave,
				Subscription: enums.SubscribeDirectories,
				ExpectedNoOf: lab.Quantities{
					Files:       0,
					Directories: 2,
				},
			},
			Description: "items that start with 'C'",
			Pattern:     "^C",
			Scope:       enums.ScopeAll,
		}),

		Entry(nil, &lab.FilterTE{
			NaviTE: lab.NaviTE{
				Given:        "directories(any scope): regex filter (negate)",
				Relative:     lab.Static.RetroWave,
				Subscription: enums.SubscribeDirectories,
				ExpectedNoOf: lab.Quantities{
					Files:       0,
					Directories: 6,
				},
			},
			Description: "items that don't start with 'C'",
			Pattern:     "^C",
			Scope:       enums.ScopeAll,
			Negate:      true,
		}),

		Entry(nil, &lab.FilterTE{
			NaviTE: lab.NaviTE{
				Given:        "directories(undefined scope): regex filter",
				Relative:     lab.Static.RetroWave,
				Subscription: enums.SubscribeDirectories,
				ExpectedNoOf: lab.Quantities{
					Files:       0,
					Directories: 2,
				},
			},
			Description: "items that start with 'C'",
			Pattern:     "^C",
		}),

		// === ifNotApplicable ===============================================

		Entry(nil, &lab.FilterTE{
			NaviTE: lab.NaviTE{
				Given:        "directories(top): regex filter (ifNotApplicable=true)",
				Relative:     "PROGRESSIVE-HOUSE",
				Subscription: enums.SubscribeDirectories,
				ExpectedNoOf: lab.Quantities{
					Files:       0,
					Directories: 10,
				},
				Mandatory: []string{"PROGRESSIVE-HOUSE"},
			},
			Description:     "top items that contain 'HOUSE'",
			Pattern:         "HOUSE",
			Scope:           enums.ScopeTop,
			IfNotApplicable: enums.TriStateBoolTrue,
		}),

		Entry(nil, &lab.FilterTE{
			NaviTE: lab.NaviTE{
				Given:        "directories(top): regex filter (ifNotApplicable=false)",
				Relative:     ".",
				Subscription: enums.SubscribeDirectories,
				Mandatory:    []string{"PROGRESSIVE-HOUSE"},
				Prohibited:   []string{"Blue Amazon", "The Javelin"},
				ExpectedNoOf: lab.Quantities{
					Files:       0,
					Directories: 1,
				},
			},
			Description:     "top items that contain 'HOUSE'",
			Pattern:         "HOUSE",
			Scope:           enums.ScopeTop,
			IfNotApplicable: enums.TriStateBoolFalse,
		}),

		// === error =========================================================

		Entry(nil, &lab.FilterTE{
			NaviTE: lab.NaviTE{
				Given:        "files(any scope): regex filter",
				Relative:     lab.Static.RetroWave,
				Subscription: enums.SubscribeFiles,
				ExpectedErr: &syntax.Error{
					Code: "missing closing )",
					Expr: "(",
				},
			},
			Description: "items that start with 'vinyl'",
			Pattern:     "(",
			Scope:       enums.ScopeAll,
		}),
	)
})
