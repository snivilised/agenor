package types_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok
	"github.com/snivilised/traverse/enums"
	lab "github.com/snivilised/traverse/internal/laboratory"
	"github.com/snivilised/traverse/internal/measure"
	"github.com/snivilised/traverse/internal/types"
)

func TestEnclave(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Enclave Suite")
}

type resultTE struct {
	lab.NaviTE
	arrange func(trig *trigger)
	assert  func(a *asserter)
}

type session struct {
	started time.Time
}

func (s *session) start() {
	s.started = time.Now()
}

func (s *session) finish(_ *types.KernelResult) {

}

func (s *session) IsComplete() bool {
	return false
}

func (s *session) StartedAt() time.Time {
	return s.started
}

func (s *session) Elapsed() time.Duration {
	return time.Since(s.started)
}

func (s *session) exec(_ context.Context) (*types.KernelResult, error) {
	return &types.KernelResult{}, nil
}

type trigger struct {
	mums measure.MutableMetrics
}

func (t *trigger) times(m enums.Metric, n uint) *trigger {
	t.mums[m].Times(n)

	return t
}

type asserter struct {
	result *types.KernelResult
}

func (a *asserter) equals(m enums.Metric, n uint) *asserter {
	Expect(a.result.Metrics().Count(
		m,
	)).To(BeEquivalentTo(n), fmt.Sprintf("💥 metric: '%v'", m))

	return a
}
