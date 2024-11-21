package age

import (
	"context"
	"time"

	"github.com/snivilised/agenor/internal/enclave"
)

type session struct {
	sync     synchroniser
	started  time.Time
	duration time.Duration
	plugins  []enclave.Plugin
}

func (s *session) start() {
	s.started = time.Now()
	s.sync.Ignite(&enclave.Ignition{
		Session: s,
	})
}

func (s *session) finish(result *enclave.KernelResult) {
	s.duration = time.Since(s.started)
	s.sync.Bye(result)
}

func (s *session) IsComplete() bool {
	return s.sync.IsComplete()
}

func (s *session) StartedAt() time.Time {
	return s.started
}

func (s *session) Elapsed() time.Duration {
	return time.Since(s.started)
}

func (s *session) exec(ctx context.Context) (*enclave.KernelResult, error) {
	return s.sync.Navigate(ctx)
}
