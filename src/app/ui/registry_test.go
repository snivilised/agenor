package ui_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/snivilised/jaywalk/src/app/report"
	"github.com/snivilised/jaywalk/src/app/ui"
	"github.com/snivilised/jaywalk/src/prism"
)

var _ = Describe("Registry", func() {

	// ------------------------------------------------------------------
	// New
	// ------------------------------------------------------------------

	Describe("New", func() {
		DescribeTable("returns a Presenter for known modes",
			func(mode string) {
				palette := prism.SystemPalette()

				presenter, err := ui.New(mode, palette)

				Expect(err).To(BeNil())
				Expect(presenter).NotTo(BeNil())
			},
			Entry("explicit linear mode", ui.ModeLinear),
			Entry("empty string defaults to linear", ""),
		)

		Context("when the mode is not registered", func() {
			It("returns an error containing the unknown mode name", func() {
				palette := prism.SystemPalette()

				_, err := ui.New("nonexistent-mode", palette)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("nonexistent-mode"))
			})
		})

		Context("when the palette contains an invalid ansi16 name", func() {
			It("returns an error propagated from prism", func() {
				palette := prism.SystemPalette()
				palette.Directory = prism.SemanticColour{ANSI16: "notacolour"}

				_, err := ui.New(ui.ModeLinear, palette)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("notacolour"))
			})
		})
	})

	// ------------------------------------------------------------------
	// RegisterMode
	// ------------------------------------------------------------------

	Describe("RegisterMode", func() {
		Context("when registering a new unique mode name", func() {
			It("succeeds and the mode becomes available via New", func() {
				const testMode = "test-mode-unique"

				err := ui.RegisterMode(testMode,
					func(palette prism.Palette) (report.Presenter, error) {
						return ui.New(ui.ModeLinear, palette)
					},
				)

				Expect(err).To(BeNil())

				presenter, err := ui.New(testMode, prism.SystemPalette())
				Expect(err).To(BeNil())
				Expect(presenter).NotTo(BeNil())
			})
		})

		Context("when registering a name that already exists", func() {
			It("returns an error containing the duplicate name", func() {
				err := ui.RegisterMode(ui.ModeLinear,
					func(palette prism.Palette) (report.Presenter, error) {
						return nil, nil
					},
				)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring(ui.ModeLinear))
			})
		})
	})
})
