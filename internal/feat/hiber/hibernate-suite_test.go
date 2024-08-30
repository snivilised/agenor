package hiber_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok
	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/internal/helpers"
)

func TestHibernate(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Hibernate Suite")
}

type hibernateTE struct {
	helpers.NaviTE
	Hibernate *core.HibernateOptions
	Mute      bool
}
