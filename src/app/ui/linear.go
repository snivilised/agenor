package ui

import (
	"fmt"
	"sync"

	"github.com/snivilised/jaywalk/src/agenor/core"
)

// linear is the default UI implementation. It writes plain text to stdout
// using fmt.Println. It is safe for concurrent use - all writes are
// serialised through a mutex so interleaved output from the run command's
// worker pool is avoided.
type linear struct {
	mu sync.Mutex
}

// OnNode prints the visited node's path to stdout.
func (l *linear) OnNode(node *core.Node) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	fmt.Printf("-> %v\n", node.Path)
	return nil
}

// Info writes a plain informational line to stdout.
func (l *linear) Info(msg string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	fmt.Println("i", msg)
}

// Warn writes a warning line to stdout.
func (l *linear) Warn(msg string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	fmt.Println("!", msg)
}

// Error writes an error line to stdout.
func (l *linear) Error(msg string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	fmt.Println("x", msg)
}
