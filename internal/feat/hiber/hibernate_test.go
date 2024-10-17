package hiber_test

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok

	"github.com/snivilised/li18ngo"

	tv "github.com/snivilised/traverse"
	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/enums"
	lab "github.com/snivilised/traverse/internal/laboratory"
	"github.com/snivilised/traverse/internal/services"
	"github.com/snivilised/traverse/internal/third/lo"
)

var _ = Describe("feature", Ordered, func() {
	var (
		FS   *lab.TestTraverseFS
		root string
	)

	BeforeAll(func() {
		const (
			verbose = false
		)

		FS, root = lab.Musico(verbose,
			lab.Static.RetroWave, "edm",
		)
		Expect(root).NotTo(BeEmpty())
		Expect(li18ngo.Use()).To(Succeed())
	})

	BeforeEach(func() {
		services.Reset()
	})

	Context("comprehension", func() {
		When("folders: wake and sleep", func() {
			It("🧪 should: invoke inside hibernation range", Label("example"),
				func(ctx SpecContext) {
					path := lab.Static.RetroWave
					result, _ := tv.Walk().Configure().Extent(tv.Prime(
						&tv.Using{
							Tree:         path,
							Subscription: enums.SubscribeFolders,
							Handler: func(servant tv.Servant) error {
								node := servant.Node()
								GinkgoWriter.Printf(
									"---> 🍯 EXAMPLE-HIBERNATE-CALLBACK: '%v'\n", node.Path,
								)
								return nil
							},
							GetTraverseFS: func(_ string) tv.TraverseFS {
								return FS
							},
						},
						tv.WithOnBegin(lab.Begin("🛡️")),
						tv.WithOnEnd(lab.End("🏁")),

						tv.WithOnWake(func(description string) {
							GinkgoWriter.Printf("===> 🔊 Wake: '%v'\n", description)
						}),

						tv.WithOnSleep(func(description string) {
							GinkgoWriter.Printf("===> 🔇 Sleep: '%v'\n", description)
						}),

						tv.WithHibernationOptions(
							&core.HibernateOptions{
								WakeAt: &core.FilterDef{
									Type:        enums.FilterTypeGlob,
									Description: "Wake At: Night Drive",
									Pattern:     "Night Drive",
								},
								SleepAt: &core.FilterDef{
									Type:        enums.FilterTypeGlob,
									Description: "Sleep At: Electric Youth",
									Pattern:     "Electric Youth",
								},
							},
						),

						// This is only required to change the default inclusivity
						// of the wake condition; by default is inclusive.
						tv.WithHibernationBehaviourExclusiveWake(),

						// This is only required to change the default inclusivity
						// of the sleep condition; by default is exclusive.
						tv.WithHibernationBehaviourInclusiveSleep(),
					)).Navigate(ctx)

					GinkgoWriter.Printf("===> 🍭 invoked '%v' folders\n",
						result.Metrics().Count(enums.MetricNoFoldersInvoked),
					)
				},
			)
		})
	})

	DescribeTable("simple hibernate",
		func(ctx SpecContext, entry *hibernateTE) {
			recording := make(lab.RecordingMap)
			once := func(node *tv.Node) error { //nolint:unparam // return nil error ok
				_, found := recording[node.Extension.Name]
				Expect(found).To(BeFalse())
				recording[node.Extension.Name] = len(node.Children)

				return nil
			}

			path := lo.Ternary(entry.NaviTE.Relative == "",
				lab.Static.RetroWave,
				entry.NaviTE.Relative,
			)

			client := func(servant tv.Servant) error {
				node := servant.Node()
				GinkgoWriter.Printf(
					"---> 🌊 HIBERNATE-CALLBACK: '%v'\n", node.Path,
				)

				return once(node)
			}

			result, err := tv.Walk().Configure().Extent(tv.Prime(
				&tv.Using{
					Tree:         path,
					Subscription: entry.Subscription,
					Handler:      client,
					GetTraverseFS: func(_ string) tv.TraverseFS {
						return FS
					},
				},
				tv.WithOnBegin(lab.Begin("🛡️")),
				tv.WithOnEnd(lab.End("🏁")),

				tv.WithOnWake(func(description string) {
					GinkgoWriter.Printf("===> 🔊 Wake: '%v'\n", description)
				}),

				tv.WithOnSleep(func(description string) {
					GinkgoWriter.Printf("===> 🔇 Sleep: '%v'\n", description)
				}),

				tv.WithHibernationOptions(
					&core.HibernateOptions{
						WakeAt:    entry.Hibernate.WakeAt,
						SleepAt:   entry.Hibernate.SleepAt,
						Behaviour: entry.Hibernate.Behaviour,
					},
				),

				tv.IfOption(entry.CaseSensitive, tv.WithHookCaseSensitiveSort()),
			)).Navigate(ctx)

			lab.AssertNavigation(&entry.NaviTE, &lab.TestOptions{
				FS:          FS,
				Recording:   recording,
				Path:        path,
				Result:      result,
				Err:         err,
				ExpectedErr: entry.ExpectedErr,
			})
		},

		func(entry *hibernateTE) string {
			return fmt.Sprintf("🧪 ===> given: '%v', should: '%v'", entry.Given, entry.Should)
		},

		// === folders =======================================================

		Entry(nil, &hibernateTE{
			NaviTE: lab.NaviTE{
				Given:        "wake and sleep (folders, inclusive:default)",
				Relative:     lab.Static.RetroWave,
				Subscription: enums.SubscribeFolders,
				Mandatory: []string{"Night Drive", "College",
					"Northern Council", "Teenage Color",
				},
				Prohibited: []string{lab.Static.RetroWave, "Chromatics",
					"Electric Youth", "Innerworld",
				},
				ExpectedNoOf: lab.Quantities{
					Folders: 4,
				},
			},
			Hibernate: &core.HibernateOptions{
				WakeAt: &core.FilterDef{
					Type:        enums.FilterTypeGlob,
					Description: "Wake At: Night Drive",
					Pattern:     "Night Drive",
				},
				SleepAt: &core.FilterDef{
					Type:        enums.FilterTypeGlob,
					Description: "Sleep At: Electric Youth",
					Pattern:     "Electric Youth",
				},
				Behaviour: core.HibernationBehaviour{
					InclusiveWake:  true,
					InclusiveSleep: false,
				},
			},
		}),

		Entry(nil, &hibernateTE{
			NaviTE: lab.NaviTE{
				Given:        "wake and sleep (folders, excl:wake, inc:sleep, mute)",
				Relative:     lab.Static.RetroWave,
				Subscription: enums.SubscribeFolders,
				Mandatory: []string{"College", "Northern Council",
					"Teenage Color", "Electric Youth",
				},
				Prohibited: []string{"Night Drive", lab.Static.RetroWave,
					"Chromatics", "Innerworld",
				},
				ExpectedNoOf: lab.Quantities{
					Folders: 4,
				},
			},
			Hibernate: &core.HibernateOptions{
				WakeAt: &core.FilterDef{
					Type:        enums.FilterTypeRegex,
					Description: "Wake At: Night Drive",
					Pattern:     "Night Drive",
				},
				SleepAt: &core.FilterDef{
					Type:        enums.FilterTypeGlob,
					Description: "Sleep At: Electric Youth",
					Pattern:     "Electric Youth",
				},
				Behaviour: core.HibernationBehaviour{
					InclusiveWake:  false,
					InclusiveSleep: true,
				},
			},
			Mute: true,
		}),

		Entry(nil, &hibernateTE{
			NaviTE: lab.NaviTE{
				Given:        "wake only (folders, inclusive:default)",
				Relative:     lab.Static.RetroWave,
				Subscription: enums.SubscribeFolders,
				Mandatory: []string{"Night Drive", "College", "Northern Council",
					"Teenage Color", "Electric Youth", "Innerworld",
				},
				Prohibited: []string{lab.Static.RetroWave, "Chromatics"},
				ExpectedNoOf: lab.Quantities{
					Folders: 6,
				},
			},
			Hibernate: &core.HibernateOptions{
				WakeAt: &core.FilterDef{
					Type:        enums.FilterTypeRegex,
					Description: "Wake At: Night Drive",
					Pattern:     "Night Drive",
				},
				Behaviour: core.HibernationBehaviour{
					InclusiveWake:  true,
					InclusiveSleep: false,
				},
			},
		}),

		Entry(nil, &hibernateTE{
			NaviTE: lab.NaviTE{
				Given:        "sleep only (folders, inclusive:default)",
				Relative:     lab.Static.RetroWave,
				Subscription: enums.SubscribeFolders,
				Mandatory: []string{lab.Static.RetroWave, "Chromatics", "Night Drive", "College",
					"Northern Council", "Teenage Color",
				},
				Prohibited: []string{"Electric Youth", "Innerworld"},
				ExpectedNoOf: lab.Quantities{
					Folders: 6,
				},
			},

			Hibernate: &core.HibernateOptions{
				SleepAt: &core.FilterDef{
					Type:        enums.FilterTypeGlob,
					Description: "Sleep At: Electric Youth",
					Pattern:     "Electric Youth",
				},
				Behaviour: core.HibernationBehaviour{
					InclusiveWake:  true,
					InclusiveSleep: false,
				},
			},
		}),

		Entry(nil, &hibernateTE{
			NaviTE: lab.NaviTE{
				Given:        "sleep only (folders, inclusive:default)",
				Relative:     lab.Static.RetroWave,
				Subscription: enums.SubscribeFolders,
				Mandatory:    []string{lab.Static.RetroWave, "Chromatics"},
				Prohibited: []string{"Night Drive", "College", "Northern Council",
					"Teenage Color", "Electric Youth", "Innerworld",
				},
				ExpectedNoOf: lab.Quantities{
					Folders: 2,
				},
			},
			Hibernate: &core.HibernateOptions{
				SleepAt: &core.FilterDef{
					Type:        enums.FilterTypeGlob,
					Description: "Sleep At: Night Drive",
					Pattern:     "Night Drive",
				},
				Behaviour: core.HibernationBehaviour{
					InclusiveWake:  true,
					InclusiveSleep: false,
				},
			},
		}),
	)
})
