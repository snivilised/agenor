package kernel_test

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
		vfs  fstest.MapFS
		root string
	)

	BeforeAll(func() {
		const (
			verbose = true
		)

		vfs, root = helpers.Musico(verbose,
			filepath.Join("MUSICO", "RETRO-WAVE"),
		)
		Expect(root).NotTo(BeEmpty())
	})

	BeforeEach(func() {
		services.Reset()
	})

	DescribeTable("glob-filter",
		func(ctx SpecContext, entry *filterTE) {
			var (
				traverseFilter core.TraverseFilter
			)

			recording := make(recordingMap)
			filterDefs := &pref.FilterOptions{
				Node: &core.FilterDef{
					Type:            enums.FilterTypeGlob,
					Description:     entry.description,
					Pattern:         entry.pattern,
					Scope:           entry.scope,
					Negate:          entry.negate,
					IfNotApplicable: entry.ifNotApplicable,
				},
				Sink: func(reply pref.FilterReply) {
					traverseFilter = reply.Node
				},
			}

			path := helpers.Path(root, entry.relative)

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
				if lo.Contains(entry.mandatory, node.Extension.Name) {
					Expect(node).Should(MatchCurrentGlobFilter(traverseFilter))
				}

				recording[node.Extension.Name] = len(node.Children)
				return nil
			}
			result, err := tv.Walk().Configure().Extent(tv.Prime(
				&tv.Using{
					Root:         path,
					Subscription: entry.subscription,
					Handler:      callback,
					GetReadDirFS: func() fs.ReadDirFS {
						return vfs
					},
					GetQueryStatusFS: func(_ fs.FS) fs.StatFS {
						return vfs
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

			assertNavigation(&entry.naviTE, &testOptions{
				vfs:       vfs,
				recording: recording,
				path:      path,
				result:    result,
				err:       err,
			})
		},
		func(entry *filterTE) string {
			return fmt.Sprintf("🧪 ===> given: '%v'", entry.given)
		},

		Entry(nil, &filterTE{
			naviTE: naviTE{
				given:        "universal(any scope): glob filter",
				relative:     "RETRO-WAVE",
				subscription: enums.SubscribeUniversal,
				expectedNoOf: quantities{
					files:   8,
					folders: 0,
				},
			},
			description: "items with '.flac' suffix",
			pattern:     "*.flac",
			scope:       enums.ScopeAll,
		}),

		Entry(nil, &filterTE{
			naviTE: naviTE{
				given:        "universal(any scope): glob filter (negate)",
				relative:     "RETRO-WAVE",
				subscription: enums.SubscribeUniversal,
				expectedNoOf: quantities{
					files:   6,
					folders: 8,
				},
			},
			description: "items without .flac suffix",
			pattern:     "*.flac",
			scope:       enums.ScopeAll,
			negate:      true,
		}),

		Entry(nil, &filterTE{
			naviTE: naviTE{
				given:        "universal(undefined scope): glob filter",
				relative:     "RETRO-WAVE",
				subscription: enums.SubscribeUniversal,
				expectedNoOf: quantities{
					files:   8,
					folders: 0,
				},
			},
			description: "items with '.flac' suffix",
			pattern:     "*.flac",
		}),

		// === ifNotApplicable ===============================================

		Entry(nil, &filterTE{
			naviTE: naviTE{
				given:        "universal(any scope): glob filter (ifNotApplicable=true)",
				relative:     "RETRO-WAVE",
				subscription: enums.SubscribeUniversal,
				expectedNoOf: quantities{
					files:   8,
					folders: 4,
				},
				mandatory: []string{"A1 - Can You Kiss Me First.flac"},
			},
			description:     "items with '.flac' suffix",
			pattern:         "*.flac",
			scope:           enums.ScopeLeaf,
			ifNotApplicable: enums.TriStateBoolTrue,
		}),

		Entry(nil, &filterTE{
			naviTE: naviTE{
				given:        "universal(leaf scope): glob filter (ifNotApplicable=false)",
				relative:     "RETRO-WAVE",
				subscription: enums.SubscribeUniversal,
				expectedNoOf: quantities{
					files:   8,
					folders: 0,
				},
				mandatory:  []string{"A1 - Can You Kiss Me First.flac"},
				prohibited: []string{"vinyl-info.teenage-color"},
			},
			description:     "items with '.flac' suffix",
			pattern:         "*.flac",
			scope:           enums.ScopeLeaf,
			ifNotApplicable: enums.TriStateBoolFalse,
		}),
	)
})
