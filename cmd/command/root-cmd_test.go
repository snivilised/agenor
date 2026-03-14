package command_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/snivilised/agenor/cmd/command"
	"github.com/snivilised/agenor/test/hanno"
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

	It("🧪 should: execute", func() {
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
