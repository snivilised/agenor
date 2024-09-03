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

var _ = Describe("filter", Ordered, func() {
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
			filepath.Join("MUSICO", "PROGRESSIVE-HOUSE"),
		)
		Expect(root).NotTo(BeEmpty())
	})

	BeforeEach(func() {
		services.Reset()
	})

	DescribeTable("regex",
		func(ctx SpecContext, entry *helpers.FilterTE) {
			var (
				traverseFilter core.TraverseFilter
			)

			recording := make(helpers.RecordingMap)
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

			path := helpers.Path(root, entry.Relative)

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
				if lo.Contains(entry.Mandatory, item.Extension.Name) {
					Expect(item).Should(MatchCurrentRegexFilter(traverseFilter))
				}

				recording[item.Extension.Name] = len(item.Children)
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

		// === files =========================================================

		Entry(nil, &helpers.FilterTE{
			NaviTE: helpers.NaviTE{
				Given:        "files(any scope): regex filter",
				Relative:     "RETRO-WAVE",
				Subscription: enums.SubscribeFiles,
				ExpectedNoOf: helpers.Quantities{
					Files:   4,
					Folders: 0,
				},
			},
			Description: "items that start with 'vinyl'",
			Pattern:     "^vinyl",
			Scope:       enums.ScopeAll,
		}),

		Entry(nil, &helpers.FilterTE{
			NaviTE: helpers.NaviTE{
				Given:        "files(any scope): regex filter (negate)",
				Relative:     "RETRO-WAVE",
				Subscription: enums.SubscribeFiles,
				ExpectedNoOf: helpers.Quantities{
					Files:   10,
					Folders: 0,
				},
			},
			Description: "items that don't start with 'vinyl'",
			Pattern:     "^vinyl",
			Scope:       enums.ScopeAll,
			Negate:      true,
		}),

		Entry(nil, &helpers.FilterTE{
			NaviTE: helpers.NaviTE{
				Given:        "files(default to any scope): regex filter",
				Relative:     "RETRO-WAVE",
				Subscription: enums.SubscribeFiles,
				ExpectedNoOf: helpers.Quantities{
					Files:   4,
					Folders: 0,
				},
			},
			Description: "items that start with 'vinyl'",
			Pattern:     "^vinyl",
		}),

		// === folders =======================================================

		Entry(nil, &helpers.FilterTE{
			NaviTE: helpers.NaviTE{
				Given:        "folders(any scope): regex filter",
				Relative:     "RETRO-WAVE",
				Subscription: enums.SubscribeFolders,
				ExpectedNoOf: helpers.Quantities{
					Files:   0,
					Folders: 2,
				},
			},
			Description: "items that start with 'C'",
			Pattern:     "^C",
			Scope:       enums.ScopeAll,
		}),

		Entry(nil, &helpers.FilterTE{
			NaviTE: helpers.NaviTE{
				Given:        "folders(any scope): regex filter (negate)",
				Relative:     "RETRO-WAVE",
				Subscription: enums.SubscribeFolders,
				ExpectedNoOf: helpers.Quantities{
					Files:   0,
					Folders: 6,
				},
			},
			Description: "items that don't start with 'C'",
			Pattern:     "^C",
			Scope:       enums.ScopeAll,
			Negate:      true,
		}),

		Entry(nil, &helpers.FilterTE{
			NaviTE: helpers.NaviTE{
				Given:        "folders(undefined scope): regex filter",
				Relative:     "RETRO-WAVE",
				Subscription: enums.SubscribeFolders,
				ExpectedNoOf: helpers.Quantities{
					Files:   0,
					Folders: 2,
				},
			},
			Description: "items that start with 'C'",
			Pattern:     "^C",
		}),

		// === ifNotApplicable ===============================================

		Entry(nil, &helpers.FilterTE{
			NaviTE: helpers.NaviTE{
				Given:        "folders(top): regex filter (ifNotApplicable=true)",
				Relative:     "PROGRESSIVE-HOUSE",
				Subscription: enums.SubscribeFolders,
				ExpectedNoOf: helpers.Quantities{
					Files:   0,
					Folders: 10,
				},
				Mandatory: []string{"PROGRESSIVE-HOUSE"},
			},
			Description:     "top items that contain 'HOUSE'",
			Pattern:         "HOUSE",
			Scope:           enums.ScopeTop,
			IfNotApplicable: enums.TriStateBoolTrue,
		}),

		Entry(nil, &helpers.FilterTE{
			NaviTE: helpers.NaviTE{
				Given:        "folders(top): regex filter (ifNotApplicable=false)",
				Relative:     "",
				Subscription: enums.SubscribeFolders,
				Mandatory:    []string{"PROGRESSIVE-HOUSE"},
				ExpectedNoOf: helpers.Quantities{
					Files:   0,
					Folders: 1,
				},
				Prohibited: []string{"Blue Amazon", "The Javelin"},
			},
			Description:     "top items that contain 'HOUSE'",
			Pattern:         "HOUSE",
			Scope:           enums.ScopeTop,
			IfNotApplicable: enums.TriStateBoolFalse,
		}),
	)
})
