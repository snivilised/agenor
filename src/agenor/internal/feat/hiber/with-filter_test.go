package hiber_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/snivilised/li18ngo"
	"github.com/snivilised/nefilim/test/luna"

	"github.com/snivilised/jaywalk/src/agenor"
	"github.com/snivilised/jaywalk/src/agenor/core"
	"github.com/snivilised/jaywalk/src/agenor/enums"
	lab "github.com/snivilised/jaywalk/src/agenor/internal/laboratory"
	"github.com/snivilised/jaywalk/src/internal/services"
	"github.com/snivilised/jaywalk/src/agenor/pref"
	"github.com/snivilised/jaywalk/src/agenor/test/hanno"
	"github.com/snivilised/jaywalk/src/agenor/tfs"
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

	DescribeTable("filter and listen both active",
		func(ctx SpecContext, entry *lab.HibernateTE) {
			path := lab.Static.RetroWave
			result, err := agenor.Walk().Configure().Extent(agenor.Prime(
				&pref.Using{
					Subscription: entry.Subscription,
					Head: pref.Head{
						Handler: entry.Callback,
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

				agenor.WithFilter(&pref.FilterOptions{
					Node: &core.FilterDef{
						Type:        enums.FilterTypeGlob,
						Description: "items with '.flac' suffix",
						Pattern:     "*.flac",
						Scope:       enums.ScopeFile,
					},
				}),

				agenor.IfOptionF(entry.Hibernate != nil && entry.Hibernate.WakeAt != nil,
					func() pref.Option {
						return agenor.WithHibernationFilterWake(
							&core.FilterDef{
								Type:        entry.Hibernate.WakeAt.Type,
								Description: entry.Hibernate.WakeAt.Description,
								Pattern:     entry.Hibernate.WakeAt.Pattern,
							},
						)
					},
				),

				agenor.IfOptionF(entry.Hibernate != nil && entry.Hibernate.SleepAt != nil,
					func() pref.Option {
						return agenor.WithHibernationFilterSleep(
							&core.FilterDef{
								Type:        entry.Hibernate.SleepAt.Type,
								Description: entry.Hibernate.SleepAt.Description,
								Pattern:     entry.Hibernate.SleepAt.Pattern,
							},
						)
					},
				),

				agenor.WithOnWake(func(description string) {
					GinkgoWriter.Printf("===> 🔆 Waking: '%v'\n", description)
				}),
				agenor.WithOnSleep(func(description string) {
					GinkgoWriter.Printf("===> 🌙 Sleeping: '%v'\n", description)
				}),
			)).Navigate(ctx)

			lab.AssertNavigation(&entry.NaviTE, &lab.TestOptions{
				FS:     fS,
				Path:   path,
				Result: result,
				Err:    err,
			})

			files := result.Metrics().Count(enums.MetricNoFilesInvoked)
			directories := result.Metrics().Count(enums.MetricNoDirectoriesInvoked)

			GinkgoWriter.Printf("---> 🍕🍕 Metrics, files:'%v', directories:'%v'\n",
				files, directories,
			)
		},
		lab.FormatHibernateTestDescription,

		Entry(nil, &lab.HibernateTE{
			DescribedTE: lab.DescribedTE{
				Given:  "File Subscription",
				Should: "wake, then apply filter until the end",
			},
			NaviTE: lab.NaviTE{
				Subscription: enums.SubscribeFiles,
				Callback: func(servant agenor.Servant) error {
					node := servant.Node()
					GinkgoWriter.Printf("---> WAKE-HIBERNATE-AND-FILTER-😵‍💫: '%v'\n", node.Path)

					return nil
				},
				ExpectedNoOf: lab.Quantities{
					Files: 6,
				},
			},
			Hibernate: &core.HibernateOptions{
				WakeAt: &core.FilterDef{
					Type:        enums.FilterTypeGlob,
					Description: "Wake At: A1 - Incident.flac",
					Pattern:     "A1 - Incident.flac",
					Scope:       enums.ScopeFile,
				},
			},
		}),
		Entry(nil, &lab.HibernateTE{
			DescribedTE: lab.DescribedTE{
				Given:  "File Subscription",
				Should: "apply filter until sleep",
			},
			NaviTE: lab.NaviTE{
				Subscription: enums.SubscribeFiles,
				Callback: func(servant agenor.Servant) error {
					node := servant.Node()
					GinkgoWriter.Printf("---> SLEEP-HIBERNATE-AND-FILTER-😴: '%v'\n", node.Path)

					return nil
				},
				ExpectedNoOf: lab.Quantities{
					Files: 2,
				},
			},
			Hibernate: &core.HibernateOptions{
				SleepAt: &core.FilterDef{
					Type:        enums.FilterTypeGlob,
					Description: "Sleep At: A1 - Incident.flac",
					Pattern:     "A1 - Incident.flac",
					Scope:       enums.ScopeFile,
				},
			},
		}),
		Entry(nil, &lab.HibernateTE{
			DescribedTE: lab.DescribedTE{
				Given:  "File Subscription",
				Should: "apply filter within hibernation range",
			},
			NaviTE: lab.NaviTE{
				Subscription: enums.SubscribeFiles,
				Callback: func(servant agenor.Servant) error {
					node := servant.Node()
					GinkgoWriter.Printf("---> WAKE/SLEEP-HIBERNATE-AND-FILTER-😴: '%v'\n", node.Path)

					return nil
				},
				ExpectedNoOf: lab.Quantities{
					Files: 4,
				},
			},
			Hibernate: &core.HibernateOptions{
				WakeAt: &core.FilterDef{
					Type:        enums.FilterTypeGlob,
					Description: "Sleep At: A1 - Incident.flac",
					Pattern:     "A1 - Incident.flac",
					Scope:       enums.ScopeDirectory,
				},
				SleepAt: &core.FilterDef{
					Type:        enums.FilterTypeGlob,
					Description: "Sleep At: A1 - Before Life.flac",
					Pattern:     "A1 - Before Life.flac",
					Scope:       enums.ScopeFile,
				},
			},
		}),
	)
})
