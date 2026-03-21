package filtering_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/snivilised/jaywalk/src/agenor"
	"github.com/snivilised/jaywalk/src/agenor/core"
	"github.com/snivilised/jaywalk/src/agenor/enums"
	lab "github.com/snivilised/jaywalk/src/agenor/internal/laboratory"
	"github.com/snivilised/jaywalk/src/internal/services"
	"github.com/snivilised/jaywalk/locale"
	"github.com/snivilised/jaywalk/src/agenor/pref"
	"github.com/snivilised/jaywalk/src/agenor/test/hanno"
	"github.com/snivilised/jaywalk/src/agenor/tfs"
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

	DescribeTable("filtering errata",
		func(ctx SpecContext, entry *lab.FilterErrataTE) {
			path := entry.Relative
			callback := func(servant agenor.Servant) error {
				_ = servant.Node()

				return nil
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

				agenor.WithFilter(entry.Filter),
			)).Navigate(ctx)

			lab.AssertNavigation(&entry.NaviTE, &lab.TestOptions{
				FS:          fS,
				Path:        path,
				Result:      result,
				Err:         err,
				ExpectedErr: entry.ExpectedErr,
			})
		},
		lab.FormatFilterErrataTestDescription,
		Entry(nil, &lab.FilterErrataTE{
			DescribedTE: lab.DescribedTE{
				Given:  "missing type",
				Should: "fail",
			},
			NaviTE: lab.NaviTE{
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
