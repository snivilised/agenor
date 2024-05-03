package enums

//go:generate stringer -type=WayPoint -linecomment -trimprefix=WayPoint -output way-point-en-auto.go

type WayPoint uint

// WayPoint used with hibernation. When active, the navigator starts of in
// sleep state. It wakes when a waypoint condition is matched and enables
// the client action to be invoked from that node.Event onwards.

const (
	WayPointUndefined WayPoint = iota // undefined-way-point
	WayPointWake                      // wake-from-hibernation
	WayPointSleep                     // sleep
	WayPointToggle                    // toggle-way-point
)
