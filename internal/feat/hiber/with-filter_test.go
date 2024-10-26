package hiber_test

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok

	"github.com/snivilised/li18ngo"
	"github.com/snivilised/nefilim/luna"
	tv "github.com/snivilised/traverse"

	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/enums"
	lab "github.com/snivilised/traverse/internal/laboratory"
	"github.com/snivilised/traverse/internal/services"
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

	DescribeTable("filter and listen both active",
		func(ctx SpecContext, entry *hibernateTE) {
			path := lab.Static.RetroWave
			result, err := tv.Walk().Configure().Extent(tv.Prime(
				&tv.Using{
					Tree:         path,
					Subscription: entry.Subscription,
					Handler:      entry.Callback,
					GetTraverseFS: func(_ string) tv.TraverseFS {
						return fS
					},
				},

				tv.WithOnBegin(lab.Begin("🛡️")),
				tv.WithOnEnd(lab.End("🏁")),

				tv.WithFilter(&pref.FilterOptions{
					Node: &core.FilterDef{
						Type:        enums.FilterTypeGlob,
						Description: "items with '.flac' suffix",
						Pattern:     "*.flac",
						Scope:       enums.ScopeFile,
					},
				}),

				tv.IfOptionF(entry.Hibernate != nil && entry.Hibernate.WakeAt != nil,
					func() pref.Option {
						return tv.WithHibernationFilterWake(
							&core.FilterDef{
								Type:        entry.Hibernate.WakeAt.Type,
								Description: entry.Hibernate.WakeAt.Description,
								Pattern:     entry.Hibernate.WakeAt.Pattern,
							},
						)
					},
				),

				tv.IfOptionF(entry.Hibernate != nil && entry.Hibernate.SleepAt != nil,
					func() pref.Option {
						return tv.WithHibernationFilterSleep(
							&core.FilterDef{
								Type:        entry.Hibernate.SleepAt.Type,
								Description: entry.Hibernate.SleepAt.Description,
								Pattern:     entry.Hibernate.SleepAt.Pattern,
							},
						)
					},
				),

				tv.WithOnWake(func(description string) {
					GinkgoWriter.Printf("===> 🔆 Waking: '%v'\n", description)
				}),
				tv.WithOnSleep(func(description string) {
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
			folders := result.Metrics().Count(enums.MetricNoFoldersInvoked)

			GinkgoWriter.Printf("---> 🍕🍕 Metrics, files:'%v', folders:'%v'\n",
				files, folders,
			)
		},
		func(entry *hibernateTE) string {
			return fmt.Sprintf("🧪 ===> given: '%v', should: '%v'", entry.Given, entry.Should)
		},

		Entry(nil, &hibernateTE{
			NaviTE: lab.NaviTE{
				Given:        "File Subscription",
				Should:       "wake, then apply filter until the end",
				Subscription: enums.SubscribeFiles,
				Callback: func(servant tv.Servant) error {
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

		Entry(nil, &hibernateTE{
			NaviTE: lab.NaviTE{
				Given:        "File Subscription",
				Should:       "apply filter until sleep",
				Subscription: enums.SubscribeFiles,
				Callback: func(servant tv.Servant) error {
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

		Entry(nil, &hibernateTE{
			NaviTE: lab.NaviTE{
				Given:        "File Subscription",
				Should:       "apply filter within hibernation range",
				Subscription: enums.SubscribeFiles,
				Callback: func(servant tv.Servant) error {
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
					Scope:       enums.ScopeFolder,
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
