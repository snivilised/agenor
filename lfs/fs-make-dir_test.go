package lfs_test

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok
	"github.com/snivilised/li18ngo"
	lab "github.com/snivilised/traverse/internal/laboratory"
	"github.com/snivilised/traverse/lfs"
)

var _ = Describe("op: make-dir/all", Ordered, func() {
	var root string

	BeforeAll(func() {
		Expect(li18ngo.Use()).To(Succeed())

		root = lab.Repo("test")
	})

	Context("tentative", func() {
		Context("fs: MakeDirFS", func() {
			var (
				fS lfs.MakeDirFS
			)

			BeforeEach(func() {
				fS = lfs.NewMakeDirFS(lfs.At{
					Root:      root,
					Overwrite: false,
				})
				scratch(root)
			})

			Context("op: MakeDir", func() {
				When("given: path does not exist", func() {
					It("🧪 should: complete ok", func() {
						path := lab.Static.FS.Scratch
						Expect(fS.MakeDir(path, lab.Perms.Dir.Perm())).To(
							Succeed(), fmt.Sprintf("failed to MakeDir %q", path),
						)

						Expect(AsDirectory(path)).To(ExistInFS(fS))
					})
				})

				When("given: path already exists", func() {
					It("🧪 should: complete ok", func() {
						path := lab.Static.FS.Existing.Directory
						Expect(fS.MakeDir(path, lab.Perms.Dir.Perm())).To(
							Succeed(), fmt.Sprintf("failed to MakeDir %q", path),
						)
					})
				})
			})

			Context("op: MakeDirAll", func() {
				When("given: path does not exist", func() {
					It("🧪 should: complete ok", func() {
						path := lab.Static.FS.MakeDir.MakeAll
						Expect(fS.MakeDirAll(path, lab.Perms.Dir.Perm())).To(
							Succeed(), fmt.Sprintf("failed to MakeDir %q", path),
						)

						Expect(AsDirectory(path)).To(ExistInFS(fS))
					})
				})

				When("given: path already exists", func() {
					It("🧪 should: complete ok", func() {
						path := lab.Static.FS.Existing.Directory
						Expect(fS.MakeDir(path, lab.Perms.Dir.Perm())).To(
							Succeed(), fmt.Sprintf("failed to MakeDir %q", path),
						)

						Expect(AsDirectory(path)).To(ExistInFS(fS))
					})
				})
			})
		})
	})
})
