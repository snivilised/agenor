package filtering_test

import (
	"fmt"
	"io/fs"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok

	"github.com/snivilised/li18ngo"
	tv "github.com/snivilised/traverse"
	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/enums"
	lab "github.com/snivilised/traverse/internal/laboratory"
	"github.com/snivilised/traverse/internal/services"
	"github.com/snivilised/traverse/internal/third/lo"
	"github.com/snivilised/traverse/lfs"
	"github.com/snivilised/traverse/pref"
)

var _ = Describe("feature", Ordered, func() {
	var (
		FS   *lab.TestTraverseFS
		root string
	)

	BeforeAll(func() {
		const (
			verbose = false
		)

		FS, root = lab.Musico(verbose,
			filepath.Join("MUSICO", "RETRO-WAVE"),
		)
		Expect(root).NotTo(BeEmpty())
		Expect(li18ngo.Use()).To(Succeed())
	})

	BeforeEach(func() {
		services.Reset()
	})

	Context("comprehension", func() {
		When("universal: filtering with poly-filter", func() {
			It("should: invoke for filtered nodes only", Label("example"),
				func(ctx SpecContext) {
					path := lab.Path(root, "RETRO-WAVE")
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
								Folder: core.FilterDef{
									Type:        enums.FilterTypeGlob,
									Description: "folders: contains i (case sensitive)",
									Pattern:     "*i*",
									Scope:       enums.ScopeFolder | enums.ScopeLeaf,
								},
							},
						},
					}
					result, _ := tv.Walk().Configure().Extent(tv.Prime(
						&tv.Using{
							Root:         path,
							Subscription: enums.SubscribeUniversal,
							Handler: func(node *core.Node) error {
								GinkgoWriter.Printf(
									"---> 🍯 EXAMPLE-POLY-FILTER-CALLBACK: '%v'\n", node.Path,
								)
								return nil
							},
							GetTraverseFS: func(_ string) lfs.TraverseFS {
								return FS
							},
						},
						tv.WithFilter(filterDefs),
						tv.WithHookQueryStatus(
							func(qsys fs.StatFS, path string) (fs.FileInfo, error) {
								return qsys.Stat(lab.TrimRoot(path))
							},
						),
						tv.WithHookReadDirectory(
							func(rfs fs.ReadDirFS, dirname string) ([]fs.DirEntry, error) {
								return rfs.ReadDir(lab.TrimRoot(dirname))
							},
						),
					)).Navigate(ctx)

					GinkgoWriter.Printf("===> 🍭 invoked '%v' folders, '%v' files.\n",
						result.Metrics().Count(enums.MetricNoFoldersInvoked),
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
						File:   entry.File,
						Folder: entry.Folder,
					},
				},
				Sink: func(reply pref.FilterReply) {
					traverseFilter = reply.Node
				},
			}

			path := lab.Path(root, entry.Relative)

			callback := func(node *core.Node) error {
				indicator := lo.Ternary(node.IsFolder(), "📁", "💠")
				GinkgoWriter.Printf(
					"===> %v Poly Filter(%v) source: '%v', item-name: '%v', item-scope(fs): '%v(%v)'\n",
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
					Root:         path,
					Subscription: entry.Subscription,
					Handler:      callback,
					GetTraverseFS: func(_ string) lfs.TraverseFS {
						return FS
					},
				},
				tv.WithFilter(filterDefs),
				tv.WithHookQueryStatus(
					func(qsys fs.StatFS, path string) (fs.FileInfo, error) {
						return qsys.Stat(lab.TrimRoot(path))
					},
				),
				tv.WithHookReadDirectory(
					func(rfs fs.ReadDirFS, dirname string) ([]fs.DirEntry, error) {
						return rfs.ReadDir(lab.TrimRoot(dirname))
					},
				),
			)).Navigate(ctx)

			lab.AssertNavigation(&entry.NaviTE, &lab.TestOptions{
				FS:        FS,
				Recording: recording,
				Path:      path,
				Result:    result,
				Err:       err,
			})
		},
		func(entry *lab.PolyTE) string {
			return fmt.Sprintf("🧪 ===> given: '%v'", entry.Given)
		},

		// === universal(file:regex; folder:glob) ============================

		Entry(nil, &lab.PolyTE{
			NaviTE: lab.NaviTE{
				Given:        "poly - files:regex; folders:glob",
				Relative:     "RETRO-WAVE",
				Subscription: enums.SubscribeUniversal,
				ExpectedNoOf: lab.Quantities{
					// file is 2 not 3 because *i* is case sensitive so Innerworld is not a match
					// The next(not this one) regex test case, fixes this because folder regex has better
					// control over case sensitivity
					Files:   2,
					Folders: 8,
				},
			},
			File: core.FilterDef{
				Type:        enums.FilterTypeRegex,
				Description: "files: starts with vinyl",
				Pattern:     "^vinyl",
				Scope:       enums.ScopeFile,
			},
			Folder: core.FilterDef{
				Type:        enums.FilterTypeGlob,
				Description: "folders: contains i (case sensitive)",
				Pattern:     "*i*",
				Scope:       enums.ScopeFolder | enums.ScopeLeaf,
			},
		}),

		// === universal(file:regex; folder:regex) ===========================

		Entry(nil, &lab.PolyTE{
			NaviTE: lab.NaviTE{
				Given:        "poly - files:regex; folders:regex",
				Relative:     "RETRO-WAVE",
				Subscription: enums.SubscribeUniversal,
				ExpectedNoOf: lab.Quantities{
					Files:   3,
					Folders: 8,
				},
			},
			File: core.FilterDef{
				Type:        enums.FilterTypeRegex,
				Description: "files: starts with vinyl",
				Pattern:     "^vinyl",
				Scope:       enums.ScopeFile,
			},
			Folder: core.FilterDef{
				Type:        enums.FilterTypeRegex,
				Description: "folders: contains i (case insensitive)",
				Pattern:     "[iI]",
				Scope:       enums.ScopeFolder | enums.ScopeLeaf,
			},
		}),

		// === universal(file:extended-glob; folder:glob) ====================

		Entry(nil, &lab.PolyTE{
			NaviTE: lab.NaviTE{
				Given:        "poly - files:extended-glob; folders:glob",
				Relative:     "RETRO-WAVE",
				Subscription: enums.SubscribeUniversal,
				ExpectedNoOf: lab.Quantities{
					// file is 2 not 3 because *i* is case sensitive so Innerworld is not a match
					// The next 2 tests regex/extended-glob test case, fixes this because they
					// have better control over case sensitivity
					//
					Files:   2,
					Folders: 8,
				},
			},
			File: core.FilterDef{
				Type:        enums.FilterTypeExtendedGlob,
				Description: "files: txt files starting with vinyl",
				Pattern:     "vinyl*|txt",
				Scope:       enums.ScopeFile,
			},
			Folder: core.FilterDef{
				Type:        enums.FilterTypeGlob,
				Description: "folders: contains i (case sensitive)",
				Pattern:     "*i*",
				Scope:       enums.ScopeFolder | enums.ScopeLeaf,
			},
		}),

		Entry(nil, &lab.PolyTE{
			NaviTE: lab.NaviTE{
				Given:        "poly - files:extended-glob; folders:regex",
				Relative:     "RETRO-WAVE",
				Subscription: enums.SubscribeUniversal,
				ExpectedNoOf: lab.Quantities{
					Files:   3,
					Folders: 8,
				},
			},
			File: core.FilterDef{
				Type:        enums.FilterTypeExtendedGlob,
				Description: "files: txt files starting with vinyl",
				Pattern:     "vinyl*|txt",
				Scope:       enums.ScopeFile,
			},
			Folder: core.FilterDef{
				Type:        enums.FilterTypeRegex,
				Description: "folders: contains i (case sensitive)",
				Pattern:     "[iI]",
				Scope:       enums.ScopeFolder | enums.ScopeLeaf,
			},
		}),

		Entry(nil, &lab.PolyTE{
			NaviTE: lab.NaviTE{
				Given:        "poly - files:extended-glob; folders:extended-glob",
				Relative:     "RETRO-WAVE",
				Subscription: enums.SubscribeUniversal,
				ExpectedNoOf: lab.Quantities{
					Files:   3,
					Folders: 8,
				},
			},
			File: core.FilterDef{
				Type:        enums.FilterTypeExtendedGlob,
				Description: "files: txt files starting with vinyl",
				Pattern:     "vinyl*|txt",
				Scope:       enums.ScopeFile,
			},
			Folder: core.FilterDef{
				Type:        enums.FilterTypeExtendedGlob,
				Description: "folders: contains i (case sensitive)",
				Pattern:     "*i*|",
				Scope:       enums.ScopeFolder | enums.ScopeLeaf,
			},
		}),

		// For the poly filter, the file/folder scopes must be set correctly, but because
		// they can be set automatically, the client is not forced to set them. This test
		// checks that when the file/folder scopes are not set, then poly filtering still works
		// properly.
		Entry(nil, &lab.PolyTE{
			NaviTE: lab.NaviTE{
				Given:        "poly(scopes omitted) - files:regex; folders:regex",
				Relative:     "RETRO-WAVE",
				Subscription: enums.SubscribeUniversal,
				ExpectedNoOf: lab.Quantities{
					Files:   3,
					Folders: 8,
				},
			},
			File: core.FilterDef{
				Type:        enums.FilterTypeRegex,
				Description: "files: starts with vinyl",
				Pattern:     "^vinyl",
				// file scope omitted
			},
			Folder: core.FilterDef{
				Type:        enums.FilterTypeRegex,
				Description: "folders: contains i (case insensitive)",
				Pattern:     "[iI]",
				Scope:       enums.ScopeLeaf, // folder scope omitted
			},
		}),

		// === files (file:regex; folder:regex) ==============================

		Entry(nil, &lab.PolyTE{
			NaviTE: lab.NaviTE{
				Given:        "poly(subscribe:files)",
				Relative:     "RETRO-WAVE",
				Subscription: enums.SubscribeFiles,
				ExpectedNoOf: lab.Quantities{
					Files:   3,
					Folders: 0,
				},
			},
			File: core.FilterDef{
				Type:        enums.FilterTypeRegex,
				Description: "files: starts with vinyl",
				Pattern:     "^vinyl",
			},
			Folder: core.FilterDef{
				Type:        enums.FilterTypeRegex,
				Description: "folders: contains i",
				Pattern:     "[iI]",
				Scope:       enums.ScopeLeaf,
			},
		}),
	)
})
