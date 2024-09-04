package hiber_test

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"testing/fstest"

	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok

	"github.com/snivilised/li18ngo"

	tv "github.com/snivilised/traverse"
	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/enums"
	"github.com/snivilised/traverse/internal/helpers"
	"github.com/snivilised/traverse/internal/services"
	"github.com/snivilised/traverse/internal/third/lo"
)

var _ = Describe("feature", Ordered, func() {
	var (
		FS   fstest.MapFS
		root string
	)

	BeforeAll(func() {
		const (
			verbose = false
		)

		FS, root = helpers.Musico(verbose,
			filepath.Join("MUSICO", "RETRO-WAVE"),
			filepath.Join("MUSICO", "edm"),
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
					path := helpers.Path(root, "RETRO-WAVE")
					result, _ := tv.Walk().Configure().Extent(tv.Prime(
						&tv.Using{
							Root:         path,
							Subscription: enums.SubscribeFolders,
							Handler: func(node *core.Node) error {
								GinkgoWriter.Printf(
									"---> 🍯 EXAMPLE-HIBERNATE-CALLBACK: '%v'\n", node.Path,
								)
								return nil
							},
							GetReadDirFS: func() fs.ReadDirFS {
								return FS
							},
							GetQueryStatusFS: func(_ fs.FS) fs.StatFS {
								return FS
							},
						},

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

						tv.WithHookQueryStatus(
							func(qsys fs.StatFS, path string) (fs.FileInfo, error) {
								return qsys.Stat(helpers.TrimRoot(path))
							},
						),

						tv.WithHookReadDirectory(
							func(rsys fs.ReadDirFS, dirname string) ([]fs.DirEntry, error) {
								// This is only required because fstest.MapFS strangely
								// can't resolve paths with a leading /. Any other program
								// using a different file system should not need to use
								// this hook for this purpose.
								//
								return rsys.ReadDir(helpers.TrimRoot(dirname))
							},
						),
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
			recording := make(helpers.RecordingMap)
			once := func(node *tv.Node) error { //nolint:unparam // return nil error ok
				_, found := recording[node.Extension.Name]
				Expect(found).To(BeFalse())
				recording[node.Extension.Name] = len(node.Children)

				return nil
			}

			path := helpers.Path(
				root,
				lo.Ternary(entry.NaviTE.Relative == "",
					"RETRO-WAVE",
					entry.NaviTE.Relative,
				),
			)

			client := func(node *tv.Node) error {
				GinkgoWriter.Printf(
					"---> 🌊 HIBERNATE-CALLBACK: '%v'\n", node.Path,
				)

				return once(node)
			}

			result, err := tv.Walk().Configure().Extent(tv.Prime(
				&tv.Using{
					Root:         path,
					Subscription: entry.Subscription,
					Handler:      client,
					GetReadDirFS: func() fs.ReadDirFS {
						return FS
					},
					GetQueryStatusFS: func(_ fs.FS) fs.StatFS {
						return FS
					},
				},

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
				tv.WithHookQueryStatus(
					func(qsys fs.StatFS, path string) (fs.FileInfo, error) {
						return qsys.Stat(helpers.TrimRoot(path))
					},
				),

				tv.WithHookReadDirectory(
					func(rsys fs.ReadDirFS, dirname string) ([]fs.DirEntry, error) {
						return rsys.ReadDir(helpers.TrimRoot(dirname))
					},
				),
			)).Navigate(ctx)

			helpers.AssertNavigation(&entry.NaviTE, &helpers.TestOptions{
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
			NaviTE: helpers.NaviTE{
				Given:        "wake and sleep (folders, inclusive:default)",
				Relative:     "RETRO-WAVE",
				Subscription: enums.SubscribeFolders,
				Mandatory: []string{"Night Drive", "College",
					"Northern Council", "Teenage Color",
				},
				Prohibited: []string{"RETRO-WAVE", "Chromatics",
					"Electric Youth", "Innerworld",
				},
				ExpectedNoOf: helpers.Quantities{
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
			NaviTE: helpers.NaviTE{
				Given:        "wake and sleep (folders, excl:wake, inc:sleep, mute)",
				Relative:     "RETRO-WAVE",
				Subscription: enums.SubscribeFolders,
				Mandatory: []string{"College", "Northern Council",
					"Teenage Color", "Electric Youth",
				},
				Prohibited: []string{"Night Drive", "RETRO-WAVE",
					"Chromatics", "Innerworld",
				},
				ExpectedNoOf: helpers.Quantities{
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
			NaviTE: helpers.NaviTE{
				Given:        "wake only (folders, inclusive:default)",
				Relative:     "RETRO-WAVE",
				Subscription: enums.SubscribeFolders,
				Mandatory: []string{"Night Drive", "College", "Northern Council",
					"Teenage Color", "Electric Youth", "Innerworld",
				},
				Prohibited: []string{"RETRO-WAVE", "Chromatics"},
				ExpectedNoOf: helpers.Quantities{
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
			NaviTE: helpers.NaviTE{
				Given:        "sleep only (folders, inclusive:default)",
				Relative:     "RETRO-WAVE",
				Subscription: enums.SubscribeFolders,
				Mandatory: []string{"RETRO-WAVE", "Chromatics", "Night Drive", "College",
					"Northern Council", "Teenage Color",
				},
				Prohibited: []string{"Electric Youth", "Innerworld"},
				ExpectedNoOf: helpers.Quantities{
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
			NaviTE: helpers.NaviTE{
				Given:        "sleep only (folders, inclusive:default)",
				Relative:     "RETRO-WAVE",
				Subscription: enums.SubscribeFolders,
				Mandatory:    []string{"RETRO-WAVE", "Chromatics"},
				Prohibited: []string{"Night Drive", "College", "Northern Council",
					"Teenage Color", "Electric Youth", "Innerworld",
				},
				ExpectedNoOf: helpers.Quantities{
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
