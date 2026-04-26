//go:build !windows

package shell_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/snivilised/jaywalk/src/agenor/enums"
	"github.com/snivilised/jaywalk/src/app/shell"
	"github.com/snivilised/li18ngo"
)

var _ = Describe("shell.Detect on Unix", Ordered, func() {
	BeforeAll(func() {
		Expect(li18ngo.Register()).To(Succeed())
	})

	Context("when running on a Unix-like platform", func() {
		It("returns KindNativeUnix", func() {
			env, err := shell.Detect()

			Expect(err).To(BeNil())
			Expect(env.Kind).To(Equal(enums.ShellKindNativeUnix))
		})

		It("returns a non-nil Locate function", func() {
			env, err := shell.Detect()

			Expect(err).To(BeNil())
			Expect(env.Locate).NotTo(BeNil())
		})
	})
})

var _ = Describe("LocateFunc on Unix", Ordered, func() {
	var locate shell.LocateFunc

	BeforeAll(func() {
		Expect(li18ngo.Register()).To(Succeed())
	})

	BeforeEach(func() {
		env, err := shell.Detect()
		Expect(err).To(BeNil())
		locate = env.Locate
	})

	Context("when locating a binary on PATH", func() {
		It("returns a non-empty path and no error for 'sh'", func() {
			// /bin/sh is guaranteed to exist on any Unix system.
			resolved, err := locate("sh")

			Expect(err).To(BeNil())
			Expect(resolved).NotTo(BeEmpty())
		})
	})

	Context("when locating a shell builtin", func() {
		It("returns a non-empty result and no error for 'echo'", func() {
			// echo is a POSIX shell builtin present in every sh implementation.
			resolved, err := locate("echo")

			Expect(err).To(BeNil())
			Expect(resolved).NotTo(BeEmpty())
		})

		It("returns a non-empty result and no error for 'cd'", func() {
			resolved, err := locate("cd")

			Expect(err).To(BeNil())
			Expect(resolved).NotTo(BeEmpty())
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
