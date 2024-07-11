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
	"github.com/snivilised/traverse/internal/lo"
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
			filepath.Join("MUSICO", "rock"),
		)
		Expect(root).NotTo(BeEmpty())
	})

	BeforeEach(func() {
		services.Reset()
	})

	DescribeTable("folders with files filtered",
		func(ctx SpecContext, entry *filterTE) {
			recording := make(recordingMap)
			filterDefs := &pref.FilterOptions{
				Node: &core.FilterDef{
					Type:            enums.FilterTypeExtendedGlob,
					Description:     entry.name,
					Pattern:         entry.pattern,
					Scope:           entry.scope,
					Negate:          entry.negate,
					IfNotApplicable: entry.ifNotApplicable,
				},
			}
			var traverseFilter core.TraverseFilter

			path := helpers.Path(root, entry.relative)

			callback := func(node *core.Node) error {
				indicator := lo.Ternary(node.IsFolder(), "📁", "💠")
				GinkgoWriter.Printf(
					"===> %v Glob Filter(%v) source: '%v', item-name: '%v', item-scope(fs): '%v(%v)'\n",
					indicator,
					traverseFilter.Description(),
					traverseFilter.Source(),
					node.Extension.Name,
					node.Extension.Scope,
					traverseFilter.Scope(),
				)
				if lo.Contains(entry.mandatory, node.Extension.Name) {
					Expect(node).Should(MatchCurrentExtendedFilter(traverseFilter))
				}

				recording[node.Extension.Name] = len(node.Children)
				return nil
			}
			result, err := tv.Walk().Configure().Extent(tv.Prime(
				&tv.Using{
					Root:         path,
					Subscription: entry.subscription,
					Handler:      callback,
					GetFS: func() fs.FS {
						return vfs
					},
				},
				tv.WithFilter(filterDefs),
				tv.WithFilterSink(func(reply pref.FilterReply) {
					traverseFilter = reply.Node
				}),
				tv.WithHookQueryStatus(func(path string) (fs.FileInfo, error) {
					return vfs.Stat(helpers.TrimRoot(path))
				}),
				tv.WithHookReadDirectory(func(_ fs.FS, dirname string) ([]fs.DirEntry, error) {
					return vfs.ReadDir(helpers.TrimRoot(dirname))
				}),
			)).Navigate(ctx)

			assertFilteredNavigation(entry, testOptions{
				vfs:       vfs,
				recording: recording,
				path:      path,
				result:    result,
				err:       err,
			})
		},

		func(entry *filterTE) string {
			return fmt.Sprintf("🧪 ===> given: '%v'", entry.given)
		},

		// === universal =====================================================

		Entry(nil, &filterTE{
			naviTE: naviTE{
				given:        "universal(any scope): extended glob filter",
				relative:     "rock/PROGRESSIVE-ROCK/Marillion",
				subscription: enums.SubscribeUniversal,
				expectedNoOf: quantities{
					files:   16,
					folders: 5,
				},
				prohibited: []string{"cover-clutching-at-straws-jpg"},
			},
			name:    "it6ems with 'flac' suffix",
			pattern: "*|flac",
			scope:   enums.ScopeAll,
		}),

		Entry(nil, &filterTE{
			naviTE: naviTE{
				given:        "universal(any scope): extended glob filter, with dot extension",
				relative:     "rock/PROGRESSIVE-ROCK/Marillion",
				subscription: enums.SubscribeUniversal,
				expectedNoOf: quantities{
					files:   16,
					folders: 5,
				},
				prohibited: []string{"cover-clutching-at-straws-jpg"},
			},
			name:    "items with 'flac' suffix",
			pattern: "*|.flac",
			scope:   enums.ScopeAll,
		}),

		Entry(nil, &filterTE{
			naviTE: naviTE{
				given:        "universal(any scope): extended glob filter, with multiple extensions",
				relative:     "rock/PROGRESSIVE-ROCK/Marillion",
				subscription: enums.SubscribeUniversal,
				expectedNoOf: quantities{
					files:   19,
					folders: 5,
				},
				mandatory:  []string{"front.jpg"},
				prohibited: []string{"cover-clutching-at-straws-jpg"},
			},
			name:    "items with 'flac' suffix",
			pattern: "*|flac,jpg",
			scope:   enums.ScopeAll,
		}),

		Entry(nil, &filterTE{
			naviTE: naviTE{
				given:        "universal(any scope): extended glob filter, without extension",
				relative:     "rock/PROGRESSIVE-ROCK/Marillion",
				subscription: enums.SubscribeUniversal,
				expectedNoOf: quantities{
					files:   3,
					folders: 5,
				},
				mandatory:  []string{"cover-clutching-at-straws-jpg"},
				prohibited: []string{"01 - Hotel Hobbies.flac"},
			},
			name:    "items with 'flac' suffix",
			pattern: "*|",
			scope:   enums.ScopeAll,
		}),

		Entry(nil, &filterTE{
			naviTE: naviTE{
				given:        "universal(file scope): extended glob filter (negate)",
				relative:     "rock/PROGRESSIVE-ROCK/Marillion",
				subscription: enums.SubscribeUniversal,
				expectedNoOf: quantities{
					files:   7,
					folders: 5,
				},
				prohibited: []string{"01 - Hotel Hobbies.flac"},
			},
			name:    "files without .flac suffix",
			pattern: "*|flac",
			scope:   enums.ScopeFile,
			negate:  true,
		}),

		Entry(nil, &filterTE{
			naviTE: naviTE{
				given:        "universal(undefined scope): extended glob filter",
				relative:     "rock/PROGRESSIVE-ROCK/Marillion",
				subscription: enums.SubscribeUniversal,
				expectedNoOf: quantities{
					files:   16,
					folders: 5,
				},
				prohibited: []string{"cover-clutching-at-straws-jpg"},
			},
			name:    "items with '.flac' suffix",
			pattern: "*|flac",
		}),

		Entry(nil, &filterTE{
			naviTE: naviTE{
				given:        "universal(any scope): extended glob filter, any extension",
				relative:     "rock/PROGRESSIVE-ROCK/Marillion",
				subscription: enums.SubscribeUniversal,
				expectedNoOf: quantities{
					files:   4,
					folders: 1,
				},
				mandatory:  []string{"cover-clutching-at-straws-jpg"},
				prohibited: []string{"01 - Hotel Hobbies.flac"},
			},
			name:    "starts with c, any extension",
			pattern: "c*|*",
			scope:   enums.ScopeAll,
		}),

		// === folders =======================================================

		Entry(nil, &filterTE{
			naviTE: naviTE{
				given:        "folders(any scope): extended glob filter",
				relative:     "rock/PROGRESSIVE-ROCK/Marillion",
				subscription: enums.SubscribeFolders,
				expectedNoOf: quantities{
					files:   0,
					folders: 2,
				},
				mandatory:  []string{"Marillion"},
				prohibited: []string{"Fugazi"},
			},
			name:    "folders starting with M",
			pattern: "M*|",
			scope:   enums.ScopeFolder,
		}),

		Entry(nil, &filterTE{
			naviTE: naviTE{
				given:        "folders(folder scope): extended glob filter (negate)",
				relative:     "rock/PROGRESSIVE-ROCK/Marillion",
				subscription: enums.SubscribeFolders,
				expectedNoOf: quantities{
					files:   0,
					folders: 3,
				},
				mandatory:  []string{"Fugazi"},
				prohibited: []string{"Marillion"},
			},
			name:    "folders NOT starting with M",
			pattern: "M*|",
			scope:   enums.ScopeFolder,
			negate:  true,
		}),

		Entry(nil, &filterTE{
			naviTE: naviTE{
				given:        "universal(undefined scope): extended glob filter",
				relative:     "rock/PROGRESSIVE-ROCK/Marillion",
				subscription: enums.SubscribeFolders,
				expectedNoOf: quantities{
					files:   0,
					folders: 2,
				},
				mandatory:  []string{"Marillion"},
				prohibited: []string{"Fugazi"},
			},
			name:    "folders starting with M",
			pattern: "M*|",
		}),

		// === files =========================================================

		Entry(nil, &filterTE{
			naviTE: naviTE{
				given:        "files(file scope): extended glob filter",
				relative:     "rock/PROGRESSIVE-ROCK/Marillion",
				subscription: enums.SubscribeFiles,
				expectedNoOf: quantities{
					files:   16,
					folders: 0,
				},
				mandatory:  []string{"01 - Hotel Hobbies.flac"},
				prohibited: []string{"cover-clutching-at-straws-jpg"},
			},
			name:    "items with 'flac' suffix",
			pattern: "*|flac",
			scope:   enums.ScopeFile,
		}),

		Entry(nil, &filterTE{
			naviTE: naviTE{
				given:        "files(any scope): extended glob filter, with dot extension",
				relative:     "rock/PROGRESSIVE-ROCK/Marillion",
				subscription: enums.SubscribeFiles,
				expectedNoOf: quantities{
					files:   16,
					folders: 0,
				},
				mandatory:  []string{"01 - Hotel Hobbies.flac"},
				prohibited: []string{"cover-clutching-at-straws-jpg"},
			},
			name:    "items with 'flac' suffix",
			pattern: "*|.flac",
			scope:   enums.ScopeFile,
		}),

		Entry(nil, &filterTE{
			naviTE: naviTE{
				given:        "files(file scope): extended glob filter, with multiple extensions",
				relative:     "rock/PROGRESSIVE-ROCK/Marillion",
				subscription: enums.SubscribeFiles,
				expectedNoOf: quantities{
					files:   19,
					folders: 0,
				},
				mandatory:  []string{"front.jpg"},
				prohibited: []string{"cover-clutching-at-straws-jpg"},
			},
			name:    "items with 'flac' suffix",
			pattern: "*|flac,jpg",
			scope:   enums.ScopeFile,
		}),

		Entry(nil, &filterTE{
			naviTE: naviTE{
				given:        "file(file scope): extended glob filter, without extension",
				relative:     "rock/PROGRESSIVE-ROCK/Marillion",
				subscription: enums.SubscribeFiles,
				expectedNoOf: quantities{
					files:   3,
					folders: 0,
				},
				mandatory:  []string{"cover-clutching-at-straws-jpg"},
				prohibited: []string{"01 - Hotel Hobbies.flac"},
			},
			name:    "items with 'flac' suffix",
			pattern: "*|",
			scope:   enums.ScopeFile,
		}),

		Entry(nil, &filterTE{
			naviTE: naviTE{
				given:        "file(file scope): extended glob filter (negate)",
				relative:     "rock/PROGRESSIVE-ROCK/Marillion",
				subscription: enums.SubscribeFiles,
				expectedNoOf: quantities{
					files:   7,
					folders: 0,
				},
				mandatory:  []string{"cover-clutching-at-straws-jpg"},
				prohibited: []string{"01 - Hotel Hobbies.flac"},
			},
			name:    "files without .flac suffix",
			pattern: "*|flac",
			scope:   enums.ScopeFile,
			negate:  true,
		}),

		Entry(nil, &filterTE{
			naviTE: naviTE{
				given:        "file(undefined scope): extended glob filter",
				relative:     "rock/PROGRESSIVE-ROCK/Marillion",
				subscription: enums.SubscribeFiles,
				expectedNoOf: quantities{
					files:   16,
					folders: 0,
				},
				mandatory:  []string{"01 - Hotel Hobbies.flac"},
				prohibited: []string{"cover-clutching-at-straws-jpg"},
			},
			name:    "items with '.flac' suffix",
			pattern: "*|flac",
		}),

		Entry(nil, &filterTE{
			naviTE: naviTE{
				given:        "file(any scope): extended glob filter, any extension",
				relative:     "rock/PROGRESSIVE-ROCK/Marillion",
				subscription: enums.SubscribeFiles,
				expectedNoOf: quantities{
					files:   4,
					folders: 0,
				},
				mandatory:  []string{"cover-clutching-at-straws-jpg"},
				prohibited: []string{"01 - Hotel Hobbies.flac"},
			},
			name:    "starts with c, any extension",
			pattern: "c*|*",
			scope:   enums.ScopeAll,
		}),

		// === ifNotApplicable ===============================================

		Entry(nil, &filterTE{
			naviTE: naviTE{
				given:        "universal(leaf scope): extended glob filter (ifNotApplicable=true)",
				relative:     "rock/PROGRESSIVE-ROCK/Marillion",
				subscription: enums.SubscribeUniversal,
				expectedNoOf: quantities{
					files:   16,
					folders: 5,
				},
				mandatory:  []string{"Marillion"},
				prohibited: []string{"cover-clutching-at-straws-jpg"},
			},
			name:            "leaf items with 'flac' suffix",
			pattern:         "*|flac",
			scope:           enums.ScopeLeaf,
			ifNotApplicable: enums.TriStateBoolTrue,
		}),

		Entry(nil, &filterTE{
			naviTE: naviTE{
				given:        "universal(leaf scope): extended glob filter (ifNotApplicable=false)",
				relative:     "rock/PROGRESSIVE-ROCK/Marillion",
				subscription: enums.SubscribeUniversal,
				expectedNoOf: quantities{
					files:   16,
					folders: 4,
				},
				prohibited: []string{"Marillion"},
			},
			name:            "items with '.flac' suffix",
			pattern:         "*|flac",
			scope:           enums.ScopeLeaf,
			ifNotApplicable: enums.TriStateBoolFalse,
		}),

		// === with-exclusion ================================================

		Entry(nil, &filterTE{
			naviTE: naviTE{
				given:        "universal(any scope): extended glob filter with exclusion",
				relative:     "rock/PROGRESSIVE-ROCK/Marillion",
				subscription: enums.SubscribeFiles,
				expectedNoOf: quantities{
					files:   12,
					folders: 0,
				},
				prohibited: []string{"01 - Hotel Hobbies.flac"},
			},
			name:    "files starting with 0, except 01 items and flac suffix",
			pattern: "0*/*01*|flac",
			scope:   enums.ScopeFile,
		}),
	)
})
