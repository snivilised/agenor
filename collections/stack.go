package collections

// NewStack creates a new empty stack.
func NewStack[T any]() *Stack[T] {
	return &Stack[T]{
		content: []T{},
	}
}

// NewStackWith creates a new stack pre-populated with the provided items.
func NewStackWith[T any](with []T) *Stack[T] {
	return &Stack[T]{
		content: with,
	}
}

// Stack is a simple last-in-first-out (LIFO) data structure that allows you to
// push items onto the stack and pop them off in reverse order. It provides methods
// to check the current item, get the size of the stack, and determine if it's
// empty.
type Stack[T any] struct {
	content []T
}

// Push adds an item to the top of the stack.
func (s *Stack[T]) Push(item T) {
	s.content = append(s.content, item)
}

// Pop removes and returns the item at the top of the stack. If the stack is empty,
// it returns an error.
func (s *Stack[T]) Pop() (T, error) {
	if s.IsEmpty() {
		var zero T
		return zero, ErrStackIsEmpty
	}

	item := s.pop()

	return item, nil
}

// MustPop removes and returns the item at the top of the stack. If the stack is empty,
// it panics.
func (s *Stack[T]) MustPop() T {
	if s.IsEmpty() {
		panic(ErrStackIsEmpty)
	}

	return s.pop()
}

// Current returns the item at the top of the stack without removing it. If the stack is empty,
// it returns an error.
func (s *Stack[T]) Current() (T, error) {
	if s.IsEmpty() {
		var zero T
		return zero, ErrStackIsEmpty
	}

	return s.content[s.top()], nil
}

// Size returns the number of items currently in the stack.
func (s *Stack[T]) Size() uint {
	return uint(len(s.content))
}

// IsEmpty returns true if the stack has no items, and false otherwise.
func (s *Stack[T]) IsEmpty() bool {
	return len(s.content) == 0
}

// Content returns a slice containing all the items in the stack.
func (s *Stack[T]) Content() []T {
	return s.content
}

func (s *Stack[T]) top() int {
	return len(s.content) - 1
}

func (s *Stack[T]) pop() T {
	t := s.top()
	item := s.content[t]
	s.content = s.content[:t]

	return item
}
