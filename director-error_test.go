package tv_test

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok

	tv "github.com/snivilised/traverse"
	"github.com/snivilised/traverse/internal/services"
)

type traverseErrorTE struct {
	given string
	using *tv.Using
	as    *tv.Was
}

var _ = Describe("director error", Ordered, func() {
	var handler tv.Client

	BeforeAll(func() {
		handler = func(_ *tv.Node) error {
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
			using: &tv.Using{
				Subscription: tv.SubscribeFiles,
				Handler:      handler,
			},
		}),

		Entry(nil, &traverseErrorTE{
			given: "using missing subscription",
			using: &tv.Using{
				Root:    "/root-traverse-path",
				Handler: handler,
			},
		}),

		Entry(nil, &traverseErrorTE{
			given: "using missing handler",
			using: &tv.Using{
				Root:         "/root-traverse-path",
				Subscription: tv.SubscribeFiles,
			},
		}),

		Entry(nil, &traverseErrorTE{
			given: "as missing restore from path",
			as: &tv.Was{
				Using: tv.Using{
					Root:         "/root-traverse-path",
					Subscription: tv.SubscribeFiles,
					Handler:      handler,
				},
				Strategy: tv.ResumeStrategySpawn,
			},
		}),

		Entry(nil, &traverseErrorTE{
			given: "as missing resume strategy",
			as: &tv.Was{
				Using: tv.Using{
					Root:         "/root-traverse-path",
					Subscription: tv.SubscribeFiles,
					Handler:      handler,
				},
				From: "/restore-from-path",
			},
		}),
	)
})
