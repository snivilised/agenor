package life_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/snivilised/jaywalk/src/agenor/internal/opts"
)

var _ = Describe("event", func() {
	var description string

	Context("wake", func() {
		Context("single", func() {
			When("listener", func() {
				It("🧪 should: invoke client's handler", func() {
					invoked := false
					o, binder, _ := opts.Get()

					o.Events.Wake.On(func(_ string) {
						invoked = true
					})
					binder.Controls.Wake.Dispatch()(description)

					Expect(invoked).To(BeTrue())
				})
			})

			When("muted then unmuted", func() {
				It("🧪 should: invoke client's handler only when not muted", func() {
					invoked := false
					o, binder, _ := opts.Get()

					o.Events.Wake.On(func(_ string) {
						invoked = true
					})
					binder.Controls.Wake.Mute()
					binder.Controls.Wake.Dispatch()(description)
					Expect(invoked).To(BeFalse(), "notification not muted")

					invoked = false

					binder.Controls.Wake.Unmute()
					binder.Controls.Wake.Dispatch()(description)
					Expect(invoked).To(BeTrue(), "notification not muted")
				})
			})
		})

		Context("multiple", func() {
			When("listener", func() {
				It("🧪 should: broadcast", func() {
					count := 0
					o, binder, _ := opts.Get()

					o.Events.Wake.On(func(_ string) {
						count++
					})
					o.Events.Wake.On(func(_ string) {
						count++
					})
					binder.Controls.Wake.Dispatch()(description)
					Expect(count).To(Equal(2), "not all listeners were invoked for first notification")

					count = 0

					o.Events.Wake.On(func(_ string) {
						count++
					})

					binder.Controls.Wake.Dispatch()(description)
					Expect(count).To(Equal(3), "not all listeners were invoked for second notification")
				})
			})

			When("muted", func() {
				It("🧪 should: not broadcast", func() {
					count := 0
					o, binder, _ := opts.Get()

					o.Events.Wake.On(func(_ string) {
						count++
					})
					o.Events.Wake.On(func(_ string) {
						count++
					})

					binder.Controls.Wake.Mute()
					binder.Controls.Wake.Dispatch()(description)

					Expect(count).To(Equal(0), "notification not muted")
				})
			})
		})

		Context("no listeners", func() {
			It("🧪 should: invoke no-op", func() {
				_, binder, _ := opts.Get()

				binder.Controls.Wake.Dispatch()(description)
			})
		})
	})
})
