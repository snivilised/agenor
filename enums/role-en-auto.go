// Code generated by "stringer -type=Role -linecomment -trimprefix=Role -output role-en-auto.go"; DO NOT EDIT.

package enums

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[RoleUndefined-0]
	_ = x[RoleDirectoryReader-1]
}

const _Role_name = "undefined-roledirectory-reader-role"

var _Role_index = [...]uint8{0, 14, 35}

func (i Role) String() string {
	if i >= Role(len(_Role_index)-1) {
		return "Role(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Role_name[_Role_index[i]:_Role_index[i+1]]
}
