package collections

import "errors"

// ❌ Stack Is Empty (internal error)

// ErrStackIsEmpty indicates stack is empty (internal error)
var ErrStackIsEmpty = errors.New("internal: stack is empty")
