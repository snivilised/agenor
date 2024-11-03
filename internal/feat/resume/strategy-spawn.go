package resume

import (
	"context"
	"errors"
	"fmt"
	"io/fs"

	"github.com/snivilised/traverse/core"
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

func (s *spawnStrategy) resume(ctx context.Context,
	_ *pref.Was,
) (*types.KernelResult, error) {
	// what is the equivalent of this?:
	// s.nc.frame.root.Set(info.ps.Active.Root)
	//
	return s.conclude(ctx, &concludeInfo{
		active:    s.active,
		tree:      s.active.Tree,
		current:   s.active.CurrentPath,
		inclusive: true,
	})
}

func (s *spawnStrategy) finish() error {
	return nil
}

func (s *spawnStrategy) ifResult() bool {
	return true // TODO: tbd...
}

func (s *spawnStrategy) conclude(ctx context.Context,
	conclusion *concludeInfo,
) (*types.KernelResult, error) {
	if conclusion.current == conclusion.active.Tree {
		// reach the top, so we're done
		//
		err := errors.New("fake spawn result")
		return s.kc.Result(ctx, err), err
	}

	// TODO: impl pending:

	return &types.KernelResult{
		// spawn not complete yet
	}, nil
}

type seedParams struct {
	parent     string
	entries    []fs.DirEntry
	conclusion *concludeInfo
}

func (s *spawnStrategy) seed(ctx context.Context,
	parent string,
	entries []fs.DirEntry,
	conclusion *concludeInfo,
) (*types.KernelResult, error) {
	s.mediator.Connect(conclusion.tree, conclusion.current)

	compoundResult := &types.KernelResult{}

	for _, entry := range entries {
		topPath := s.forest.T.Calc().Join(parent, entry.Name())

		result, err := s.mediator.Spawn(ctx, &core.ActiveState{
			Tree: topPath,
			// other members tbd
		})

		// We can do away with the concept of merging results

		_, _ = compoundResult.Merge(result)

		if err != nil {
			return compoundResult, err
		}
	}

	// TODO: check wether its ok to do this weird thing with the embedded error...
	//
	return compoundResult, compoundResult.Error() //nolint:gocritic // fuck off
}

func (s *spawnStrategy) following(parent, anchor string,
	filesFirst, inclusive bool,
) *shard {
	// nil here is the fS(fs.ReadDirFS)
	_, _ = s.o.Hooks.ReadDirectory.Invoke()(nil, parent)

	// would it be better to just call a method on the mediator to
	// read directory contents? That way, we can be abstracted away
	// from the details of reading a directory, like having to
	// need access to the fS ...
	//
	// s.mediator.Read(path)
	//

	_ = shard{
		siblings: nil, // TODO: create new contents
	}

	panic(
		fmt.Sprintf("NOT-IMPL: spawnStrategy.following, %v, %v, %v, %v",
			parent, anchor, filesFirst, inclusive),
	)
}
