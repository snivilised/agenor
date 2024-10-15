package filtering_test

import (
	"fmt"
	"io/fs"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok

	"github.com/snivilised/li18ngo"
	nef "github.com/snivilised/nefilim"
	tv "github.com/snivilised/traverse"
	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/enums"
	lab "github.com/snivilised/traverse/internal/laboratory"
	"github.com/snivilised/traverse/internal/services"
	"github.com/snivilised/traverse/internal/third/lo"
	"github.com/snivilised/traverse/pref"
)

var _ = Describe("NavigatorFilterGlob", Ordered, func() {
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
		When("universal: filtering with glob", func() {
			It("should: invoke for filtered nodes only", Label("example"),
				func(ctx SpecContext) {
					path := lab.Path(root, "RETRO-WAVE")
					filterDefs := &pref.FilterOptions{
						Node: &core.FilterDef{
							Type:        enums.FilterTypeGlob,
							Description: "items with '.flac' suffix",
							Pattern:     "*.flac",
							Scope:       enums.ScopeAll,
						},
					}
					result, _ := tv.Walk().Configure().Extent(tv.Prime(
						&tv.Using{
							Root:         path,
							Subscription: enums.SubscribeUniversal,
							Handler: func(node *core.Node) error {
								GinkgoWriter.Printf(
									"---> 🍯 EXAMPLE-GLOB-FILTER-CALLBACK: '%v'\n", node.Path,
								)
								return nil
							},
							GetTraverseFS: func(_ string) nef.TraverseFS {
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

	DescribeTable("glob-filter",
		func(ctx SpecContext, entry *lab.FilterTE) {
			var (
				traverseFilter core.TraverseFilter
			)

			recording := make(lab.RecordingMap)
			filterDefs := &pref.FilterOptions{
				Node: &core.FilterDef{
					Type:            enums.FilterTypeGlob,
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
					GetTraverseFS: func(_ string) nef.TraverseFS {
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

		Entry(nil, &lab.FilterTE{
			NaviTE: lab.NaviTE{
				Given:        "universal(any scope): glob filter",
				Relative:     "RETRO-WAVE",
				Subscription: enums.SubscribeUniversal,
				ExpectedNoOf: lab.Quantities{
					Files:   8,
					Folders: 0,
				},
			},
			Description: "items with '.flac' suffix",
			Pattern:     "*.flac",
			Scope:       enums.ScopeAll,
		}),

		Entry(nil, &lab.FilterTE{
			NaviTE: lab.NaviTE{
				Given:        "universal(any scope): glob filter (negate)",
				Relative:     "RETRO-WAVE",
				Subscription: enums.SubscribeUniversal,
				ExpectedNoOf: lab.Quantities{
					Files:   6,
					Folders: 8,
				},
			},
			Description: "items without .flac suffix",
			Pattern:     "*.flac",
			Scope:       enums.ScopeAll,
			Negate:      true,
		}),

		Entry(nil, &lab.FilterTE{
			NaviTE: lab.NaviTE{
				Given:        "universal(undefined scope): glob filter",
				Relative:     "RETRO-WAVE",
				Subscription: enums.SubscribeUniversal,
				ExpectedNoOf: lab.Quantities{
					Files:   8,
					Folders: 0,
				},
			},
			Description: "items with '.flac' suffix",
			Pattern:     "*.flac",
		}),

		// === ifNotApplicable ===============================================

		Entry(nil, &lab.FilterTE{
			NaviTE: lab.NaviTE{
				Given:        "universal(any scope): glob filter (ifNotApplicable=true)",
				Relative:     "RETRO-WAVE",
				Subscription: enums.SubscribeUniversal,
				ExpectedNoOf: lab.Quantities{
					Files:   8,
					Folders: 4,
				},
				Mandatory: []string{"A1 - Can You Kiss Me First.flac"},
			},
			Description:     "items with '.flac' suffix",
			Pattern:         "*.flac",
			Scope:           enums.ScopeLeaf,
			IfNotApplicable: enums.TriStateBoolTrue,
		}),

		Entry(nil, &lab.FilterTE{
			NaviTE: lab.NaviTE{
				Given:        "universal(leaf scope): glob filter (ifNotApplicable=false)",
				Relative:     "RETRO-WAVE",
				Subscription: enums.SubscribeUniversal,
				ExpectedNoOf: lab.Quantities{
					Files:   8,
					Folders: 0,
				},
				Mandatory:  []string{"A1 - Can You Kiss Me First.flac"},
				Prohibited: []string{"vinyl-info.teenage-color"},
			},
			Description:     "items with '.flac' suffix",
			Pattern:         "*.flac",
			Scope:           enums.ScopeLeaf,
			IfNotApplicable: enums.TriStateBoolFalse,
		}),
	)
})
