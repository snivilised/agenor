package enums

//go:generate stringer -type=Subscription -linecomment -trimprefix=Subscribe -output subscription-en-auto.go

// Subscription type to define traversal subscription (for which file system
// items the client defined callback are invoked for).
type Subscription uint

const (
	SubscribeUndefined        Subscription = iota
	SubscribeFiles                         // subscribe-files
	SubscribeFolders                       // subscribe-folders
	SubscribeFoldersWithFiles              // subscribe-folders-with-files
	SubscribeUniversal                     // subscribe-to-everything
)
