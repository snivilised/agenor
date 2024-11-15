package hiber_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok
	"github.com/snivilised/agenor/core"
	lab "github.com/snivilised/agenor/internal/laboratory"
)

func TestHibernate(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Hibernate Suite")
}

type hibernateTE struct {
	lab.NaviTE
	Hibernate *core.HibernateOptions
	Mute      bool
}
