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
	"github.com/snivilised/traverse/locale"
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

	// TODO: need to return error from multiple places
	XDescribeTable("filtering errata", Label("BROKEN"),
		func(ctx SpecContext, entry *lab.FilterErrataTE) {
			path := entry.Relative
			callback := func(servant tv.Servant) error {
				_ = servant.Node()

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

				tv.WithFilter(entry.Filter),
			)).Navigate(ctx)

			lab.AssertNavigation(&entry.NaviTE, &lab.TestOptions{
				FS:          fS,
				Path:        path,
				Result:      result,
				Err:         err,
				ExpectedErr: entry.ExpectedErr,
			})
		},
		func(entry *lab.FilterErrataTE) string {
			return fmt.Sprintf("🧪 ===> given: '%v', should: '%v'",
				entry.Given, entry.Should,
			)
		},

		Entry(nil, &lab.FilterErrataTE{
			NaviTE: lab.NaviTE{
				Given:        "missing type",
				Should:       "fail",
				Relative:     lab.Static.RetroWave,
				Subscription: enums.SubscribeFiles,
				ExpectedErr:  locale.ErrFilterMissingType,
			},
			Filter: &pref.FilterOptions{
				Node: &core.FilterDef{
					Description: "filter missing type",
					Pattern:     "*",
					Scope:       enums.ScopeAll,
				},
			},
		}),
	)
})
