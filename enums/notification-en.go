package enums

import (
	"math"
)

//go:generate stringer -type=Notification -linecomment -trimprefix=Notification -output notification-en-auto.go

type Notification uint32

const (
	NotificationUndefined Notification     = 0               // undefined-notification
	NotificationBegin                      = 1 << (iota - 1) // begin-notification
	NotificationEnd                                          // end-notification
	NotificationDescend                                      // descend-notification
	NotificationAscend                                       // ascend-notification
	NotificationStart                                        // start-notification
	NotificationStop                                         // stop-notification
	NotificationAll       = math.MaxUint32                   // all-notification
)
