package resume

import (
	"context"
	"fmt"
	"io/fs"

	nef "github.com/snivilised/nefilim"
	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/enums"
	"github.com/snivilised/traverse/internal/enclave"
	"github.com/snivilised/traverse/internal/kernel"
	"github.com/snivilised/traverse/internal/opts"
	"github.com/snivilised/traverse/internal/third/lo"
	"github.com/snivilised/traverse/pref"
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

func (s *spawnStrategy) resume(ctx context.Context,
	_ *pref.Was,
) (result *enclave.KernelResult, err error) {
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
		return s.kc.Result(ctx, nil), nil
	}

	fmt.Printf("\t🟢 conclude, current: '%v'\n",
		conc.current)

	parent, child := s.calc.Split(conc.current)
	following, err := s.following(parent,
		child,
		conc.inclusive,
	)
	if err != nil {
		return s.kc.Result(ctx, err), err
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

	result := s.kc.Result(ctx, nil)

	for _, entry := range entries {
		top := s.calc.Join(parent, entry.Name())

		intermediate, err := s.mediator.Spawn(ctx, &core.ActiveState{
			Tree: top,
		})

		if err != nil {
			return intermediate, err
		}
	}

	return result, result.Error() //nolint:gocritic // baa
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
