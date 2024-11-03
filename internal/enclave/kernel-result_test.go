package enclave_test

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	"github.com/snivilised/traverse/enums"
	"github.com/snivilised/traverse/internal/enclave"
	lab "github.com/snivilised/traverse/internal/laboratory"
	"github.com/snivilised/traverse/internal/measure"
)

var _ = Describe("KernelResult", func() {
	Context("Metrics", Ordered, func() {
		var (
			sess     *session
			reporter *measure.Supervisor
			trig     *trigger
			err      error
			complete bool
		)

		BeforeEach(func() {
			sess = &session{}
			reporter = measure.New()
			trig = &trigger{
				mums: reporter.Many(
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
				arrange: func(trig *trigger) {
					trig.times(
						enums.MetricNoFilesInvoked, 10).times(
						enums.MetricNoFilesFilteredOut, 20).times(
						enums.MetricNoDirectoriesInvoked, 30).times(
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
