package age_test

import (
	"context"
	"path/filepath"
	"sync"

	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok

	age "github.com/snivilised/agenor"
	"github.com/snivilised/agenor/core"
	"github.com/snivilised/agenor/enums"
	"github.com/snivilised/agenor/internal/enclave"
	lab "github.com/snivilised/agenor/internal/laboratory"
	"github.com/snivilised/agenor/internal/services"
	"github.com/snivilised/agenor/life"
	"github.com/snivilised/agenor/locale"
	"github.com/snivilised/agenor/pref"
	"github.com/snivilised/agenor/test/hanno"
	"github.com/snivilised/agenor/tfs"
	"github.com/snivilised/li18ngo"
	"github.com/snivilised/nefilim/test/luna"
)

const (
	ResumeAtTeenageColor = "RETRO-WAVE/College/Teenage Color"
)

var _ = Describe("Director(Resume)", Ordered, func() {
	var (
		fS      *luna.MemFS
		restore pref.Option

		jsonPath, resumeAt, tree string
	)

	BeforeAll(func() {
		restore = func(o *age.Options) error {
			o.Events.Begin.On(func(_ *life.BeginState) {})

			return nil
		}
		fS = luna.NewMemFS()
		jsonPath = lab.GetJSONPath()
		tree = hanno.Repo("test")
		resumeAt = filepath.Join(tree, "hydra")

		Expect(li18ngo.Use(
			func(o *li18ngo.UseOptions) {
				o.From.Sources = li18ngo.TranslationFiles{
					locale.SourceID: li18ngo.TranslationSource{Name: "agenor"},
				}
			},
		)).To(Succeed())
	})

	BeforeEach(func() {
		services.Reset()
	})

	Context("simple", func() {
		Context("Walk", func() {
			It("🧪 should: walk resume navigation successfully", func(specCtx SpecContext) {
				lab.WithTestContext(specCtx, func(ctx context.Context, _ context.CancelFunc) {
					const depth = 2

					_, err := age.Walk().Configure(enclave.Loader(func(active *core.ActiveState) {
						active.Tree = tree
						active.Depth = depth
						active.TraverseDescription.IsRelative = true
						active.ResumeDescription.IsRelative = false
						active.Subscription = enums.SubscribeUniversal
						active.CurrentPath = ResumeAtTeenageColor
					})).Extent(age.Resume(
						&pref.Relic{
							Head: pref.Head{
								Handler: noOpHandler,
								GetForest: func(_ string) *core.Forest {
									return &core.Forest{
										T: fS,
										R: tfs.New(),
									}
								},
							},
							From:     jsonPath,
							Strategy: age.ResumeStrategyFastward,
						},
						age.WithDepth(depth),
						age.WithOnDescend(func(_ *core.Node) {}),
						age.WithFaultHandler(age.Accepter(lab.IgnoreFault)),
						restore,
					)).Navigate(ctx)

					Expect(err).To(Succeed())
				})
			})
		})

		Context("Run", func() {
			XIt("🧪 should: perform run navigation successfully", func(specCtx SpecContext) {
				lab.WithTestContext(specCtx, func(ctx context.Context, _ context.CancelFunc) {
					var wg sync.WaitGroup

					_, err := age.Run(&wg).Configure(enclave.Loader(func(active *core.ActiveState) {
						active.Tree = tree
						active.Depth = 2
						active.TraverseDescription.IsRelative = true
						active.ResumeDescription.IsRelative = false
						active.Subscription = enums.SubscribeUniversal
						active.CurrentPath = resumeAt
					})).Extent(age.Resume(
						&pref.Relic{
							Head: pref.Head{
								Handler: noOpHandler,
								GetForest: func(_ string) *core.Forest {
									return &core.Forest{
										T: fS,
										R: tfs.New(),
									}
								},
							},
							From:     jsonPath,
							Strategy: age.ResumeStrategySpawn,
						},
						age.WithOnDescend(func(_ *core.Node) {}),
						restore,
					)).Navigate(ctx)

					wg.Wait()
					Expect(err).To(Succeed())
				})
			})
		})
	})

	Context("features", func() {
		Context("Run", func() {
			When("filter", func() {
				It("🧪 should: register ok", func(specCtx SpecContext) {
					lab.WithTestContext(specCtx, func(ctx context.Context, _ context.CancelFunc) {
						var wg sync.WaitGroup

						_, err := age.Run(&wg).Configure().Extent(age.Resume(
							&pref.Relic{
								Head: pref.Head{
									Handler: noOpHandler,
								},
								From:     jsonPath,
								Strategy: age.ResumeStrategyFastward,
							},
							age.WithFilter(&pref.FilterOptions{}),
							age.WithFaultHandler(age.Accepter(lab.IgnoreFault)),
							restore,
						)).Navigate(ctx)

						wg.Wait()
						Expect(err).To(Succeed())

					})
				})
			})

			When("hibernate", func() {
				It("🧪 should: register ok", func(specCtx SpecContext) {
					lab.WithTestContext(specCtx, func(ctx context.Context, _ context.CancelFunc) {
						var wg sync.WaitGroup

						_, err := age.Run(&wg).Configure().Extent(age.Resume(
							&pref.Relic{
								Head: pref.Head{
									Handler: noOpHandler,
								},
								From:     jsonPath,
								Strategy: age.ResumeStrategyFastward,
							},
							age.WithHibernationFilterWake(&core.FilterDef{
								Description: "nonsense",
								Type:        enums.FilterTypeGlob,
								Pattern:     "*",
							}),
							age.WithFaultHandler(age.Accepter(lab.IgnoreFault)),
							restore,
						)).Navigate(ctx)

						wg.Wait()
						Expect(err).To(Succeed())
					})
				})
			})

			When("sample", func() {
				It("🧪 should: register ok", func(specCtx SpecContext) {
					lab.WithTestContext(specCtx, func(ctx context.Context, _ context.CancelFunc) {
						var wg sync.WaitGroup

						_, err := age.Run(&wg).Configure().Extent(age.Resume(
							&pref.Relic{
								Head: pref.Head{
									Handler: noOpHandler,
								},
								From:     jsonPath,
								Strategy: age.ResumeStrategyFastward,
							},
							age.WithSamplingOptions(&pref.SamplingOptions{
								NoOf: pref.EntryQuantities{
									Files:       files,
									Directories: directories,
								},
								Type: enums.SampleTypeSlice,
								Iteration: pref.SamplingIterationOptions{
									Each:  func(_ *core.Node) bool { return false },
									While: func(_ *pref.FilteredInfo) bool { return false },
								},
							}),
							age.WithFaultHandler(age.Accepter(lab.IgnoreFault)),
							restore,
						)).Navigate(ctx)

						wg.Wait()
						Expect(err).To(Succeed())
					})
				})
			})
		})
	})
})
