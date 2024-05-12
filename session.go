package traverse

import (
	"context"
	"time"

	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/pref"
)

// Session represents a traversal session and keeps tracks of
// timing.
type Session interface {
	StartedAt() time.Time
	Elapsed() time.Duration
}

type session struct {
	started  time.Time
	duration time.Duration
	ctx      context.Context
	cancel   context.CancelFunc
	nav      core.Navigator
	o        *pref.Options
}

func (s *session) start() {
	s.started = time.Now()
}

func (s *session) finish(_ core.TraverseResult) {
	// I wonder if the traverse result should become available
	// as a result of a message sent of the bus. Any component
	// needing access to the result should handle the message. This
	// way, we don't have to explicitly pass it around.
	//
	s.duration = time.Since(s.started)
}

func (s *session) StartedAt() time.Time {
	return s.started
}

func (s *session) Elapsed() time.Duration {
	return time.Since(s.started)
}

func (s *session) exec() (core.TraverseResult, error) {
	return s.nav.Navigate()
}
