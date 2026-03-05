package enclave_test

import (
	"fmt"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/snivilised/agenor/enums"
	"github.com/snivilised/agenor/internal/enclave"
	lab "github.com/snivilised/agenor/internal/laboratory"
)

func TestEnclave(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Enclave Suite")
}

type resultTE struct {
	lab.DescribedTE
	arrange func(trig *lab.Trigger)
	assert  func(a *asserter)
}

func FormatResultTestDescription(entry *resultTE) string {
	return fmt.Sprintf("Given: %v 🧪 should: %v", entry.Given, entry.Should)
}

type session struct {
	started time.Time
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

type asserter struct {
	result *enclave.KernelResult
}

func (a *asserter) equals(m enums.Metric, n uint) *asserter {
	Expect(a.result.Metrics().Count(
		m,
	)).To(BeEquivalentTo(n), fmt.Sprintf("💥 metric: '%v'", m))

	return a
}
