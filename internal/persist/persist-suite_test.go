package persist_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok
	"github.com/snivilised/traverse/internal/opts/json"
	"github.com/snivilised/traverse/pref"
)

func TestPersist(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Persist Suite")
}

const (
	NoOverwrite = false
	from        = "json/unmarshal"
	to          = "json/marshal"
	permDir     = 0o777
	permFile    = 0o666
	tempFile    = "test-state-marshal.TEMP.json"
	home        = "/home"
)

type (
	persistTE struct {
		given string
	}
	marshalTE struct {
		persistTE
		// option defines a single option to be defined for the unit test. When
		// a test case wants to test an optional option in pref.Options (ie it
		// is a pointer), then that test case will not define this option. Instead
		// it will define the tweak function to contain the corresponding member
		// on the json instance, such that the pref.member is nil and json.member
		// is noy til, thereby triggering an unequal error.
		option func() pref.Option

		// tweak allows a test case to change json.Options to provoke unequal error
		tweak func(jo *json.Options)
	}
)
