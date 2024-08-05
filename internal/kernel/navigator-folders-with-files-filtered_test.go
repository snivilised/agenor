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
	"github.com/snivilised/traverse/pref"
)

var _ = Describe("NavigatorFoldersWithFiles", Ordered, func() {
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

	DescribeTable("folders with files filtered",
		func(ctx SpecContext, entry *filterTE) {
			var (
				childFilter core.ChildTraverseFilter
			)

			recording := make(recordingMap)
			filterDefs := &pref.FilterOptions{
				Child: &core.ChildFilterDef{
					Type:        enums.FilterTypeGlob,
					Description: entry.name,
					Pattern:     entry.pattern,
					Negate:      entry.negate,
				},
				Sink: func(reply pref.FilterReply) {
					childFilter = reply.Child
				},
			}

			path := helpers.Path(root, entry.relative)

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

			Expect(err).To(Succeed())

			if entry.mandatory != nil {
				for _, name := range entry.mandatory {
					_, found := recording[name]
					Expect(found).To(BeTrue(), helpers.Reason(name))
				}
			}

			if entry.prohibited != nil {
				for _, name := range entry.prohibited {
					_, found := recording[name]
					Expect(found).To(BeFalse(), helpers.Reason(name))
				}
			}

			for n, actualNoChildren := range entry.expectedNoOf.children {
				expected := recording[n]

				Expect(expected).To(Equal(actualNoChildren),
					helpers.BecauseQuantity("Incorrect no of children",
						expected,
						actualNoChildren,
					),
				)
			}

			Expect(result.Metrics().Count(enums.MetricNoFilesInvoked)).To(
				Equal(entry.expectedNoOf.files),
				helpers.BecauseQuantity("Incorrect no of files",
					int(entry.expectedNoOf.files),
					int(result.Metrics().Count(enums.MetricNoFilesInvoked)),
				),
			)

			Expect(result.Metrics().Count(enums.MetricNoFoldersInvoked)).To(
				Equal(entry.expectedNoOf.folders),
				helpers.BecauseQuantity("Incorrect no of folders",
					int(entry.expectedNoOf.folders),
					int(result.Metrics().Count(enums.MetricNoFoldersInvoked)),
				),
			)
		},
		func(entry *filterTE) string {
			return fmt.Sprintf("🧪 ===> given: '%v'", entry.given)
		},

		// folders with files not implemented yet
		XEntry(nil, &filterTE{
			naviTE: naviTE{
				given:        "folder(with files): glob filter",
				relative:     "RETRO-WAVE",
				subscription: enums.SubscribeFoldersWithFiles,
				expectedNoOf: quantities{
					files:   0,
					folders: 8,
					children: map[string]int{
						"Night Drive":      2,
						"Northern Council": 2,
						"Teenage Color":    2,
						"Innerworld":       2,
					},
				},
			},
			name:    "items with '.flac' suffix",
			pattern: "*.flac",
		}),

		// folders with files not implemented yet
		XEntry(nil, &filterTE{
			naviTE: naviTE{
				given:        "folder(with files): glob filter (negate)",
				relative:     "RETRO-WAVE",
				subscription: enums.SubscribeFoldersWithFiles,
				expectedNoOf: quantities{
					files:   0,
					folders: 8,
					children: map[string]int{
						"Night Drive":      3,
						"Northern Council": 3,
						"Teenage Color":    2,
						"Innerworld":       2,
					},
				},
			},
			name:    "items without '.txt' suffix",
			pattern: "*.txt",
			negate:  true,
		}),
	)
})
