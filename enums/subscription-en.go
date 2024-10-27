package enums

//go:generate stringer -type=Subscription -linecomment -trimprefix=Subscribe -output subscription-en-auto.go

// Subscription type to define traversal subscription (for which file system
// items the client defined callback are invoked for).
type Subscription uint

const (
	SubscribeUndefined            Subscription = iota
	SubscribeFiles                             // subscribe-files
	SubscribeDirectories                       // subscribe-directories
	SubscribeDirectoriesWithFiles              // subscribe-directories-with-files
	SubscribeUniversal                         // subscribe-to-everything
)
