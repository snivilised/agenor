package lfs_test

import (
	"os"

	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok
	"github.com/snivilised/li18ngo"
	lab "github.com/snivilised/traverse/internal/laboratory"
	"github.com/snivilised/traverse/lfs"
)

var _ = Describe("op: create", Ordered, func() {
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

			Context("op: Create", func() {
				When("given: file does not already exist", func() {
					It("🧪 should: create successfully", func() {
						Expect(require(root, lab.Static.FS.Scratch)).To(Succeed())
						name := lab.Static.FS.Create.Destination
						file, err := fS.Create(name)
						Expect(err).To(Succeed())
						defer file.Close()

						Expect(AsFile(name)).To(ExistInFS(fS))
					})
				})

				When("given: file exists", func() {
					It("🧪 should: create successfully", func() {
						Expect(require(
							root, lab.Static.FS.Scratch, lab.Static.FS.Create.Destination,
						)).To(Succeed())
						name := lab.Static.FS.Create.Destination
						file, err := fS.Create(name)
						Expect(err).To(Succeed())
						defer file.Close()

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

			Context("op: Create", func() {
				When("given: file exists", func() {
					It("🧪 should: fail", func() {
						file := lab.Static.FS.Create.Destination
						Expect(require(
							root, lab.Static.FS.Scratch, file,
						)).To(Succeed())
						_, err := fS.Create(file)
						Expect(err).To(MatchError(os.ErrExist))
					})
				})
			})
		})
	})
})
