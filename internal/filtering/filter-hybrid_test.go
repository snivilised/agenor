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
	"github.com/snivilised/traverse/locale"
	"github.com/snivilised/traverse/pref"
)

var _ = Describe("feature", Ordered, func() {
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
		Expect(li18ngo.Use(
			func(o *li18ngo.UseOptions) {
				o.From.Sources = li18ngo.TranslationFiles{
					locale.SourceID: li18ngo.TranslationSource{Name: "traverse"},
				}
			},
		)).To(Succeed())
	})

	BeforeEach(func() {
		services.Reset()
	})

	DescribeTable("folders with files filtered",
		func(ctx SpecContext, entry *helpers.HybridFilterTE) {
			var (
				childFilter core.ChildTraverseFilter
			)

			recording := make(helpers.RecordingMap)
			filterDefs := &pref.FilterOptions{
				Node:  &entry.NodeDef,
				Child: &entry.ChildDef,
				Sink: func(reply pref.FilterReply) {
					childFilter = reply.Child
				},
			}

			path := helpers.Path(root, entry.Relative)
			callback := func(item *core.Node) error {
				actualNoChildren := len(item.Children)
				GinkgoWriter.Printf(
					"===> 💠 Child Glob Filter(%v, children: %v)"+
						"source: '%v', node-name: '%v', node-scope: '%v', depth: '%v'\n",
					childFilter.Description(),
					actualNoChildren,
					childFilter.Source(),
					item.Extension.Name,
					item.Extension.Scope,
					item.Extension.Depth,
				)

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
				tv.WithOnBegin(helpers.Begin("🛡️")),
				tv.WithOnEnd(helpers.End("🏁")),
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
		func(entry *helpers.HybridFilterTE) string {
			return fmt.Sprintf("🧪 ===> given: '%v'", entry.Given)
		},

		Entry(nil, &helpers.HybridFilterTE{
			NaviTE: helpers.NaviTE{
				Given:        "folder(with files): glob filter",
				Relative:     "RETRO-WAVE",
				Subscription: enums.SubscribeFoldersWithFiles,
				ExpectedNoOf: helpers.Quantities{
					Folders: 6,
					Children: map[string]int{
						"Northern Council": 2,
						"Teenage Color":    2,
						"Innerworld":       2,
					},
				},
			},
			NodeDef: core.FilterDef{
				Type:        enums.FilterTypeGlob,
				Description: "folders contains o",
				Pattern:     "*o*",
				Scope:       enums.ScopeFolder,
			},
			ChildDef: core.ChildFilterDef{
				Type:        enums.FilterTypeGlob,
				Description: "items with '.flac' suffix",
				Pattern:     "*.flac",
			},
		}),

		Entry(nil, &helpers.HybridFilterTE{
			NaviTE: helpers.NaviTE{
				Given:        "folder(with files): glob filter (negate)",
				Relative:     "RETRO-WAVE",
				Subscription: enums.SubscribeFoldersWithFiles,
				ExpectedNoOf: helpers.Quantities{
					Folders: 2,
					Children: map[string]int{
						"Night Drive": 3,
					},
				},
			},
			NodeDef: core.FilterDef{
				Type:        enums.FilterTypeGlob,
				Description: "folders don't contain o",
				Pattern:     "*o*",
				Scope:       enums.ScopeFolder,
				Negate:      true,
			},
			ChildDef: core.ChildFilterDef{
				Type:        enums.FilterTypeGlob,
				Description: "items without '.txt' suffix",
				Pattern:     "*.txt",
				Negate:      true,
			},
		}),
	)
})
