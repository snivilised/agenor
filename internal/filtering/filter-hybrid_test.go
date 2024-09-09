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

		FS, root = lab.Musico(verbose,
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
		func(ctx SpecContext, entry *lab.HybridFilterTE) {
			var (
				childFilter core.ChildTraverseFilter
			)

			recording := make(lab.RecordingMap)
			filterDefs := &pref.FilterOptions{
				Node:  &entry.NodeDef,
				Child: &entry.ChildDef,
				Sink: func(reply pref.FilterReply) {
					childFilter = reply.Child
				},
			}

			path := lab.Path(root, entry.Relative)
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
				tv.WithOnBegin(lab.Begin("🛡️")),
				tv.WithOnEnd(lab.End("🏁")),
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
		func(entry *lab.HybridFilterTE) string {
			return fmt.Sprintf("🧪 ===> given: '%v'", entry.Given)
		},

		Entry(nil, &lab.HybridFilterTE{
			NaviTE: lab.NaviTE{
				Given:        "folder(with files): glob child filter",
				Relative:     "RETRO-WAVE",
				Subscription: enums.SubscribeFoldersWithFiles,
				ExpectedNoOf: lab.Quantities{
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

		Entry(nil, &lab.HybridFilterTE{
			NaviTE: lab.NaviTE{
				Given:        "folder(with files): glob child filter (negate)",
				Relative:     "RETRO-WAVE",
				Subscription: enums.SubscribeFoldersWithFiles,
				ExpectedNoOf: lab.Quantities{
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

		// =====

		XEntry(nil, &lab.HybridFilterTE{
			NaviTE: lab.NaviTE{
				Given:        "folder(with files): regex child filter",
				Relative:     "RETRO-WAVE",
				Subscription: enums.SubscribeFoldersWithFiles,
				ExpectedNoOf: lab.Quantities{
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
				Pattern:     `\.flac*`,
			},
		}),

		XEntry(nil, &lab.HybridFilterTE{
			NaviTE: lab.NaviTE{
				Given:        "folder(with files): regex child filter (negate)",
				Relative:     "RETRO-WAVE",
				Subscription: enums.SubscribeFoldersWithFiles,
				ExpectedNoOf: lab.Quantities{
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
				Pattern:     `\.txt$`,
				Negate:      true,
			},
		}),
	)
})
