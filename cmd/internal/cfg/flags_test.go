package cfg_test

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"

	"github.com/snivilised/agenor/cmd/internal/cfg"
)

var _ = Describe("FlagResolver", func() {

	// Helper: build a minimal cobra.Command with int flags.
	newCmd := func(name string, defaults map[string]int) *cobra.Command {
		cmd := &cobra.Command{Use: name}
		for flag, def := range defaults {
			cmd.Flags().Int(flag, def, "")
		}
		return cmd
	}

	// Helper: simulate user supplying a flag explicitly.
	setFlag := func(cmd *cobra.Command, flag string, val int) {
		_ = cmd.Flags().Set(flag, intToString(val))
	}

	Describe("ResolveInt", func() {
		var (
			resolver *cfg.FlagResolver
			flags    cfg.FlagsConfig
		)

		BeforeEach(func() {
			flags = cfg.FlagsConfig{
				Invoke: cfg.FlagInvokeDefaults{
					"walk": {"files": 10},
					"any":  {"files": 5, "folders": 3},
				},
				Component: cfg.FlagComponentDefaults{
					"sampler": {"files": 7, "folders": 2},
				},
			}
			resolver = cfg.NewFlagResolver(flags)
		})

		Context("when the flag is explicitly set on the CLI", func() {
			It("returns the CLI value regardless of config", func() {
				cmd := newCmd("walk", map[string]int{"files": 1})
				setFlag(cmd, "files", 99)

				val, ok := resolver.ResolveInt(cmd, "files", "sampler")
				Expect(ok).To(BeTrue())
				Expect(val).To(Equal(99))
			})
		})

		Context("when the flag is not set on the CLI", func() {
			Context("and a command-specific invoke default exists", func() {
				It("uses the command-specific invoke default", func() {
					cmd := newCmd("walk", map[string]int{"files": 1})

					val, ok := resolver.ResolveInt(cmd, "files", "sampler")
					Expect(ok).To(BeTrue())
					Expect(val).To(Equal(10)) // walk.files from invoke
				})
			})

			Context("and only the 'any' wildcard invoke default exists", func() {
				It("uses the 'any' wildcard default", func() {
					cmd := newCmd("walk", map[string]int{"folders": 1})

					val, ok := resolver.ResolveInt(cmd, "folders", "sampler")
					Expect(ok).To(BeTrue())
					Expect(val).To(Equal(3)) // any.folders from invoke
				})
			})

			Context("and no invoke default exists but a component default does", func() {
				It("uses the component default", func() {
					// No invoke entry for "run"
					flags2 := cfg.FlagsConfig{
						Invoke: cfg.FlagInvokeDefaults{},
						Component: cfg.FlagComponentDefaults{
							"sampler": {"files": 7},
						},
					}
					r2 := cfg.NewFlagResolver(flags2)
					cmd := newCmd("run", map[string]int{"files": 1})

					val, ok := r2.ResolveInt(cmd, "files", "sampler")
					Expect(ok).To(BeTrue())
					Expect(val).To(Equal(7))
				})
			})

			Context("and no config defaults exist at all", func() {
				It("falls back to the cobra default", func() {
					r2 := cfg.NewFlagResolver(cfg.FlagsConfig{})
					cmd := newCmd("run", map[string]int{"files": 42})

					val, ok := r2.ResolveInt(cmd, "files", "")
					Expect(ok).To(BeTrue())
					Expect(val).To(Equal(42))
				})
			})
		})
	})

	Describe("ApplyShortOverrides", func() {
		It("remaps the shorthand for the named command's flags", func() {
			flags := cfg.FlagsConfig{
				Short: cfg.FlagShortOverride{
					"walk": {"foo": "F"},
				},
			}
			resolver := cfg.NewFlagResolver(flags)

			cmd := &cobra.Command{Use: "walk"}
			cmd.Flags().StringP("foo", "f", "", "a flag")

			resolver.ApplyShortOverrides(cmd)

			f := cmd.Flags().Lookup("foo")
			Expect(f).NotTo(BeNil())
			Expect(f.Shorthand).To(Equal("F"))
		})

		It("is a no-op for commands not in the short overrides", func() {
			flags := cfg.FlagsConfig{
				Short: cfg.FlagShortOverride{
					"walk": {"foo": "F"},
				},
			}
			resolver := cfg.NewFlagResolver(flags)

			cmd := &cobra.Command{Use: "run"}
			cmd.Flags().StringP("foo", "f", "", "a flag")

			resolver.ApplyShortOverrides(cmd)

			f := cmd.Flags().Lookup("foo")
			Expect(f.Shorthand).To(Equal("f")) // unchanged
		})
	})
})

// intToString converts an int to its string representation for pflag.Set.
func intToString(n int) string {
	return fmt.Sprintf("%d", n)
}
