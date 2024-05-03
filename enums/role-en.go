package enums

//go:generate stringer -type=Role -linecomment -trimprefix=Role -output role-en-auto.go

type Role uint32

const (
	RoleUndefined       Role = iota // undefined-role
	RoleDirectoryReader             // directory-reader-role
)
