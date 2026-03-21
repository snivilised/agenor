package sampling_test

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/snivilised/jaywalk/src/agenor"
	"github.com/snivilised/jaywalk/src/agenor/core"
	"github.com/snivilised/jaywalk/src/agenor/enums"
	lab "github.com/snivilised/jaywalk/src/agenor/internal/laboratory"
	"github.com/snivilised/jaywalk/src/internal/services"
	"github.com/snivilised/jaywalk/src/internal/third/lo"
	"github.com/snivilised/jaywalk/locale"
	"github.com/snivilised/jaywalk/src/agenor/pref"
	"github.com/snivilised/jaywalk/src/agenor/test/hanno"
	"github.com/snivilised/jaywalk/src/agenor/tfs"
	"github.com/snivilised/li18ngo"
	"github.com/snivilised/nefilim/test/luna"
)

var _ = Describe("feature", Ordered, func() {
	var (
		fS *luna.MemFS
	)

	BeforeAll(func() {
		const (
			verbose = false
		)

		fS = hanno.Nuxx(verbose, lab.Static.RetroWave, "edm")

		Expect(li18ngo.Use()).To(Succeed())
	})

	BeforeEach(func() {
		services.Reset()
	})

	Context("comprehension", func() {
		When("universal: slice sample", func() {
			It("🧪 should: foo", Label("example"), func(ctx SpecContext) {
				path := lab.Static.RetroWave
				result, _ := agenor.Walk().Configure().Extent(agenor.Prime(
					&pref.Using{
						Subscription: enums.SubscribeUniversal,
						Head: pref.Head{
							Handler: func(servant agenor.Servant) error {
								node := servant.Node()
								GinkgoWriter.Printf(
									"---> 🍯 EXAMPLE-SAMPLE-CALLBACK: '%v'\n", node.Path,
								)

								return nil
							},
							GetForest: func(_ string) *core.Forest {
								return &core.Forest{
									T: fS,
									R: tfs.New(),
								}
							},
						},
						Tree: path,
					},
					agenor.WithOnBegin(lab.Begin("🛡️")),
					agenor.WithOnEnd(lab.End("🏁")),

					agenor.WithSamplingOptions(&pref.SamplingOptions{
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
			recall := make(lab.Recall)
			once := func(node *agenor.Node) error { //nolint:unparam // return nil error ok
				_, found := recall[node.Extension.Name]
				Expect(found).To(BeFalse())

				recall[node.Extension.Name] = len(node.Children)

				return nil
			}

			path := lo.Ternary(entry.Relative == "",
				lab.Static.RetroWave,
				entry.Relative,
			)

			callback := func(servant agenor.Servant) error {
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

			result, err := agenor.Walk().Configure().Extent(agenor.Prime(
				&pref.Using{
					Subscription: entry.Subscription,
					Head: pref.Head{
						Handler: callback,
						GetForest: func(_ string) *core.Forest {
							return &core.Forest{
								T: fS,
								R: tfs.New(),
							}
						},
					},
					Tree: path,
				},
				agenor.WithOnBegin(lab.Begin("🛡️")),
				agenor.WithOnEnd(lab.End("🏁")),

				agenor.WithSamplingOptions(&pref.SamplingOptions{
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
				agenor.IfOptionF(entry.Filter != nil, func() pref.Option {
					return agenor.WithFilter(&pref.FilterOptions{
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
				agenor.IfOption(entry.CaseSensitive, agenor.WithHookCaseSensitiveSort()),
			)).Navigate(ctx)

			lab.AssertNavigation(&entry.NaviTE, &lab.TestOptions{
				FS:          fS,
				Recording:   recall,
				Path:        path,
				Result:      result,
				Err:         err,
				ExpectedErr: entry.ExpectedErr,
			})
		},
		lab.FormatSampleTestDescription,

		// === universal =====================================================

		Entry(nil, &lab.SampleTE{
			DescribedTE: lab.DescribedTE{
				Given:  "universal(slice): first, with 2 files",
				Should: "invoke for at most 2 files per directory",
			},
			NaviTE: lab.NaviTE{
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
			DescribedTE: lab.DescribedTE{
				Given:  "universal(slice): first, with 2 directories",
				Should: "invoke for at most 2 directories per directory",
			},
			NaviTE: lab.NaviTE{
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
			DescribedTE: lab.DescribedTE{
				Given:  "universal(slice): first, with 2 files and 2 directories",
				Should: "invoke for at most 2 files and 2 directories per directory",
			},
			NaviTE: lab.NaviTE{
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
			DescribedTE: lab.DescribedTE{
				Given:  "universal(filter): first, single file, first 2 directories",
				Should: "invoke for at most single file per directory",
			},
			NaviTE: lab.NaviTE{
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
			DescribedTE: lab.DescribedTE{
				Given:  "universal(filter): first, single file, first 2 directories",
				Should: "invoke for at most single file per directory",
			},
			NaviTE: lab.NaviTE{
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
			DescribedTE: lab.DescribedTE{
				Given:  "universal(filter): last, last single files, last 2 directories",
				Should: "invoke for at most single file per directory",
			},
			NaviTE: lab.NaviTE{
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
			DescribedTE: lab.DescribedTE{
				Given:  "universal(filter): last, last single files, last 2 directories",
				Should: "invoke for at most single file per directory",
			},
			NaviTE: lab.NaviTE{
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
			DescribedTE: lab.DescribedTE{
				Given:  "directories(slice): first, with 2 directories",
				Should: "invoke for at most 2 directories per directory",
			},
			NaviTE: lab.NaviTE{
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
			DescribedTE: lab.DescribedTE{
				Given:  "directories(slice): last, with last single directory",
				Should: "invoke for only last directory per directory",
			},
			NaviTE: lab.NaviTE{
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
			DescribedTE: lab.DescribedTE{
				Given:  "filtered directories(filter): first, with 2 directories that start with A",
				Should: "invoke for at most 2 directories per directory",
			},
			NaviTE: lab.NaviTE{
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
			DescribedTE: lab.DescribedTE{
				Given:  "filtered directories(filter): first, with 2 directories that start with A",
				Should: "invoke for at most 2 directories per directory",
			},
			NaviTE: lab.NaviTE{
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
			DescribedTE: lab.DescribedTE{
				Given:  "filtered directories(filter): last, with single directory that start with A",
				Should: "invoke for at most a single directory per directory",
			},
			NaviTE: lab.NaviTE{
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
			DescribedTE: lab.DescribedTE{
				Given:  "filtered directories(filter): last, with single directory that start with A",
				Should: "invoke for at most a single directory per directory",
			},
			NaviTE: lab.NaviTE{
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
			DescribedTE: lab.DescribedTE{
				Given:  "directories with files(slice): first, with 2 directories",
				Should: "invoke for at most 2 directories per directory",
			},
			NaviTE: lab.NaviTE{
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
			DescribedTE: lab.DescribedTE{
				Given:  "directories with files(slice): last, with last single directory",
				Should: "invoke for only last directory per directory",
			},
			NaviTE: lab.NaviTE{
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
			DescribedTE: lab.DescribedTE{
				Given:  "filtered directories with files(filter): last, with single directory that start with A",
				Should: "invoke for at most a single directory per directory",
			},
			NaviTE: lab.NaviTE{
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
			DescribedTE: lab.DescribedTE{
				Given:  "filtered directories with files(filter): last, with single directory that start with A",
				Should: "invoke for at most a single directory per directory",
			},
			NaviTE: lab.NaviTE{
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
			DescribedTE: lab.DescribedTE{
				Given:  "files(slice): first, with 2 files",
				Should: "invoke for at most 2 files per directory",
			},
			NaviTE: lab.NaviTE{
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
			DescribedTE: lab.DescribedTE{
				Given:  "files(slice): last, with last single file",
				Should: "invoke for only last file per directory",
			},
			NaviTE: lab.NaviTE{
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
			DescribedTE: lab.DescribedTE{
				Given:  "filtered files(filter): first, 2 files",
				Should: "invoke for at most 2 files per directory",
			},
			NaviTE: lab.NaviTE{
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
			DescribedTE: lab.DescribedTE{
				Given:  "filtered files(filter): first, 2 files",
				Should: "invoke for at most 2 files per directory",
			},
			NaviTE: lab.NaviTE{
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
			DescribedTE: lab.DescribedTE{
				Given:  "filtered files(filter): last, last 2 files",
				Should: "invoke for at most 2 files per directory",
			},
			NaviTE: lab.NaviTE{
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
			DescribedTE: lab.DescribedTE{
				Given:  "filtered files(filter): last, last 2 files",
				Should: "invoke for at most 2 files per directory",
			},
			NaviTE: lab.NaviTE{
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
			DescribedTE: lab.DescribedTE{
				Given:  "universal(custom): first, single file, 2 directories",
				Should: "invoke for at most single file per directory",
			},
			NaviTE: lab.NaviTE{
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
					Sample:      agenor.NewCustomSampleFilter(enums.ScopeFile),
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
			DescribedTE: lab.DescribedTE{
				Given:  "filtered directories(custom): last, single directory that starts with A",
				Should: "invoke for at most a single directory per directory",
			},
			NaviTE: lab.NaviTE{
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
					Sample:      agenor.NewCustomSampleFilter(enums.ScopeDirectory),
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
			DescribedTE: lab.DescribedTE{
				Given:  "filtered files(custom): last, last 2 files",
				Should: "invoke for at most 2 files per directory",
			},
			NaviTE: lab.NaviTE{
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
					Sample:      agenor.NewCustomSampleFilter(enums.ScopeFile),
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
			DescribedTE: lab.DescribedTE{
				Given:  "directory spec, without no of directories",
				Should: "return invalid directory spec error",
			},
			NaviTE: lab.NaviTE{
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
			DescribedTE: lab.DescribedTE{
				Given:  "file spec, without no of files",
				Should: "return invalid file spec error",
			},
			NaviTE: lab.NaviTE{
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
			DescribedTE: lab.DescribedTE{
				Given:  "custom filter not defined",
				Should: "fail",
			},
			NaviTE: lab.NaviTE{
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
			DescribedTE: lab.DescribedTE{
				Given:  "filter missing type",
				Should: "fail",
			},
			NaviTE: lab.NaviTE{
				Relative:     "edm",
				Subscription: enums.SubscribeUniversal,
				ExpectedErr:  locale.ErrFilterMissingType,
			},
			Filter:     &lab.FilterTE{},
			SampleType: enums.SampleTypeCustom,
		}),
	)
})
