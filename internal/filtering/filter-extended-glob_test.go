package filtering_test

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"testing/fstest"

	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok
	"github.com/snivilised/li18ngo"
	tv "github.com/snivilised/traverse"
	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/enums"
	lab "github.com/snivilised/traverse/internal/laboratory"
	"github.com/snivilised/traverse/internal/services"
	"github.com/snivilised/traverse/internal/third/lo"
	"github.com/snivilised/traverse/pref"
)

var _ = Describe("filtering", Ordered, func() {
	var (
		FS   fstest.MapFS
		root string
	)

	BeforeAll(func() {
		const (
			verbose = false
		)

		FS, root = lab.Musico(verbose,
			filepath.Join("MUSICO", "rock"),
		)
		Expect(root).NotTo(BeEmpty())
		Expect(li18ngo.Use()).To(Succeed())
	})

	BeforeEach(func() {
		services.Reset()
	})

	Context("comprehension", func() {
		When("universal: filtering with extended glob", func() {
			It("should: invoke for filtered nodes only", Label("example"),
				func(ctx SpecContext) {
					path := lab.Path(root, "RETRO-WAVE")
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
							Root:         path,
							Subscription: enums.SubscribeUniversal,
							Handler: func(node *core.Node) error {
								GinkgoWriter.Printf(
									"---> 🍯 EXAMPLE-EXTENDED-GLOB-FILTER-CALLBACK: '%v'\n", node.Path,
								)
								return nil
							},
							GetReadDirFS: func() fs.ReadDirFS {
								return FS
							},
							GetQueryStatusFS: func(_ fs.FS) fs.StatFS {
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

	DescribeTable("folders with files",
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

			path := lab.Path(root, entry.Relative)

			callback := func(node *core.Node) error {
				indicator := lo.Ternary(node.IsFolder(), "📁", "💠")
				GinkgoWriter.Printf(
					"===> %v Glob Filter(%v) source: '%v', item-name: '%v', item-scope(fs): '%v(%v)'\n",
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
					Root:         path,
					Subscription: entry.Subscription,
					Handler:      callback,
					GetReadDirFS: func() fs.ReadDirFS {
						return FS
					},
					GetQueryStatusFS: func(_ fs.FS) fs.StatFS {
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
					Files:   16,
					Folders: 5,
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
					Files:   16,
					Folders: 5,
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
					Files:   19,
					Folders: 5,
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
					Files:   3,
					Folders: 5,
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
					Files:   7,
					Folders: 5,
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
					Files:   16,
					Folders: 5,
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
					Files:   4,
					Folders: 1,
				},
				Mandatory:  []string{"cover-clutching-at-straws-jpg"},
				Prohibited: []string{"01 - Hotel Hobbies.flac"},
			},
			Description: "starts with c, any extension",
			Pattern:     "c*|*",
			Scope:       enums.ScopeAll,
		}),

		// === folders =======================================================

		Entry(nil, &lab.FilterTE{
			NaviTE: lab.NaviTE{
				Given:        "folders(any scope): extended glob filter",
				Relative:     "rock/PROGRESSIVE-ROCK/Marillion",
				Subscription: enums.SubscribeFolders,
				ExpectedNoOf: lab.Quantities{
					Files:   0,
					Folders: 2,
				},
				Mandatory:  []string{"Marillion"},
				Prohibited: []string{"Fugazi"},
			},
			Description: "folders starting with M",
			Pattern:     "M*|",
			Scope:       enums.ScopeFolder,
		}),

		Entry(nil, &lab.FilterTE{
			NaviTE: lab.NaviTE{
				Given:        "folders(folder scope): extended glob filter (negate)",
				Relative:     "rock/PROGRESSIVE-ROCK/Marillion",
				Subscription: enums.SubscribeFolders,
				ExpectedNoOf: lab.Quantities{
					Files:   0,
					Folders: 3,
				},
				Mandatory:  []string{"Fugazi"},
				Prohibited: []string{"Marillion"},
			},
			Description: "folders NOT starting with M",
			Pattern:     "M*|",
			Scope:       enums.ScopeFolder,
			Negate:      true,
		}),

		Entry(nil, &lab.FilterTE{
			NaviTE: lab.NaviTE{
				Given:        "universal(undefined scope): extended glob filter",
				Relative:     "rock/PROGRESSIVE-ROCK/Marillion",
				Subscription: enums.SubscribeFolders,
				ExpectedNoOf: lab.Quantities{
					Files:   0,
					Folders: 2,
				},
				Mandatory:  []string{"Marillion"},
				Prohibited: []string{"Fugazi"},
			},
			Description: "folders starting with M",
			Pattern:     "M*|",
		}),

		// === files =========================================================

		Entry(nil, &lab.FilterTE{
			NaviTE: lab.NaviTE{
				Given:        "files(file scope): extended glob filter",
				Relative:     "rock/PROGRESSIVE-ROCK/Marillion",
				Subscription: enums.SubscribeFiles,
				ExpectedNoOf: lab.Quantities{
					Files:   16,
					Folders: 0,
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
					Files:   16,
					Folders: 0,
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
					Files:   19,
					Folders: 0,
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
					Files:   3,
					Folders: 0,
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
					Files:   7,
					Folders: 0,
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
					Files:   16,
					Folders: 0,
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
					Files:   4,
					Folders: 0,
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
					Files:   16,
					Folders: 5,
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
					Files:   16,
					Folders: 4,
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
					Files:   12,
					Folders: 0,
				},
				Prohibited: []string{"01 - Hotel Hobbies.flac"},
			},
			Description: "files starting with 0, except 01 items and flac suffix",
			Pattern:     "0*/*01*|flac",
			Scope:       enums.ScopeFile,
		}),
	)
})
