package ui

import (
	"fmt"
	"sync"

	"github.com/snivilised/jaywalk/src/app/report"
	"github.com/snivilised/jaywalk/src/locale"
	"github.com/snivilised/li18ngo"
)

// linear is the default UI implementation. It writes plain text to stdout
// using fmt.Println. It is safe for concurrent use - all writes are
// serialised through a mutex so interleaved output from the run command's
// worker pool is avoided.
type linear struct {
	mu sync.Mutex
}

// OnNodeEvent prints the visited node's path to stdout.
func (l *linear) OnNodeEvent(e *report.NeutralEvent) {
	l.mu.Lock()
	defer l.mu.Unlock()

	fmt.Println(li18ngo.Text(locale.NodeVisitedTemplData{
		Path: e.Node.Path,
	}))
}

// OnActionEvent prints the action name and node path to stdout.
func (l *linear) OnActionEvent(e *report.ActionEvent) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if e.Err != nil {
		fmt.Println(li18ngo.Text(locale.ActionFailedTemplData{
			Name: e.Name,
			Path: e.Node.Path,
			Err:  e.Err.Error(),
		}))
		return
	}

	fmt.Println(li18ngo.Text(locale.ActionVisitedTemplData{
		Name: e.Name,
		Path: e.Node.Path,
	}))
}

// OnPipelineEvent prints the pipeline name and node path to stdout.
func (l *linear) OnPipelineEvent(e *report.PipelineEvent) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if e.Err != nil {
		fmt.Println(li18ngo.Text(locale.PipelineFailedTemplData{
			Name: e.Name,
			Path: e.Node.Path,
			Err:  e.Err.Error(),
		}))
		return
	}

	fmt.Println(li18ngo.Text(locale.PipelineVisitedTemplData{
		Name: e.Name,
		Path: e.Node.Path,
	}))
}

// OnComplete renders the traversal outcome as plain text.
func (l *linear) OnComplete(t *report.Traversal) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if t.Err != nil {
		fmt.Println(li18ngo.Text(locale.TraversalFailedTemplData{
			Err: t.Err.Error(),
		}))
		return
	}

	fmt.Println(li18ngo.Text(locale.TraversalCompleteTemplData{
		Files:   t.FilesVisited,
		Dirs:    t.DirsVisited,
		Elapsed: t.Elapsed.Round(1e6).String(),
	}))
}
