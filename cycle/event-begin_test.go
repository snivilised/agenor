package cycle_test

import (
	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok

	"github.com/snivilised/traverse/cycle"
	"github.com/snivilised/traverse/pref"
)

var _ = Describe("event", func() {
	Context("begin", func() {
		Context("single", func() {
			When("listener", func() {
				It("🧪 should: invoke client's handler", func() {
					invoked := false
					o, _ := pref.Get()

					o.Events.Begin.On(func(_ *cycle.BeginState) {
						invoked = true
					})
					o.Binder.Controls.Begin.Dispatch()(&cycle.BeginState{
						Root: traversalRoot,
					})

					Expect(invoked).To(BeTrue())
				})
			})

			When("muted then unmuted", func() {
				It("🧪 should: invoke client's handler only when not muted", func() {
					invoked := false
					o, _ := pref.Get()

					o.Events.Begin.On(func(_ *cycle.BeginState) {
						invoked = true
					})
					o.Binder.Controls.Begin.Mute()
					o.Binder.Controls.Begin.Dispatch()(&cycle.BeginState{
						Root: traversalRoot,
					})
					Expect(invoked).To(BeFalse(), "notification not muted")

					invoked = false
					o.Binder.Controls.Begin.Unmute()
					o.Binder.Controls.Begin.Dispatch()(&cycle.BeginState{
						Root: traversalRoot,
					})
					Expect(invoked).To(BeTrue(), "notification not muted")
				})
			})
		})

		Context("multiple", func() {
			When("listener", func() {
				It("🧪 should: broadcast", func() {
					count := 0
					o, _ := pref.Get()

					o.Events.Begin.On(func(_ *cycle.BeginState) {
						count++
					})
					o.Events.Begin.On(func(_ *cycle.BeginState) {
						count++
					})
					o.Binder.Controls.Begin.Dispatch()(&cycle.BeginState{
						Root: traversalRoot,
					})
					Expect(count).To(Equal(2), "not all listeners were invoked for first notification")

					count = 0
					o.Events.Begin.On(func(_ *cycle.BeginState) {
						count++
					})

					o.Binder.Controls.Begin.Dispatch()(&cycle.BeginState{
						Root: anotherRoot,
					})
					Expect(count).To(Equal(3), "not all listeners were invoked for second notification")
				})
			})

			When("muted", func() {
				It("🧪 should: not broadcast", func() {
					count := 0
					o, _ := pref.Get()

					o.Events.Begin.On(func(_ *cycle.BeginState) {
						count++
					})
					o.Events.Begin.On(func(_ *cycle.BeginState) {
						count++
					})

					o.Binder.Controls.Begin.Mute()
					o.Binder.Controls.Begin.Dispatch()(&cycle.BeginState{
						Root: anotherRoot,
					})

					Expect(count).To(Equal(0), "notification not muted")
				})
			})
		})

		Context("no listeners", func() {
			It("🧪 should: invoke no-op", func() {
				o, _ := pref.Get()

				o.Binder.Controls.Begin.Dispatch()(&cycle.BeginState{
					Root: traversalRoot,
				})
			})
		})
	})
})
