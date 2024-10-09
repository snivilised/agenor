package lfs_test

import (
	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok
	"github.com/snivilised/li18ngo"
	lab "github.com/snivilised/traverse/internal/laboratory"
	"github.com/snivilised/traverse/lfs"
)

var _ = Describe("op: copy/all", Ordered, func() {
	var (
		root   string
		fS     lfs.UniversalFS
		single string
	)

	BeforeAll(func() {
		Expect(li18ngo.Use()).To(Succeed())

		root = lab.Repo("test")
	})

	BeforeEach(func() {
		fS = lfs.NewUniversalFS(lfs.At{
			Root:      root,
			Overwrite: false,
		})
		scratch(root)
	})

	Context("op: Copy", func() {
		When("given: ", func() {
			It("🧪 should: ", func() {
				_ = fS
				_ = single
			})
		})
	})

	Context("op: CopyAll", func() {
		When("given: ", func() {
			It("🧪 should: ", func() {

			})
		})
	})
})
