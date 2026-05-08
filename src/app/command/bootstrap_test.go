package command_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/snivilised/jaywalk/src/agenor/test/hanno"
	"github.com/snivilised/jaywalk/src/app/command"
	"github.com/snivilised/li18ngo"
	nef "github.com/snivilised/nefilim"
	"github.com/spf13/cobra"
	"golang.org/x/text/language"
)

const (
	configName = "jay-test"
	configPath = "../../test/data/configuration"
)

// DetectorStub is a test double that satisfies the LocaleDetector
// interface and always returns British English.
type DetectorStub struct{}

func (j *DetectorStub) Scan() language.Tag {
	return language.BritishEnglish
}

func applyTestConfig(co *command.ConfigureOptions) {
	co.Detector = &DetectorStub{}
	co.Config.Name = configName
	co.Config.ConfigPath = configPath
}

func buildRoot() *cobra.Command {
	return (&command.Bootstrap{}).Root(applyTestConfig)
}

var _ = Describe("Bootstrap", Ordered, func() {
	var (
		repo     string
		l10nPath string
	)

	BeforeAll(func() {
		Expect(li18ngo.Register()).To(Succeed())
	})

	BeforeAll(func() {
		repo = hanno.Repo("")
		l10nPath = hanno.Combine(repo, "test/data/l10n")
		fS := nef.NewUniversalABS()
		Expect(fS.DirectoryExists(l10nPath)).To(BeTrue())
	})

	Context("given: root defined", func() {
		It("🧪 should: build command tree without error", func() {
			Expect(buildRoot()).NotTo(BeNil())
		})
	})

	Context("given: command tree built", func() {

		Context("ghost commands", func() {
			It("🧪 should: hide nav from user-facing help", func() {
				root := buildRoot()
				for _, sub := range root.Commands() {
					if sub.Name() == "nav" {
						Expect(sub.Hidden).To(BeTrue())
						return
					}
				}
			})

			It("🧪 should: hide exec from user-facing help", func() {
				root := buildRoot()
				for _, sub := range root.Commands() {
					if sub.Name() == "nav" {
						for _, navSub := range sub.Commands() {
							if navSub.Name() == "exec" {
								Expect(navSub.Hidden).To(BeTrue())
								return
							}
						}
					}
				}
			})
		})

		Context("flag inheritance - nav level", func() {
			It("🧪 should: expose --subscribe on walk", func() {
				walkCmd, _, err := buildRoot().Find([]string{"nav", "exec", "walk"})
				Expect(err).To(BeNil())
				walkCmd.InheritedFlags()
				Expect(walkCmd.Flags().Lookup("subscribe")).NotTo(BeNil())
			})

			It("🧪 should: expose --subscribe on sprint", func() {
				sprintCmd, _, err := buildRoot().Find([]string{"nav", "exec", "sprint"})
				Expect(err).To(BeNil())
				sprintCmd.InheritedFlags()
				Expect(sprintCmd.Flags().Lookup("subscribe")).NotTo(BeNil())
			})

			It("🧪 should: expose --subscribe on query", func() {
				queryCmd, _, err := buildRoot().Find([]string{"nav", "query"})
				Expect(err).To(BeNil())
				queryCmd.InheritedFlags()
				Expect(queryCmd.Flags().Lookup("subscribe")).NotTo(BeNil())
			})

			It("🧪 should: expose --action on walk", func() {
				walkCmd, _, err := buildRoot().Find([]string{"nav", "exec", "walk"})
				Expect(err).To(BeNil())
				walkCmd.InheritedFlags()
				Expect(walkCmd.Flags().Lookup("action")).NotTo(BeNil())
			})

			It("🧪 should: expose --action on query", func() {
				queryCmd, _, err := buildRoot().Find([]string{"nav", "query"})
				Expect(err).To(BeNil())
				queryCmd.InheritedFlags()
				Expect(queryCmd.Flags().Lookup("action")).NotTo(BeNil())
			})

			It("🧪 should: expose --pipeline on walk", func() {
				walkCmd, _, err := buildRoot().Find([]string{"nav", "exec", "walk"})
				Expect(err).To(BeNil())
				walkCmd.InheritedFlags()
				Expect(walkCmd.Flags().Lookup("pipeline")).NotTo(BeNil())
			})
		})

		Context("flag inheritance - exec level", func() {
			It("🧪 should: expose --resume on walk", func() {
				walkCmd, _, err := buildRoot().Find([]string{"nav", "exec", "walk"})
				Expect(err).To(BeNil())
				walkCmd.InheritedFlags()
				Expect(walkCmd.Flags().Lookup("resume")).NotTo(BeNil())
			})

			It("🧪 should: expose --resume on sprint", func() {
				sprintCmd, _, err := buildRoot().Find([]string{"nav", "exec", "sprint"})
				Expect(err).To(BeNil())
				sprintCmd.InheritedFlags()
				Expect(sprintCmd.Flags().Lookup("resume")).NotTo(BeNil())
			})

			It("🧪 should: not expose --resume on query", func() {
				queryCmd, _, err := buildRoot().Find([]string{"nav", "query"})
				Expect(err).To(BeNil())
				queryCmd.InheritedFlags()
				Expect(queryCmd.Flags().Lookup("resume")).To(BeNil())
			})
		})

		Context("flag inheritance - root level", func() {
			It("🧪 should: expose --theme on walk", func() {
				walkCmd, _, err := buildRoot().Find([]string{"nav", "exec", "walk"})
				Expect(err).To(BeNil())
				walkCmd.InheritedFlags()
				Expect(walkCmd.Flags().Lookup("theme")).NotTo(BeNil())
			})

			It("🧪 should: expose --theme on query", func() {
				queryCmd, _, err := buildRoot().Find([]string{"nav", "query"})
				Expect(err).To(BeNil())
				queryCmd.InheritedFlags()
				Expect(queryCmd.Flags().Lookup("theme")).NotTo(BeNil())
			})

			It("🧪 should: not register --subscribe as a root persistent flag", func() {
				Expect(buildRoot().PersistentFlags().Lookup("subscribe")).To(BeNil())
			})

			It("🧪 should: not register --resume as a root persistent flag", func() {
				Expect(buildRoot().PersistentFlags().Lookup("resume")).To(BeNil())
			})
		})

		Context("worker-pool flags - sprint exclusive", func() {
			It("🧪 should: expose --cpu as a local flag on sprint", func() {
				sprintCmd, _, err := buildRoot().Find([]string{"nav", "exec", "sprint"})
				Expect(err).To(BeNil())
				Expect(sprintCmd.Flags().Lookup("cpu")).NotTo(BeNil())
			})

			It("🧪 should: not expose --cpu on walk", func() {
				walkCmd, _, err := buildRoot().Find([]string{"nav", "exec", "walk"})
				Expect(err).To(BeNil())
				Expect(walkCmd.Flags().Lookup("cpu")).To(BeNil())
				Expect(walkCmd.InheritedFlags().Lookup("cpu")).To(BeNil())
			})

			It("🧪 should: not expose --cpu on query", func() {
				queryCmd, _, err := buildRoot().Find([]string{"nav", "query"})
				Expect(err).To(BeNil())
				Expect(queryCmd.Flags().Lookup("cpu")).To(BeNil())
				Expect(queryCmd.InheritedFlags().Lookup("cpu")).To(BeNil())
			})
		})
	})
})
