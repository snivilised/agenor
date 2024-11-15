package enclave_test

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	"github.com/snivilised/agenor/core"
	"github.com/snivilised/agenor/enums"
	"github.com/snivilised/agenor/internal/enclave"
	lab "github.com/snivilised/agenor/internal/laboratory"
)

var _ = Describe("KernelResult", func() {
	Context("Metrics", Ordered, func() {
		var (
			sess     *session
			reporter *core.Supervisor
			trig     *lab.Trigger
			err      error
			complete bool
		)

		BeforeEach(func() {
			sess = &session{}
			reporter = core.NewSupervisor()
			trig = &lab.Trigger{
				Metrics: reporter.Many(
					enums.MetricNoFilesInvoked,
					enums.MetricNoFilesFilteredOut,
					enums.MetricNoDirectoriesInvoked,
					enums.MetricNoDirectoriesFilteredOut,
				),
			}
			complete = false
		})

		DescribeTable("Times",
			func(entry *resultTE) {
				entry.arrange(trig)
				result := enclave.NewResult(sess,
					reporter,
					err,
					complete,
				)
				entry.assert(&asserter{
					result: result,
				})
			},
			func(entry *resultTE) string {
				return fmt.Sprintf("🧪 ===> given: '%v', should: '%v'", entry.Given, entry.Should)
			},

			Entry(nil, &resultTE{
				NaviTE: lab.NaviTE{
					Given:  "metrics populated",
					Should: "count metrics",
				},
				arrange: func(trig *lab.Trigger) {
					trig.Times(
						enums.MetricNoFilesInvoked, 10).Times(
						enums.MetricNoFilesFilteredOut, 20).Times(
						enums.MetricNoDirectoriesInvoked, 30).Times(
						enums.MetricNoDirectoriesFilteredOut, 40,
					)
				},
				assert: func(a *asserter) {
					a.equals(
						enums.MetricNoFilesInvoked, 10).equals(
						enums.MetricNoFilesFilteredOut, 20).equals(
						enums.MetricNoDirectoriesInvoked, 30).equals(
						enums.MetricNoDirectoriesFilteredOut, 40,
					)
				},
			}),
		)
	})
})
