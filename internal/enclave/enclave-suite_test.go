package enclave_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok
	"github.com/snivilised/traverse/enums"
	"github.com/snivilised/traverse/internal/enclave"
	lab "github.com/snivilised/traverse/internal/laboratory"
)

func TestEnclave(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Enclave Suite")
}

type resultTE struct {
	lab.NaviTE
	arrange func(trig *lab.Trigger)
	assert  func(a *asserter)
}

type session struct {
	started time.Time
}

func (s *session) start() {
	s.started = time.Now()
}

func (s *session) finish(_ *enclave.KernelResult) {

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

func (s *session) exec(_ context.Context) (*enclave.KernelResult, error) {
	return &enclave.KernelResult{}, nil
}

type asserter struct {
	result *enclave.KernelResult
}

func (a *asserter) equals(m enums.Metric, n uint) *asserter {
	Expect(a.result.Metrics().Count(
		m,
	)).To(BeEquivalentTo(n), fmt.Sprintf("💥 metric: '%v'", m))

	return a
}
