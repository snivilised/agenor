package kernel_test

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
	"testing/fstest"

	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok
	tv "github.com/snivilised/traverse"
	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/enums"
	"github.com/snivilised/traverse/i18n"
	"github.com/snivilised/traverse/internal/helpers"
	"github.com/snivilised/traverse/internal/lo"
	"github.com/snivilised/traverse/internal/services"
	"github.com/snivilised/traverse/pref"
)

var _ = Describe("Sampling", Ordered, func() {
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
			filepath.Join("MUSICO", "edm"),
		)
		Expect(root).NotTo(BeEmpty())

		_ = vfs
	})

	BeforeEach(func() {
		services.Reset()
	})

	DescribeTable("sample",
		func(ctx SpecContext, entry *sampleTE) {
			recording := make(recordingMap)
			once := func(node *tv.Node) error { //nolint:unparam // return nil error ok
				_, found := recording[node.Extension.Name]
				Expect(found).To(BeFalse())
				recording[node.Extension.Name] = len(node.Children)

				return nil
			}

			path := helpers.Path(
				root,
				lo.Ternary(entry.naviTE.relative == "",
					"RETRO-WAVE",
					entry.naviTE.relative,
				),
			)

			callback := func(node *tv.Node) error {
				GinkgoWriter.Printf(
					"---> 🌊 SAMPLE-CALLBACK: '%v'\n", node.Path,
				)
				prohibited := fmt.Sprintf("%v, was invoked, but does not satisfy sample criteria",
					helpers.Reason(node.Extension.Name),
				)
				Expect(entry.prohibited).ToNot(ContainElement(node.Extension.Name), prohibited)

				return once(node)
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
				tv.WithSampling(&pref.SamplingOptions{
					SampleType:      entry.sampleType,
					SampleInReverse: entry.reverse,
					NoOf: pref.EntryQuantities{
						Files:   entry.noOf.Files,
						Folders: entry.noOf.Folders,
					},
					Iteration: lo.TernaryF(entry.each != nil,
						func() pref.SamplingIterationOptions {
							return pref.SamplingIterationOptions{
								Each:  entry.each,
								While: entry.while,
							}
						},
						func() pref.SamplingIterationOptions {
							return pref.SamplingIterationOptions{}
						},
					),
				}),
				tv.IfOptionF(entry.filter != nil, func() pref.Option {
					return tv.WithFilter(&pref.FilterOptions{
						Sample: &core.SampleFilterDef{
							Type:        enums.FilterTypeGlob,
							Description: entry.filter.name,
							Scope:       entry.filter.scope,
							Pattern:     entry.filter.pattern,
						},
					})
				}),
				tv.IfOption(entry.caseSensitive, tv.WithHookCaseSensitiveSort()),
				tv.WithHookQueryStatus(
					func(qsys fs.StatFS, path string) (fs.FileInfo, error) {
						return qsys.Stat(helpers.TrimRoot(path))
					},
				),
				tv.WithHookReadDirectory(
					func(rsys fs.ReadDirFS, dirname string) ([]fs.DirEntry, error) {
						return rsys.ReadDir(helpers.TrimRoot(dirname))
					},
				),
			)).Navigate(ctx)

			assertNavigation(&entry.naviTE, &testOptions{
				vfs:         vfs,
				recording:   recording,
				path:        path,
				result:      result,
				err:         err,
				expectedErr: entry.expectedErr,
			})
		},
		func(entry *sampleTE) string {
			return fmt.Sprintf("🧪 ===> given: '%v', should: '%v'", entry.given, entry.should)
		},
		// === universal =====================================================

		Entry(nil, &sampleTE{
			naviTE: naviTE{
				given:        "universal(slice): first, with 2 files",
				should:       "invoke for at most 2 files per directory",
				subscription: enums.SubscribeUniversal,
				prohibited:   []string{"cover.night-drive.jpg"},
				expectedNoOf: quantities{
					files:   8,
					folders: 8,
				},
			},
			sampleType: enums.SampleTypeSlice,
			noOf: pref.EntryQuantities{
				Files: 2,
			},
		}),

		Entry(nil, &sampleTE{
			naviTE: naviTE{
				given:        "universal(slice): first, with 2 folders",
				should:       "invoke for at most 2 folders per directory",
				subscription: enums.SubscribeUniversal,
				prohibited:   []string{"Electric Youth"},
				expectedNoOf: quantities{
					files:   11,
					folders: 6,
				},
			},
			sampleType: enums.SampleTypeSlice,
			noOf: pref.EntryQuantities{
				Folders: 2,
			},
		}),

		Entry(nil, &sampleTE{
			naviTE: naviTE{
				given:        "universal(slice): first, with 2 files and 2 folders",
				should:       "invoke for at most 2 files and 2 folders per directory",
				subscription: enums.SubscribeUniversal,
				prohibited:   []string{"cover.night-drive.jpg", "Electric Youth"},
				expectedNoOf: quantities{
					files:   6,
					folders: 6,
				},
			},
			sampleType: enums.SampleTypeSlice,
			noOf: pref.EntryQuantities{
				Files:   2,
				Folders: 2,
			},
		}),

		Entry(nil, &sampleTE{
			naviTE: naviTE{
				given:        "universal(filter): first, single file, first 2 folders",
				should:       "invoke for at most single file per directory",
				relative:     "edm",
				subscription: enums.SubscribeUniversal,
				prohibited:   []string{"02 - Swab.flac"},
				expectedNoOf: quantities{
					files:   7,
					folders: 14,
				},
			},
			filter: &filterTE{
				name: "items with .flac suffix",

				pattern: "*.flac",
				scope:   enums.ScopeFile,
			},
			sampleType: enums.SampleTypeFilter,
			noOf: pref.EntryQuantities{
				Files:   1,
				Folders: 2,
			},
		}),

		Entry(nil, &sampleTE{
			naviTE: naviTE{
				given:        "universal(filter): last, last single files, last 2 folders",
				should:       "invoke for at most single file per directory",
				relative:     "edm",
				subscription: enums.SubscribeUniversal,
				prohibited:   []string{"01 - Dre.flac"},
				expectedNoOf: quantities{
					files:   8,
					folders: 15,
				},
			},
			filter: &filterTE{
				name:    "items with .flac suffix",
				pattern: "*.flac",
				scope:   enums.ScopeFile,
			},
			sampleType: enums.SampleTypeFilter,
			reverse:    true,
			noOf: pref.EntryQuantities{
				Files:   1,
				Folders: 2,
			},
		}),

		// === folders =======================================================

		Entry(nil, &sampleTE{
			naviTE: naviTE{
				given:        "folders(slice): first, with 2 folders",
				should:       "invoke for at most 2 folders per directory",
				subscription: enums.SubscribeFolders,
				prohibited:   []string{"Electric Youth"},
				expectedNoOf: quantities{
					folders: 6,
				},
			},
			sampleType: enums.SampleTypeSlice,
			noOf: pref.EntryQuantities{
				Folders: 2,
			},
		}),

		Entry(nil, &sampleTE{
			naviTE: naviTE{
				given:        "folders(slice): last, with last single folder",
				should:       "invoke for only last folder per directory",
				subscription: enums.SubscribeFolders,
				prohibited:   []string{"Chromatics"},
				expectedNoOf: quantities{
					folders: 3,
				},
			},
			sampleType: enums.SampleTypeSlice,
			reverse:    true,
			noOf: pref.EntryQuantities{
				Folders: 1,
			},
		}),

		Entry(nil, &sampleTE{
			naviTE: naviTE{
				given:        "filtered folders(filter): first, with 2 folders that start with A",
				should:       "invoke for at most 2 folders per directory",
				relative:     "edm",
				subscription: enums.SubscribeFolders,
				prohibited:   []string{"Tales Of Ephidrina"},
				expectedNoOf: quantities{
					// AMBIENT-TECHNO, Amorphous Androgynous, Aphex Twin
					folders: 3,
				},
			},
			filter: &filterTE{
				name:    "items with that start with A",
				pattern: "A*",
				scope:   enums.ScopeFolder,
			},
			sampleType: enums.SampleTypeFilter,
			noOf: pref.EntryQuantities{
				Folders: 2,
			},
		}),

		Entry(nil, &sampleTE{
			naviTE: naviTE{
				given:        "filtered folders(filter): last, with single folder that start with A",
				should:       "invoke for at most a single folder per directory",
				relative:     "edm",
				subscription: enums.SubscribeFolders,
				prohibited:   []string{"Amorphous Androgynous"},
				expectedNoOf: quantities{
					folders: 2,
				},
			},
			filter: &filterTE{
				name:    "items with that start with A",
				pattern: "A*",
				scope:   enums.ScopeFolder,
			},
			sampleType: enums.SampleTypeFilter,
			reverse:    true,
			noOf: pref.EntryQuantities{
				Folders: 1,
			},
		}),

		// === folders with files ============================================

		Entry(nil, &sampleTE{
			naviTE: naviTE{
				given:        "folders with files(slice): first, with 2 folders",
				should:       "invoke for at most 2 folders per directory",
				subscription: enums.SubscribeFoldersWithFiles,
				prohibited:   []string{"Electric Youth"},
				expectedNoOf: quantities{
					folders: 6,
					children: map[string]int{
						"Night Drive":      4,
						"Northern Council": 4,
						"Teenage Color":    3,
					},
				},
			},
			sampleType: enums.SampleTypeSlice,
			noOf: pref.EntryQuantities{
				Folders: 2,
			},
		}),

		Entry(nil, &sampleTE{
			naviTE: naviTE{
				given:        "folders with files(slice): last, with last single folder",
				should:       "invoke for only last folder per directory",
				subscription: enums.SubscribeFoldersWithFiles,
				prohibited:   []string{"Chromatics"},
				expectedNoOf: quantities{
					folders: 3,
					children: map[string]int{
						"Innerworld": 3,
					},
				},
			},
			sampleType: enums.SampleTypeSlice,
			reverse:    true,
			noOf: pref.EntryQuantities{
				Folders: 1,
			},
		}),

		// child filter not implemented yet
		Entry(nil, &sampleTE{
			naviTE: naviTE{
				given:        "filtered folders with files(filter): last, with single folder that start with A",
				should:       "invoke for at most a single folder per directory",
				relative:     "edm",
				subscription: enums.SubscribeFoldersWithFiles,
				prohibited:   []string{"Amorphous Androgynous"},
				expectedNoOf: quantities{
					folders:  2,
					children: map[string]int{},
				},
			},
			filter: &filterTE{ // this is folder filter, not child filter
				name:    "items that start with A",
				pattern: "A*",
				scope:   enums.ScopeFolder,
			},
			sampleType: enums.SampleTypeFilter,
			reverse:    true,
			noOf: pref.EntryQuantities{
				Folders: 1,
			},
		}),

		// === files =========================================================

		Entry(nil, &sampleTE{
			naviTE: naviTE{
				given:        "files(slice): first, with 2 files",
				should:       "invoke for at most 2 files per directory",
				subscription: enums.SubscribeFiles,
				prohibited:   []string{"cover.night-drive.jpg"},
				expectedNoOf: quantities{
					files: 8,
				},
			},
			sampleType: enums.SampleTypeSlice,
			noOf: pref.EntryQuantities{
				Files: 2,
			},
		}),

		Entry(nil, &sampleTE{
			naviTE: naviTE{
				given:        "files(slice): last, with last single file",
				should:       "invoke for only last file per directory",
				subscription: enums.SubscribeFiles,
				prohibited:   []string{"A1 - The Telephone Call.flac"},
				expectedNoOf: quantities{
					files: 4,
				},
			},
			sampleType: enums.SampleTypeSlice,
			reverse:    true,
			noOf: pref.EntryQuantities{
				Files: 1,
			},
		}),

		// ScopeLeaf is not supported. Sampling filters only support
		// file/folder scopes because a node's scope is determined after
		// a directory's contents are read, but sampling filter is
		// applied at the point the contents are read. Any scopes other
		// than file/folder are ignored.
		Entry(nil, &sampleTE{
			naviTE: naviTE{
				given:        "filtered files(filter): first, 2 files",
				should:       "invoke for at most 2 files per directory",
				relative:     "edm/ELECTRONICA",
				subscription: enums.SubscribeFiles,
				prohibited:   []string{"03 - Mountain Goat.flac"},
				expectedNoOf: quantities{
					files: 24,
				},
			},
			filter: &filterTE{
				name:    "items with .flac suffix",
				pattern: "*.flac",
				scope:   enums.ScopeFile,
			},
			sampleType: enums.SampleTypeFilter,
			noOf: pref.EntryQuantities{
				Files: 2,
			},
		}),

		Entry(nil, &sampleTE{
			naviTE: naviTE{
				given:        "filtered files(filter): last, last 2 files",
				should:       "invoke for at most 2 files per directory",
				relative:     "edm",
				subscription: enums.SubscribeFiles,
				prohibited:   []string{"01 - Liquid Insects.flac"},
				expectedNoOf: quantities{
					files: 42,
				},
			},
			filter: &filterTE{
				name:    "items with .flac suffix",
				pattern: "*.flac",
				scope:   enums.ScopeFile,
			},
			sampleType: enums.SampleTypeFilter,
			reverse:    true,
			noOf: pref.EntryQuantities{
				Files: 2,
			},
		}),

		// === custom ========================================================

		// custom not implemented yet
		XEntry(nil, &sampleTE{
			naviTE: naviTE{
				given:        "universal(custom): first, single file, 2 folders",
				should:       "invoke for at most single file per directory",
				relative:     "edm",
				subscription: enums.SubscribeUniversal,
				prohibited:   []string{"02 - Swab.flac"},
				expectedNoOf: quantities{
					files:   7,
					folders: 14,
				},
			},
			filter: &filterTE{},
			each: func(node *core.Node) bool { // convert to child filter
				if node.IsFolder() {
					return true
				}

				return strings.HasPrefix(node.Extension.Name, "cover")
			},
			while: func(fi *pref.FilteredInfo) bool {
				fi.Enough.Files = fi.Counts.Files == 1
				fi.Enough.Folders = fi.Counts.Folders == 2

				return !fi.Enough.Files && !fi.Enough.Folders
			},
			sampleType: enums.SampleTypeCustom,
			noOf: pref.EntryQuantities{
				Files:   1,
				Folders: 2,
			},
		}),

		// custom not implemented yet
		XEntry(nil, &sampleTE{
			naviTE: naviTE{
				given:        "filtered folders(custom): last, single folder that starts with A",
				should:       "invoke for at most a single folder per directory",
				relative:     "edm",
				subscription: enums.SubscribeFolders,
				prohibited:   []string{"Amorphous Androgynous"},
				expectedNoOf: quantities{
					folders: 3,
				},
			},
			each: func(node *core.Node) bool {
				return strings.HasPrefix(node.Extension.Name, "A")
			},
			while: func(fi *pref.FilteredInfo) bool {
				return fi.Counts.Folders < 1
			},
			sampleType: enums.SampleTypeCustom,
			reverse:    true,
		}),

		// custom filter not implemented yet
		XEntry(nil, &sampleTE{
			naviTE: naviTE{
				given:        "filtered files(custom): last, last 2 files",
				should:       "invoke for at most 2 files per directory",
				relative:     "edm",
				subscription: enums.SubscribeFiles,
				prohibited:   []string{"01 - Liquid Insects.flac"},
				expectedNoOf: quantities{
					files: 42,
				},
			},
			each: func(node *core.Node) bool {
				return strings.HasSuffix(node.Extension.Name, ".flac")
			},
			while: func(fi *pref.FilteredInfo) bool {
				return fi.Counts.Files != 2
			},
			sampleType: enums.SampleTypeCustom,
			reverse:    true,
		}),

		// === errors ========================================================

		Entry(nil, &sampleTE{
			naviTE: naviTE{
				given:        "folder spec, without no of folders",
				should:       "return invalid folder spec error",
				relative:     "edm/ELECTRONICA",
				subscription: enums.SubscribeFiles,
				prohibited:   []string{"03 - Mountain Goat.flac"},
				expectedNoOf: quantities{
					files: 24,
				},
				expectedErr: i18n.ErrInvalidFolderSamplingSpecification,
			},
			filter: &filterTE{
				name:    "items with .flac suffix",
				pattern: "*.flac",
				scope:   enums.ScopeFolder,
			},
			sampleType: enums.SampleTypeFilter,
			noOf: pref.EntryQuantities{
				Files: 2,
			},
		}),

		Entry(nil, &sampleTE{
			naviTE: naviTE{
				given:        "file spec, without no of files",
				should:       "return invalid file spec error",
				relative:     "edm/ELECTRONICA",
				subscription: enums.SubscribeFiles,
				prohibited:   []string{"03 - Mountain Goat.flac"},
				expectedNoOf: quantities{
					files: 24,
				},
				expectedErr: i18n.ErrInvalidFileSamplingSpecification,
			},
			filter: &filterTE{
				name:    "items with .flac suffix",
				pattern: "*.flac",
				scope:   enums.ScopeFile,
			},
			sampleType: enums.SampleTypeFilter,
			noOf: pref.EntryQuantities{
				Folders: 2,
			},
		}),
	)
})
