package sampling_test

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok
	"github.com/snivilised/li18ngo"
	nef "github.com/snivilised/nefilim"
	"github.com/snivilised/nefilim/test/luna"
	tv "github.com/snivilised/traverse"
	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/enums"
	lab "github.com/snivilised/traverse/internal/laboratory"
	"github.com/snivilised/traverse/internal/services"
	"github.com/snivilised/traverse/internal/third/lo"
	"github.com/snivilised/traverse/locale"
	"github.com/snivilised/traverse/pref"
	"github.com/snivilised/traverse/test/hydra"
)

var _ = Describe("feature", Ordered, func() {
	var (
		fS *luna.MemFS
	)

	BeforeAll(func() {
		const (
			verbose = false
		)

		fS = hydra.Nuxx(verbose, lab.Static.RetroWave, "edm")
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
						GetForest: func(_ string) *core.Forest {
							return &core.Forest{
								T: fS,
								R: nef.NewTraverseABS(),
							}
						},
					},
					tv.WithOnBegin(lab.Begin("🛡️")),
					tv.WithOnEnd(lab.End("🏁")),

					tv.WithSamplingOptions(&pref.SamplingOptions{
						Type: enums.SampleTypeSlice,
						NoOf: pref.EntryQuantities{
							Files:       2,
							Directories: 2,
						},
					}),
				)).Navigate(ctx)

				GinkgoWriter.Printf("===> 🍭 invoked '%v' directories, '%v' files.\n",
					result.Metrics().Count(enums.MetricNoDirectoriesInvoked),
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
					GetForest: func(_ string) *core.Forest {
						return &core.Forest{
							T: fS,
							R: nef.NewTraverseABS(),
						}
					},
				},
				tv.WithOnBegin(lab.Begin("🛡️")),
				tv.WithOnEnd(lab.End("🏁")),

				tv.WithSamplingOptions(&pref.SamplingOptions{
					Type:      entry.SampleType,
					InReverse: entry.Reverse,
					NoOf: pref.EntryQuantities{
						Files:       entry.NoOf.Files,
						Directories: entry.NoOf.Directories,
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
				FS:          fS,
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
					Files:       8,
					Directories: 8,
				},
			},
			SampleType: enums.SampleTypeSlice,
			NoOf: pref.EntryQuantities{
				Files: 2,
			},
		}),

		Entry(nil, &lab.SampleTE{
			NaviTE: lab.NaviTE{
				Given:        "universal(slice): first, with 2 directories",
				Should:       "invoke for at most 2 directories per directory",
				Subscription: enums.SubscribeUniversal,
				Prohibited:   []string{"Electric Youth"},
				ExpectedNoOf: lab.Quantities{
					Files:       11,
					Directories: 6,
				},
			},
			SampleType: enums.SampleTypeSlice,
			NoOf: pref.EntryQuantities{
				Directories: 2,
			},
		}),

		Entry(nil, &lab.SampleTE{
			NaviTE: lab.NaviTE{
				Given:        "universal(slice): first, with 2 files and 2 directories",
				Should:       "invoke for at most 2 files and 2 directories per directory",
				Subscription: enums.SubscribeUniversal,
				Prohibited:   []string{"cover.night-drive.jpg", "Electric Youth"},
				ExpectedNoOf: lab.Quantities{
					Files:       6,
					Directories: 6,
				},
			},
			SampleType: enums.SampleTypeSlice,
			NoOf: pref.EntryQuantities{
				Files:       2,
				Directories: 2,
			},
		}),

		Entry(nil, &lab.SampleTE{
			NaviTE: lab.NaviTE{
				Given:        "universal(filter): first, single file, first 2 directories",
				Should:       "invoke for at most single file per directory",
				Relative:     "edm",
				Subscription: enums.SubscribeUniversal,
				Prohibited:   []string{"02 - Swab.flac"},
				ExpectedNoOf: lab.Quantities{
					Files:       7,
					Directories: 14,
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
				Files:       1,
				Directories: 2,
			},
		}),

		Entry(nil, &lab.SampleTE{
			NaviTE: lab.NaviTE{
				Given:        "universal(filter): first, single file, first 2 directories",
				Should:       "invoke for at most single file per directory",
				Relative:     "edm",
				Subscription: enums.SubscribeUniversal,
				Prohibited:   []string{"02 - Swab.flac"},
				ExpectedNoOf: lab.Quantities{
					Files:       7,
					Directories: 14,
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
				Files:       1,
				Directories: 2,
			},
		}),

		Entry(nil, &lab.SampleTE{
			NaviTE: lab.NaviTE{
				Given:        "universal(filter): last, last single files, last 2 directories",
				Should:       "invoke for at most single file per directory",
				Relative:     "edm",
				Subscription: enums.SubscribeUniversal,
				Prohibited:   []string{"01 - Dre.flac"},
				ExpectedNoOf: lab.Quantities{
					Files:       8,
					Directories: 15,
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
				Files:       1,
				Directories: 2,
			},
		}),

		Entry(nil, &lab.SampleTE{
			NaviTE: lab.NaviTE{
				Given:        "universal(filter): last, last single files, last 2 directories",
				Should:       "invoke for at most single file per directory",
				Relative:     "edm",
				Subscription: enums.SubscribeUniversal,
				Prohibited:   []string{"01 - Dre.flac"},
				ExpectedNoOf: lab.Quantities{
					Files:       8,
					Directories: 15,
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
				Files:       1,
				Directories: 2,
			},
		}),

		// === directories ===================================================

		Entry(nil, &lab.SampleTE{
			NaviTE: lab.NaviTE{
				Given:        "directories(slice): first, with 2 directories",
				Should:       "invoke for at most 2 directories per directory",
				Subscription: enums.SubscribeDirectories,
				Prohibited:   []string{"Electric Youth"},
				ExpectedNoOf: lab.Quantities{
					Directories: 6,
				},
			},
			SampleType: enums.SampleTypeSlice,
			NoOf: pref.EntryQuantities{
				Directories: 2,
			},
		}),

		Entry(nil, &lab.SampleTE{
			NaviTE: lab.NaviTE{
				Given:        "directories(slice): last, with last single directory",
				Should:       "invoke for only last directory per directory",
				Subscription: enums.SubscribeDirectories,
				Prohibited:   []string{"Chromatics"},
				ExpectedNoOf: lab.Quantities{
					Directories: 3,
				},
			},
			SampleType: enums.SampleTypeSlice,
			Reverse:    true,
			NoOf: pref.EntryQuantities{
				Directories: 1,
			},
		}),

		Entry(nil, &lab.SampleTE{
			NaviTE: lab.NaviTE{
				Given:        "filtered directories(filter): first, with 2 directories that start with A",
				Should:       "invoke for at most 2 directories per directory",
				Relative:     "edm",
				Subscription: enums.SubscribeDirectories,
				Prohibited:   []string{"Tales Of Ephidrina"},
				ExpectedNoOf: lab.Quantities{
					// AMBIENT-TECHNO, Amorphous Androgynous, Aphex Twin
					Directories: 3,
				},
			},
			Filter: &lab.FilterTE{ // 🧄
				Description: "glob: items with that start with A",
				Type:        enums.FilterTypeGlob,
				Pattern:     "A*",
				Scope:       enums.ScopeDirectory,
			},
			SampleType: enums.SampleTypeFilter,
			NoOf: pref.EntryQuantities{
				Directories: 2,
			},
		}),

		Entry(nil, &lab.SampleTE{
			NaviTE: lab.NaviTE{
				Given:        "filtered directories(filter): first, with 2 directories that start with A",
				Should:       "invoke for at most 2 directories per directory",
				Relative:     "edm",
				Subscription: enums.SubscribeDirectories,
				Prohibited:   []string{"Tales Of Ephidrina"},
				ExpectedNoOf: lab.Quantities{
					// AMBIENT-TECHNO, Amorphous Androgynous, Aphex Twin
					Directories: 3,
				},
			},
			Filter: &lab.FilterTE{ // 🚀
				Description: "regex: items with that start with A",
				Type:        enums.FilterTypeRegex,
				Pattern:     "^A",
				Scope:       enums.ScopeDirectory,
			},
			SampleType: enums.SampleTypeFilter,
			NoOf: pref.EntryQuantities{
				Directories: 2,
			},
		}),

		Entry(nil, &lab.SampleTE{
			NaviTE: lab.NaviTE{
				Given:        "filtered directories(filter): last, with single directory that start with A",
				Should:       "invoke for at most a single directory per directory",
				Relative:     "edm",
				Subscription: enums.SubscribeDirectories,
				Prohibited:   []string{"Amorphous Androgynous"},
				ExpectedNoOf: lab.Quantities{
					Directories: 2,
				},
			},
			Filter: &lab.FilterTE{ // 🧄
				Description: "glob: items with that start with A",
				Type:        enums.FilterTypeGlob,
				Pattern:     "A*",
				Scope:       enums.ScopeDirectory,
			},
			SampleType: enums.SampleTypeFilter,
			Reverse:    true,
			NoOf: pref.EntryQuantities{
				Directories: 1,
			},
		}),

		Entry(nil, &lab.SampleTE{
			NaviTE: lab.NaviTE{
				Given:        "filtered directories(filter): last, with single directory that start with A",
				Should:       "invoke for at most a single directory per directory",
				Relative:     "edm",
				Subscription: enums.SubscribeDirectories,
				Prohibited:   []string{"Amorphous Androgynous"},
				ExpectedNoOf: lab.Quantities{
					Directories: 2,
				},
			},
			Filter: &lab.FilterTE{ // 🚀
				Description: "regex: items with that start with A",
				Type:        enums.FilterTypeRegex,
				Pattern:     "^A",
				Scope:       enums.ScopeDirectory,
			},
			SampleType: enums.SampleTypeFilter,
			Reverse:    true,
			NoOf: pref.EntryQuantities{
				Directories: 1,
			},
		}),

		// === directories with files ========================================

		Entry(nil, &lab.SampleTE{
			NaviTE: lab.NaviTE{
				Given:        "directories with files(slice): first, with 2 directories",
				Should:       "invoke for at most 2 directories per directory",
				Subscription: enums.SubscribeDirectoriesWithFiles,
				Prohibited:   []string{"Electric Youth"},
				ExpectedNoOf: lab.Quantities{
					Directories: 6,
					Children: map[string]int{
						"Night Drive":      4,
						"Northern Council": 4,
						"Teenage Color":    3,
					},
				},
			},
			SampleType: enums.SampleTypeSlice,
			NoOf: pref.EntryQuantities{
				Directories: 2,
			},
		}),

		Entry(nil, &lab.SampleTE{
			NaviTE: lab.NaviTE{
				Given:        "directories with files(slice): last, with last single directory",
				Should:       "invoke for only last directory per directory",
				Subscription: enums.SubscribeDirectoriesWithFiles,
				Prohibited:   []string{"Chromatics"},
				ExpectedNoOf: lab.Quantities{
					Directories: 3,
					Children: map[string]int{
						"Innerworld": 3,
					},
				},
			},
			SampleType: enums.SampleTypeSlice,
			Reverse:    true,
			NoOf: pref.EntryQuantities{
				Directories: 1,
			},
		}),

		Entry(nil, &lab.SampleTE{
			NaviTE: lab.NaviTE{
				Given:        "filtered directories with files(filter): last, with single directory that start with A",
				Should:       "invoke for at most a single directory per directory",
				Relative:     "edm",
				Subscription: enums.SubscribeDirectoriesWithFiles,
				Prohibited:   []string{"Amorphous Androgynous"},
				ExpectedNoOf: lab.Quantities{
					Directories: 2,
					Children:    map[string]int{},
				},
			},
			Filter: &lab.FilterTE{ // 🧄 this is directory filter, not child filter
				Description: "glob: items that start with A",
				Type:        enums.FilterTypeGlob,
				Pattern:     "A*",
				Scope:       enums.ScopeDirectory,
			},
			SampleType: enums.SampleTypeFilter,
			Reverse:    true,
			NoOf: pref.EntryQuantities{
				Directories: 1,
			},
		}),

		Entry(nil, &lab.SampleTE{
			NaviTE: lab.NaviTE{
				Given:        "filtered directories with files(filter): last, with single directory that start with A",
				Should:       "invoke for at most a single directory per directory",
				Relative:     "edm",
				Subscription: enums.SubscribeDirectoriesWithFiles,
				Prohibited:   []string{"Amorphous Androgynous"},
				ExpectedNoOf: lab.Quantities{
					Directories: 2,
					Children:    map[string]int{},
				},
			},
			Filter: &lab.FilterTE{ // 🚀
				Description: "regex: items that start with A",
				Type:        enums.FilterTypeRegex,
				Pattern:     "^A",
				Scope:       enums.ScopeDirectory,
			},
			SampleType: enums.SampleTypeFilter,
			Reverse:    true,
			NoOf: pref.EntryQuantities{
				Directories: 1,
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
		// file/directory scopes because a node's scope is determined after
		// a directory's contents are read, but sampling filter is
		// applied at the point the contents are read. Any scopes other
		// than file/directory are ignored.
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
				Given:        "universal(custom): first, single file, 2 directories",
				Should:       "invoke for at most single file per directory",
				Relative:     "edm",
				Subscription: enums.SubscribeUniversal,
				Prohibited:   []string{"02 - Swab.flac"},
				ExpectedNoOf: lab.Quantities{
					Files:       7,
					Directories: 14,
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
				Files:       1,
				Directories: 2,
			},
		}),

		Entry(nil, &lab.SampleTE{
			NaviTE: lab.NaviTE{
				Given:        "filtered directories(custom): last, single directory that starts with A",
				Should:       "invoke for at most a single directory per directory",
				Relative:     "edm",
				Subscription: enums.SubscribeDirectories,
				Prohibited:   []string{"Amorphous Androgynous"},
				ExpectedNoOf: lab.Quantities{
					Directories: 2,
				},
			},
			Filter: &lab.FilterTE{ // 🍒
				Type: enums.FilterTypeCustom,
				Sample: &customSamplingFilter{
					Sample:      tv.NewCustomSampleFilter(enums.ScopeDirectory),
					description: "custom(glob): items with A prefix",
					pattern:     "A*",
				},
			},
			SampleType: enums.SampleTypeCustom,
			NoOf: pref.EntryQuantities{
				Directories: 1,
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
				Given:        "directory spec, without no of directories",
				Should:       "return invalid directory spec error",
				Relative:     "edm/ELECTRONICA",
				Subscription: enums.SubscribeFiles,
				Prohibited:   []string{"03 - Mountain Goat.flac"},
				ExpectedNoOf: lab.Quantities{
					Files: 24,
				},
				ExpectedErr: locale.ErrInvalidSamplingSpecMissingDirectories,
			},
			Filter: &lab.FilterTE{ // 🧄
				Description: "glob: items with .flac suffix",
				Type:        enums.FilterTypeGlob,
				Pattern:     "*.flac",
				Scope:       enums.ScopeDirectory,
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
				Directories: 2,
			},
		}),

		Entry(nil, &lab.SampleTE{
			NaviTE: lab.NaviTE{
				Given:        "custom filter not defined",
				Should:       "fail",
				Relative:     "edm",
				Subscription: enums.SubscribeUniversal,
				ExpectedErr:  locale.ErrFilterIsNil,
			},
			Filter: &lab.FilterTE{ // 🍒
				Type: enums.FilterTypeCustom,
			},
			SampleType: enums.SampleTypeCustom,
		}),

		Entry(nil, &lab.SampleTE{
			NaviTE: lab.NaviTE{
				Given:        "filter missing type",
				Should:       "fail",
				Relative:     "edm",
				Subscription: enums.SubscribeUniversal,
				ExpectedErr:  locale.ErrFilterMissingType,
			},
			Filter:     &lab.FilterTE{},
			SampleType: enums.SampleTypeCustom,
		}),
	)
})
