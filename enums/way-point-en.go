package enums

//go:generate stringer -type=WayPoint -linecomment -trimprefix=WayPoint -output way-point-en-auto.go

// WayPoint used with hibernation. When active, the navigator starts of in
// sleep state. It wakes when a waypoint condition is matched and enables
// the client action to be invoked from that node.Event onwards.
type WayPoint uint

// WayPoint used with hibernation. When active, the navigator starts of in
// sleep state. It wakes when a waypoint condition is matched and enables
// the client action to be invoked from that node.Event onwards.

const (
	// WayPointUndefined represents the undefined way point
	WayPointUndefined WayPoint = iota // undefined-way-point

	// WayPointWake represents the wake way point
	WayPointWake // wake-from-hibernation

	// WayPointSleep represents the sleep way point
	WayPointSleep // sleep

	// WayPointToggle represents the toggle way point
	WayPointToggle // toggle-way-point
)
