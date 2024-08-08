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

var _ = Describe("NavigatorFilterPoly", Ordered, func() {
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

	DescribeTable("poly-filter",
		func(ctx SpecContext, entry *polyTE) {
			var (
				traverseFilter core.TraverseFilter
			)

			recording := make(recordingMap)
			filterDefs := &pref.FilterOptions{
				Node: &core.FilterDef{
					Type: enums.FilterTypePoly,
					Poly: &core.PolyFilterDef{
						File:   entry.file,
						Folder: entry.folder,
					},
				},
				Sink: func(reply pref.FilterReply) {
					traverseFilter = reply.Node
				},
			}

			path := helpers.Path(root, entry.relative)

			callback := func(node *core.Node) error {
				indicator := lo.Ternary(node.IsFolder(), "📁", "💠")
				GinkgoWriter.Printf(
					"===> %v Poly Filter(%v) source: '%v', item-name: '%v', item-scope(fs): '%v(%v)'\n",
					indicator,
					traverseFilter.Description(),
					traverseFilter.Source(),
					node.Extension.Name,
					node.Extension.Scope,
					traverseFilter.Scope(),
				)
				if lo.Contains(entry.mandatory, node.Extension.Name) {
					Expect(node).Should(MatchCurrentGlobFilter(traverseFilter))
				}

				recording[node.Extension.Name] = len(node.Children)
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

			assertNavigation(&entry.naviTE, &testOptions{
				vfs:       vfs,
				recording: recording,
				path:      path,
				result:    result,
				err:       err,
			})
		},
		func(entry *polyTE) string {
			return fmt.Sprintf("🧪 ===> given: '%v'", entry.given)
		},

		// === universal(file:regex; folder:glob) ============================

		Entry(nil, &polyTE{
			naviTE: naviTE{
				given:        "poly - files:regex; folders:glob",
				relative:     "RETRO-WAVE",
				subscription: enums.SubscribeUniversal,
				expectedNoOf: quantities{
					// file is 2 not 3 because *i* is case sensitive so Innerworld is not a match
					// The next(not this one) regex test case, fixes this because folder regex has better
					// control over case sensitivity
					files:   2,
					folders: 8,
				},
			},
			file: core.FilterDef{
				Type:        enums.FilterTypeRegex,
				Description: "files: starts with vinyl",
				Pattern:     "^vinyl",
				Scope:       enums.ScopeFile,
			},
			folder: core.FilterDef{
				Type:        enums.FilterTypeGlob,
				Description: "folders: contains i (case sensitive)",
				Pattern:     "*i*",
				Scope:       enums.ScopeFolder | enums.ScopeLeaf,
			},
		}),

		// === universal(file:regex; folder:regex) ===========================

		Entry(nil, &polyTE{
			naviTE: naviTE{
				given:        "poly - files:regex; folders:regex",
				relative:     "RETRO-WAVE",
				subscription: enums.SubscribeUniversal,
				expectedNoOf: quantities{
					files:   3,
					folders: 8,
				},
			},
			file: core.FilterDef{
				Type:        enums.FilterTypeRegex,
				Description: "files: starts with vinyl",
				Pattern:     "^vinyl",
				Scope:       enums.ScopeFile,
			},
			folder: core.FilterDef{
				Type:        enums.FilterTypeRegex,
				Description: "folders: contains i (case insensitive)",
				Pattern:     "[iI]",
				Scope:       enums.ScopeFolder | enums.ScopeLeaf,
			},
		}),

		// === universal(file:extended-glob; folder:glob) ====================

		Entry(nil, &polyTE{
			naviTE: naviTE{
				given:        "poly - files:extended-glob; folders:glob",
				relative:     "RETRO-WAVE",
				subscription: enums.SubscribeUniversal,
				expectedNoOf: quantities{
					// file is 2 not 3 because *i* is case sensitive so Innerworld is not a match
					// The next 2 tests regex/extended-glob test case, fixes this because they
					// have better control over case sensitivity
					//
					files:   2,
					folders: 8,
				},
			},
			file: core.FilterDef{
				Type:        enums.FilterTypeExtendedGlob,
				Description: "files: txt files starting with vinyl",
				Pattern:     "vinyl*|txt",
				Scope:       enums.ScopeFile,
			},
			folder: core.FilterDef{
				Type:        enums.FilterTypeGlob,
				Description: "folders: contains i (case sensitive)",
				Pattern:     "*i*",
				Scope:       enums.ScopeFolder | enums.ScopeLeaf,
			},
		}),

		Entry(nil, &polyTE{
			naviTE: naviTE{
				given:        "poly - files:extended-glob; folders:regex",
				relative:     "RETRO-WAVE",
				subscription: enums.SubscribeUniversal,
				expectedNoOf: quantities{
					files:   3,
					folders: 8,
				},
			},
			file: core.FilterDef{
				Type:        enums.FilterTypeExtendedGlob,
				Description: "files: txt files starting with vinyl",
				Pattern:     "vinyl*|txt",
				Scope:       enums.ScopeFile,
			},
			folder: core.FilterDef{
				Type:        enums.FilterTypeRegex,
				Description: "folders: contains i (case sensitive)",
				Pattern:     "[iI]",
				Scope:       enums.ScopeFolder | enums.ScopeLeaf,
			},
		}),

		Entry(nil, &polyTE{
			naviTE: naviTE{
				given:        "poly - files:extended-glob; folders:extended-glob",
				relative:     "RETRO-WAVE",
				subscription: enums.SubscribeUniversal,
				expectedNoOf: quantities{
					files:   3,
					folders: 8,
				},
			},
			file: core.FilterDef{
				Type:        enums.FilterTypeExtendedGlob,
				Description: "files: txt files starting with vinyl",
				Pattern:     "vinyl*|txt",
				Scope:       enums.ScopeFile,
			},
			folder: core.FilterDef{
				Type:        enums.FilterTypeExtendedGlob,
				Description: "folders: contains i (case sensitive)",
				Pattern:     "*i*|",
				Scope:       enums.ScopeFolder | enums.ScopeLeaf,
			},
		}),

		// For the poly filter, the file/folder scopes must be set correctly, but because
		// they can be set automatically, the client is not forced to set them. This test
		// checks that when the file/folder scopes are not set, then poly filtering still works
		// properly.
		Entry(nil, &polyTE{
			naviTE: naviTE{
				given:        "poly(scopes omitted) - files:regex; folders:regex",
				relative:     "RETRO-WAVE",
				subscription: enums.SubscribeUniversal,
				expectedNoOf: quantities{
					files:   3,
					folders: 8,
				},
			},
			file: core.FilterDef{
				Type:        enums.FilterTypeRegex,
				Description: "files: starts with vinyl",
				Pattern:     "^vinyl",
				// file scope omitted
			},
			folder: core.FilterDef{
				Type:        enums.FilterTypeRegex,
				Description: "folders: contains i (case insensitive)",
				Pattern:     "[iI]",
				Scope:       enums.ScopeLeaf, // folder scope omitted
			},
		}),

		// === files (file:regex; folder:regex) ==============================

		Entry(nil, &polyTE{
			naviTE: naviTE{
				given:        "poly(subscribe:files)",
				relative:     "RETRO-WAVE",
				subscription: enums.SubscribeFiles,
				expectedNoOf: quantities{
					files:   3,
					folders: 0,
				},
			},
			file: core.FilterDef{
				Type:        enums.FilterTypeRegex,
				Description: "files: starts with vinyl",
				Pattern:     "^vinyl",
			},
			folder: core.FilterDef{
				Type:        enums.FilterTypeRegex,
				Description: "folders: contains i",
				Pattern:     "[iI]",
				Scope:       enums.ScopeLeaf,
			},
		}),
	)
})
