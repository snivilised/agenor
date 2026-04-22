package ui_test

import (
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/snivilised/jaywalk/src/agenor/core"
	"github.com/snivilised/jaywalk/src/app/report"
	"github.com/snivilised/jaywalk/src/app/ui"
	"github.com/snivilised/jaywalk/src/locale"
	"github.com/snivilised/li18ngo"
)

// ---------------------------------------------------------------------------
// Specs
// ---------------------------------------------------------------------------

var _ = Describe("ui.New", Ordered, func() {
	BeforeAll(func() {
		Expect(li18ngo.Register(
			func(o *li18ngo.UseOptions) {
				o.From.Sources = li18ngo.TranslationFiles{
					locale.SourceID: li18ngo.TranslationSource{Name: "agenor"},
				}
			},
		)).To(Succeed())
	})

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

var _ = Describe("linear Manager", Ordered, func() {
	var (
		m    ui.Manager
		node *core.Node
	)

	BeforeAll(func() {
		Expect(li18ngo.Register(
			func(o *li18ngo.UseOptions) {
				o.From.Sources = li18ngo.TranslationFiles{
					locale.SourceID: li18ngo.TranslationSource{Name: "agenor"},
				}
			},
		)).To(Succeed())
	})

	BeforeEach(func() {
		var err error
		m, err = ui.New(ui.ModeLinear)
		Expect(err).To(BeNil())
		node = &core.Node{Path: "/some/path/file.txt"}
	})

	Describe("OnNodeEvent", func() {
		It("does not panic for a valid node", func() {
			Expect(func() {
				m.OnNodeEvent(&report.NeutralEvent{
					DisplayEvent: report.DisplayEvent{Node: node},
				})
			}).NotTo(Panic())
		})
	})

	Describe("OnActionEvent", func() {
		It("does not panic on success", func() {
			Expect(func() {
				m.OnActionEvent(&report.ActionEvent{
					DisplayEvent: report.DisplayEvent{Node: node, Name: "my-action"},
				})
			}).NotTo(Panic())
		})

		It("does not panic on failure", func() {
			Expect(func() {
				m.OnActionEvent(&report.ActionEvent{
					DisplayEvent: report.DisplayEvent{Node: node, Name: "my-action"},
					Err:          errors.New("action failed"),
				})
			}).NotTo(Panic())
		})
	})

	Describe("OnPipelineEvent", func() {
		It("does not panic on success", func() {
			Expect(func() {
				m.OnPipelineEvent(&report.PipelineEvent{
					DisplayEvent: report.DisplayEvent{Node: node, Name: "my-pipeline"},
				})
			}).NotTo(Panic())
		})

		It("does not panic on failure", func() {
			Expect(func() {
				m.OnPipelineEvent(&report.PipelineEvent{
					DisplayEvent: report.DisplayEvent{Node: node, Name: "my-pipeline"},
					Err:          errors.New("pipeline failed"),
				})
			}).NotTo(Panic())
		})
	})

	Describe("OnComplete", func() {
		It("does not panic on a successful traversal", func() {
			Expect(func() {
				m.OnComplete(&report.Traversal{
					FilesVisited: 10,
					DirsVisited:  3,
				})
			}).NotTo(Panic())
		})

		It("does not panic when the traversal contains an error", func() {
			Expect(func() {
				m.OnComplete(&report.Traversal{
					Err: errors.New("something broke"),
				})
			}).NotTo(Panic())
		})
	})
})

// ---------------------------------------------------------------------------
// Test double - satisfies ui.Manager for registration tests
// ---------------------------------------------------------------------------

type stubManager struct{}

func (s *stubManager) OnNodeEvent(_ *report.NeutralEvent)      {}
func (s *stubManager) OnActionEvent(_ *report.ActionEvent)     {}
func (s *stubManager) OnPipelineEvent(_ *report.PipelineEvent) {}
func (s *stubManager) OnComplete(_ *report.Traversal)          {}
