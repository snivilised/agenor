package collections

import (
	"github.com/snivilised/traverse/internal/lo"
	"golang.org/x/exp/constraints"
)

type Orderable interface {
	constraints.Integer | string
}

type PositionalSet[T Orderable] struct {
	order     []T        // Defines the valid elements and their order
	items     map[T]bool // Tracks which items are in the set
	positions map[T]int  // Maps each item to its position in the order
	anchor    T
}

func NewPositionalSet[T Orderable](order []T, anchor T) *PositionalSet[T] {
	o := lo.Reject(lo.Uniq(order), func(
		item T, _ int,
	) bool {
		return item == anchor
	})
	o = append(o, anchor)

	ps := &PositionalSet[T]{
		order:     o,
		items:     make(map[T]bool),
		positions: make(map[T]int),
		anchor:    anchor,
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
		return true
	}

	return false
}

// Add insert multiple items into the set under the same conditions as
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

// Delete removes an item from the set
func (ps *PositionalSet[T]) Delete(item T) {
	if item == ps.anchor {
		return
	}

	ps.items[item] = false
	delete(ps.items, item)
}

// Contains checks if an item is in the set
func (ps *PositionalSet[T]) Contains(item T) bool {
	return ps.items[item]
}

// Items returns all items in the set, in the defined order
func (ps *PositionalSet[T]) Items() []T {
	result := make([]T, 0, len(ps.items))
	for _, item := range ps.order {
		if ps.items[item] {
			result = append(result, item)
		}
	}
	return result
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
