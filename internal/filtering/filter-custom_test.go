package filtering_test

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok

	age "github.com/snivilised/agenor"
	"github.com/snivilised/agenor/core"
	"github.com/snivilised/agenor/enums"
	lab "github.com/snivilised/agenor/internal/laboratory"
	"github.com/snivilised/agenor/internal/services"
	"github.com/snivilised/agenor/internal/third/lo"
	"github.com/snivilised/agenor/pref"
	"github.com/snivilised/agenor/test/hanno"
	"github.com/snivilised/agenor/tfs"
	"github.com/snivilised/li18ngo"
	"github.com/snivilised/nefilim/test/luna"
)

var _ = Describe("NavigatorFilterCustom", Ordered, func() {
	var (
		fS *luna.MemFS
	)

	BeforeAll(func() {
		const (
			verbose = false
		)

		fS = hanno.Nuxx(verbose, lab.Static.RetroWave)
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
			callback := func(servant age.Servant) error {
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
			result, err := age.Walk().Configure().Extent(age.Prime(
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
				age.WithOnBegin(lab.Begin("🛡️")),
				age.WithOnEnd(lab.End("🏁")),

				age.WithFilter(&pref.FilterOptions{
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
