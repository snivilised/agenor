package enums

import (
	"math"
)

//go:generate stringer -type=Notification -linecomment -trimprefix=Notification -output notification-en-auto.go

// Notification represents the type of notification to send
type Notification uint32

const (
	// NotificationUndefined represents the undefined notification
	NotificationUndefined Notification = 0 // undefined-notification

	// NotificationBegin represents the begin notification
	NotificationBegin = 1 << (iota - 1) // begin-notification

	// NotificationEnd represents the end notification
	NotificationEnd // end-notification

	// NotificationDescend represents the descend notification
	NotificationDescend // descend-notification

	// NotificationAscend represents the ascend notification
	NotificationAscend // ascend-notification

	// NotificationWake represents the wake notification
	NotificationWake // wake-notification

	// NotificationSleep represents the sleep notification
	NotificationSleep // sleep-notification

	// NotificationAll represents all notifications
	NotificationAll = math.MaxUint32 // all-notification
)
