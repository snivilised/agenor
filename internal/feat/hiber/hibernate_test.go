package hiber_test

import (
	"regexp/syntax"

	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok

	"github.com/snivilised/li18ngo"
	"github.com/snivilised/nefilim/test/luna"

	age "github.com/snivilised/agenor"
	"github.com/snivilised/agenor/core"
	"github.com/snivilised/agenor/enums"
	lab "github.com/snivilised/agenor/internal/laboratory"
	"github.com/snivilised/agenor/internal/services"
	"github.com/snivilised/agenor/internal/third/lo"
	"github.com/snivilised/agenor/pref"
	"github.com/snivilised/agenor/test/hanno"
	"github.com/snivilised/agenor/tfs"
)

var _ = Describe("feature", Ordered, func() {
	var (
		fS *luna.MemFS
	)

	BeforeAll(func() {
		const (
			verbose = false
		)

		var (
			err error
		)

		repo := hanno.Repo("")
		index := hanno.Combine(repo, "test/data/musico-index.xml")
		fS, err = hanno.CustomTree(index, "MUSICO", verbose, lab.Static.RetroWave, "edm")

		Expect(err).To(Succeed(), "Failed to initialise custom tree with MUSICO data")
		Expect(li18ngo.Use()).To(Succeed())
	})

	BeforeEach(func() {
		services.Reset()
	})

	Context("comprehension", func() {
		When("directories: wake and sleep", func() {
			It("🧪 should: invoke inside hibernation range", Label("example"),
				func(ctx SpecContext) {
					path := lab.Static.RetroWave
					result, _ := age.Walk().Configure().Extent(age.Prime(
						&pref.Using{
							Subscription: enums.SubscribeDirectories,
							Head: pref.Head{
								Handler: func(servant age.Servant) error {
									node := servant.Node()
									GinkgoWriter.Printf(
										"---> 🍯 EXAMPLE-HIBERNATE-CALLBACK: '%v'\n", node.Path,
									)
									return nil
								},
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

						age.WithOnWake(func(description string) {
							GinkgoWriter.Printf("===> 🔊 Wake: '%v'\n", description)
						}),

						age.WithOnSleep(func(description string) {
							GinkgoWriter.Printf("===> 🔇 Sleep: '%v'\n", description)
						}),

						age.WithHibernationOptions(
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
						age.WithHibernationBehaviourExclusiveWake(),

						// This is only required to change the default inclusivity
						// of the sleep condition; by default is exclusive.
						age.WithHibernationBehaviourInclusiveSleep(),
					)).Navigate(ctx)

					GinkgoWriter.Printf("===> 🍭 invoked '%v' directories\n",
						result.Metrics().Count(enums.MetricNoDirectoriesInvoked),
					)
				},
			)
		})
	})

	DescribeTable("simple hibernate",
		func(ctx SpecContext, entry *lab.HibernateTE) {
			recall := make(lab.Recall)
			once := func(node *age.Node) error { //nolint:unparam // return nil error ok
				_, found := recall[node.Extension.Name]
				Expect(found).To(BeFalse())
				recall[node.Extension.Name] = len(node.Children)

				return nil
			}

			path := lo.Ternary(entry.NaviTE.Relative == "",
				lab.Static.RetroWave,
				entry.NaviTE.Relative,
			)

			client := func(servant age.Servant) error {
				node := servant.Node()
				GinkgoWriter.Printf(
					"---> 🌊 HIBERNATE-CALLBACK: '%v'\n", node.Path,
				)

				return once(node)
			}

			result, err := age.Walk().Configure().Extent(age.Prime(
				&pref.Using{
					Subscription: entry.Subscription,
					Head: pref.Head{
						Handler: client,
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

				age.WithOnWake(func(description string) {
					GinkgoWriter.Printf("===> 🔊 Wake: '%v'\n", description)
				}),

				age.WithOnSleep(func(description string) {
					GinkgoWriter.Printf("===> 🔇 Sleep: '%v'\n", description)
				}),

				age.WithHibernationOptions(
					&core.HibernateOptions{
						WakeAt:    entry.Hibernate.WakeAt,
						SleepAt:   entry.Hibernate.SleepAt,
						Behaviour: entry.Hibernate.Behaviour,
					},
				),

				age.IfOption(entry.CaseSensitive, age.WithHookCaseSensitiveSort()),
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
		lab.FormatHibernateTestDescription,

		// === directories ===================================================

		Entry(nil, &lab.HibernateTE{
			DescribedTE: lab.DescribedTE{
				Given:  "wake and sleep (directories, inclusive:default)",
				Should: "fail",
			},
			NaviTE: lab.NaviTE{
				Relative:     lab.Static.RetroWave,
				Subscription: enums.SubscribeDirectories,
				Mandatory: []string{"Night Drive", "College",
					"Northern Council", "Teenage Color",
				},
				Prohibited: []string{lab.Static.RetroWave, "Chromatics",
					"Electric Youth", "Innerworld",
				},
				ExpectedNoOf: lab.Quantities{
					Directories: 4,
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

		Entry(nil, &lab.HibernateTE{
			DescribedTE: lab.DescribedTE{
				Given:  "wake and sleep (directories, excl:wake, inc:sleep, mute)",
				Should: "fail",
			},

			NaviTE: lab.NaviTE{
				Relative:     lab.Static.RetroWave,
				Subscription: enums.SubscribeDirectories,
				Mandatory: []string{"College", "Northern Council",
					"Teenage Color", "Electric Youth",
				},
				Prohibited: []string{"Night Drive", lab.Static.RetroWave,
					"Chromatics", "Innerworld",
				},
				ExpectedNoOf: lab.Quantities{
					Directories: 4,
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

		Entry(nil, &lab.HibernateTE{
			DescribedTE: lab.DescribedTE{
				Given: "wake only (directories, inclusive:default)",
			},
			NaviTE: lab.NaviTE{
				Relative:     lab.Static.RetroWave,
				Subscription: enums.SubscribeDirectories,
				Mandatory: []string{"Night Drive", "College", "Northern Council",
					"Teenage Color", "Electric Youth", "Innerworld",
				},
				Prohibited: []string{lab.Static.RetroWave, "Chromatics"},
				ExpectedNoOf: lab.Quantities{
					Directories: 6,
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

		Entry(nil, &lab.HibernateTE{
			DescribedTE: lab.DescribedTE{
				Given: "sleep only (directories, inclusive:default)",
			},
			NaviTE: lab.NaviTE{
				Relative:     lab.Static.RetroWave,
				Subscription: enums.SubscribeDirectories,
				Mandatory: []string{lab.Static.RetroWave, "Chromatics", "Night Drive", "College",
					"Northern Council", "Teenage Color",
				},
				Prohibited: []string{"Electric Youth", "Innerworld"},
				ExpectedNoOf: lab.Quantities{
					Directories: 6,
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

		Entry(nil, &lab.HibernateTE{
			DescribedTE: lab.DescribedTE{
				Given: "sleep only (directories, inclusive:default)",
			},
			NaviTE: lab.NaviTE{

				Relative:     lab.Static.RetroWave,
				Subscription: enums.SubscribeDirectories,
				Mandatory:    []string{lab.Static.RetroWave, "Chromatics"},
				Prohibited: []string{"Night Drive", "College", "Northern Council",
					"Teenage Color", "Electric Youth", "Innerworld",
				},
				ExpectedNoOf: lab.Quantities{
					Directories: 2,
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

		// error ==================================================================

		Entry(nil, &lab.HibernateTE{
			DescribedTE: lab.DescribedTE{
				Given:  "wake only (directories, inclusive:default)",
				Should: "fail",
			},
			NaviTE: lab.NaviTE{
				Relative:     lab.Static.RetroWave,
				Subscription: enums.SubscribeDirectories,
				Mandatory: []string{"Night Drive", "College", "Northern Council",
					"Teenage Color", "Electric Youth", "Innerworld",
				},
				Prohibited: []string{lab.Static.RetroWave, "Chromatics"},
				ExpectedErr: &syntax.Error{
					Code: "missing closing )",
					Expr: "(",
				},
			},
			Hibernate: &core.HibernateOptions{
				WakeAt: &core.FilterDef{
					Type:        enums.FilterTypeRegex,
					Description: "Wake At: Night Drive",
					Pattern:     "(",
				},
			},
		}),

		Entry(nil, &lab.HibernateTE{
			DescribedTE: lab.DescribedTE{
				Given:  "sleep only (directories, inclusive:default)",
				Should: "fail",
			},
			NaviTE: lab.NaviTE{
				Relative:     lab.Static.RetroWave,
				Subscription: enums.SubscribeDirectories,
				Mandatory: []string{lab.Static.RetroWave, "Chromatics", "Night Drive", "College",
					"Northern Council", "Teenage Color",
				},
				Prohibited: []string{"Electric Youth", "Innerworld"},
				ExpectedErr: &syntax.Error{
					Code: "missing closing )",
					Expr: "(",
				},
			},

			Hibernate: &core.HibernateOptions{
				SleepAt: &core.FilterDef{
					Type:        enums.FilterTypeRegex,
					Description: "Sleep At: Electric Youth",
					Pattern:     "(",
				},
			},
		}),
	)
})
