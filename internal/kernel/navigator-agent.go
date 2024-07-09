package kernel

import (
	"context"
	"errors"
	"io/fs"
	"path/filepath"

	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/internal/lo"
	"github.com/snivilised/traverse/internal/services"
	"github.com/snivilised/traverse/internal/types"
	"github.com/snivilised/traverse/pref"
	"github.com/snivilised/traverse/tapable"
)

type readHooks struct {
	read tapable.Hook[core.ReadDirectoryHook]
	sort tapable.Hook[core.SortHook]
}

type readOptions struct {
	hooks     readHooks
	behaviour *pref.SortBehaviour
}

type agentOptions struct {
	hooks   *tapable.Hooks
	defects *pref.DefectOptions
}

// navigatorAgent does work on behalf of the navigator. It is distinct
// from navigatorBase and should only be used when the limited polymorphism
// on base is inadequate. The agent functions performs generic tasks that
// apply to all navigators. agent is really an abstract concept that isn't
// represented by state, just functions that take state,
// typically navigationStatic.
type navigatorAgent struct {
	o         *agentOptions
	ro        *readOptions
	using     *pref.Using
	resources *types.Resources
	session   core.Session
}

/*
func (n *navigator) descend(navi *NavigationInfo) bool {
	if !navi.frame.periscope.descend(n.o.Store.Behaviours.Cascade.Depth) {
		return false
	}

	navi.frame.notifiers.descend.invoke(navi.Item)

	return true
}

func (n *navigator) ascend(navi *NavigationInfo, permit bool) {
	if permit {
		navi.frame.periscope.ascend()
		navi.frame.notifiers.ascend.invoke(navi.Item)
	}
}
*/

func (n *navigatorAgent) descend(*navigationInfo) bool {
	return true
}

func (n *navigatorAgent) ascend(navi *navigationInfo, permit bool) {
	_, _ = navi, permit
}

func (n *navigatorAgent) Ignite(ignition *types.Ignition) {
	n.session = ignition.Session
}

func (n *navigatorAgent) top(ctx context.Context,
	ns *navigationStatic,
) (*types.KernelResult, error) {
	info, ie := n.o.hooks.QueryStatus.Invoke()(ns.root)
	err := lo.TernaryF(ie != nil,
		func() error {
			return n.o.defects.Fault.Accept(&pref.NavigationFault{
				Err:  ie,
				Path: ns.root,
				Info: info,
			})
		},
		func() error {
			_, te := ns.mediator.impl.Travel(ctx, ns,
				core.Root(ns.root, info),
			)

			return te
		},
	)

	return ns.mediator.impl.Result(ctx, err), err
}

func (n *navigatorAgent) Travel(context.Context,
	*navigationStatic,
	*core.Node,
) (bool, error) {
	return continueTraversal, nil
}

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
		parent = vapour.current()
	)

	for _, entry := range vapour.entries() {
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
		if progress, err := ns.mediator.impl.Travel(
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
