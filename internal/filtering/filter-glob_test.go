package filtering_test

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"testing/fstest"

	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok

	tv "github.com/snivilised/traverse"
	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/enums"
	"github.com/snivilised/traverse/internal/helpers"
	"github.com/snivilised/traverse/internal/services"
	"github.com/snivilised/traverse/internal/third/lo"
	"github.com/snivilised/traverse/pref"
)

var _ = Describe("NavigatorFilterGlob", Ordered, func() {
	var (
		FS   fstest.MapFS
		root string
	)

	BeforeAll(func() {
		const (
			verbose = false
		)

		FS, root = helpers.Musico(verbose,
			filepath.Join("MUSICO", "RETRO-WAVE"),
		)
		Expect(root).NotTo(BeEmpty())
	})

	BeforeEach(func() {
		services.Reset()
	})

	DescribeTable("glob-filter",
		func(ctx SpecContext, entry *helpers.FilterTE) {
			var (
				traverseFilter core.TraverseFilter
			)

			recording := make(helpers.RecordingMap)
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

			path := helpers.Path(root, entry.Relative)

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
						return qsys.Stat(helpers.TrimRoot(path))
					},
				),
				tv.WithHookReadDirectory(
					func(rfs fs.ReadDirFS, dirname string) ([]fs.DirEntry, error) {
						return rfs.ReadDir(helpers.TrimRoot(dirname))
					},
				),
			)).Navigate(ctx)

			helpers.AssertNavigation(&entry.NaviTE, &helpers.TestOptions{
				FS:        FS,
				Recording: recording,
				Path:      path,
				Result:    result,
				Err:       err,
			})
		},
		func(entry *helpers.FilterTE) string {
			return fmt.Sprintf("🧪 ===> given: '%v'", entry.Given)
		},

		Entry(nil, &helpers.FilterTE{
			NaviTE: helpers.NaviTE{
				Given:        "universal(any scope): glob filter",
				Relative:     "RETRO-WAVE",
				Subscription: enums.SubscribeUniversal,
				ExpectedNoOf: helpers.Quantities{
					Files:   8,
					Folders: 0,
				},
			},
			Description: "items with '.flac' suffix",
			Pattern:     "*.flac",
			Scope:       enums.ScopeAll,
		}),

		Entry(nil, &helpers.FilterTE{
			NaviTE: helpers.NaviTE{
				Given:        "universal(any scope): glob filter (negate)",
				Relative:     "RETRO-WAVE",
				Subscription: enums.SubscribeUniversal,
				ExpectedNoOf: helpers.Quantities{
					Files:   6,
					Folders: 8,
				},
			},
			Description: "items without .flac suffix",
			Pattern:     "*.flac",
			Scope:       enums.ScopeAll,
			Negate:      true,
		}),

		Entry(nil, &helpers.FilterTE{
			NaviTE: helpers.NaviTE{
				Given:        "universal(undefined scope): glob filter",
				Relative:     "RETRO-WAVE",
				Subscription: enums.SubscribeUniversal,
				ExpectedNoOf: helpers.Quantities{
					Files:   8,
					Folders: 0,
				},
			},
			Description: "items with '.flac' suffix",
			Pattern:     "*.flac",
		}),

		// === ifNotApplicable ===============================================

		Entry(nil, &helpers.FilterTE{
			NaviTE: helpers.NaviTE{
				Given:        "universal(any scope): glob filter (ifNotApplicable=true)",
				Relative:     "RETRO-WAVE",
				Subscription: enums.SubscribeUniversal,
				ExpectedNoOf: helpers.Quantities{
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

		Entry(nil, &helpers.FilterTE{
			NaviTE: helpers.NaviTE{
				Given:        "universal(leaf scope): glob filter (ifNotApplicable=false)",
				Relative:     "RETRO-WAVE",
				Subscription: enums.SubscribeUniversal,
				ExpectedNoOf: helpers.Quantities{
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
