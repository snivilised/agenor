package command_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/snivilised/jaywalk/src/app/command"
	"github.com/snivilised/jaywalk/src/agenor/test/hanno"

	nef "github.com/snivilised/nefilim"

	"golang.org/x/text/language"
)

const (
	configName = "jay-test"
	configPath = "../../test/data/configuration"
)

// DetectorStub is a test double that satisfies the LocaleDetector
// interface and always returns British English.
type DetectorStub struct {
}

func (j *DetectorStub) Scan() language.Tag {
	return language.BritishEnglish
}

var _ = Describe("Bootstrap", Ordered, func() {

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

	Context("given: root defined", func() {
		It("🧪 should: setup command without error", func() {
			bootstrap := command.Bootstrap{}
			rootCmd := bootstrap.Root(func(co *command.ConfigureOptions) {
				co.Detector = &DetectorStub{}
				co.Config.Name = configName
				co.Config.ConfigPath = configPath
			})
			Expect(rootCmd).NotTo(BeNil())
		})
	})
})
