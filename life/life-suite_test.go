package life_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok
)

func TestLife(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Life Suite")
}
