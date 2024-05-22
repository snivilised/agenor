package kernel_test

import (
	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok

	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/enums"
	"github.com/snivilised/traverse/internal/kernel"
	"github.com/snivilised/traverse/pref"
)

var _ = Describe("NavigatorUniversal", func() {
	var o *pref.Options

	BeforeEach(func() {
		o, _ = pref.Get()
	})

	Context("nav", func() {
		When("foo", func() {
			It("🧪 should: not fail", func() {
				nav, err := kernel.PrimeNav(
					pref.Using{
						Root:         RootPath,
						Subscription: enums.SubscribeUniversal,
						Handler: func(_ *core.Node) error {
							return nil
						},
					},
					o,
				)

				Expect(err).To(Succeed())
				Expect(nav).NotTo(BeNil())
			})
		})
	})
})
