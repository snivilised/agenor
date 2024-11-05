package filtering_test

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
	"github.com/snivilised/traverse/pref"
	"github.com/snivilised/traverse/test/hydra"
)

var _ = Describe("NavigatorFilterCustom", Ordered, func() {
	var (
		fS *luna.MemFS
	)

	BeforeAll(func() {
		const (
			verbose = false
		)

		fS = hydra.Nuxx(verbose, lab.Static.RetroWave)
		Expect(li18ngo.Use()).To(Succeed())
	})

	BeforeEach(func() {
		services.Reset()
	})

	DescribeTable("custom-filter (glob)",
		func(ctx SpecContext, entry *lab.FilterTE) {
			recall := make(lab.Recall)
			customFilter := &customFilter{
				name:    entry.Description,
				pattern: entry.Pattern,
				scope:   entry.Scope,
				negate:  entry.Negate,
			}

			path := entry.Relative
			callback := func(servant tv.Servant) error {
				node := servant.Node()
				indicator := lo.Ternary(node.IsDirectory(), "📁", "💠")
				GinkgoWriter.Printf(
					"===> %v Glob Filter(%v) source: '%v', node-name: '%v', node-scope(fs): '%v(%v)'\n",
					indicator,
					customFilter.Description(),
					customFilter.Source(),
					node.Extension.Name,
					node.Extension.Scope,
					customFilter.Scope(),
				)
				if lo.Contains(entry.Mandatory, node.Extension.Name) {
					Expect(node).Should(MatchCurrentCustomFilter(customFilter))
				}

				recall[node.Extension.Name] = len(node.Children)
				return nil
			}
			result, err := tv.Walk().Configure().Extent(tv.Prime(
				&pref.Using{
					Head: pref.Head{
						Subscription: entry.Subscription,
						Handler:      callback,
						GetForest: func(_ string) *core.Forest {
							return &core.Forest{
								T: fS,
								R: nef.NewTraverseABS(),
							}
						},
					},
					Tree: path,
				},
				tv.WithOnBegin(lab.Begin("🛡️")),
				tv.WithOnEnd(lab.End("🏁")),

				tv.WithFilter(&pref.FilterOptions{
					Custom: customFilter,
				}),
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
		func(entry *lab.FilterTE) string {
			return fmt.Sprintf("🧪 ===> given: '%v'", entry.Given)
		},

		// === universal =====================================================

		Entry(nil, &lab.FilterTE{
			NaviTE: lab.NaviTE{
				Given:        "universal(any scope): custom filter",
				Relative:     lab.Static.RetroWave,
				Subscription: enums.SubscribeUniversal,
				ExpectedNoOf: lab.Quantities{
					Files:       8,
					Directories: 0,
				},
			},
			Description: "items with '.flac' suffix",
			Pattern:     "*.flac",
			Scope:       enums.ScopeFile,
		}),

		Entry(nil, &lab.FilterTE{
			NaviTE: lab.NaviTE{
				Given:        "universal(any scope): custom filter (negate)",
				Relative:     lab.Static.RetroWave,
				Subscription: enums.SubscribeUniversal,
				ExpectedNoOf: lab.Quantities{
					Files:       6,
					Directories: 8,
				},
			},
			Description: "items without .flac suffix",
			Pattern:     "*.flac",
			Scope:       enums.ScopeAll,
			Negate:      true,
		}),

		Entry(nil, &lab.FilterTE{
			NaviTE: lab.NaviTE{
				Given:        "universal(undefined scope): custom filter",
				Relative:     lab.Static.RetroWave,
				Subscription: enums.SubscribeUniversal,
				ExpectedNoOf: lab.Quantities{
					Files:       8,
					Directories: 0,
				},
			},
			Description: "items with '.flac' suffix",
			Pattern:     "*.flac",
		}),

		// === error =========================================================

		Entry(nil, &lab.FilterTE{
			NaviTE: lab.NaviTE{
				Given:        "universal(any scope): missing filter type",
				Relative:     lab.Static.RetroWave,
				Subscription: enums.SubscribeUniversal,
				ExpectedNoOf: lab.Quantities{
					Files:       8,
					Directories: 0,
				},
			},
			Description: "items with '.flac' suffix",
			Pattern:     "*.flac",
			Scope:       enums.ScopeFile,
		}),
	)
})
