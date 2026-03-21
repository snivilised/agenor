package life_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/snivilised/jaywalk/src/agenor/core"
	"github.com/snivilised/jaywalk/src/agenor/internal/opts"
)

var _ = Describe("event", func() {
	var node core.Node

	Context("descend", func() {
		Context("single", func() {
			When("listener", func() {
				It("🧪 should: invoke client's handler", func() {
					invoked := false
					o, binder, _ := opts.Get()

					o.Events.Descend.On(func(_ *core.Node) {
						invoked = true
					})
					binder.Controls.Descend.Dispatch()(&node)

					Expect(invoked).To(BeTrue())
				})
			})

			When("muted then unmuted", func() {
				It("🧪 should: invoke client's handler only when not muted", func() {
					invoked := false
					o, binder, _ := opts.Get()

					o.Events.Descend.On(func(_ *core.Node) {
						invoked = true
					})
					binder.Controls.Descend.Mute()
					binder.Controls.Descend.Dispatch()(&node)
					Expect(invoked).To(BeFalse(), "notification not muted")

					invoked = false

					binder.Controls.Descend.Unmute()
					binder.Controls.Descend.Dispatch()(&node)
					Expect(invoked).To(BeTrue(), "notification not muted")
				})
			})
		})

		Context("multiple", func() {
			When("listener", func() {
				It("🧪 should: broadcast", func() {
					count := 0
					o, binder, _ := opts.Get()

					o.Events.Descend.On(func(_ *core.Node) {
						count++
					})
					o.Events.Descend.On(func(_ *core.Node) {
						count++
					})
					binder.Controls.Descend.Dispatch()(&node)
					Expect(count).To(Equal(2), "not all listeners were invoked for first notification")

					count = 0

					o.Events.Descend.On(func(_ *core.Node) {
						count++
					})

					binder.Controls.Descend.Dispatch()(&node)
					Expect(count).To(Equal(3), "not all listeners were invoked for second notification")
				})
			})

			When("muted", func() {
				It("🧪 should: not broadcast", func() {
					count := 0
					o, binder, _ := opts.Get()

					o.Events.Descend.On(func(_ *core.Node) {
						count++
					})
					o.Events.Descend.On(func(_ *core.Node) {
						count++
					})

					binder.Controls.Descend.Mute()
					binder.Controls.Descend.Dispatch()(&node)

					Expect(count).To(Equal(0), "notification not muted")
				})
			})
		})

		Context("no listeners", func() {
			It("🧪 should: invoke no-op", func() {
				_, binder, _ := opts.Get()

				binder.Controls.Descend.Dispatch()(&node)
			})
		})
	})
})
