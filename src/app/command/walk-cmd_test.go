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
			// nav is hidden but still reachable via its Use name. When
			// invoked, it should print the not-invocable message and
			// root help without returning an error.
			bootstrap := command.Bootstrap{}
			tester := hanno.CommandTester{
				Args: command.InjectGhostAncestors([]string{
					"walk", "RETRO-WAVE", "--action", "echo", "--theme", "system",
				}),
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
})
