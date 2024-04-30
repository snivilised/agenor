package enums

type WayPoint uint

// WayPoint used with hibernation. When active, the navigator starts of in
// sleep state. It wakes when a waypoint condition is matched and enables
// the client action t be invoke from that node.Event onwards

const (
	WayPointUndefined WayPoint = iota
	WayPointWake
	WayPointSleep
	WayPointToggle
)
