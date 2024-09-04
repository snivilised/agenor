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
	"github.com/snivilised/traverse/internal/helpers"
	"github.com/snivilised/traverse/internal/services"
	"github.com/snivilised/traverse/internal/third/lo"
	"github.com/snivilised/traverse/pref"
)

var _ = Describe("NavigatorFilterCustom", Ordered, func() {
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
		Expect(li18ngo.Use()).To(Succeed())
	})

	BeforeEach(func() {
		services.Reset()
	})

	DescribeTable("custom-filter (glob)",
		func(ctx SpecContext, entry *helpers.FilterTE) {
			recording := make(helpers.RecordingMap)
			customFilter := &customFilter{
				name:    entry.Description,
				pattern: entry.Pattern,
				scope:   entry.Scope,
				negate:  entry.Negate,
			}

			path := helpers.Path(root, entry.Relative)
			callback := func(item *core.Node) error {
				indicator := lo.Ternary(item.IsFolder(), "📁", "💠")
				GinkgoWriter.Printf(
					"===> %v Glob Filter(%v) source: '%v', item-name: '%v', item-scope(fs): '%v(%v)'\n",
					indicator,
					customFilter.Description(),
					customFilter.Source(),
					item.Extension.Name,
					item.Extension.Scope,
					customFilter.Scope(),
				)
				if lo.Contains(entry.Mandatory, item.Extension.Name) {
					Expect(item).Should(MatchCurrentCustomFilter(customFilter))
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
				tv.WithFilter(&pref.FilterOptions{
					Custom: customFilter,
				}),
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

		// === universal =====================================================

		Entry(nil, &helpers.FilterTE{
			NaviTE: helpers.NaviTE{
				Given:        "universal(any scope): custom filter",
				Relative:     "RETRO-WAVE",
				Subscription: enums.SubscribeUniversal,
				ExpectedNoOf: helpers.Quantities{
					Files:   8,
					Folders: 0,
				},
			},
			Description: "items with '.flac' suffix",
			Pattern:     "*.flac",
			Scope:       enums.ScopeFile,
		}),

		Entry(nil, &helpers.FilterTE{
			NaviTE: helpers.NaviTE{
				Given:        "universal(any scope): custom filter (negate)",
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
				Given:        "universal(undefined scope): custom filter",
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
	)
})
