package cycle_test

import (
	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok

	"github.com/snivilised/traverse/pref"
)

var _ = Describe("event", func() {
	var description string

	Context("wake", func() {
		Context("single", func() {
			When("listener", func() {
				It("🧪 should: invoke client's handler", func() {
					invoked := false
					o, _ := pref.Get()

					o.Events.Wake.On(func(_ string) {
						invoked = true
					})
					o.Binder.Controls.Wake.Dispatch()(description)

					Expect(invoked).To(BeTrue())
				})
			})

			When("muted then unmuted", func() {
				It("🧪 should: invoke client's handler only when not muted", func() {
					invoked := false
					o, _ := pref.Get()

					o.Events.Wake.On(func(_ string) {
						invoked = true
					})
					o.Binder.Controls.Wake.Mute()
					o.Binder.Controls.Wake.Dispatch()(description)
					Expect(invoked).To(BeFalse(), "notification not muted")

					invoked = false
					o.Binder.Controls.Wake.Unmute()
					o.Binder.Controls.Wake.Dispatch()(description)
					Expect(invoked).To(BeTrue(), "notification not muted")
				})
			})
		})

		Context("multiple", func() {
			When("listener", func() {
				It("🧪 should: broadcast", func() {
					count := 0
					o, _ := pref.Get()

					o.Events.Wake.On(func(_ string) {
						count++
					})
					o.Events.Wake.On(func(_ string) {
						count++
					})
					o.Binder.Controls.Wake.Dispatch()(description)
					Expect(count).To(Equal(2), "not all listeners were invoked for first notification")

					count = 0
					o.Events.Wake.On(func(_ string) {
						count++
					})

					o.Binder.Controls.Wake.Dispatch()(description)
					Expect(count).To(Equal(3), "not all listeners were invoked for second notification")
				})
			})

			When("muted", func() {
				It("🧪 should: not broadcast", func() {
					count := 0
					o, _ := pref.Get()

					o.Events.Wake.On(func(_ string) {
						count++
					})
					o.Events.Wake.On(func(_ string) {
						count++
					})

					o.Binder.Controls.Wake.Mute()
					o.Binder.Controls.Wake.Dispatch()(description)

					Expect(count).To(Equal(0), "notification not muted")
				})
			})
		})

		Context("no listeners", func() {
			It("🧪 should: invoke no-op", func() {
				o, _ := pref.Get()

				o.Binder.Controls.Wake.Dispatch()(description)
			})
		})
	})
})
