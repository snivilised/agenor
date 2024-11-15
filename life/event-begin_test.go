package life_test

import (
	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok

	"github.com/snivilised/agenor/internal/opts"
	"github.com/snivilised/agenor/life"
)

var _ = Describe("event", func() {
	Context("begin", func() {
		Context("single", func() {
			When("listener", func() {
				It("🧪 should: invoke client's handler", func() {
					invoked := false
					o, binder, _ := opts.Get()

					o.Events.Begin.On(func(_ *life.BeginState) {
						invoked = true
					})
					binder.Controls.Begin.Dispatch()(&life.BeginState{
						Tree: traversalRoot,
					})

					Expect(invoked).To(BeTrue())
				})
			})

			When("muted then unmuted", func() {
				It("🧪 should: invoke client's handler only when not muted", func() {
					invoked := false
					o, binder, _ := opts.Get()

					o.Events.Begin.On(func(_ *life.BeginState) {
						invoked = true
					})
					binder.Controls.Begin.Mute()
					binder.Controls.Begin.Dispatch()(&life.BeginState{
						Tree: traversalRoot,
					})
					Expect(invoked).To(BeFalse(), "notification not muted")

					invoked = false
					binder.Controls.Begin.Unmute()
					binder.Controls.Begin.Dispatch()(&life.BeginState{
						Tree: traversalRoot,
					})
					Expect(invoked).To(BeTrue(), "notification not muted")
				})
			})
		})

		Context("multiple", func() {
			When("listener", func() {
				It("🧪 should: broadcast", func() {
					count := 0
					o, binder, _ := opts.Get()

					o.Events.Begin.On(func(_ *life.BeginState) {
						count++
					})
					o.Events.Begin.On(func(_ *life.BeginState) {
						count++
					})
					binder.Controls.Begin.Dispatch()(&life.BeginState{
						Tree: traversalRoot,
					})
					Expect(count).To(Equal(2), "not all listeners were invoked for first notification")

					count = 0
					o.Events.Begin.On(func(_ *life.BeginState) {
						count++
					})

					binder.Controls.Begin.Dispatch()(&life.BeginState{
						Tree: anotherRoot,
					})
					Expect(count).To(Equal(3), "not all listeners were invoked for second notification")
				})
			})

			When("muted", func() {
				It("🧪 should: not broadcast", func() {
					count := 0
					o, binder, _ := opts.Get()

					o.Events.Begin.On(func(_ *life.BeginState) {
						count++
					})
					o.Events.Begin.On(func(_ *life.BeginState) {
						count++
					})

					binder.Controls.Begin.Mute()
					binder.Controls.Begin.Dispatch()(&life.BeginState{
						Tree: anotherRoot,
					})

					Expect(count).To(Equal(0), "notification not muted")
				})
			})
		})

		Context("no listeners", func() {
			It("🧪 should: invoke no-op", func() {
				_, binder, _ := opts.Get()

				binder.Controls.Begin.Dispatch()(&life.BeginState{
					Tree: traversalRoot,
				})
			})
		})
	})
})
