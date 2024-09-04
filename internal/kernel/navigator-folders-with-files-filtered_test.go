package kernel_test

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

var _ = Describe("NavigatorFoldersWithFiles", Ordered, func() {
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
		func(ctx SpecContext, entry *helpers.FilterTE) {
			var (
				childFilter core.ChildTraverseFilter
			)

			recording := make(helpers.RecordingMap)
			filterDefs := &pref.FilterOptions{
				Child: &core.ChildFilterDef{
					Type:        enums.FilterTypeGlob,
					Description: entry.Description,
					Pattern:     entry.Pattern,
					Negate:      entry.Negate,
				},
				Sink: func(reply pref.FilterReply) {
					childFilter = reply.Child
				},
			}

			path := helpers.Path(root, entry.Relative)

			callback := func(item *core.Node) error {
				actualNoChildren := len(item.Children)
				GinkgoWriter.Printf(
					"===> 💠 Compound Glob Filter(%v, children: %v) source: '%v', node-name: '%v', node-scope: '%v', depth: '%v'\n",
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

			Expect(err).To(Succeed())

			if entry.Mandatory != nil {
				for _, name := range entry.Mandatory {
					_, found := recording[name]
					Expect(found).To(BeTrue(), helpers.Reason(name))
				}
			}

			if entry.Prohibited != nil {
				for _, name := range entry.Prohibited {
					_, found := recording[name]
					Expect(found).To(BeFalse(), helpers.Reason(name))
				}
			}

			for n, actualNoChildren := range entry.ExpectedNoOf.Children {
				expected := recording[n]

				Expect(expected).To(Equal(actualNoChildren),
					helpers.BecauseQuantity("Incorrect no of children",
						expected,
						actualNoChildren,
					),
				)
			}

			Expect(result.Metrics().Count(enums.MetricNoFilesInvoked)).To(
				Equal(entry.ExpectedNoOf.Files),
				helpers.BecauseQuantity("Incorrect no of files",
					int(entry.ExpectedNoOf.Files),
					int(result.Metrics().Count(enums.MetricNoFilesInvoked)),
				),
			)

			Expect(result.Metrics().Count(enums.MetricNoFoldersInvoked)).To(
				Equal(entry.ExpectedNoOf.Folders),
				helpers.BecauseQuantity("Incorrect no of folders",
					int(entry.ExpectedNoOf.Folders),
					int(result.Metrics().Count(enums.MetricNoFoldersInvoked)),
				),
			)
		},
		func(entry *helpers.FilterTE) string {
			return fmt.Sprintf("🧪 ===> given: '%v'", entry.Given)
		},

		// folders with files not implemented yet
		XEntry(nil, &helpers.FilterTE{
			NaviTE: helpers.NaviTE{
				Given:        "folder(with files): glob filter",
				Relative:     "RETRO-WAVE",
				Subscription: enums.SubscribeFoldersWithFiles,
				ExpectedNoOf: helpers.Quantities{
					Files:   0,
					Folders: 8,
					Children: map[string]int{
						"Night Drive":      2,
						"Northern Council": 2,
						"Teenage Color":    2,
						"Innerworld":       2,
					},
				},
			},
			Description: "items with '.flac' suffix",
			Pattern:     "*.flac",
		}),

		// folders with files not implemented yet
		XEntry(nil, &helpers.FilterTE{
			NaviTE: helpers.NaviTE{
				Given:        "folder(with files): glob filter (negate)",
				Relative:     "RETRO-WAVE",
				Subscription: enums.SubscribeFoldersWithFiles,
				ExpectedNoOf: helpers.Quantities{
					Files:   0,
					Folders: 8,
					Children: map[string]int{
						"Night Drive":      3,
						"Northern Council": 3,
						"Teenage Color":    2,
						"Innerworld":       2,
					},
				},
			},
			Description: "items without '.txt' suffix",
			Pattern:     "*.txt",
			Negate:      true,
		}),
	)
})
