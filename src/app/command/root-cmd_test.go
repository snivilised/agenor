package command_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/snivilised/jaywalk/src/agenor/test/hanno"
	"github.com/snivilised/jaywalk/src/app/command"
	nef "github.com/snivilised/nefilim"
)

var _ = Describe("RootCmd", Ordered, func() {
	var (
		repo     string
		l10nPath string
	)

	BeforeAll(func() {
		repo = hanno.Repo("")
		l10nPath = hanno.Combine(repo, "test/data/l10n")
		fS := nef.NewUniversalABS()
		Expect(fS.DirectoryExists(l10nPath)).To(BeTrue())
	})

	Context("given: no arguments", func() {
		It("🧪 should: execute root without error", func() {
			bootstrap := command.Bootstrap{}
			tester := hanno.CommandTester{
				Args: []string{},
				Root: bootstrap.Root(func(co *command.ConfigureOptions) {
					co.Detector = &DetectorStub{}
					co.Config.Name = configName
					co.Config.ConfigPath = configPath
				}),
			}
			_, err := tester.Execute()
			Expect(err).Error().To(BeNil())
		})
	})

	Context("given: nav invoked directly", func() {
		It("🧪 should: not error and display help", func() {
			// nav is hidden but still reachable via its Use name. When
			// invoked, it should print the not-invocable message and
			// root help without returning an error.
			bootstrap := command.Bootstrap{}
			tester := hanno.CommandTester{
				Args: []string{"nav"},
				Root: bootstrap.Root(func(co *command.ConfigureOptions) {
					co.Detector = &DetectorStub{}
					co.Config.Name = configName
					co.Config.ConfigPath = configPath
				}),
			}
			_, err := tester.Execute()
			Expect(err).Error().To(BeNil())
		})
	})

	Context("given: exec invoked directly", func() {
		It("🧪 should: not error and display help", func() {
			bootstrap := command.Bootstrap{}
			tester := hanno.CommandTester{
				Args: []string{"nav", "exec"},
				Root: bootstrap.Root(func(co *command.ConfigureOptions) {
					co.Detector = &DetectorStub{}
					co.Config.Name = configName
					co.Config.ConfigPath = configPath
				}),
			}
			_, err := tester.Execute()
			Expect(err).Error().To(BeNil())
		})
	})

	Context("given: --theme flag", func() {
		It("🧪 should: be accepted on root", func() {
			bootstrap := command.Bootstrap{}
			tester := hanno.CommandTester{
				Args: []string{"--theme", "system"},
				Root: bootstrap.Root(func(co *command.ConfigureOptions) {
					co.Detector = &DetectorStub{}
					co.Config.Name = configName
					co.Config.ConfigPath = configPath
				}),
			}
			_, err := tester.Execute()
			Expect(err).Error().To(BeNil())
		})
	})
})
