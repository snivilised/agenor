package pref_test

import (
	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok

	"github.com/snivilised/traverse/pref"
)

var _ = Describe("Options", func() {
	Context("Init", func() {
		Context("RequestOptions", func() {
			Context("Notification", func() {
				When("client listens", func() {
					It("should: invoke client's handler", func() {
						begun := false
						reg := pref.NewRegistry()
						o := pref.RequestOptions(reg)

						o.Events.Begin.On(func(_ string) {
							begun = true
						})
						reg.Notification.Begin.Dispatch.Invoke("/traversal-root")

						Expect(begun).To(BeTrue())
					})
				})

				When("multiple listeners", func() {
					It("should: broadcast", func() {
						count := 0
						reg := pref.NewRegistry()
						o := pref.RequestOptions(reg)

						o.Events.Begin.On(func(_ string) {
							count++
						})
						o.Events.Begin.On(func(_ string) {
							count++
						})
						reg.Notification.Begin.Dispatch.Invoke("/traversal-root")
						Expect(count).To(Equal(2), "not all listeners were invoked for first notification")

						count = 0
						o.Events.Begin.On(func(_ string) {
							count++
						})

						reg.Notification.Begin.Dispatch.Invoke("/another-root")
						Expect(count).To(Equal(3), "not all listeners were invoked for second notification")
					})
				})

				When("no subscription", func() {
					It("should: ...", func() {
						reg := pref.NewRegistry()
						_ = pref.RequestOptions(reg)

						reg.Notification.Begin.Dispatch.Invoke("/traversal-root")
					})
				})
			})
		})
	})
})
