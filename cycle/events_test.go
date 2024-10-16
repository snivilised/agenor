package cycle_test

import (
	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok

	"github.com/snivilised/li18ngo"
	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/cycle"
)

var _ = Describe("controls", Ordered, func() {
	BeforeAll(func() {
		Expect(li18ngo.Use()).To(Succeed())
	})

	When("bind", func() {
		It("🧪 should: dispatch notification to event handler", func() {
			const path = "/traversal-tree"

			var (
				controls cycle.Controls
				events   cycle.Events
				begun    bool
				ended    bool
			)

			// init(binder->options):
			//
			events.Bind(&controls)

			// client:
			//
			events.Begin.On(func(state *cycle.BeginState) {
				begun = true
				Expect(state.Tree).To(Equal(path))
			})

			events.End.On(func(_ core.TraverseResult) {
				ended = true
			})

			// component side:
			//
			controls.Begin.Dispatch()(&cycle.BeginState{
				Tree: path,
			})
			controls.End.Dispatch()(nil)

			Expect(begun).To(BeTrue(), "begin notification handler not invoked")
			Expect(ended).To(BeTrue(), "end notification handler not invoked")
		})
	})
})
