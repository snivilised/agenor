package command_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/snivilised/jaywalk/src/agenor/core"
	"github.com/snivilised/jaywalk/src/agenor/test/hanno"
	"github.com/snivilised/jaywalk/src/agenor/tfs"
	"github.com/snivilised/jaywalk/src/app/command"
	"github.com/snivilised/jaywalk/src/internal/services"
	"github.com/snivilised/jaywalk/src/locale"
	lab "github.com/snivilised/jaywalk/test/laboratory"
	"github.com/snivilised/li18ngo"
	"github.com/snivilised/nefilim/test/luna"
)

var _ = Describe("NavigatorUniversal", Ordered, func() {
	var (
		fS                *luna.MemFS
		configurationPath string
		bootstrap         command.Bootstrap
	)

	BeforeAll(func() {
		const (
			verbose = false
		)

		fS = hanno.Nuxx(verbose, lab.Static.RetroWave)
		configurationPath = hanno.Repo("test/data/")

		Expect(li18ngo.Register(
			func(o *li18ngo.UseOptions) {
				o.From.Sources = li18ngo.TranslationFiles{
					locale.SourceID: li18ngo.TranslationSource{Name: "agenor"},
				}
			},
		)).To(Succeed())
	})

	BeforeEach(func() {
		bootstrap = command.Bootstrap{}
		services.Reset()
	})

	DescribeTable("Walk command",
		func(ctx SpecContext, entry *lab.GeneralTE) {
			tester := hanno.CommandTester{
				Args: entry.Args,
				Root: bootstrap.Root(func(co *command.ConfigureAppOptions) {
					co.Detector = &DetectorStub{}
					co.ConfigInfo.Name = configName
					co.ConfigInfo.ConfigPath = configurationPath
					co.GetForest = func(_ string) *core.Forest {
						return &core.Forest{
							T: fS,
							R: tfs.New(),
						}
					}
				}),
			}

			_, err := tester.Execute()
			entry.Asserter(err)
		},
		lab.FormatGeneralTestDescription,

		// === regular =============================================================

		Entry(nil, &lab.GeneralTE{
			DescribedTE: lab.DescribedTE{
				Given:  "walk invoked with action",
				Should: "result in no error",
			},
			NaviTE: lab.NaviTE{
				Args: []string{
					"walk", "RETRO-WAVE", "--action", "echo", "--theme", "system",
				},
				Asserter: func(err error) {
					Expect(err).Error().To(BeNil())
				},
			},
		}),

		// === errors ==============================================================

		Entry(nil, &lab.GeneralTE{
			DescribedTE: lab.DescribedTE{
				Given:  "walk invoked missing action and pipeline",
				Should: "🧪 result in one of flags constraint error",
			},
			NaviTE: lab.NaviTE{
				Args: []string{
					"walk", "RETRO-WAVE", "--theme", "system",
				},
				Asserter: func(err error) {
					Expect(err).Error().NotTo(BeNil())
					Expect(err.Error()).To(ContainSubstring("at least one of the flags in the group"))
				},
			},
		}),

		Entry(nil, &lab.GeneralTE{
			DescribedTE: lab.DescribedTE{
				Given:  "walk invoked with filter",
				Should: "result in no error",
			},
			NaviTE: lab.NaviTE{
				Args: []string{
					"walk", "RETRO-WAVE", "--action", "echo", "--theme", "system", "--files", "*|.flac",
				},
				Asserter: func(err error) {
					Expect(err).Error().To(BeNil())
				},
			},
		}),
	)
})
