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
		It("returns the default linear presenter", func() {
			p, err := ui.New("")
			Expect(err).To(BeNil())
			Expect(p).NotTo(BeNil())
		})
	})

	Context("given mode 'linear'", func() {
		It("returns a Presenter without error", func() {
			p, err := ui.New(ui.ModeLinear)
			Expect(err).To(BeNil())
			Expect(p).NotTo(BeNil())
		})
	})

	Context("given an unknown mode", func() {
		It("returns an error containing the unknown mode name", func() {
			p, err := ui.New("flashy")
			Expect(p).To(BeNil())
			Expect(err).NotTo(BeNil())
			// TODO: once lingo generates UnknownModeError, replace the
			// ContainSubstring check with:
			//   var target *locale.UnknownModeError
			//   Expect(errors.As(err, &target)).To(BeTrue())
			Expect(err.Error()).To(ContainSubstring("flashy"))
		})
	})
})

var _ = Describe("RegisterMode", func() {
	Context("registering a new mode", func() {
		It("makes the mode available via New", func() {
			err := ui.RegisterMode("test-stub", func() report.Presenter {
				return &stubPresenter{}
			})
			Expect(err).To(BeNil())

			p, err := ui.New("test-stub")
			Expect(err).To(BeNil())
			Expect(p).NotTo(BeNil())
		})
	})

	Context("registering a duplicate mode", func() {
		It("returns an error", func() {
			// "test-stub" was registered in the previous spec; registering
			// it again must return an error, not panic.
			err := ui.RegisterMode("test-stub", func() report.Presenter {
				return &stubPresenter{}
			})
			Expect(err).NotTo(BeNil())
			// TODO: once lingo generates DuplicateModeError, add:
			//   var target *locale.DuplicateModeError
			//   Expect(errors.As(err, &target)).To(BeTrue())
			Expect(err.Error()).To(ContainSubstring("test-stub"))
		})
	})
})

var _ = Describe("linear Presenter", Ordered, func() {
	var (
		p    report.Presenter
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
		p, err = ui.New(ui.ModeLinear)
		Expect(err).To(BeNil())
		node = &core.Node{Path: "/some/path/file.txt"}
	})

	Describe("OnNodeEvent", func() {
		It("does not panic for a valid node", func() {
			Expect(func() {
				p.OnNodeEvent(&report.NeutralEvent{
					DisplayEvent: report.DisplayEvent{Node: node},
				})
			}).NotTo(Panic())
		})
	})

	Describe("OnActionEvent", func() {
		It("does not panic on success", func() {
			Expect(func() {
				p.OnActionEvent(&report.ActionEvent{
					DisplayEvent: report.DisplayEvent{Node: node, Name: "my-action"},
				})
			}).NotTo(Panic())
		})

		It("does not panic on failure", func() {
			Expect(func() {
				p.OnActionEvent(&report.ActionEvent{
					DisplayEvent: report.DisplayEvent{Node: node, Name: "my-action"},
					Err:          errors.New("action failed"),
				})
			}).NotTo(Panic())
		})
	})

	Describe("OnPipelineEvent", func() {
		It("does not panic on success", func() {
			Expect(func() {
				p.OnPipelineEvent(&report.PipelineEvent{
					DisplayEvent: report.DisplayEvent{Node: node, Name: "my-pipeline"},
				})
			}).NotTo(Panic())
		})

		It("does not panic on failure", func() {
			Expect(func() {
				p.OnPipelineEvent(&report.PipelineEvent{
					DisplayEvent: report.DisplayEvent{Node: node, Name: "my-pipeline"},
					Err:          errors.New("pipeline failed"),
				})
			}).NotTo(Panic())
		})
	})

	Describe("OnComplete", func() {
		It("does not panic on a successful traversal", func() {
			Expect(func() {
				p.OnComplete(&report.Traversal{
					FilesVisited: 10,
					DirsVisited:  3,
				})
			}).NotTo(Panic())
		})

		It("does not panic when the traversal contains an error", func() {
			Expect(func() {
				p.OnComplete(&report.Traversal{
					Err: errors.New("something broke"),
				})
			}).NotTo(Panic())
		})
	})

	Describe("OnSkipEvent", func() {
		It("does not panic for a populated skip event", func() {
			Expect(func() {
				p.OnSkipEvent(&report.SkipEvent{
					DisplayEvent: report.DisplayEvent{Node: node, Name: "my-action"},
					Placeholder:  "{{.grand}}",
					ResolvedPath: "/some",
				})
			}).NotTo(Panic())
		})
	})
})

// ---------------------------------------------------------------------------
// Test double - satisfies report.Presenter for registration tests
// ---------------------------------------------------------------------------

type stubPresenter struct{}

// OnBegin is called once before any traversal events, with the
// opening metadata. Implementations should use this to render
// any opening banner or header.
func (s *stubPresenter) OnBegin(_ *report.BeginEvent)            {}
func (s *stubPresenter) OnNodeEvent(_ *report.NeutralEvent)      {}
func (s *stubPresenter) OnActionEvent(_ *report.ActionEvent)     {}
func (s *stubPresenter) OnPipelineEvent(_ *report.PipelineEvent) {}
func (s *stubPresenter) OnSkipEvent(_ *report.SkipEvent)         {}
func (s *stubPresenter) OnComplete(_ *report.Traversal)          {}
