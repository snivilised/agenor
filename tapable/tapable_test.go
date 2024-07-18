package tapable_test

import (
	"io/fs"
	"os"
	"testing/fstest"

	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok

	tv "github.com/snivilised/traverse"
	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/pref"
)

const (
	root = "/traversal-root-path"
)

var fakeSubPath = &core.SubPathInfo{
	Root: root,
	Node: &core.Node{
		Extension: core.Root("/root", nil).Extension,
	},
}

var _ = Describe("Tapable", Ordered, func() {
	var (
		invoked bool
		o       *pref.Options
		err     error
		emptyFS fstest.MapFS
	)

	BeforeAll(func() {
		emptyFS = fstest.MapFS{
			".": &fstest.MapFile{
				Mode: os.ModeDir,
			},
		}
	})

	BeforeEach(func() {
		invoked = false
		o, err = pref.Get()
		Expect(err).To(Succeed())
	})

	Context("hooks", func() {
		When("FileSubPath hooked", func() {
			It("🧪 should: invoke hook", func() {
				o.Hooks.FileSubPath.Tap(func(_ *core.SubPathInfo) string {
					invoked = true
					return ""
				})
				o.Hooks.FileSubPath.Default()(fakeSubPath)
				o.Hooks.FileSubPath.Invoke()(fakeSubPath)

				Expect(invoked).To(BeTrue(), "FileSubPath hook not invoked")
			})
		})

		When("FolderSubPath hooked", func() {
			It("🧪 should: invoke hook", func() {
				o.Hooks.FolderSubPath.Tap(func(_ *core.SubPathInfo) string {
					invoked = true
					return ""
				})
				o.Hooks.FolderSubPath.Default()(fakeSubPath)
				o.Hooks.FolderSubPath.Invoke()(fakeSubPath)

				Expect(invoked).To(BeTrue(), "FolderSubPath hook not invoked")
			})
		})

		When("ReadDirectory hooked", func() {
			It("🧪 should: invoke hook", func() {
				o.Hooks.ReadDirectory.Tap(func(_ fs.ReadDirFS, _ string) ([]fs.DirEntry, error) {
					invoked = true
					return []fs.DirEntry{}, nil
				})

				sys := tv.NewNativeFS(root)
				_, _ = o.Hooks.ReadDirectory.Default()(sys, root)
				_, _ = o.Hooks.ReadDirectory.Invoke()(sys, root)

				Expect(invoked).To(BeTrue(), "ReadDirectory hook not invoked")
			})
		})

		When("QueryStatus hooked", func() {
			It("🧪 should: invoke hook", func() {
				o.Hooks.QueryStatus.Tap(func(_ fs.StatFS, _ string) (fs.FileInfo, error) {
					invoked = true
					return nil, nil
				})
				_, _ = o.Hooks.QueryStatus.Default()(emptyFS, root)
				_, _ = o.Hooks.QueryStatus.Invoke()(emptyFS, root)

				Expect(invoked).To(BeTrue(), "QueryStatus hook not invoked")
			})
		})

		When("Sort hooked", func() {
			It("🧪 should: invoke hook", func() {
				o.Hooks.Sort.Tap(func(_ []fs.DirEntry, _ ...any) {
					invoked = true
				})
				o.Hooks.Sort.Default()([]fs.DirEntry{})
				o.Hooks.Sort.Invoke()([]fs.DirEntry{})

				Expect(invoked).To(BeTrue(), "Sort hook not invoked")
			})
		})
	})
})
