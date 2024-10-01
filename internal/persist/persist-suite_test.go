package persist_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok
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
	marshalTE struct {
		given  string
		option func() pref.Option
	}

	errorTE struct {
		marshalTE
	}

	conversionTE struct {
		marshalTE
	}
)
