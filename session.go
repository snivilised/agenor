package tv

import (
	"context"
	"time"

	"github.com/snivilised/traverse/internal/types"
)

type session struct {
	sync     synchroniser
	started  time.Time
	duration time.Duration
	plugins  []types.Plugin
}

func (s *session) start() {
	s.started = time.Now()
	s.sync.Ignite(&types.Ignition{
		Session: s,
	})
}

func (s *session) finish(result *types.KernelResult) {
	s.duration = time.Since(s.started)
	s.sync.Conclude(result)
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

func (s *session) exec(ctx context.Context) (*types.KernelResult, error) {
	return s.sync.Navigate(ctx)
}
