package resume

import (
	"context"
	"fmt"
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

func (s *spawnStrategy) init(load *opts.LoadInfo) error {
	s.calc = s.forest.T.Calc()

	fmt.Printf("===> 🍭 RESTORED '%v' directories, '%v' files.\n",
		load.State.Metrics[enums.MetricNoDirectoriesInvoked].Counter,
		load.State.Metrics[enums.MetricNoFilesInvoked].Counter,
	)

	return nil
}

func (s *spawnStrategy) resume(ctx context.Context) (result *enclave.KernelResult, err error) {
	fmt.Printf("\t💙 resume, tree: '%v', current path: '%v'\n",
		s.active.Tree, s.active.CurrentPath)

	result, err = s.conclude(ctx, &conclusion{
		active:    s.active,
		tree:      s.active.Tree,
		current:   s.active.CurrentPath,
		inclusive: true,
	})

	fmt.Printf("===> 🍭 invoked '%v' directories, '%v' files.\n",
		result.Metrics().Count(enums.MetricNoDirectoriesInvoked),
		result.Metrics().Count(enums.MetricNoFilesInvoked),
	)

	return result, err
}

func (s *spawnStrategy) ifResult() bool {
	return s.complete
}

func (s *spawnStrategy) conclude(ctx context.Context,
	conc *conclusion,
) (*enclave.KernelResult, error) {
	if conc.current == conc.active.Tree {
		fmt.Printf("\t💎 conclude(COMPLETE), current: '%v'\n",
			conc.current)

		if s.complete {
			return nil, core.ErrDetectedSpawnStackOverflow
		}

		// reach the top, so we're done
		//
		s.complete = true
		return s.kc.Result(ctx), nil
	}

	fmt.Printf("\t🟢 conclude, current: '%v'\n",
		conc.current)

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

	result, err := s.seed(ctx, parent, following.siblings.All(), conc)

	if err != nil {
		return result, err
	}

	conc.current = parent
	conc.inclusive = false

	return s.conclude(ctx, conc)
}

func (s *spawnStrategy) seed(ctx context.Context,
	parent string,
	entries []fs.DirEntry,
	conc *conclusion,
) (*enclave.KernelResult, error) {
	fmt.Printf("\t🔊 seed, current: '%v'\n", conc.current)
	s.mediator.Bridge(conc.tree, conc.current)

	result := s.kc.Result(ctx)

	for _, entry := range entries {
		top := s.calc.Join(parent, entry.Name())

		intermediate, err := s.mediator.Spawn(ctx, &core.ActiveState{
			Tree: top,
			TraverseDescription: core.FsDescription{
				IsRelative: s.forest.T.IsRelative(),
			},
			ResumeDescription: core.FsDescription{
				IsRelative: s.forest.R.IsRelative(),
			},
			// Subscription: tbd,
		})

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
