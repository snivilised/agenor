package lfs_test

import (
	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok
	"github.com/snivilised/li18ngo"
	lab "github.com/snivilised/traverse/internal/laboratory"
	"github.com/snivilised/traverse/lfs"
)

var _ = Describe("op: write-file", Ordered, func() {
	var root string

	BeforeAll(func() {
		Expect(li18ngo.Use()).To(Succeed())

		root = lab.Repo("test")
	})

	Context("fs: WriteFileFS", func() {
		BeforeEach(func() {
			scratch(root)
		})

		Context("overwrite", func() {
			var fS lfs.WriteFileFS

			BeforeEach(func() {
				fS = lfs.NewWriteFileFS(lfs.At{
					Root:      root,
					Overwrite: true,
				})
			})

			Context("op: WriteFile", func() {
				When("given: file does not already exist", func() {
					It("🧪 should: write successfully", func() {
						Expect(require(root, lab.Static.FS.Scratch)).To(Succeed())
						name := lab.Static.FS.Write.Destination
						Expect(fS.WriteFile(
							name, lab.Static.FS.Write.Content, lab.Perms.File.Perm(),
						)).To(Succeed())
						Expect(AsFile(name)).To(ExistInFS(fS))
					})
				})
			})
		})

		Context("tentative", func() {
			var fS lfs.WriteFileFS

			BeforeEach(func() {
				fS = lfs.NewWriteFileFS(lfs.At{
					Root:      root,
					Overwrite: false,
				})
			})

			Context("op: WriteFile", func() {
				When("given: file does not already exist", func() {
					It("🧪 should: write successfully", func() {
						file := lab.Static.FS.Write.Destination
						Expect(require(
							root, lab.Static.FS.Scratch, file,
						)).To(Succeed())
						Expect(fS.WriteFile(
							file, lab.Static.FS.Write.Content, lab.Perms.File.Perm(),
						)).To(Succeed())
					})
				})
			})
		})
	})
})
