package resume_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestResume(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Resume Suite")
}
