package command_test

import (
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"

	"github.com/snivilised/jaywalk/src/app/command"
	"github.com/snivilised/jaywalk/src/locale"
)

// buildTestHierarchy constructs a minimal parent/child cobra command pair
// with the named flags registered as persistent flags on the parent,
// mirroring the nav->exec->walk/run inheritance pattern.
func buildTestHierarchy(persistentFlags ...string) (child *cobra.Command) {
	parent := &cobra.Command{Use: "parent"}
	child = &cobra.Command{
		Use:  "child",
		RunE: func(_ *cobra.Command, _ []string) error { return nil },
	}

	for _, name := range persistentFlags {
		parent.PersistentFlags().String(name, "", "test flag: "+name)
	}

	parent.AddCommand(child)

	return child
}

// setFlag simulates the user supplying a flag value by marking it as
// changed on the child's merged flag set.
func setFlag(cmd *cobra.Command, name, value string) {
	cmd.InheritedFlags()
	_ = cmd.Flags().Set(name, value)
}

var _ = Describe("CobraUtils", func() {

	// -----------------------------------------------------------------------
	// MarkInheritedFlagRequired
	// -----------------------------------------------------------------------

	Describe("MarkInheritedFlagRequired", func() {
		Context("given: required flag is supplied", func() {
			It("🧪 should: return nil", func() {
				child := buildTestHierarchy("action")
				setFlag(child, "action", "convert")

				err := command.MarkInheritedFlagRequired(child, "action")
				Expect(err).To(BeNil())
			})
		})

		Context("given: required flag is not supplied", func() {
			It("🧪 should: return error naming the command and flag", func() {
				child := buildTestHierarchy("action")

				err := command.MarkInheritedFlagRequired(child, "action")
				Expect(err).NotTo(BeNil())

				var target *locale.MarkInheritedFlagsRequiredError
				Expect(errors.As(err, &target)).To(BeTrue())
				Expect(target.Command).To(Equal("child"))
				Expect(target.Flag).To(Equal("action"))
			})
		})

		Context("given: flag is defined locally rather than inherited", func() {
			It("🧪 should: return nil when supplied", func() {
				cmd := &cobra.Command{Use: "solo"}
				cmd.Flags().String("action", "", "local flag")
				_ = cmd.Flags().Set("action", "convert")

				err := command.MarkInheritedFlagRequired(cmd, "action")
				Expect(err).To(BeNil())
			})
		})
	})

	// -----------------------------------------------------------------------
	// MarkInheritedFlagsOneRequired
	// -----------------------------------------------------------------------

	Describe("MarkInheritedFlagsOneRequired", func() {
		Context("given: first flag is supplied", func() {
			It("🧪 should: return nil", func() {
				child := buildTestHierarchy("action", "pipeline")
				setFlag(child, "action", "convert")

				err := command.MarkInheritedFlagsOneRequired(child, "action", "pipeline")
				Expect(err).To(BeNil())
			})
		})

		Context("given: second flag is supplied", func() {
			It("🧪 should: return nil", func() {
				child := buildTestHierarchy("action", "pipeline")
				setFlag(child, "pipeline", "encode")

				err := command.MarkInheritedFlagsOneRequired(child, "action", "pipeline")
				Expect(err).To(BeNil())
			})
		})

		Context("given: both flags are supplied", func() {
			It("🧪 should: return nil", func() {
				child := buildTestHierarchy("action", "pipeline")
				setFlag(child, "action", "convert")
				setFlag(child, "pipeline", "encode")

				err := command.MarkInheritedFlagsOneRequired(child, "action", "pipeline")
				Expect(err).To(BeNil())
			})
		})

		Context("given: neither flag is supplied", func() {
			It("🧪 should: return error naming the command and missing flags", func() {
				child := buildTestHierarchy("action", "pipeline")

				err := command.MarkInheritedFlagsOneRequired(child, "action", "pipeline")
				Expect(err).NotTo(BeNil())

				var target *locale.MarkInheritedFlagsOneRequiredError
				Expect(errors.As(err, &target)).To(BeTrue())
				Expect(target.Command).To(Equal("child"))
				Expect(target.Flags).To(ContainSubstring("action"))
				Expect(target.Flags).To(ContainSubstring("pipeline"))
			})
		})
	})

	// -----------------------------------------------------------------------
	// MarkInheritedFlagsMutuallyExclusive
	// -----------------------------------------------------------------------

	Describe("MarkInheritedFlagsMutuallyExclusive", func() {
		Context("given: neither flag is supplied", func() {
			It("🧪 should: return nil", func() {
				child := buildTestHierarchy("cpu", "now")

				err := command.MarkInheritedFlagsMutuallyExclusive(child, "cpu", "now")
				Expect(err).To(BeNil())
			})
		})

		Context("given: first flag only is supplied", func() {
			It("🧪 should: return nil", func() {
				child := buildTestHierarchy("cpu", "now")
				setFlag(child, "cpu", "true")

				err := command.MarkInheritedFlagsMutuallyExclusive(child, "cpu", "now")
				Expect(err).To(BeNil())
			})
		})

		Context("given: second flag only is supplied", func() {
			It("🧪 should: return nil", func() {
				child := buildTestHierarchy("cpu", "now")
				setFlag(child, "now", "4")

				err := command.MarkInheritedFlagsMutuallyExclusive(child, "cpu", "now")
				Expect(err).To(BeNil())
			})
		})

		Context("given: both flags are supplied", func() {
			It("🧪 should: return error naming the command and conflicting flags", func() {
				child := buildTestHierarchy("cpu", "now")
				setFlag(child, "cpu", "true")
				setFlag(child, "now", "4")

				err := command.MarkInheritedFlagsMutuallyExclusive(child, "cpu", "now")
				Expect(err).NotTo(BeNil())

				var target *locale.MutuallyExclusiveFlagsPresentError
				Expect(errors.As(err, &target)).To(BeTrue())
				Expect(target.Command).To(Equal("child"))
				Expect(target.Flags).To(ContainSubstring("cpu"))
				Expect(target.Flags).To(ContainSubstring("now"))
			})
		})

		Context("given: three flags where two are supplied", func() {
			It("🧪 should: return error naming only the conflicting flags", func() {
				child := buildTestHierarchy("a", "b", "c")
				setFlag(child, "a", "foo")
				setFlag(child, "c", "bar")

				err := command.MarkInheritedFlagsMutuallyExclusive(child, "a", "b", "c")
				Expect(err).NotTo(BeNil())

				var target *locale.MutuallyExclusiveFlagsPresentError
				Expect(errors.As(err, &target)).To(BeTrue())
				Expect(target.Command).To(Equal("child"))
				Expect(target.Flags).To(ContainSubstring("a"))
				Expect(target.Flags).To(ContainSubstring("c"))
				Expect(target.Flags).NotTo(ContainSubstring("b"))
			})
		})
	})
})
