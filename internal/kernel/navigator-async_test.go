package kernel_test

import (
	"context"
	"fmt"
	"sync"
	"time"

	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok

	age "github.com/snivilised/agenor"
	"github.com/snivilised/agenor/core"
	"github.com/snivilised/agenor/enums"
	"github.com/snivilised/agenor/internal/enclave"
	lab "github.com/snivilised/agenor/internal/laboratory"
	"github.com/snivilised/agenor/internal/services"
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

func PrimeBuilder(entry *lab.AsyncOkTE, path string, fS *luna.MemFS) *age.Builders {
	return age.Prime(
		&pref.Using{
			Tree:         path,
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
		},
		Settings(entry)...,
	)
}

func ResumeBuilder(entry *lab.AsyncOkTE, path string, fS *luna.MemFS) *age.Builders {
	return age.Resume(
		&pref.Relic{
			From:     path,
			Strategy: enums.ResumeStrategyFastward,
			Head: pref.Head{
				Handler: entry.Callback,
				GetForest: func(_ string) *core.Forest {
					return &core.Forest{
						T: fS,
						R: tfs.New(),
					}
				},
			},
		},
		Settings(entry)...,
	)
}

func Settings(entry *lab.AsyncOkTE) []pref.Option {
	return []pref.Option{
		age.WithOnBegin(lab.Begin("🛡️")),
		age.WithOnEnd(lab.End("🏁")),
		age.IfOptionF(entry.NoWorkers > 0, func() pref.Option {
			return age.WithNoW(entry.NoWorkers)
		}),
		age.IfOptionF(entry.CPU, func() pref.Option {
			return age.WithCPU()
		}),
	}
}

var _ = Describe("Navigator", Ordered, func() {
	var (
		from string
		fS   *luna.MemFS
	)

	BeforeAll(func() {
		Expect(li18ngo.Use(
			func(o *li18ngo.UseOptions) {
				o.From.Sources = li18ngo.TranslationFiles{
					locale.SourceID: li18ngo.TranslationSource{Name: "agenor"},
				}
			},
		)).To(Succeed())

		fS = hanno.Nuxx(verbose, lab.Static.RetroWave)
		from = lab.GetJSONPath()
	})

	BeforeEach(func() {
		services.Reset()
	})

	DescribeTable("run",
		func(specCtx SpecContext, entry *lab.AsyncOkTE) {
			lab.WithTestContext(specCtx, func(ctx context.Context) {
				var wg sync.WaitGroup

				path := entry.Path()

				result, err := age.Run(&wg).Configure(enclave.Loader(func(active *core.ActiveState) {
					GinkgoWriter.Printf("===> 🐚 restoring state: resume at=%v, subscription=%v\n",
						entry.Resume.At, entry.Subscription,
					)
					active.Tree = lab.Static.RetroWave
					active.Depth = 2
					active.TraverseDescription.IsRelative = true
					active.ResumeDescription.IsRelative = false
					active.Subscription = entry.Subscription
					active.CurrentPath = entry.Resume.At
				})).Extent(
					entry.Builder(entry, path, fS),
				).Navigate(ctx)

				wg.Wait()
				Expect(err).To(Succeed())
				Expect(result).NotTo(BeNil())
			})
		},
		func(entry *lab.AsyncOkTE) string {
			return fmt.Sprintf("🧪 ===> given: '%v', should: '%v'",
				entry.Given, entry.Should,
			)
		},

		Entry(nil, &lab.AsyncOkTE{
			AsyncTE: lab.AsyncTE{
				Given:  "Primary Session WithCPUPool",
				Should: "run with context",
				Callback: func(servant core.Servant) error {
					node := servant.Node()
					name := node.Extension.Name
					GinkgoWriter.Printf("---> 🌀 ASYNC//%v-PRIME-CALLBACK(CPU): '%v'\n", name, node.Path)

					return nil
				},
				Builder:      PrimeBuilder,
				Path:         func() string { return lab.Static.RetroWave },
				Subscription: enums.SubscribeUniversal,
				CPU:          true,
			},
		}, SpecTimeout(time.Second*2)),

		Entry(nil, &lab.AsyncOkTE{
			AsyncTE: lab.AsyncTE{
				Given:  "Primary Session NoW=3",
				Should: "run with context",
				Callback: func(servant core.Servant) error {
					node := servant.Node()
					name := node.Extension.Name
					GinkgoWriter.Printf("---> 🌀 ASYNC//%v-PRIME-CALLBACK(NoW=3): '%v'\n", name, node.Path)

					return nil
				},
				Builder:      PrimeBuilder,
				Path:         func() string { return lab.Static.RetroWave },
				Subscription: enums.SubscribeUniversal,
				NoWorkers:    3,
			},
		}, SpecTimeout(time.Second*2)),

		Entry(nil, &lab.AsyncOkTE{
			AsyncTE: lab.AsyncTE{
				Given:  "Resume Session NoW=3",
				Should: "run with context",
				Callback: func(servant core.Servant) error {
					node := servant.Node()
					name := node.Extension.Name
					GinkgoWriter.Printf("---> 🌀 ASYNC//%v-RESUME-CALLBACK(NoW=3): '%v'\n", name, node.Path)

					return nil
				},
				Builder:      ResumeBuilder,
				Path:         func() string { return from },
				Subscription: enums.SubscribeUniversal,
				NoWorkers:    3,
				Resume: lab.AsyncResumeTE{
					At:       ResumeAtTeenageColor,
					Strategy: enums.ResumeStrategyFastward,
				},
			},
		}, SpecTimeout(time.Second*2)),
	)
})
