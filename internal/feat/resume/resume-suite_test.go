package resume_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok
	"github.com/snivilised/agenor/enums"
	lab "github.com/snivilised/agenor/internal/laboratory"
)

func TestResume(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Resume Suite")
}

type (
	activeTE struct {
		resumeAt    string
		listenState enums.Hibernation // rename listenState...
	}

	resumeTE struct {
		lab.NaviTE
		active         activeTE
		clientListenAt string // rename clientListenAt
		profile        string
	}
)
