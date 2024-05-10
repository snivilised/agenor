package traverse

import (
	"context"
	"time"

	"github.com/snivilised/traverse/pref"
)

type Session interface {
	StartedAt() time.Time
	Elapsed() time.Duration
}

type session struct {
	started             time.Time
	ctx                 context.Context
	cancel              context.CancelFunc
	accelerationOptions []pref.Option
}

func (s *session) StartedAt() time.Time {
	return s.started
}

func (s *session) Elapsed() time.Duration {
	return time.Since(s.started)
}
