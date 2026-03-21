package ui

import (
	"fmt"
	"sync"

	"github.com/snivilised/jaywalk/src/app/report"
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
	fmt.Printf("-> %v\n", e.Node.Path)
}

// OnActionEvent prints the action name and node path to stdout.
func (l *linear) OnActionEvent(e *report.ActionEvent) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if e.Err != nil {
		fmt.Printf("x action %q failed on %v: %v\n", e.Name, e.Node.Path, e.Err)
		return
	}

	fmt.Printf("a [%v] %v\n", e.Name, e.Node.Path)
}

// OnPipelineEvent prints the pipeline name and node path to stdout.
func (l *linear) OnPipelineEvent(e *report.PipelineEvent) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if e.Err != nil {
		fmt.Printf("x pipeline %q failed on %v: %v\n", e.Name, e.Node.Path, e.Err)
		return
	}

	fmt.Printf("p [%v] %v\n", e.Name, e.Node.Path)
}

// OnComplete renders the traversal outcome as plain text.
func (l *linear) OnComplete(t *report.Traversal) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if t.Err != nil {
		fmt.Printf("x traversal failed: %v\n", t.Err)
		return
	}

	fmt.Printf(
		"i complete: %d files, %d dirs visited in %s\n",
		t.FilesVisited,
		t.DirsVisited,
		t.Elapsed.Round(1e6),
	)
}
