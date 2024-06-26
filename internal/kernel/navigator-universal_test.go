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
	"github.com/snivilised/traverse/internal/services"
)

var _ = Describe("NavigatorUniversal", func() {
	var memFS fstest.MapFS

	BeforeEach(func() {
		services.Reset()
		memFS = fstest.MapFS{
			filepath.Join(RootPath, "foo.txt"): {},
			RootPath:                           {},
		}
	})

	Context("nav", func() {
		When("foo", func() {
			It("🧪 should: not fail", func(specCtx SpecContext) {
				defer leaktest.Check(GinkgoT())()

				ctx, cancel := context.WithCancel(specCtx)
				defer cancel()

				_, err := tv.Walk().Configure().Extent(tv.Prime(
					&tv.Using{
						Root:         RootPath,
						Subscription: tv.SubscribeUniversal,
						Handler: func(_ *tv.Node) error {
							return nil
						},
						GetFS: func() fs.FS {
							return memFS
						},
					},

					tv.WithHookQueryStatus(memFS.Stat),
				),
				).Navigate(ctx)

				Expect(err).To(Succeed())
			})
		})
	})
})
