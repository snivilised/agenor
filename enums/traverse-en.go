package enums

// TraverseSubscription type to define traversal subscription (for which file system
// items the client defined callback are invoked for).
type TraverseSubscription uint

const (
	_                         TraverseSubscription = iota
	SubscribeAny                                   // invoke callback for files and folders
	SubscribeFolders                               // invoke callback for folders only
	SubscribeFoldersWithFiles                      // invoke callback for folders only but include files
	SubscribeFiles                                 // invoke callback for files only
)

// Role represents the role of an application entity (like a plugin role) The
// key element of a role is that there should be just a single entity that can take up
// the role which is bound to a service.
//
// For example, there can only be 1 logger. which means there can only be 1 entity
// that claims to provide this service, typically when the client invoke WithLogger
// option. This is similar to plugin architectures that allows plugins to register
// to provide a particular service.
//
// The mediator knows about Roles and manage registration requests
type Role uint

const (
	RoleUndefined Role = iota
	RoleLogger         // WithLogger
	RoleSampler        // WithSampler (need a specific sampler interface)
	RoleResume         // this is not an option; so might not be a valid role
)

// do we need to distinguish between internal and external entities. It looks
// like external entities are interested in traversal events, and internal
// entities are interested in initialisation events.
//
// --> external
// * logger
//
// --> internal
// * resume
