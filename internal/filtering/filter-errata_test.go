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
	"github.com/snivilised/agenor/locale"
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

	// TODO: need to return error from multiple places
	DescribeTable("filtering errata", Label("BROKEN"),
		func(ctx SpecContext, entry *lab.FilterErrataTE) {
			path := entry.Relative
			callback := func(servant age.Servant) error {
				_ = servant.Node()

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

				age.WithFilter(entry.Filter),
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

		XEntry(nil, &lab.FilterErrataTE{
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
