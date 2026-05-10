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
		services.Reset()
	})

	Context("given: walk invoked with action", func() {
		It("🧪 should: not error", func() {
			bootstrap := command.Bootstrap{}
			tester := hanno.CommandTester{
				Args: []string{
					"walk", "RETRO-WAVE", "--action", "echo", "--theme", "system",
				},
				Root: bootstrap.Root(func(co *command.ConfigureOptions) {
					co.Detector = &DetectorStub{}
					co.Config.Name = configName
					co.Config.ConfigPath = configurationPath
					co.GetForest = func(_ string) *core.Forest {
						return &core.Forest{
							T: fS,
							R: tfs.New(),
						}
					}
				}),
			}
			_, err := tester.Execute()
			Expect(err).Error().To(BeNil())
		})
	})

	Context("given: walk invoked missing action and pipeline", func() {
		It("🧪 should: result in one of flags constraint error", func() {
			bootstrap := command.Bootstrap{}
			tester := hanno.CommandTester{
				Args: []string{
					"walk", "RETRO-WAVE", "--theme", "system",
				},
				Root: bootstrap.Root(func(co *command.ConfigureOptions) {
					co.Detector = &DetectorStub{}
					co.Config.Name = configName
					co.Config.ConfigPath = configurationPath
					co.GetForest = func(_ string) *core.Forest {
						return &core.Forest{
							T: fS,
							R: tfs.New(),
						}
					}
				}),
			}
			_, err := tester.Execute()
			Expect(err).Error().NotTo(BeNil())
			Expect(err.Error()).To(ContainSubstring("at least one of the flags in the group"))
		})
	})
})
