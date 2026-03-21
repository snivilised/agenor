package collections

import (
	"github.com/snivilised/jaywalk/src/internal/third/lo"
	"golang.org/x/exp/constraints"
)

// Orderable is a constraint that allows only types that can be ordered
// (e.g., integers, strings). This is used to ensure that the elements
// in the PositionalSet can be compared and ordered.
type Orderable interface {
	constraints.Integer | string
}

// PositionalSet is a set that maintains the order of its elements based on a
// predefined list. It allows insertion of elements in a specific order and
// provides methods to check for containment, retrieve items in order, and get
// the position of items.
type PositionalSet[T Orderable] struct {
	order       []T        // Defines the valid elements and their order
	items       map[T]bool // Tracks which items are in the set
	positions   map[T]int  // Maps each item to its position in the order
	anchor      T
	invalidated bool
	cache       []T
}

// NewPositionalSet creates a new PositionalSet with the given order and anchor.
// The order slice defines the valid elements and their order, while the anchor
// is a special element that is always present in the set and cannot be removed.
// The function ensures that the anchor is included in the order and initializes
// the internal structures of the PositionalSet.
func NewPositionalSet[T Orderable](order []T, anchor T) *PositionalSet[T] {
	o := lo.Reject(lo.Uniq(order), func(
		item T, _ int,
	) bool {
		return item == anchor
	})
	o = append(o, anchor)

	ps := &PositionalSet[T]{
		order:       o,
		items:       make(map[T]bool),
		positions:   make(map[T]int),
		anchor:      anchor,
		invalidated: true,
	}

	for i, item := range o {
		ps.positions[item] = i
	}

	ps.items[anchor] = true

	return ps
}

// Insert adds an item to the set if it's in the order and is not present
func (ps *PositionalSet[T]) Insert(item T) bool {
	if item == ps.anchor {
		return false
	}

	if _, exists := ps.positions[item]; exists {
		if _, found := ps.items[item]; found {
			return false
		}

		ps.items[item] = true
		ps.invalidated = true

		return true
	}

	return false
}

// All inserts multiple items into the set under the same conditions as
// Insert
func (ps *PositionalSet[T]) All(items ...T) bool {
	result := true

	for _, item := range items {
		inserted := ps.Insert(item)

		if result {
			result = inserted
		}
	}

	return result
}

// Delete removes an item from the set. Removing the anchor is prohibited.
func (ps *PositionalSet[T]) Delete(item T) {
	if item == ps.anchor {
		return
	}

	ps.items[item] = false
	ps.invalidated = true
	delete(ps.items, item)
}

// Contains checks if an item is in the set
func (ps *PositionalSet[T]) Contains(item T) bool {
	return ps.items[item]
}

// Top returns the first item in the set, in the defined order
func (ps *PositionalSet[T]) Top() T {
	return ps.Items()[0]
}

// Items returns all items in the set, in the defined order
func (ps *PositionalSet[T]) Items() []T {
	if !ps.invalidated {
		return ps.cache
	}

	ps.cache = make([]T, 0, len(ps.items))
	for _, item := range ps.order {
		if ps.items[item] {
			ps.cache = append(ps.cache, item)
		}
	}

	ps.invalidated = false

	return ps.cache
}

// Position returns the position of an item in the order
func (ps *PositionalSet[T]) Position(item T) (int, bool) {
	pos, exists := ps.positions[item]
	return pos, exists
}

// Count returns the elements in the set
func (ps *PositionalSet[T]) Count() int {
	return len(ps.items)
}
