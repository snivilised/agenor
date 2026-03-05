package enums

//go:generate stringer -type=TriStateBool -linecomment -trimprefix=TriStateBool -output tri-state-bool-en-auto.go

// TriStateBool boolean value with an undefined state
type TriStateBool uint

const (
	// TriStateBoolUndefined represent bool value noty set
	TriStateBoolUndefined TriStateBool = iota // undefined-bool

	// TriStateBoolTrue boolean true vlaue
	TriStateBoolTrue // true

	// TriStateBoolFalse boolean false value
	TriStateBoolFalse // false
)
