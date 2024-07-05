package enums

//go:generate stringer -type=Role -linecomment -trimprefix=Role -output role-en-auto.go

type Role uint32

const (
	RoleUndefined        Role = iota // undefined-role
	RoleAnchor                       // anchor-role
	RoleClientFilter                 // client-filter-role
	RoleClientHiberWake              // client-hiber-wake-role
	RoleClientHiberSleep             // client-hiber-sleep-role
	RoleFastward                     // fastward-role
	RoleSampler                      // sampler-role
)
