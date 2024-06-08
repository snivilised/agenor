package tv

import (
	"context"
	"time"

	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/internal/types"
)

// Session represents a traversal session and keeps tracks of
// timing.
type Session interface {
	StartedAt() time.Time
	Elapsed() time.Duration
}

type session struct {
	sync     synchroniser
	started  time.Time
	duration time.Duration
	plugins  []types.Plugin
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

func (s *session) exec(ctx context.Context) (core.TraverseResult, error) {
	return s.sync.Navigate(ctx)
}
