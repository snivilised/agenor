package hiber_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok
	"github.com/snivilised/traverse/core"
	lab "github.com/snivilised/traverse/internal/laboratory"
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
