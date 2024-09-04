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
	"github.com/snivilised/traverse/pref"
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

	DescribeTable("filter and listen both active",
		func(ctx SpecContext, entry *hibernateTE) {
			path := helpers.Path(root, "RETRO-WAVE")
			result, err := tv.Walk().Configure().Extent(tv.Prime(
				&tv.Using{
					Root:         path,
					Subscription: entry.Subscription,
					Handler:      entry.Callback,
					GetReadDirFS: func() fs.ReadDirFS {
						return FS
					},
					GetQueryStatusFS: func(_ fs.FS) fs.StatFS {
						return FS
					},
				},

				tv.WithOnBegin(helpers.Begin("🛡️")),
				tv.WithOnEnd(helpers.End("🏁")),

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
				tv.WithOnStart(func(description string) {
					GinkgoWriter.Printf("===> 🔆 Waking: '%v'\n", description)
				}),
				tv.WithOnStop(func(description string) {
					GinkgoWriter.Printf("===> 🌙 Sleeping: '%v'\n", description)
				}),
			)).Navigate(ctx)

			helpers.AssertNavigation(&entry.NaviTE, &helpers.TestOptions{
				FS:     FS,
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
			NaviTE: helpers.NaviTE{
				Given:        "File Subscription",
				Should:       "wake, then apply filter until the end",
				Subscription: enums.SubscribeFiles,
				Callback: func(node *core.Node) error {
					GinkgoWriter.Printf("---> WAKE-HIBERNATE-AND-FILTER-😵‍💫: '%v'\n", node.Path)

					return nil
				},
				ExpectedNoOf: helpers.Quantities{
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
			NaviTE: helpers.NaviTE{
				Given:        "File Subscription",
				Should:       "apply filter until sleep",
				Subscription: enums.SubscribeFiles,
				Callback: func(node *core.Node) error {
					GinkgoWriter.Printf("---> SLEEP-HIBERNATE-AND-FILTER-😴: '%v'\n", node.Path)

					return nil
				},
				ExpectedNoOf: helpers.Quantities{
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
			NaviTE: helpers.NaviTE{
				Given:        "File Subscription",
				Should:       "apply filter within hibernation range",
				Subscription: enums.SubscribeFiles,
				Callback: func(node *core.Node) error {
					GinkgoWriter.Printf("---> WAKE/SLEEP-HIBERNATE-AND-FILTER-😴: '%v'\n", node.Path)

					return nil
				},
				ExpectedNoOf: helpers.Quantities{
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
