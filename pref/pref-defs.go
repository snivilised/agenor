package pref

// 📦 pkg: pref - contains user option definitions; do not use anything
// in kernel (cyclic).

const (
	badge = "badge: option-requester"
)

type (
	RescueData interface {
		Data() interface{}
	}

	Recovery interface {
		Save(data RescueData) (string, error)
	}
)
