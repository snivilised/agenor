package cycle_test

import (
	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok

	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/cycle"
)

var _ = Describe("Notify", func() {
	Context("foo", func() {
		It("should:", func() {
			const path = "/traversal-root"

			var (
				notifications cycle.Controls
				taps          cycle.Events
				begun         bool
				ended         bool
			)

			// init(registry->options):
			//
			taps.Bind(&notifications)

			// client:
			//
			taps.Begin.On(func(root string) {
				begun = true
				Expect(root).To(Equal(path))
			})

			taps.End.On(func(_ core.TraverseResult) {
				ended = true
			})

			// component side:
			//
			notifications.Begin.Dispatch.Invoke(path)
			notifications.End.Dispatch.Invoke(nil)

			Expect(begun).To(BeTrue())
			Expect(ended).To(BeTrue())
		})
	})
})
