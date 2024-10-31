package resume

import (
	"context"

	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/enums"
	"github.com/snivilised/traverse/internal/opts"
	"github.com/snivilised/traverse/internal/types"
	"github.com/snivilised/traverse/pref"
)

type spawnStrategy struct {
	baseStrategy
}

func (s *spawnStrategy) init(*opts.LoadInfo) error {
	return nil
}

// Next invokes this decorator which returns true if
// next link in the chain can be run or false to stop
// execution of subsequent links.
func (s *spawnStrategy) Next(_ core.Servant,
	_ types.Inspection,
) (bool, error) {
	return true, nil
}

// Role indicates the identity of the link
func (s *spawnStrategy) Role() enums.Role {
	// what role?
	return enums.RoleFastward
}

func (s *spawnStrategy) resume(_ context.Context,
	_ *pref.Was,
) (*types.KernelResult, error) {
	return nil, nil
}

func (s *spawnStrategy) finish() error {
	return nil
}

func (s *spawnStrategy) ifResult() bool {
	return true // TODO: tbd...
}
