package enums

//go:generate stringer -type=Role -linecomment -trimprefix=Role -output role-en-auto.go

// Role represents the role of a subscription
type Role uint32

const (
	// RoleUndefined represents the undefined role
	RoleUndefined Role = iota // undefined-role

	// RoleAnchor represents the anchor role
	RoleAnchor // anchor-role

	// RoleClientFilter represents the client filter role
	RoleClientFilter // client-filter-role

	// RoleHibernate represents the hibernate role
	RoleHibernate // hibernate-role

	// RoleSampler represents the sampler role
	RoleSampler // sampler-role

	// RoleNanny represents the nanny role
	RoleNanny // nanny-role

	// RoleFastward represents the fastward role
	RoleFastward // fastward-role
)
