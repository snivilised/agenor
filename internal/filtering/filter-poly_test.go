package filtering_test

import (
	"fmt"
	"regexp/syntax"

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
					result, _ := tv.Walk().Configure().Extent(tv.Prime(
						&tv.Using{
							Tree:         path,
							Subscription: enums.SubscribeUniversal,
							Handler: func(servant tv.Servant) error {
								node := servant.Node()
								GinkgoWriter.Printf(
									"---> 🍯 EXAMPLE-POLY-FILTER-CALLBACK: '%v'\n", node.Path,
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

	DescribeTable("poly-filter",
		func(ctx SpecContext, entry *lab.PolyTE) {
			var (
				traverseFilter core.TraverseFilter
			)

			recording := make(lab.RecordingMap)
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
			callback := func(servant tv.Servant) error {
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
				FS:          fS,
				Recording:   recording,
				Path:        path,
				Result:      result,
				Err:         err,
				ExpectedErr: entry.ExpectedErr,
			})
		},
		func(entry *lab.PolyTE) string {
			return fmt.Sprintf("🧪 ===> given: '%v'", entry.Given)
		},

		// === universal(file:regex; directory:glob) =========================

		Entry(nil, &lab.PolyTE{
			NaviTE: lab.NaviTE{
				Given:        "poly - files:regex; directories:glob",
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
			NaviTE: lab.NaviTE{
				Given:        "poly - files:regex; directories:regex",
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
			NaviTE: lab.NaviTE{
				Given:        "poly - files:extended-glob; directories:glob",
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
				Type:        enums.FilterTypeExtendedGlob,
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
			NaviTE: lab.NaviTE{
				Given:        "poly - files:extended-glob; directories:regex",
				Relative:     lab.Static.RetroWave,
				Subscription: enums.SubscribeUniversal,
				ExpectedNoOf: lab.Quantities{
					Files:       3,
					Directories: 8,
				},
			},
			File: core.FilterDef{
				Type:        enums.FilterTypeExtendedGlob,
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
			NaviTE: lab.NaviTE{
				Given:        "poly - files:extended-glob; directories:extended-glob",
				Relative:     lab.Static.RetroWave,
				Subscription: enums.SubscribeUniversal,
				ExpectedNoOf: lab.Quantities{
					Files:       3,
					Directories: 8,
				},
			},
			File: core.FilterDef{
				Type:        enums.FilterTypeExtendedGlob,
				Description: "files: txt files starting with vinyl",
				Pattern:     "vinyl*|txt",
				Scope:       enums.ScopeFile,
			},
			Directory: core.FilterDef{
				Type:        enums.FilterTypeExtendedGlob,
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
			NaviTE: lab.NaviTE{
				Given:        "poly(scopes omitted) - files:regex; directories:regex",
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
			NaviTE: lab.NaviTE{
				Given:        "poly(subscribe:files)",
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
			NaviTE: lab.NaviTE{
				Given:        "invalid poly: constituent is also poly",
				Should:       "fail",
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
			NaviTE: lab.NaviTE{
				Given:        "poly - files:regex; directories:glob",
				Should:       "fail",
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
