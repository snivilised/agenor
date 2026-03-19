package ui_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/snivilised/agenor/cmd/ui"
	"github.com/snivilised/agenor/core"
)

// ---------------------------------------------------------------------------
// Specs
// ---------------------------------------------------------------------------

var _ = Describe("ui.New", func() {

	Context("given an empty mode string", func() {
		It("returns the default linear manager", func() {
			m, err := ui.New("")
			Expect(err).To(BeNil())
			Expect(m).NotTo(BeNil())
		})
	})

	Context("given mode 'linear'", func() {
		It("returns a Manager without error", func() {
			m, err := ui.New(ui.ModeLinear)
			Expect(err).To(BeNil())
			Expect(m).NotTo(BeNil())
		})
	})

	Context("given an unknown mode", func() {
		It("returns an ErrUnknownMode error", func() {
			m, err := ui.New("flashy")
			Expect(m).To(BeNil())
			Expect(err).NotTo(BeNil())

			var unknownErr *ui.ErrUnknownMode
			Expect(err).To(BeAssignableToTypeOf(unknownErr))
			Expect(err.Error()).To(ContainSubstring("flashy"))
		})
	})
})

var _ = Describe("RegisterMode", func() {
	Context("registering a new mode", func() {
		It("makes the mode available via New", func() {
			ui.RegisterMode("test-stub", func() ui.Manager {
				return &stubManager{}
			})
			m, err := ui.New("test-stub")
			Expect(err).To(BeNil())
			Expect(m).NotTo(BeNil())
		})
	})

	Context("registering a duplicate mode", func() {
		It("panics", func() {
			Expect(func() {
				ui.RegisterMode("test-stub", func() ui.Manager {
					return &stubManager{}
				})
			}).To(Panic())
		})
	})
})

var _ = Describe("linear Manager", func() {
	var m ui.Manager

	BeforeEach(func() {
		var err error
		m, err = ui.New(ui.ModeLinear)
		Expect(err).To(BeNil())
	})

	Describe("OnNode", func() {
		It("does not return an error for a valid node", func() {
			node := &core.Node{Path: "/some/path/file.txt"}
			Expect(m.OnNode(node)).To(BeNil())
		})
	})

	Describe("Info / Warn / Error", func() {
		It("Info does not panic", func() {
			Expect(func() { m.Info("all good") }).NotTo(Panic())
		})

		It("Warn does not panic", func() {
			Expect(func() { m.Warn("something odd") }).NotTo(Panic())
		})

		It("Error does not panic", func() {
			Expect(func() { m.Error("something broke") }).NotTo(Panic())
		})
	})
})

// ---------------------------------------------------------------------------
// Test double - satisfies ui.Manager for registration tests
// ---------------------------------------------------------------------------

type stubManager struct{}

func (s *stubManager) OnNode(_ *core.Node) error { return nil }
func (s *stubManager) Info(_ string)             {}
func (s *stubManager) Warn(_ string)             {}
func (s *stubManager) Error(_ string)            {}
