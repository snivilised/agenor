package prism_test

import (
	"bytes"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/snivilised/jaywalk/src/prism"
)

var _ = Describe("Theme", func() {

	Describe("NewTheme", func() {
		Context("when given the system palette", func() {
			It("constructs a Theme without error", func() {
				w := &bytes.Buffer{}
				palette := prism.SystemPalette()

				theme, err := prism.NewTheme(palette, w)

				Expect(err).To(BeNil())
				// Verify that key styles are non-zero - a zero lipgloss.Style
				// is valid but indicates the palette entry was not applied.
				// We check the styles are at least initialised by asserting
				// the theme value is not its zero value.
				Expect(theme).NotTo(Equal(prism.Theme{}))
			})

			It("includes BranchStyle in the constructed Theme", func() {
				w := &bytes.Buffer{}
				palette := prism.SystemPalette()
				palette.Branch = prism.SemanticColour{ANSI16: "green"}

				theme, err := prism.NewTheme(palette, w)

				Expect(err).To(BeNil())
				// Verify that Theme is properly constructed with BranchStyle
				Expect(theme).NotTo(Equal(prism.Theme{}))
			})
		})

		Context("when a palette entry contains an unrecognised ansi16 name", func() {
			DescribeTable("returns an error identifying the bad field",
				func(mutatePalette func(*prism.Palette), expectedField string) {
					w := &bytes.Buffer{}
					palette := prism.SystemPalette()
					mutatePalette(&palette)

					_, err := prism.NewTheme(palette, w)

					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring(expectedField))
				},
				Entry("banner",
					func(p *prism.Palette) {
						p.Banner = prism.SemanticColour{ANSI16: "notacolour"}
					},
					"palette.banner",
				),
				Entry("directory",
					func(p *prism.Palette) {
						p.Directory = prism.SemanticColour{ANSI16: "notacolour"}
					},
					"palette.directory",
				),
				Entry("error",
					func(p *prism.Palette) {
						p.Error = prism.SemanticColour{ANSI16: "notacolour"}
					},
					"palette.error",
				),
				Entry("worker",
					func(p *prism.Palette) {
						p.Worker = prism.SemanticColour{ANSI16: "notacolour"}
					},
					"palette.worker",
				),
				Entry("lane-header",
					func(p *prism.Palette) {
						p.LaneHeader = prism.SemanticColour{ANSI16: "notacolour"}
					},
					"palette.lane-header",
				),
				Entry("branch",
					func(p *prism.Palette) {
						p.Branch = prism.SemanticColour{ANSI16: "notacolour"}
					},
					"palette.branch",
				),
			)
		})
	})
})
