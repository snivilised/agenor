package enums

//go:generate stringer -type=TriStateBool -linecomment -trimprefix=TriStateBool -output tri-state-bool-en-auto.go

type TriStateBool uint

const (
	TriStateBoolUndefined TriStateBool = iota // undefined-bool
	TriStateBoolTrue                          // true
	TriStateBoolFalse                         // false
)
