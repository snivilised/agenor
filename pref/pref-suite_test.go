package pref_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok
	"github.com/snivilised/agenor/core"
	"github.com/snivilised/agenor/enums"
	"github.com/snivilised/agenor/pref"
)

func TestPref(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Pref Suite")
}

type (
	testFaultHandler struct{}
	testPanicHandler struct{}
	testSkipHandler  struct{}
)

func (*testFaultHandler) Accept(*pref.NavigationFault) error                        { return nil }
func (*testPanicHandler) Rescue(_ pref.Recovery, _ pref.RescueData) (string, error) { return "", nil }
func (*testSkipHandler) Ask(*core.Node, core.DirectoryContents, error) (enums.SkipTraversal, error) {
	return enums.SkipAllTraversal, nil
}
