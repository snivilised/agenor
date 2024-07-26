package tv

import (
	"slices"

	"github.com/snivilised/traverse/enums"
	"github.com/snivilised/traverse/internal/lo"
)

type (
	rule          func(current enums.Role, active, all []enums.Role) bool
	manifestRules map[string]rule
)

var (
	manifestOrder = []enums.Role{
		enums.RoleFastward,
		enums.RoleClientHiberWake,
		enums.RoleClientHiberSleep,
		enums.RoleClientFilter,
		enums.RoleSampler,
	}
)

// manifest defines the order of roles and which roles can be
// activated at a time over the top of the client callback
// function.
func manifest(active []enums.Role) []enums.Role {
	all := manifestOrder
	rules := manifestRules{
		"contained-in-all": func(current enums.Role, _, all []enums.Role) bool {
			return slices.Contains(all, current)
		},
		"filter-defers-to-sampler": func(current enums.Role, active, _ []enums.Role) bool {
			if current == enums.RoleClientFilter &&
				slices.Contains(active, enums.RoleSampler) {
				return false
			}
			return true
		},
	}

	// only roles that satisfy all rules are returned
	//
	return lo.Reduce(active,
		func(acc []enums.Role, role enums.Role, _ int) []enums.Role {
			if lo.EveryBy(lo.Keys(rules), func(name string) bool {
				return rules[name](role, acc, all)
			}) {
				acc = append(acc, role)
			}
			return acc
		},
		make([]enums.Role, 0, len(active)+1),
	)
}
