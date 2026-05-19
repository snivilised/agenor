package controller

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/snivilised/pants"
)

type shellResult struct {
	output []byte
	err    error
}

type shellPoolExecutor struct {
	pool    *pants.ShellPool
	counter uint64
	once    sync.Once
	done    chan struct{}
	mux     sync.Mutex
	pending map[string]chan shellResult
}

func newShellPoolExecutor(pool *pants.ShellPool) *shellPoolExecutor {
	return &shellPoolExecutor{
		pool:    pool,
		done:    make(chan struct{}),
		pending: make(map[string]chan shellResult),
	}
}

func (e *shellPoolExecutor) Execute(
	ctx context.Context,
	command string,
) ([]byte, error) {
	id := strconv.FormatUint(atomic.AddUint64(&e.counter, 1), 36)
	marker := fmt.Sprintf("__JAYWALK_SHELL_STATUS_%s__:", id)
	resultCh := make(chan shellResult, 1)

	e.once.Do(func() {
		go e.observe()
	})

	e.mux.Lock()
	e.pending[marker] = resultCh
	e.mux.Unlock()

	if err := e.pool.Post(ctx, wrapShellCommand(command, marker)); err != nil {
		e.mux.Lock()
		delete(e.pending, marker)
		e.mux.Unlock()

		return nil, err
	}

	select {
	case result := <-resultCh:
		return result.output, result.err

	case <-ctx.Done():
		e.mux.Lock()
		delete(e.pending, marker)
		e.mux.Unlock()

		return nil, ctx.Err()

	case <-e.done:
		return nil, nil
	}
}

func (e *shellPoolExecutor) observe() {
	defer close(e.done)

	for output := range e.pool.Observe() {
		payload := output.Payload

		e.mux.Lock()
		for marker, resultCh := range e.pending {
			if result, ok := parseShellResult(payload, marker, output.Error); ok {
				delete(e.pending, marker)
				resultCh <- result
				close(resultCh)

				break
			}
		}
		e.mux.Unlock()
	}
}

func wrapShellCommand(command, marker string) string {
	return fmt.Sprintf("{\n%s\n} 2>&1\n__jaywalk_status=$?\nprintf '\\n%s%%d\\n' \"$__jaywalk_status\"",
		command,
		marker,
	)
}

func parseShellResult(
	payload, marker string,
	err error,
) (shellResult, bool) {
	index := strings.LastIndex(payload, marker)
	if index < 0 {
		return shellResult{}, false
	}

	body := strings.TrimSuffix(payload[:index], "\n")
	statusText := strings.TrimSpace(payload[index+len(marker):])
	if fields := strings.Fields(statusText); len(fields) > 0 {
		statusText = fields[0]
	}

	if status, convErr := strconv.Atoi(statusText); convErr == nil && status != 0 && err == nil {
		err = fmt.Errorf("shell command exited with status %d", status)
	}

	return shellResult{
		output: []byte(body),
		err:    err,
	}, true
}
