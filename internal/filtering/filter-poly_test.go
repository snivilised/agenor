package filtering_test

import (
	"regexp/syntax"

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

var _ = Describe("feature", Ordered, func() {
	var (
		fS *luna.MemFS
	)

	BeforeAll(func() {
		const (
			verbose = false
		)

		fS = hanno.Nuxx(verbose, lab.Static.RetroWave)
		Expect(li18ngo.Use()).To(Succeed())
	})

	BeforeEach(func() {
		services.Reset()
	})

	Context("comprehension", func() {
		When("universal: filtering with poly-filter", func() {
			It("🧪 should: invoke for filtered nodes only", Label("example"),
				func(ctx SpecContext) {
					path := lab.Static.RetroWave
					filterDefs := &pref.FilterOptions{
						Node: &core.FilterDef{
							Type: enums.FilterTypePoly,
							Poly: &core.PolyFilterDef{
								File: core.FilterDef{
									Type:        enums.FilterTypeRegex,
									Description: "files: starts with vinyl",
									Pattern:     "^vinyl",
									Scope:       enums.ScopeFile,
								},
								Directory: core.FilterDef{
									Type:        enums.FilterTypeGlob,
									Description: "directories: contains i (case sensitive)",
									Pattern:     "*i*",
									Scope:       enums.ScopeDirectory | enums.ScopeLeaf,
								},
							},
						},
					}
					result, _ := age.Walk().Configure().Extent(age.Prime(
						&pref.Using{
							Subscription: enums.SubscribeUniversal,
							Head: pref.Head{
								Handler: func(servant age.Servant) error {
									node := servant.Node()
									GinkgoWriter.Printf(
										"---> 🍯 EXAMPLE-POLY-FILTER-CALLBACK: '%v'\n", node.Path,
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

	DescribeTable("poly-filter",
		func(ctx SpecContext, entry *lab.PolyTE) {
			var (
				traverseFilter core.TraverseFilter
			)

			recall := make(lab.Recall)
			filterDefs := &pref.FilterOptions{
				Node: &core.FilterDef{
					Type: enums.FilterTypePoly,
					Poly: &core.PolyFilterDef{
						File:      entry.File,
						Directory: entry.Directory,
					},
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
					"===> %v Poly Filter(%v) source: '%v', node-name: '%v', node-scope(fs): '%v(%v)'\n",
					indicator,
					traverseFilter.Description(),
					traverseFilter.Source(),
					node.Extension.Name,
					node.Extension.Scope,
					traverseFilter.Scope(),
				)
				if lo.Contains(entry.Mandatory, node.Extension.Name) {
					Expect(node).Should(MatchCurrentGlobFilter(traverseFilter))
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
				FS:          fS,
				Recording:   recall,
				Path:        path,
				Result:      result,
				Err:         err,
				ExpectedErr: entry.ExpectedErr,
			})
		},
		lab.FormatPolyFilterTestDescription,

		// === universal(file:regex; directory:glob) =========================

		Entry(nil, &lab.PolyTE{
			DescribedTE: lab.DescribedTE{
				Given: "poly - files:regex; directories:glob",
			},
			NaviTE: lab.NaviTE{
				Relative:     lab.Static.RetroWave,
				Subscription: enums.SubscribeUniversal,
				ExpectedNoOf: lab.Quantities{
					// file is 2 not 3 because *i* is case sensitive so Innerworld is not a match
					// The next(not this one) regex test case, fixes this because directory regex has better
					// control over case sensitivity
					Files:       2,
					Directories: 8,
				},
			},
			File: core.FilterDef{
				Type:        enums.FilterTypeRegex,
				Description: "files: starts with vinyl",
				Pattern:     "^vinyl",
				Scope:       enums.ScopeFile,
			},
			Directory: core.FilterDef{
				Type:        enums.FilterTypeGlob,
				Description: "directories: contains i (case sensitive)",
				Pattern:     "*i*",
				Scope:       enums.ScopeDirectory | enums.ScopeLeaf,
			},
		}),

		// === universal(file:regex; directory:regex) ========================

		Entry(nil, &lab.PolyTE{
			DescribedTE: lab.DescribedTE{
				Given: "poly - files:regex; directories:regex",
			},
			NaviTE: lab.NaviTE{
				Relative:     lab.Static.RetroWave,
				Subscription: enums.SubscribeUniversal,
				ExpectedNoOf: lab.Quantities{
					Files:       3,
					Directories: 8,
				},
			},
			File: core.FilterDef{
				Type:        enums.FilterTypeRegex,
				Description: "files: starts with vinyl",
				Pattern:     "^vinyl",
				Scope:       enums.ScopeFile,
			},
			Directory: core.FilterDef{
				Type:        enums.FilterTypeRegex,
				Description: "directories: contains i (case insensitive)",
				Pattern:     "[iI]",
				Scope:       enums.ScopeDirectory | enums.ScopeLeaf,
			},
		}),

		// === universal(file:extended-glob; directory:glob) =================

		Entry(nil, &lab.PolyTE{
			DescribedTE: lab.DescribedTE{
				Given: "poly - files:extended-glob; directories:glob",
			},
			NaviTE: lab.NaviTE{
				Relative:     lab.Static.RetroWave,
				Subscription: enums.SubscribeUniversal,
				ExpectedNoOf: lab.Quantities{
					// file is 2 not 3 because *i* is case sensitive so Innerworld is not a match
					// The next 2 tests regex/extended-glob test case, fixes this because they
					// have better control over case sensitivity
					//
					Files:       2,
					Directories: 8,
				},
			},
			File: core.FilterDef{
				Type:        enums.FilterTypeGlobEx,
				Description: "files: txt files starting with vinyl",
				Pattern:     "vinyl*|txt",
				Scope:       enums.ScopeFile,
			},
			Directory: core.FilterDef{
				Type:        enums.FilterTypeGlob,
				Description: "directories: contains i (case sensitive)",
				Pattern:     "*i*",
				Scope:       enums.ScopeDirectory | enums.ScopeLeaf,
			},
		}),

		Entry(nil, &lab.PolyTE{
			DescribedTE: lab.DescribedTE{
				Given: "poly - files:extended-glob; directories:regex",
			},
			NaviTE: lab.NaviTE{
				Relative:     lab.Static.RetroWave,
				Subscription: enums.SubscribeUniversal,
				ExpectedNoOf: lab.Quantities{
					Files:       3,
					Directories: 8,
				},
			},
			File: core.FilterDef{
				Type:        enums.FilterTypeGlobEx,
				Description: "files: txt files starting with vinyl",
				Pattern:     "vinyl*|txt",
				Scope:       enums.ScopeFile,
			},
			Directory: core.FilterDef{
				Type:        enums.FilterTypeRegex,
				Description: "directories: contains i (case sensitive)",
				Pattern:     "[iI]",
				Scope:       enums.ScopeDirectory | enums.ScopeLeaf,
			},
		}),

		Entry(nil, &lab.PolyTE{
			DescribedTE: lab.DescribedTE{
				Given: "poly - files:extended-glob; directories:extended-glob",
			},
			NaviTE: lab.NaviTE{
				Relative:     lab.Static.RetroWave,
				Subscription: enums.SubscribeUniversal,
				ExpectedNoOf: lab.Quantities{
					Files:       3,
					Directories: 8,
				},
			},
			File: core.FilterDef{
				Type:        enums.FilterTypeGlobEx,
				Description: "files: txt files starting with vinyl",
				Pattern:     "vinyl*|txt",
				Scope:       enums.ScopeFile,
			},
			Directory: core.FilterDef{
				Type:        enums.FilterTypeGlobEx,
				Description: "directories: contains i (case sensitive)",
				Pattern:     "*i*|",
				Scope:       enums.ScopeDirectory | enums.ScopeLeaf,
			},
		}),

		// For the poly filter, the file/directory scopes must be set correctly, but because
		// they can be set automatically, the client is not forced to set them. This test
		// checks that when the file/directory scopes are not set, then poly filtering still works
		// properly.
		Entry(nil, &lab.PolyTE{
			DescribedTE: lab.DescribedTE{
				Given: "poly(scopes omitted) - files:regex; directories:regex",
			},
			NaviTE: lab.NaviTE{
				Relative:     lab.Static.RetroWave,
				Subscription: enums.SubscribeUniversal,
				ExpectedNoOf: lab.Quantities{
					Files:       3,
					Directories: 8,
				},
			},
			File: core.FilterDef{
				Type:        enums.FilterTypeRegex,
				Description: "files: starts with vinyl",
				Pattern:     "^vinyl",
				// file scope omitted
			},
			Directory: core.FilterDef{
				Type:        enums.FilterTypeRegex,
				Description: "directories: contains i (case insensitive)",
				Pattern:     "[iI]",
				Scope:       enums.ScopeLeaf, // directory scope omitted
			},
		}),

		// === files (file:regex; directory:regex) ==============================

		Entry(nil, &lab.PolyTE{
			DescribedTE: lab.DescribedTE{
				Given: "poly(subscribe:files)",
			},
			NaviTE: lab.NaviTE{
				Relative:     lab.Static.RetroWave,
				Subscription: enums.SubscribeFiles,
				ExpectedNoOf: lab.Quantities{
					Files:       3,
					Directories: 0,
				},
			},
			File: core.FilterDef{
				Type:        enums.FilterTypeRegex,
				Description: "files: starts with vinyl",
				Pattern:     "^vinyl",
			},
			Directory: core.FilterDef{
				Type:        enums.FilterTypeRegex,
				Description: "directories: contains i",
				Pattern:     "[iI]",
				Scope:       enums.ScopeLeaf,
			},
		}),

		// === errors ========================================================

		Entry(nil, &lab.PolyTE{
			DescribedTE: lab.DescribedTE{
				Given:  "invalid poly: constituent is also poly",
				Should: "fail",
			},
			NaviTE: lab.NaviTE{
				Relative:     lab.Static.RetroWave,
				Subscription: enums.SubscribeFiles,
				ExpectedErr:  locale.ErrPolyFilterIsInvalid,
			},
			File: core.FilterDef{
				Type:        enums.FilterTypePoly,
				Description: "files: constituent is poly",
			},
			Directory: core.FilterDef{
				Type:        enums.FilterTypePoly,
				Description: "directories: constituent is poly",
			},
		}),

		Entry(nil, &lab.PolyTE{
			DescribedTE: lab.DescribedTE{
				Given:  "poly - files:regex; directories:glob",
				Should: "fail",
			},
			NaviTE: lab.NaviTE{
				Relative:     lab.Static.RetroWave,
				Subscription: enums.SubscribeUniversal,
				ExpectedErr: &syntax.Error{
					Code: "missing closing )",
					Expr: "(",
				},
			},
			File: core.FilterDef{
				Type:        enums.FilterTypeRegex,
				Description: "files: starts with vinyl",
				Pattern:     "(",
				Scope:       enums.ScopeFile,
			},
			Directory: core.FilterDef{
				Type:        enums.FilterTypeGlob,
				Description: "directories: contains i (case sensitive)",
				Pattern:     "*i*",
				Scope:       enums.ScopeDirectory | enums.ScopeLeaf,
			},
		}),
	)
})
