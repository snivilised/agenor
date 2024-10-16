package kernel

import (
	"context"
	"errors"
	"io/fs"
	"path/filepath"

	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/internal/services"
	"github.com/snivilised/traverse/internal/third/lo"
	"github.com/snivilised/traverse/internal/types"
	"github.com/snivilised/traverse/pref"
	"github.com/snivilised/traverse/tapable"
)

type readHooks struct {
	read tapable.Hook[core.ReadDirectoryHook, core.ChainReadDirectoryHook]
	sort tapable.Hook[core.SortHook, core.ChainSortHook]
}

type readOptions struct {
	hooks     readHooks
	behaviour *pref.SortBehaviour
}

type agentOptions struct {
	hooks   *tapable.Hooks
	defects *pref.DefectOptions
}

// navigatorAgent does work on behalf of the navigator. The agent performs
// generic tasks that apply to all navigators.
type navigatorAgent struct {
	ao        *agentOptions
	ro        *readOptions
	using     *pref.Using
	resources *types.Resources
	session   core.Session
}

func (n *navigatorAgent) Ignite(ignition *types.Ignition) {
	n.session = ignition.Session
}

func (n *navigatorAgent) top(ctx context.Context,
	ns *navigationStatic,
) (*types.KernelResult, error) {
	info, ie := n.ao.hooks.QueryStatus.Invoke()(
		ns.mediator.resources.FS.T, ns.tree,
	)

	err := lo.TernaryF(ie != nil,
		func() error {
			return n.ao.defects.Fault.Accept(&pref.NavigationFault{
				Err:  ie,
				Path: ns.tree,
				Info: info,
			})
		},
		func() error {
			_, te := ns.mediator.impl.Traverse(ctx, ns,
				core.Top(ns.tree, info),
			)

			return te
		},
	)

	return ns.mediator.impl.Result(ctx, err), err
}

func (n *navigatorAgent) Traverse(context.Context,
	*navigationStatic,
	*core.Node,
) (bool, error) {
	return continueTraversal, nil
}

// Result is the single point at which a Result is constructed. Due to
// the spawn resume strategy, a Result may occur more than once during
// a navigation session. The session knows when completion occurs. Any
// Result that occurs prior to completion are as a result of child
// navigation whose result should be combined in the final Result. This
// is all handled by the strategy.
func (n *navigatorAgent) Result(ctx context.Context,
	err error,
) *types.KernelResult {
	complete := n.session.IsComplete()
	result := types.NewResult(n.session,
		n.resources.Supervisor,
		err,
		complete,
	)

	if complete {
		_ = services.Broker.Emit(ctx, services.TopicNavigationComplete, result)
	}

	return result
}

const (
	continueTraversal = true
	skipTraversal     = false
)

// travel is the general recursive navigation function which returns a bool
// indicating whether we continue travelling or not in response to an
// error.
// true: success path; continue/progress
// false: skip (all, dir)
//
// When an error occurs for this node, we return false (skipTraversal) indicating
// a skip. A skip can mean skip the entire navigation process (fs.SkipAll),
// or just skip all remaining sibling nodes in this directory (fs.SkipDir).
func (n *navigatorAgent) travel(ctx context.Context,
	ns *navigationStatic,
	vapour inspection,
) (bool, error) {
	var (
		parent = vapour.Current()
	)

	for _, entry := range vapour.Entries() {
		path := filepath.Join(parent.Path, entry.Name())
		info, e := entry.Info()

		// TODO: check sampling; should happen transparently, by plugin

		current := core.New(
			path,
			entry,
			info,
			parent,
			e,
		)

		// TODO: ok for Travel to by-pass mediator?
		//
		if progress, err := ns.mediator.impl.Traverse(
			ctx, ns, current,
		); !progress {
			if err != nil {
				if errors.Is(err, fs.SkipDir) {
					// The returning of skipTraversal by the child, denotes
					// a skip. So when a child node returns a SkipDir error and
					// skipTraversal, what we're saying is that we want to skip
					// processing all successive siblings but continue traversal.
					// The !progress indicates we're skipping the remaining
					// processing of all of the parent item's remaining children.
					// (see the ✨ below ...)
					//
					return skipTraversal, err
				}

				return continueTraversal, err
			}
		} else if err != nil {
			// ✨ ... we skip processing all the remaining children for
			// this node, but still continue the overall traversal.
			//
			switch {
			case errors.Is(err, fs.SkipDir):
				continue
			case errors.Is(err, fs.SkipAll):
				break
			default:
				return continueTraversal, err
			}
		}
	}

	return continueTraversal, nil
}
