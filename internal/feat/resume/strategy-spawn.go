package resume

import (
	"context"
	"io/fs"

	"github.com/snivilised/agenor/core"
	"github.com/snivilised/agenor/enums"
	"github.com/snivilised/agenor/internal/enclave"
	"github.com/snivilised/agenor/internal/kernel"
	"github.com/snivilised/agenor/internal/opts"
	"github.com/snivilised/agenor/internal/third/lo"
	nef "github.com/snivilised/nefilim"
)

type spawnStrategy struct {
	baseStrategy
	calc     nef.PathCalc
	complete bool
}

func (s *spawnStrategy) init(_ *opts.LoadInfo) error {
	s.calc = s.forest.T.Calc()

	return nil
}

func (s *spawnStrategy) ignite() {
}

func (s *spawnStrategy) resume(ctx context.Context) (result *enclave.KernelResult, err error) {
	// Bridge announces the availability of the ActiveState and acts as a conduit
	// between the previous navigation session and this resume session.
	s.mediator.Bridge(s.active)

	result, err = s.crown(ctx, &conclusion{
		active:    s.active,
		tree:      s.active.Tree,
		current:   s.active.CurrentPath,
		inclusive: true,
	})

	return result, err
}

func (s *spawnStrategy) ifResult() bool {
	return s.complete
}

// crown finishes off the navigation of a directory's contents as the
// Spawn resume has to deal with fractured ancestors; this is to say when a
// navigation is halted and saved, the state of the tree is that it is only
// partly navigated. The point at which it stops, marks a dividing line in
// that directory, where the parent's children have not all been visited;
// ie fractured ancestor and this fracturing occurs all the way up to the root
// of the tree. For the current node, crown identifies the following
// siblings and invokes seed for each one and this progresses on a recursive
// basis via crown.
func (s *spawnStrategy) crown(ctx context.Context,
	conc *conclusion,
) (*enclave.KernelResult, error) {
	if conc.current == conc.active.Tree {
		if s.complete {
			return nil, core.ErrDetectedSpawnStackOverflow
		}

		// reach the top, so we're done
		//
		s.complete = true
		return s.kc.Result(ctx), nil
	}

	parent, child := s.calc.Split(conc.current)
	following, err := s.following(parent,
		child,
		conc.inclusive,
	)
	if err != nil {
		return s.kc.Result(ctx), err
	}

	following.siblings.Sort(enums.EntryTypeFile)
	following.siblings.Sort(enums.EntryTypeDirectory)

	result, err := s.seed(ctx, parent, following.siblings.All())

	if err != nil {
		return result, err
	}

	conc.current = parent
	conc.inclusive = false

	return s.crown(ctx, conc)
}

// seed invokes a Spawn for each new siblings it is presented with by crown.
func (s *spawnStrategy) seed(ctx context.Context,
	parent string,
	entries []fs.DirEntry,
) (*enclave.KernelResult, error) {
	result := s.kc.Result(ctx)

	for _, entry := range entries {
		top := s.calc.Join(parent, entry.Name())
		intermediate, err := s.mediator.Spawn(ctx, top)

		if err != nil {
			return intermediate, err
		}
	}

	return result, nil
}

func (s *spawnStrategy) following(parent, anchor string,
	inclusive bool,
) (*shard, error) {
	entries, err := s.mediator.Read(parent)

	if err != nil {
		return nil, err
	}

	groups := lo.GroupBy(entries, func(entry fs.DirEntry) bool {
		if inclusive {
			return entry.Name() >= anchor
		}

		return entry.Name() > anchor
	})

	return &shard{
		siblings: kernel.NewContents(
			&s.o.Behaviours.Sort,
			s.o.Hooks.Sort,
			groups[followingSiblings],
		),
	}, nil
}
