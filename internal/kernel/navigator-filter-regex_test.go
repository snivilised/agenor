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
	"github.com/snivilised/traverse/internal/lo"
	"github.com/snivilised/traverse/internal/services"
	"github.com/snivilised/traverse/pref"
)

var _ = Describe("NavigatorFilterRegex", Ordered, func() {
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
			filepath.Join("MUSICO", "PROGRESSIVE-HOUSE"),
		)
		Expect(root).NotTo(BeEmpty())
	})

	BeforeEach(func() {
		services.Reset()
	})

	DescribeTable("regex-filter",
		func(ctx SpecContext, entry *filterTE) {
			var (
				traverseFilter core.TraverseFilter
			)

			recording := make(recordingMap)
			filterDefs := &pref.FilterOptions{
				Node: &core.FilterDef{
					Type:            enums.FilterTypeRegex,
					Description:     entry.name,
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

			callback := func(item *core.Node) error {
				indicator := lo.Ternary(item.IsFolder(), "📁", "💠")
				GinkgoWriter.Printf(
					"===> %v Glob Filter(%v) source: '%v', item-name: '%v', item-scope(fs): '%v(%v)'\n",
					indicator,
					traverseFilter.Description(),
					traverseFilter.Source(),
					item.Extension.Name,
					item.Extension.Scope,
					traverseFilter.Scope(),
				)
				if lo.Contains(entry.mandatory, item.Extension.Name) {
					Expect(item).Should(MatchCurrentRegexFilter(traverseFilter))
				}

				recording[item.Extension.Name] = len(item.Children)
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

		// === files =========================================================

		Entry(nil, &filterTE{
			naviTE: naviTE{
				given:        "files(any scope): regex filter",
				relative:     "RETRO-WAVE",
				subscription: enums.SubscribeFiles,
				expectedNoOf: quantities{
					files:   4,
					folders: 0,
				},
			},
			name:    "items that start with 'vinyl'",
			pattern: "^vinyl",
			scope:   enums.ScopeAll,
		}),

		Entry(nil, &filterTE{
			naviTE: naviTE{
				given:        "files(any scope): regex filter (negate)",
				relative:     "RETRO-WAVE",
				subscription: enums.SubscribeFiles,
				expectedNoOf: quantities{
					files:   10,
					folders: 0,
				},
			},
			name:    "items that don't start with 'vinyl'",
			pattern: "^vinyl",
			scope:   enums.ScopeAll,
			negate:  true,
		}),

		Entry(nil, &filterTE{
			naviTE: naviTE{
				given:        "files(default to any scope): regex filter",
				relative:     "RETRO-WAVE",
				subscription: enums.SubscribeFiles,
				expectedNoOf: quantities{
					files:   4,
					folders: 0,
				},
			},
			name:    "items that start with 'vinyl'",
			pattern: "^vinyl",
		}),

		// === folders =======================================================

		Entry(nil, &filterTE{
			naviTE: naviTE{
				given:        "folders(any scope): regex filter",
				relative:     "RETRO-WAVE",
				subscription: enums.SubscribeFolders,
				expectedNoOf: quantities{
					files:   0,
					folders: 2,
				},
			},
			name:    "items that start with 'C'",
			pattern: "^C",
			scope:   enums.ScopeAll,
		}),

		Entry(nil, &filterTE{
			naviTE: naviTE{
				given:        "folders(any scope): regex filter (negate)",
				relative:     "RETRO-WAVE",
				subscription: enums.SubscribeFolders,
				expectedNoOf: quantities{
					files:   0,
					folders: 6,
				},
			},
			name:    "items that don't start with 'C'",
			pattern: "^C",
			scope:   enums.ScopeAll,
			negate:  true,
		}),

		Entry(nil, &filterTE{
			naviTE: naviTE{
				given:        "folders(undefined scope): regex filter",
				relative:     "RETRO-WAVE",
				subscription: enums.SubscribeFolders,
				expectedNoOf: quantities{
					files:   0,
					folders: 2,
				},
			},
			name:    "items that start with 'C'",
			pattern: "^C",
		}),

		// === ifNotApplicable ===============================================

		Entry(nil, &filterTE{
			naviTE: naviTE{
				given:        "folders(top): regex filter (ifNotApplicable=true)",
				relative:     "PROGRESSIVE-HOUSE",
				subscription: enums.SubscribeFolders,
				expectedNoOf: quantities{
					files:   0,
					folders: 10,
				},
				mandatory: []string{"PROGRESSIVE-HOUSE"},
			},
			name:            "top items that contain 'HOUSE'",
			pattern:         "HOUSE",
			scope:           enums.ScopeTop,
			ifNotApplicable: enums.TriStateBoolTrue,
		}),

		Entry(nil, &filterTE{
			naviTE: naviTE{
				given:        "folders(top): regex filter (ifNotApplicable=false)",
				relative:     "",
				subscription: enums.SubscribeFolders,
				mandatory:    []string{"PROGRESSIVE-HOUSE"},
				expectedNoOf: quantities{
					files:   0,
					folders: 1,
				},
				prohibited: []string{"Blue Amazon", "The Javelin"},
			},
			name:            "top items that contain 'HOUSE'",
			pattern:         "HOUSE",
			scope:           enums.ScopeTop,
			ifNotApplicable: enums.TriStateBoolFalse,
		}),
	)
})
