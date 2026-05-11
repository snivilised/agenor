package prism_test

import (
	"bytes"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/snivilised/jaywalk/src/prism"
	"github.com/snivilised/jaywalk/src/prism/flow"
)

var _ = Describe("New", func() {
	BeforeEach(func() {
		flow.Register()
	})

	Context("when given a valid palette", func() {
		DescribeTable("returns a non-nil Renderer for known view kinds",
			func(kind prism.ViewKind) {
				w := &bytes.Buffer{}
				palette := prism.SystemPalette()

				renderer, err := prism.New(kind, palette, w)

				Expect(err).To(BeNil())
				Expect(renderer).NotTo(BeNil())
			},
			Entry("linear view", prism.LinearView),
		)

		Context("when given an unknown view kind", func() {
			It("falls back to the linear renderer without error", func() {
				w := &bytes.Buffer{}
				palette := prism.SystemPalette()

				renderer, err := prism.New(prism.ViewKind("unknown"), palette, w)

				Expect(err).To(BeNil())
				Expect(renderer).NotTo(BeNil())
			})
		})
	})

	Context("when the palette contains an invalid ansi16 name", func() {
		It("returns an error and a nil renderer", func() {
			w := &bytes.Buffer{}
			palette := prism.SystemPalette()
			palette.Directory = prism.SemanticColour{ANSI16: "notacolour"}

			renderer, err := prism.New(prism.LinearView, palette, w)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("prism.New"))
			Expect(renderer).To(BeNil())
		})
	})
})
