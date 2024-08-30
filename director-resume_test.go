package tv_test

import (
	"context"
	"io/fs"
	"os"
	"sync"
	"testing/fstest"

	"github.com/fortytw2/leaktest"
	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok

	"github.com/snivilised/li18ngo"
	tv "github.com/snivilised/traverse"
	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/cycle"
	"github.com/snivilised/traverse/enums"
	"github.com/snivilised/traverse/internal/services"
	"github.com/snivilised/traverse/locale"
	"github.com/snivilised/traverse/pref"
)

var _ = Describe("Director(Resume)", Ordered, func() {
	var (
		emptyFS fstest.MapFS
		restore pref.Option
	)

	BeforeAll(func() {
		restore = func(o *tv.Options) error {
			o.Events.Begin.On(func(_ *cycle.BeginState) {})

			return nil
		}
		emptyFS = fstest.MapFS{
			".": &fstest.MapFile{
				Mode: os.ModeDir,
			},
		}

		Expect(li18ngo.Use(
			func(o *li18ngo.UseOptions) {
				o.From.Sources = li18ngo.TranslationFiles{
					locale.SourceID: li18ngo.TranslationSource{Name: "traverse"},
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
				defer leaktest.Check(GinkgoT())()

				ctx, cancel := context.WithCancel(specCtx)
				defer cancel()

				const depth = 2

				_, err := tv.Walk().Configure().Extent(tv.Resume(
					&tv.Was{
						Using: tv.Using{
							Subscription: tv.SubscribeFiles,
							Handler:      noOpHandler,
							GetReadDirFS: func() fs.ReadDirFS {
								return emptyFS
							},
							GetQueryStatusFS: func(_ fs.FS) fs.StatFS {
								return emptyFS
							},
						},
						From:     RestorePath,
						Strategy: tv.ResumeStrategyFastward,
					},
					tv.WithDepth(depth),
					tv.WithOnDescend(func(_ *core.Node) {}),
					restore,
				)).Navigate(ctx)

				Expect(err).To(Succeed())
			})
		})

		Context("Run", func() {
			It("🧪 should: perform run navigation successfully", func(specCtx SpecContext) {
				defer leaktest.Check(GinkgoT())()

				ctx, cancel := context.WithCancel(specCtx)
				defer cancel()

				var wg sync.WaitGroup

				_, err := tv.Run(&wg).Configure().Extent(tv.Resume(
					&tv.Was{
						Using: tv.Using{
							Subscription: tv.SubscribeFiles,
							Handler:      noOpHandler,
						},
						From:     RestorePath,
						Strategy: tv.ResumeStrategySpawn,
					},
					tv.WithOnDescend(func(_ *core.Node) {}),
					restore,
				)).Navigate(ctx)

				wg.Wait()
				Expect(err).To(Succeed())
			})
		})
	})

	Context("features", func() {
		Context("Run", func() {
			When("filter", func() {
				It("🧪 should: register ok", func(specCtx SpecContext) {
					defer leaktest.Check(GinkgoT())()

					ctx, cancel := context.WithCancel(specCtx)
					defer cancel()

					var wg sync.WaitGroup

					_, err := tv.Run(&wg).Configure().Extent(tv.Resume(
						&tv.Was{
							Using: tv.Using{
								Subscription: tv.SubscribeFiles,
								Handler:      noOpHandler,
							},
							From:     RestorePath,
							Strategy: tv.ResumeStrategySpawn,
						},
						tv.WithFilter(&pref.FilterOptions{}),
						restore,
					)).Navigate(ctx)

					wg.Wait()
					Expect(err).To(Succeed())
				})
			})

			When("hibernate", func() {
				It("🧪 should: register ok", func(specCtx SpecContext) {
					defer leaktest.Check(GinkgoT())()

					ctx, cancel := context.WithCancel(specCtx)
					defer cancel()

					var wg sync.WaitGroup

					_, err := tv.Run(&wg).Configure().Extent(tv.Resume(
						&tv.Was{
							Using: tv.Using{
								Subscription: tv.SubscribeFiles,
								Handler:      noOpHandler,
							},
							From:     RestorePath,
							Strategy: tv.ResumeStrategySpawn,
						},
						tv.WithHibernationFilterWake(&core.FilterDef{
							Description: "nonsense",
							Type:        enums.FilterTypeGlob,
							Pattern:     "*",
						}),
						restore,
					)).Navigate(ctx)

					wg.Wait()
					Expect(err).To(Succeed())
				})
			})

			When("sample", func() {
				It("🧪 should: register ok", func(specCtx SpecContext) {
					defer leaktest.Check(GinkgoT())()

					ctx, cancel := context.WithCancel(specCtx)
					defer cancel()

					var wg sync.WaitGroup

					_, err := tv.Run(&wg).Configure().Extent(tv.Resume(
						&tv.Was{
							Using: tv.Using{
								Subscription: tv.SubscribeFiles,
								Handler:      noOpHandler,
							},
							From:     RestorePath,
							Strategy: tv.ResumeStrategySpawn,
						},
						tv.WithSampling(&pref.SamplingOptions{
							NoOf: pref.EntryQuantities{
								Files:   files,
								Folders: folders,
							},
							SampleType: enums.SampleTypeSlice,
							Iteration: pref.SamplingIterationOptions{
								Each:  func(_ *core.Node) bool { return false },
								While: func(_ *pref.FilteredInfo) bool { return false },
							},
						}),
						restore,
					)).Navigate(ctx)

					wg.Wait()
					Expect(err).To(Succeed())
				})
			})
		})
	})
})
