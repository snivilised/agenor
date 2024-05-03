package enums

//go:generate stringer -type=InternalRole -linecomment -trimprefix=InternalRole -output internal-role-en-auto.go

// InternalRole represents the role of an application entity (like a plugin role) The
// key element of a role is that there should be just a single entity that can take up
// the role which is bound to a service.
//
// For example, there can only be 1 logger. which means there can only be 1 entity
// that claims to provide this service, typically when the client invoke WithLogger
// option. This is similar to plugin architectures that allows plugins to register
// to provide a particular service.
//
// The mediator knows about Roles and manage registration requests
type InternalRole uint

const (
	InternalRoleRoleUndefined InternalRole = iota // undefined-role
	InternalRoleLogger                            // logger-role
	InternalRoleSampler                           // sampler-role
	InternalRoleResume                            // resume-role
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

/*
	InternalRoleRoleUndefined InternalRole = iota
	InternalRoleLogger                     // WithLogger
	InternalRoleSampler                    // WithSampler (need a specific sampler interface)
	InternalRoleResume                     // this is not an option; so might not be a valid role

*/
