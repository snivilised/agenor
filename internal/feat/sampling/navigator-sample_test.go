package sampling_test

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok
	"github.com/snivilised/li18ngo"
	tv "github.com/snivilised/traverse"
	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/enums"
	lab "github.com/snivilised/traverse/internal/laboratory"
	"github.com/snivilised/traverse/internal/services"
	"github.com/snivilised/traverse/internal/third/lo"
	"github.com/snivilised/traverse/locale"
	"github.com/snivilised/traverse/pref"
)

var _ = Describe("feature", Ordered, func() {
	var (
		FS   *lab.TestTraverseFS
		root string
	)

	BeforeAll(func() {
		const (
			verbose = false
		)

		FS, root = lab.Musico(verbose,
			lab.Static.RetroWave, "edm",
		)
		Expect(root).NotTo(BeEmpty())
		Expect(li18ngo.Use()).To(Succeed())
	})

	BeforeEach(func() {
		services.Reset()
	})

	Context("comprehension", func() {
		When("universal: slice sample", func() {
			It("🧪 should: foo", Label("example"), func(ctx SpecContext) {
				path := lab.Static.RetroWave
				result, _ := tv.Walk().Configure().Extent(tv.Prime(
					&tv.Using{
						Tree:         path,
						Subscription: enums.SubscribeUniversal,
						Handler: func(servant tv.Servant) error {
							node := servant.Node()
							GinkgoWriter.Printf(
								"---> 🍯 EXAMPLE-SAMPLE-CALLBACK: '%v'\n", node.Path,
							)
							return nil
						},
						GetTraverseFS: func(_ string) tv.TraverseFS {
							return FS
						},
					},
					tv.WithOnBegin(lab.Begin("🛡️")),
					tv.WithOnEnd(lab.End("🏁")),

					tv.WithSamplingOptions(&pref.SamplingOptions{
						Type: enums.SampleTypeSlice,
						NoOf: pref.EntryQuantities{
							Files:   2,
							Folders: 2,
						},
					}),
				)).Navigate(ctx)

				GinkgoWriter.Printf("===> 🍭 invoked '%v' folders, '%v' files.\n",
					result.Metrics().Count(enums.MetricNoFoldersInvoked),
					result.Metrics().Count(enums.MetricNoFilesInvoked),
				)
			})
		})
	})

	DescribeTable("sample",
		func(ctx SpecContext, entry *lab.SampleTE) {
			recording := make(lab.RecordingMap)
			once := func(node *tv.Node) error { //nolint:unparam // return nil error ok
				_, found := recording[node.Extension.Name]
				Expect(found).To(BeFalse())
				recording[node.Extension.Name] = len(node.Children)

				return nil
			}

			path := lo.Ternary(entry.NaviTE.Relative == "",
				lab.Static.RetroWave,
				entry.NaviTE.Relative,
			)

			callback := func(servant tv.Servant) error {
				node := servant.Node()
				GinkgoWriter.Printf(
					"---> 🌊 SAMPLE-CALLBACK: '%v'\n", node.Path,
				)
				prohibited := fmt.Sprintf("%v, was invoked, but does not satisfy sample criteria",
					lab.Reason(node.Extension.Name),
				)
				Expect(entry.Prohibited).ToNot(ContainElement(node.Extension.Name), prohibited)

				return once(node)
			}

			result, err := tv.Walk().Configure().Extent(tv.Prime(
				&tv.Using{
					Tree:         path,
					Subscription: entry.Subscription,
					Handler:      callback,
					GetTraverseFS: func(_ string) tv.TraverseFS {
						return FS
					},
				},
				tv.WithOnBegin(lab.Begin("🛡️")),
				tv.WithOnEnd(lab.End("🏁")),

				tv.WithSamplingOptions(&pref.SamplingOptions{
					Type:      entry.SampleType,
					InReverse: entry.Reverse,
					NoOf: pref.EntryQuantities{
						Files:   entry.NoOf.Files,
						Folders: entry.NoOf.Folders,
					},
					Iteration: lo.TernaryF(entry.Each != nil,
						func() pref.SamplingIterationOptions {
							return pref.SamplingIterationOptions{
								Each:  entry.Each,
								While: entry.While,
							}
						},
						func() pref.SamplingIterationOptions {
							return pref.SamplingIterationOptions{}
						},
					),
				}),
				tv.IfOptionF(entry.Filter != nil, func() pref.Option {
					return tv.WithFilter(&pref.FilterOptions{
						Sample: &core.SampleFilterDef{
							Type:        entry.Filter.Type,
							Description: entry.Filter.Description,
							Scope:       entry.Filter.Scope,
							Pattern:     entry.Filter.Pattern,
							Custom:      entry.Filter.Sample,
						},
						Custom: entry.Filter.Custom,
					})
				}),
				tv.IfOption(entry.CaseSensitive, tv.WithHookCaseSensitiveSort()),
			)).Navigate(ctx)

			lab.AssertNavigation(&entry.NaviTE, &lab.TestOptions{
				FS:          FS,
				Recording:   recording,
				Path:        path,
				Result:      result,
				Err:         err,
				ExpectedErr: entry.ExpectedErr,
			})
		},
		func(entry *lab.SampleTE) string {
			return fmt.Sprintf("🧪 ===> given: '%v', should: '%v'", entry.Given, entry.Should)
		},
		// === universal =====================================================

		Entry(nil, &lab.SampleTE{
			NaviTE: lab.NaviTE{
				Given:        "universal(slice): first, with 2 files",
				Should:       "invoke for at most 2 files per directory",
				Subscription: enums.SubscribeUniversal,
				Prohibited:   []string{"cover.night-drive.jpg"},
				ExpectedNoOf: lab.Quantities{
					Files:   8,
					Folders: 8,
				},
			},
			SampleType: enums.SampleTypeSlice,
			NoOf: pref.EntryQuantities{
				Files: 2,
			},
		}),

		Entry(nil, &lab.SampleTE{
			NaviTE: lab.NaviTE{
				Given:        "universal(slice): first, with 2 folders",
				Should:       "invoke for at most 2 folders per directory",
				Subscription: enums.SubscribeUniversal,
				Prohibited:   []string{"Electric Youth"},
				ExpectedNoOf: lab.Quantities{
					Files:   11,
					Folders: 6,
				},
			},
			SampleType: enums.SampleTypeSlice,
			NoOf: pref.EntryQuantities{
				Folders: 2,
			},
		}),

		Entry(nil, &lab.SampleTE{
			NaviTE: lab.NaviTE{
				Given:        "universal(slice): first, with 2 files and 2 folders",
				Should:       "invoke for at most 2 files and 2 folders per directory",
				Subscription: enums.SubscribeUniversal,
				Prohibited:   []string{"cover.night-drive.jpg", "Electric Youth"},
				ExpectedNoOf: lab.Quantities{
					Files:   6,
					Folders: 6,
				},
			},
			SampleType: enums.SampleTypeSlice,
			NoOf: pref.EntryQuantities{
				Files:   2,
				Folders: 2,
			},
		}),

		Entry(nil, &lab.SampleTE{
			NaviTE: lab.NaviTE{
				Given:        "universal(filter): first, single file, first 2 folders",
				Should:       "invoke for at most single file per directory",
				Relative:     "edm",
				Subscription: enums.SubscribeUniversal,
				Prohibited:   []string{"02 - Swab.flac"},
				ExpectedNoOf: lab.Quantities{
					Files:   7,
					Folders: 14,
				},
			},
			Filter: &lab.FilterTE{ // 🧄
				Description: "glob: items with .flac suffix",
				Type:        enums.FilterTypeGlob,
				Pattern:     "*.flac",
				Scope:       enums.ScopeFile,
			},
			SampleType: enums.SampleTypeFilter,
			NoOf: pref.EntryQuantities{
				Files:   1,
				Folders: 2,
			},
		}),

		Entry(nil, &lab.SampleTE{
			NaviTE: lab.NaviTE{
				Given:        "universal(filter): first, single file, first 2 folders",
				Should:       "invoke for at most single file per directory",
				Relative:     "edm",
				Subscription: enums.SubscribeUniversal,
				Prohibited:   []string{"02 - Swab.flac"},
				ExpectedNoOf: lab.Quantities{
					Files:   7,
					Folders: 14,
				},
			},
			Filter: &lab.FilterTE{ // 🚀
				Description: "regex: items with .flac suffix",
				Type:        enums.FilterTypeRegex,
				Pattern:     "\\.flac$",
				Scope:       enums.ScopeFile,
			},
			SampleType: enums.SampleTypeFilter,
			NoOf: pref.EntryQuantities{
				Files:   1,
				Folders: 2,
			},
		}),

		Entry(nil, &lab.SampleTE{
			NaviTE: lab.NaviTE{
				Given:        "universal(filter): last, last single files, last 2 folders",
				Should:       "invoke for at most single file per directory",
				Relative:     "edm",
				Subscription: enums.SubscribeUniversal,
				Prohibited:   []string{"01 - Dre.flac"},
				ExpectedNoOf: lab.Quantities{
					Files:   8,
					Folders: 15,
				},
			},
			Filter: &lab.FilterTE{ // 🧄
				Description: "glob: items with .flac suffix",
				Type:        enums.FilterTypeGlob,
				Pattern:     "*.flac",
				Scope:       enums.ScopeFile,
			},
			SampleType: enums.SampleTypeFilter,
			Reverse:    true,
			NoOf: pref.EntryQuantities{
				Files:   1,
				Folders: 2,
			},
		}),

		Entry(nil, &lab.SampleTE{
			NaviTE: lab.NaviTE{
				Given:        "universal(filter): last, last single files, last 2 folders",
				Should:       "invoke for at most single file per directory",
				Relative:     "edm",
				Subscription: enums.SubscribeUniversal,
				Prohibited:   []string{"01 - Dre.flac"},
				ExpectedNoOf: lab.Quantities{
					Files:   8,
					Folders: 15,
				},
			},
			Filter: &lab.FilterTE{ // 🚀
				Description: "regex: items with .flac suffix",
				Type:        enums.FilterTypeRegex,
				Pattern:     "\\.flac$",
				Scope:       enums.ScopeFile,
			},
			SampleType: enums.SampleTypeFilter,
			Reverse:    true,
			NoOf: pref.EntryQuantities{
				Files:   1,
				Folders: 2,
			},
		}),

		// === folders =======================================================

		Entry(nil, &lab.SampleTE{
			NaviTE: lab.NaviTE{
				Given:        "folders(slice): first, with 2 folders",
				Should:       "invoke for at most 2 folders per directory",
				Subscription: enums.SubscribeFolders,
				Prohibited:   []string{"Electric Youth"},
				ExpectedNoOf: lab.Quantities{
					Folders: 6,
				},
			},
			SampleType: enums.SampleTypeSlice,
			NoOf: pref.EntryQuantities{
				Folders: 2,
			},
		}),

		Entry(nil, &lab.SampleTE{
			NaviTE: lab.NaviTE{
				Given:        "folders(slice): last, with last single folder",
				Should:       "invoke for only last folder per directory",
				Subscription: enums.SubscribeFolders,
				Prohibited:   []string{"Chromatics"},
				ExpectedNoOf: lab.Quantities{
					Folders: 3,
				},
			},
			SampleType: enums.SampleTypeSlice,
			Reverse:    true,
			NoOf: pref.EntryQuantities{
				Folders: 1,
			},
		}),

		Entry(nil, &lab.SampleTE{
			NaviTE: lab.NaviTE{
				Given:        "filtered folders(filter): first, with 2 folders that start with A",
				Should:       "invoke for at most 2 folders per directory",
				Relative:     "edm",
				Subscription: enums.SubscribeFolders,
				Prohibited:   []string{"Tales Of Ephidrina"},
				ExpectedNoOf: lab.Quantities{
					// AMBIENT-TECHNO, Amorphous Androgynous, Aphex Twin
					Folders: 3,
				},
			},
			Filter: &lab.FilterTE{ // 🧄
				Description: "glob: items with that start with A",
				Type:        enums.FilterTypeGlob,
				Pattern:     "A*",
				Scope:       enums.ScopeFolder,
			},
			SampleType: enums.SampleTypeFilter,
			NoOf: pref.EntryQuantities{
				Folders: 2,
			},
		}),

		Entry(nil, &lab.SampleTE{
			NaviTE: lab.NaviTE{
				Given:        "filtered folders(filter): first, with 2 folders that start with A",
				Should:       "invoke for at most 2 folders per directory",
				Relative:     "edm",
				Subscription: enums.SubscribeFolders,
				Prohibited:   []string{"Tales Of Ephidrina"},
				ExpectedNoOf: lab.Quantities{
					// AMBIENT-TECHNO, Amorphous Androgynous, Aphex Twin
					Folders: 3,
				},
			},
			Filter: &lab.FilterTE{ // 🚀
				Description: "regex: items with that start with A",
				Type:        enums.FilterTypeRegex,
				Pattern:     "^A",
				Scope:       enums.ScopeFolder,
			},
			SampleType: enums.SampleTypeFilter,
			NoOf: pref.EntryQuantities{
				Folders: 2,
			},
		}),

		Entry(nil, &lab.SampleTE{
			NaviTE: lab.NaviTE{
				Given:        "filtered folders(filter): last, with single folder that start with A",
				Should:       "invoke for at most a single folder per directory",
				Relative:     "edm",
				Subscription: enums.SubscribeFolders,
				Prohibited:   []string{"Amorphous Androgynous"},
				ExpectedNoOf: lab.Quantities{
					Folders: 2,
				},
			},
			Filter: &lab.FilterTE{ // 🧄
				Description: "glob: items with that start with A",
				Type:        enums.FilterTypeGlob,
				Pattern:     "A*",
				Scope:       enums.ScopeFolder,
			},
			SampleType: enums.SampleTypeFilter,
			Reverse:    true,
			NoOf: pref.EntryQuantities{
				Folders: 1,
			},
		}),

		Entry(nil, &lab.SampleTE{
			NaviTE: lab.NaviTE{
				Given:        "filtered folders(filter): last, with single folder that start with A",
				Should:       "invoke for at most a single folder per directory",
				Relative:     "edm",
				Subscription: enums.SubscribeFolders,
				Prohibited:   []string{"Amorphous Androgynous"},
				ExpectedNoOf: lab.Quantities{
					Folders: 2,
				},
			},
			Filter: &lab.FilterTE{ // 🚀
				Description: "regex: items with that start with A",
				Type:        enums.FilterTypeRegex,
				Pattern:     "^A",
				Scope:       enums.ScopeFolder,
			},
			SampleType: enums.SampleTypeFilter,
			Reverse:    true,
			NoOf: pref.EntryQuantities{
				Folders: 1,
			},
		}),

		// === folders with files ============================================

		Entry(nil, &lab.SampleTE{
			NaviTE: lab.NaviTE{
				Given:        "folders with files(slice): first, with 2 folders",
				Should:       "invoke for at most 2 folders per directory",
				Subscription: enums.SubscribeFoldersWithFiles,
				Prohibited:   []string{"Electric Youth"},
				ExpectedNoOf: lab.Quantities{
					Folders: 6,
					Children: map[string]int{
						"Night Drive":      4,
						"Northern Council": 4,
						"Teenage Color":    3,
					},
				},
			},
			SampleType: enums.SampleTypeSlice,
			NoOf: pref.EntryQuantities{
				Folders: 2,
			},
		}),

		Entry(nil, &lab.SampleTE{
			NaviTE: lab.NaviTE{
				Given:        "folders with files(slice): last, with last single folder",
				Should:       "invoke for only last folder per directory",
				Subscription: enums.SubscribeFoldersWithFiles,
				Prohibited:   []string{"Chromatics"},
				ExpectedNoOf: lab.Quantities{
					Folders: 3,
					Children: map[string]int{
						"Innerworld": 3,
					},
				},
			},
			SampleType: enums.SampleTypeSlice,
			Reverse:    true,
			NoOf: pref.EntryQuantities{
				Folders: 1,
			},
		}),

		Entry(nil, &lab.SampleTE{
			NaviTE: lab.NaviTE{
				Given:        "filtered folders with files(filter): last, with single folder that start with A",
				Should:       "invoke for at most a single folder per directory",
				Relative:     "edm",
				Subscription: enums.SubscribeFoldersWithFiles,
				Prohibited:   []string{"Amorphous Androgynous"},
				ExpectedNoOf: lab.Quantities{
					Folders:  2,
					Children: map[string]int{},
				},
			},
			Filter: &lab.FilterTE{ // 🧄 this is folder filter, not child filter
				Description: "glob: items that start with A",
				Type:        enums.FilterTypeGlob,
				Pattern:     "A*",
				Scope:       enums.ScopeFolder,
			},
			SampleType: enums.SampleTypeFilter,
			Reverse:    true,
			NoOf: pref.EntryQuantities{
				Folders: 1,
			},
		}),

		Entry(nil, &lab.SampleTE{
			NaviTE: lab.NaviTE{
				Given:        "filtered folders with files(filter): last, with single folder that start with A",
				Should:       "invoke for at most a single folder per directory",
				Relative:     "edm",
				Subscription: enums.SubscribeFoldersWithFiles,
				Prohibited:   []string{"Amorphous Androgynous"},
				ExpectedNoOf: lab.Quantities{
					Folders:  2,
					Children: map[string]int{},
				},
			},
			Filter: &lab.FilterTE{ // 🚀
				Description: "regex: items that start with A",
				Type:        enums.FilterTypeRegex,
				Pattern:     "^A",
				Scope:       enums.ScopeFolder,
			},
			SampleType: enums.SampleTypeFilter,
			Reverse:    true,
			NoOf: pref.EntryQuantities{
				Folders: 1,
			},
		}),

		// === files =========================================================

		Entry(nil, &lab.SampleTE{
			NaviTE: lab.NaviTE{
				Given:        "files(slice): first, with 2 files",
				Should:       "invoke for at most 2 files per directory",
				Subscription: enums.SubscribeFiles,
				Prohibited:   []string{"cover.night-drive.jpg"},
				ExpectedNoOf: lab.Quantities{
					Files: 8,
				},
			},
			SampleType: enums.SampleTypeSlice,
			NoOf: pref.EntryQuantities{
				Files: 2,
			},
		}),

		Entry(nil, &lab.SampleTE{
			NaviTE: lab.NaviTE{
				Given:        "files(slice): last, with last single file",
				Should:       "invoke for only last file per directory",
				Subscription: enums.SubscribeFiles,
				Prohibited:   []string{"A1 - The Telephone Call.flac"},
				ExpectedNoOf: lab.Quantities{
					Files: 4,
				},
			},
			SampleType: enums.SampleTypeSlice,
			Reverse:    true,
			NoOf: pref.EntryQuantities{
				Files: 1,
			},
		}),

		// ScopeLeaf is not supported. Sampling filters only support
		// file/folder scopes because a node's scope is determined after
		// a directory's contents are read, but sampling filter is
		// applied at the point the contents are read. Any scopes other
		// than file/folder are ignored.
		Entry(nil, &lab.SampleTE{
			NaviTE: lab.NaviTE{
				Given:        "filtered files(filter): first, 2 files",
				Should:       "invoke for at most 2 files per directory",
				Relative:     "edm/ELECTRONICA",
				Subscription: enums.SubscribeFiles,
				Prohibited:   []string{"03 - Mountain Goat.flac"},
				ExpectedNoOf: lab.Quantities{
					Files: 24,
				},
			},
			Filter: &lab.FilterTE{ // 🧄
				Description: "glob: items with .flac suffix",
				Type:        enums.FilterTypeGlob,
				Pattern:     "*.flac",
				Scope:       enums.ScopeFile,
			},
			SampleType: enums.SampleTypeFilter,
			NoOf: pref.EntryQuantities{
				Files: 2,
			},
		}),

		Entry(nil, &lab.SampleTE{
			NaviTE: lab.NaviTE{
				Given:        "filtered files(filter): first, 2 files",
				Should:       "invoke for at most 2 files per directory",
				Relative:     "edm/ELECTRONICA",
				Subscription: enums.SubscribeFiles,
				Prohibited:   []string{"03 - Mountain Goat.flac"},
				ExpectedNoOf: lab.Quantities{
					Files: 24,
				},
			},
			Filter: &lab.FilterTE{ // 🚀
				Description: "regex: items with .flac suffix",
				Type:        enums.FilterTypeRegex,
				Pattern:     "\\.flac$",
				Scope:       enums.ScopeFile,
			},
			SampleType: enums.SampleTypeFilter,
			NoOf: pref.EntryQuantities{
				Files: 2,
			},
		}),

		Entry(nil, &lab.SampleTE{
			NaviTE: lab.NaviTE{
				Given:        "filtered files(filter): last, last 2 files",
				Should:       "invoke for at most 2 files per directory",
				Relative:     "edm",
				Subscription: enums.SubscribeFiles,
				Prohibited:   []string{"01 - Liquid Insects.flac"},
				ExpectedNoOf: lab.Quantities{
					Files: 42,
				},
			},
			Filter: &lab.FilterTE{ // 🧄
				Description: "glob: items with .flac suffix",
				Type:        enums.FilterTypeGlob,
				Pattern:     "*.flac",
				Scope:       enums.ScopeFile,
			},
			SampleType: enums.SampleTypeFilter,
			Reverse:    true,
			NoOf: pref.EntryQuantities{
				Files: 2,
			},
		}),

		Entry(nil, &lab.SampleTE{
			NaviTE: lab.NaviTE{
				Given:        "filtered files(filter): last, last 2 files",
				Should:       "invoke for at most 2 files per directory",
				Relative:     "edm",
				Subscription: enums.SubscribeFiles,
				Prohibited:   []string{"01 - Liquid Insects.flac"},
				ExpectedNoOf: lab.Quantities{
					Files: 42,
				},
			},
			Filter: &lab.FilterTE{ // 🚀
				Description: "regex: items with .flac suffix",
				Type:        enums.FilterTypeRegex,
				Pattern:     "\\.flac$",
				Scope:       enums.ScopeFile,
			},
			SampleType: enums.SampleTypeFilter,
			Reverse:    true,
			NoOf: pref.EntryQuantities{
				Files: 2,
			},
		}),

		// === custom ========================================================

		Entry(nil, &lab.SampleTE{
			NaviTE: lab.NaviTE{
				Given:        "universal(custom): first, single file, 2 folders",
				Should:       "invoke for at most single file per directory",
				Relative:     "edm",
				Subscription: enums.SubscribeUniversal,
				Prohibited:   []string{"02 - Swab.flac"},
				ExpectedNoOf: lab.Quantities{
					Files:   7,
					Folders: 14,
				},
			},
			Filter: &lab.FilterTE{ // 🍒
				Type: enums.FilterTypeCustom,
				Sample: &customSamplingFilter{
					Sample:      tv.NewCustomSampleFilter(enums.ScopeFile),
					description: "custom(glob): items with cover prefix",
					pattern:     "cover*",
				},
			},
			SampleType: enums.SampleTypeCustom,
			NoOf: pref.EntryQuantities{
				Files:   1,
				Folders: 2,
			},
		}),

		Entry(nil, &lab.SampleTE{
			NaviTE: lab.NaviTE{
				Given:        "filtered folders(custom): last, single folder that starts with A",
				Should:       "invoke for at most a single folder per directory",
				Relative:     "edm",
				Subscription: enums.SubscribeFolders,
				Prohibited:   []string{"Amorphous Androgynous"},
				ExpectedNoOf: lab.Quantities{
					Folders: 2,
				},
			},
			Filter: &lab.FilterTE{ // 🍒
				Type: enums.FilterTypeCustom,
				Sample: &customSamplingFilter{
					Sample:      tv.NewCustomSampleFilter(enums.ScopeFolder),
					description: "custom(glob): items with A prefix",
					pattern:     "A*",
				},
			},
			SampleType: enums.SampleTypeCustom,
			NoOf: pref.EntryQuantities{
				Folders: 1,
			},
			Reverse: true,
		}),

		Entry(nil, &lab.SampleTE{
			NaviTE: lab.NaviTE{
				Given:        "filtered files(custom): last, last 2 files",
				Should:       "invoke for at most 2 files per directory",
				Relative:     "edm",
				Subscription: enums.SubscribeFiles,
				Prohibited:   []string{"01 - Liquid Insects.flac"},
				ExpectedNoOf: lab.Quantities{
					Files: 42,
				},
			},
			Filter: &lab.FilterTE{ // 🍒
				Type: enums.FilterTypeCustom,
				Sample: &customSamplingFilter{
					Sample:      tv.NewCustomSampleFilter(enums.ScopeFile),
					description: "custom(glob): items with .flac suffix",
					pattern:     "*.flac",
				},
			},
			SampleType: enums.SampleTypeCustom,
			NoOf: pref.EntryQuantities{
				Files: 2,
			},
			Reverse: true,
		}),

		// === errors ========================================================

		Entry(nil, &lab.SampleTE{
			NaviTE: lab.NaviTE{
				Given:        "folder spec, without no of folders",
				Should:       "return invalid folder spec error",
				Relative:     "edm/ELECTRONICA",
				Subscription: enums.SubscribeFiles,
				Prohibited:   []string{"03 - Mountain Goat.flac"},
				ExpectedNoOf: lab.Quantities{
					Files: 24,
				},
				ExpectedErr: locale.ErrInvalidFolderSamplingSpecMissingFolders,
			},
			Filter: &lab.FilterTE{ // 🧄
				Description: "glob: items with .flac suffix",
				Type:        enums.FilterTypeGlob,
				Pattern:     "*.flac",
				Scope:       enums.ScopeFolder,
			},
			SampleType: enums.SampleTypeFilter,
			NoOf: pref.EntryQuantities{
				Files: 2,
			},
		}),

		Entry(nil, &lab.SampleTE{
			NaviTE: lab.NaviTE{
				Given:        "file spec, without no of files",
				Should:       "return invalid file spec error",
				Relative:     "edm/ELECTRONICA",
				Subscription: enums.SubscribeFiles,
				Prohibited:   []string{"03 - Mountain Goat.flac"},
				ExpectedNoOf: lab.Quantities{
					Files: 24,
				},
				ExpectedErr: locale.ErrInvalidFileSamplingSpecMissingFiles,
			},
			Filter: &lab.FilterTE{ // 🧄
				Description: "glob: items with .flac suffix",
				Type:        enums.FilterTypeGlob,
				Pattern:     "*.flac",
				Scope:       enums.ScopeFile,
			},
			SampleType: enums.SampleTypeFilter,
			NoOf: pref.EntryQuantities{
				Folders: 2,
			},
		}),
	)
})
