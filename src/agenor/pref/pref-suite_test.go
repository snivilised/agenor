package pref_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/snivilised/jaywalk/src/agenor/core"
	"github.com/snivilised/jaywalk/src/agenor/enums"
	"github.com/snivilised/jaywalk/src/agenor/pref"
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
