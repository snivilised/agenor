package life_test

import (
	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok

	"github.com/snivilised/agenor/core"
	"github.com/snivilised/agenor/internal/enclave"
	"github.com/snivilised/agenor/internal/opts"
)

var _ = Describe("event", func() {
	var result enclave.KernelResult

	Context("end", func() {
		Context("single", func() {
			When("listener", func() {
				It("🧪 should: invoke client's handler", func() {
					invoked := false
					o, binder, _ := opts.Get()

					o.Events.End.On(func(_ core.TraverseResult) {
						invoked = true
					})
					binder.Controls.End.Dispatch()(&result)

					Expect(invoked).To(BeTrue())
				})
			})

			When("muted then unmuted", func() {
				It("🧪 should: invoke client's handler only when not muted", func() {
					invoked := false
					o, binder, _ := opts.Get()

					o.Events.End.On(func(_ core.TraverseResult) {
						invoked = true
					})
					binder.Controls.End.Mute()
					binder.Controls.End.Dispatch()(&result)
					Expect(invoked).To(BeFalse(), "notification not muted")

					invoked = false
					binder.Controls.End.Unmute()
					binder.Controls.End.Dispatch()(&result)
					Expect(invoked).To(BeTrue(), "notification not muted")
				})
			})
		})

		Context("multiple", func() {
			When("listener", func() {
				It("🧪 should: broadcast", func() {
					count := 0
					o, binder, _ := opts.Get()

					o.Events.End.On(func(_ core.TraverseResult) {
						count++
					})
					o.Events.End.On(func(_ core.TraverseResult) {
						count++
					})
					binder.Controls.End.Dispatch()(&result)
					Expect(count).To(Equal(2), "not all listeners were invoked for first notification")

					count = 0
					o.Events.End.On(func(_ core.TraverseResult) {
						count++
					})

					binder.Controls.End.Dispatch()(&result)
					Expect(count).To(Equal(3), "not all listeners were invoked for second notification")
				})
			})

			When("muted", func() {
				It("🧪 should: not broadcast", func() {
					count := 0
					o, binder, _ := opts.Get()

					o.Events.End.On(func(_ core.TraverseResult) {
						count++
					})
					o.Events.End.On(func(_ core.TraverseResult) {
						count++
					})

					binder.Controls.End.Mute()
					binder.Controls.End.Dispatch()(&result)

					Expect(count).To(Equal(0), "notification not muted")
				})
			})
		})

		Context("no listeners", func() {
			It("🧪 should: invoke no-op", func() {
				_, binder, _ := opts.Get()

				binder.Controls.End.Dispatch()(&result)
			})
		})
	})
})
