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
	"github.com/snivilised/pants"
)

const (
	ResumeAtTeenageColor = "RETRO-WAVE/College/Teenage Color"
)

func PrimeBuilder(post *lab.AsyncPostage) *age.Builders {
	return age.Prime(
		&pref.Using{
			Tree:         post.Path,
			Subscription: post.Entry.Subscription,
			Head: pref.Head{
				Handler: post.Entry.Callback,
				GetForest: func(_ string) *core.Forest {
					return &core.Forest{
						T: post.FS,
						R: tfs.New(),
					}
				},
			},
		},
		Settings(post)...,
	)
}

func ResumeBuilder(post *lab.AsyncPostage) *age.Builders {
	return age.Resume(
		&pref.Relic{
			From:     post.Path,
			Strategy: enums.ResumeStrategyFastward,
			Head: pref.Head{
				Handler: post.Entry.Callback,
				GetForest: func(_ string) *core.Forest {
					return &core.Forest{
						T: post.FS,
						R: tfs.New(),
					}
				},
			},
		},
		Settings(post)...,
	)
}

func Settings(post *lab.AsyncPostage) []pref.Option {
	return []pref.Option{
		age.WithOnBegin(lab.Begin("🛡️")),
		age.WithOnEnd(lab.End("🏁")),
		age.IfOptionF(post.Entry.NoWorkers > 0, func() pref.Option {
			return age.WithNoW(post.Entry.NoWorkers)
		}),
		age.IfOptionF(post.Entry.CPU, func() pref.Option {
			return age.WithCPU()
		}),
		age.IfOptionF(post.Entry.Consume, func() pref.Option {
			return age.WithOutput(&pref.OutputOptions{
				CheckCloseInterval: time.Second / 10,
				TimeoutOnSend:      time.Second * 3,
				On:                 post.On,
			})
		}),
	}
}

func consumeOk(_ context.Context,
	outs core.OutputStream,
	wg pants.WaitGroup,
) {
	defer wg.Done()

	// We don't need to use a timeout on the observe channel
	// because the navigator invokes Conclude, which results in
	// the observe channel being closed, terminating the range.
	// This aspect is specific to this example and to test
	// cancellation, we'll have to use a select.
	//
	for output := range outs {
		GinkgoWriter.Printf("🍒 payload: '%v', id: '%v', seq: '%v' (e: '%v')\n",
			output.Payload, output.ID, output.SequenceNo, output.Error,
		)
	}

	GinkgoWriter.Println("===> 🍒 finished consuming output channel.")
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
			lab.WithTestContext(specCtx, func(ctx context.Context, _ context.CancelFunc) {
				var (
					wg sync.WaitGroup
					on core.OutputFunc
				)

				if entry.Consume {
					on = func(outs core.OutputStream) {
						wg.Add(1)
						go consumeOk(ctx, outs, &wg)
					}
				}

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
					entry.Builder(&lab.AsyncPostage{
						Entry: entry,
						Path:  entry.Path(),
						FS:    fS,
						On:    on,
					}),
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

		Entry(nil, &lab.AsyncOkTE{
			AsyncTE: lab.AsyncTE{
				Given:  "Primary Session With Output",
				Should: "consume output",
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
			Consume: true,
		}, SpecTimeout(time.Second*2)),

		Entry(nil, &lab.AsyncOkTE{
			AsyncTE: lab.AsyncTE{
				Given:  "Resume Session With Output",
				Should: "consume output",
				Callback: func(servant core.Servant) error {
					node := servant.Node()
					name := node.Extension.Name
					GinkgoWriter.Printf("---> 🌀 ASYNC//%v-PRIME-CALLBACK(NoW=3): '%v'\n", name, node.Path)

					return nil
				},
				Builder: ResumeBuilder,
				Resume: lab.AsyncResumeTE{
					At:       lab.Static.TeenageColor,
					Strategy: enums.ResumeStrategyFastward,
				},
				Path:         lab.GetJSONPath,
				Subscription: enums.SubscribeUniversal,
				NoWorkers:    3,
			},
			Consume: true,
		}, SpecTimeout(time.Second*2)),
	)
})
