package kernel_test

import (
	"context"
	"io/fs"
	"path/filepath"
	"testing/fstest"

	"github.com/fortytw2/leaktest"
	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok

	tv "github.com/snivilised/traverse"
	"github.com/snivilised/traverse/internal/helpers"
	"github.com/snivilised/traverse/internal/services"
)

var _ = Describe("NavigatorUniversal", Ordered, func() {
	var (
		memFS fstest.MapFS
		root  string
	)

	BeforeAll(func() {
		const (
			verbose = true
		)
		var portion = filepath.Join("MUSICO", "bass")
		memFS, root = helpers.Musico(portion, verbose)
		Expect(root).NotTo(BeEmpty())
	})

	BeforeEach(func() {
		services.Reset()
	})

	Context("nav", func() {
		When("foo", func() {
			It("🧪 should: not fail", func(specCtx SpecContext) {
				defer leaktest.Check(GinkgoT())()

				ctx, cancel := context.WithCancel(specCtx)
				defer cancel()

				_, err := tv.Walk().Configure().Extent(tv.Prime(
					&tv.Using{
						Root:         root,
						Subscription: tv.SubscribeUniversal,
						Handler: func(_ *tv.Node) error {
							return nil
						},
						GetFS: func() fs.FS {
							return memFS
						},
					},

					tv.WithHookQueryStatus(func(path string) (fs.FileInfo, error) {
						return memFS.Stat(helpers.TrimRoot(path))
					}),
					tv.WithHookReadDirectory(func(_ fs.FS, dirname string) ([]fs.DirEntry, error) {
						return memFS.ReadDir(helpers.TrimRoot(dirname))
					}),
				),
				).Navigate(ctx)

				Expect(err).To(Succeed())
			})
		})
	})
})
