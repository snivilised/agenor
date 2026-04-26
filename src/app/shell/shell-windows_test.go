//go:build windows

package shell_test

import (
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/snivilised/jaywalk/src/agenor/enums"
	"github.com/snivilised/jaywalk/src/app/shell"
	"github.com/snivilised/li18ngo"
)

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

// withEnv temporarily sets an environment variable for the duration of
// a test and restores the original value in DeferCleanup.
func withEnv(key, value string) {
	original, wasSet := os.LookupEnv(key)

	if value == "" {
		os.Unsetenv(key)
	} else {
		os.Setenv(key, value)
	}

	DeferCleanup(func() {
		if wasSet {
			os.Setenv(key, original)
		} else {
			os.Unsetenv(key)
		}
	})
}

// ---------------------------------------------------------------------------
// Detect() specs
// ---------------------------------------------------------------------------

var _ = Describe("shell.Detect on Windows", Ordered, func() {

	BeforeAll(func() {
		li18ngo.Register()
	})

	Context("when CYGWIN env var is set", func() {
		BeforeEach(func() {
			withEnv("CYGWIN", "nodosfilewarning")
			withEnv("MSYSTEM", "")
			withEnv("PSModulePath", "")
		})

		It("returns KindCygwin", func() {
			env, err := shell.Detect()

			Expect(err).To(BeNil())
			Expect(env.Kind).To(Equal(enums.ShellKindCygwin))
		})
	})

	Context("when MSYSTEM env var is set", func() {
		BeforeEach(func() {
			withEnv("CYGWIN", "")
			withEnv("MSYSTEM", "MINGW64")
			withEnv("PSModulePath", "")
		})

		It("returns KindMSYS2", func() {
			env, err := shell.Detect()

			Expect(err).To(BeNil())
			Expect(env.Kind).To(Equal(enums.ShellKindMSYS2))
		})
	})

	Context("when PSModulePath env var is set", func() {
		BeforeEach(func() {
			withEnv("CYGWIN", "")
			withEnv("MSYSTEM", "")
			// PSModulePath is already set in a real PowerShell session;
			// we preserve it here and rely on the CI runner being PowerShell.
		})

		It("returns KindPowerShell", func() {
			// This test is only meaningful when running inside a PowerShell
			// session on the CI Windows runner, where PSModulePath is set.
			if os.Getenv("PSModulePath") == "" {
				Skip("PSModulePath not set - not running inside PowerShell")
			}

			env, err := shell.Detect()

			Expect(err).To(BeNil())
			Expect(env.Kind).To(Equal(enums.ShellKindPowerShell))
		})
	})

	Context("when no shell-specific env vars are set", func() {
		BeforeEach(func() {
			withEnv("CYGWIN", "")
			withEnv("MSYSTEM", "")
			withEnv("PSModulePath", "")
		})

		It("returns KindCmdExe", func() {
			env, err := shell.Detect()

			Expect(err).To(BeNil())
			Expect(env.Kind).To(Equal(enums.ShellKindCmdExe))
		})
	})
})

// ---------------------------------------------------------------------------
// LocateFunc specs - cmd.exe environment
// ---------------------------------------------------------------------------

var _ = Describe("LocateFunc in cmd.exe environment", Ordered, func() {
	var locate shell.LocateFunc

	BeforeAll(func() {
		li18ngo.Register()
	})

	BeforeEach(func() {
		withEnv("CYGWIN", "")
		withEnv("MSYSTEM", "")
		withEnv("PSModulePath", "")

		env, err := shell.Detect()
		Expect(err).To(BeNil())
		Expect(env.Kind).To(Equal(enums.ShellKindCmdExe))

		locate = env.Locate
	})

	Context("when locating a binary on PATH", func() {
		It("resolves cmd.exe itself via where.exe", func() {
			// cmd.exe is always on PATH on Windows.
			resolved, err := locate("cmd.exe")

			Expect(err).To(BeNil())
			Expect(resolved).NotTo(BeEmpty())
		})
	})

	Context("when locating a cmd.exe builtin", func() {
		DescribeTable("returns success for known builtins",
			func(builtin string) {
				resolved, err := locate(builtin)

				Expect(err).To(BeNil())
				Expect(resolved).NotTo(BeEmpty())
			},
			Entry("echo", "echo"),
			Entry("dir", "dir"),
			Entry("cd", "cd"),
			Entry("set", "set"),
			Entry("copy", "copy"),
			Entry("del", "del"),
			Entry("md", "md"),
			Entry("rd", "rd"),
			Entry("type", "type"),
			Entry("cls", "cls"),
		)
	})

	Context("when locating a token that does not exist", func() {
		It("returns an empty string and an error", func() {
			resolved, err := locate("__jay_nonexistent_binary_xyz__")

			Expect(err).To(HaveOccurred())
			Expect(resolved).To(BeEmpty())
			Expect(err.Error()).To(ContainSubstring("__jay_nonexistent_binary_xyz__"))
		})
	})
})

// ---------------------------------------------------------------------------
// LocateFunc specs - PowerShell environment
// ---------------------------------------------------------------------------

var _ = Describe("LocateFunc in PowerShell environment", Ordered, func() {
	var locate shell.LocateFunc

	BeforeAll(func() {
		li18ngo.Register()
	})

	BeforeEach(func() {
		if os.Getenv("PSModulePath") == "" {
			Skip("PSModulePath not set - not running inside PowerShell")
		}

		withEnv("CYGWIN", "")
		withEnv("MSYSTEM", "")

		env, err := shell.Detect()
		Expect(err).To(BeNil())
		Expect(env.Kind).To(Equal(enums.ShellKindPowerShell))

		locate = env.Locate
	})

	Context("when locating a PowerShell cmdlet", func() {
		It("resolves Get-Command itself", func() {
			// Get-Command is a core cmdlet present in all PowerShell versions.
			resolved, err := locate("Get-Command")

			Expect(err).To(BeNil())
			Expect(resolved).NotTo(BeEmpty())
		})

		It("resolves Write-Output", func() {
			resolved, err := locate("Write-Output")

			Expect(err).To(BeNil())
			Expect(resolved).NotTo(BeEmpty())
		})
	})

	Context("when locating a binary on PATH", func() {
		It("resolves powershell.exe or pwsh.exe", func() {
			// At least one of these must exist if we are inside PowerShell.
			resolvedPwsh, errPwsh := locate("pwsh.exe")
			resolvedPs, errPs := locate("powershell.exe")

			Expect(errPwsh == nil || errPs == nil).To(BeTrue(),
				"expected at least one of pwsh.exe or powershell.exe to be locatable",
			)
			Expect(resolvedPwsh != "" || resolvedPs != "").To(BeTrue())
		})
	})

	Context("when locating a token that does not exist", func() {
		It("returns an empty string and an error", func() {
			resolved, err := locate("__jay_nonexistent_binary_xyz__")

			Expect(err).To(HaveOccurred())
			Expect(resolved).To(BeEmpty())
			Expect(err.Error()).To(ContainSubstring("__jay_nonexistent_binary_xyz__"))
		})
	})
})
