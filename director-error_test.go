package traverse_test

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok

	"github.com/snivilised/traverse"
	"github.com/snivilised/traverse/internal/services"
)

type traverseErrorTE struct {
	given string
	using *traverse.Using
	as    *traverse.As
}

var _ = Describe("director error", Ordered, func() {
	var handler traverse.Client

	BeforeAll(func() {
		handler = func(_ *traverse.Node) error {
			return nil
		}
	})

	BeforeEach(func() {
		services.Reset()
	})

	DescribeTable("Validate",
		func(entry *traverseErrorTE) {
			if entry.using != nil {
				Expect(entry.using.Validate()).NotTo(Succeed())

				return
			}

			if entry.as != nil {
				Expect(entry.as.Validate()).NotTo(Succeed())

				return
			}
		},
		func(entry *traverseErrorTE) string {
			return fmt.Sprintf("given: %v, 🧪 should fail", entry.given)
		},
		Entry(nil, &traverseErrorTE{
			given: "using missing root path",
			using: &traverse.Using{
				Subscription: traverse.SubscribeFiles,
				Handler:      handler,
			},
		}),

		Entry(nil, &traverseErrorTE{
			given: "using missing subscription",
			using: &traverse.Using{
				Root:    "/root-traverse-path",
				Handler: handler,
			},
		}),

		Entry(nil, &traverseErrorTE{
			given: "using missing handler",
			using: &traverse.Using{
				Root:         "/root-traverse-path",
				Subscription: traverse.SubscribeFiles,
			},
		}),

		Entry(nil, &traverseErrorTE{
			given: "as missing restore from path",
			as: &traverse.As{
				Using: traverse.Using{
					Root:         "/root-traverse-path",
					Subscription: traverse.SubscribeFiles,
					Handler:      handler,
				},
				Strategy: traverse.ResumeStrategySpawn,
			},
		}),

		Entry(nil, &traverseErrorTE{
			given: "as missing resume strategy",
			as: &traverse.As{
				Using: traverse.Using{
					Root:         "/root-traverse-path",
					Subscription: traverse.SubscribeFiles,
					Handler:      handler,
				},
				From: "/restore-from-path",
			},
		}),
	)
})
