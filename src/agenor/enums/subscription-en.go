package enums

//go:generate stringer -type=Subscription -linecomment -trimprefix=Subscribe -output subscription-en-auto.go

// Subscription type to define traversal subscription (for which file system
// items the client defined callback are invoked for).
type Subscription uint

const (
	// SubscribeUndefined represents the undefined subscription
	SubscribeUndefined Subscription = iota // undefined-subscription

	// SubscribeFiles represents the files subscription
	SubscribeFiles // subscribe-files

	// SubscribeDirectories represents the directories subscription
	SubscribeDirectories // subscribe-directories

	// SubscribeDirectoriesWithFiles represents the directories with files subscription
	SubscribeDirectoriesWithFiles // subscribe-directories-with-files

	// SubscribeUniversal represents the universal subscription
	SubscribeUniversal // subscribe-to-everything
)
