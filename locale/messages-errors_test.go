package locale_test

import (
	"errors"
	"fmt"
	"os"

	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok

	"github.com/snivilised/li18ngo"
	"github.com/snivilised/traverse/locale"
	"github.com/snivilised/traverse/test/hydra"
)

var _ = Describe("error messages", Ordered, func() {
	var (
		repo                string
		l10nPath            string
		testTranslationFile li18ngo.TranslationFiles
	)

	BeforeAll(func() {
		repo = hydra.Repo("")
		l10nPath = hydra.Combine(repo, "test/data/l10n")

		_, err := os.Stat(l10nPath)
		Expect(err).To(Succeed(),
			fmt.Sprintf("l10n '%v' path does not exist", l10nPath),
		)

		testTranslationFile = li18ngo.TranslationFiles{
			locale.SourceID: li18ngo.TranslationSource{Name: "foo"},
		}
	})

	BeforeEach(func() {
		if err := li18ngo.Use(func(o *li18ngo.UseOptions) {
			o.Tag = li18ngo.DefaultLanguage
			o.From.Sources = testTranslationFile
		}); err != nil {
			Fail(err.Error())
		}
	})

	Context("InvalidExtGlobFilterMissingSeparator error", func() { // PENDING
		When("variant error created", func() {
			It("should: render translated content", func() {
				const (
					expected = "invalid glob ex filter definition; pattern is missing separator, pattern: foo"
				)
				text := locale.NewInvalidExtGlobFilterMissingSeparatorError(
					"foo",
				).Error()
				Expect(text).To(Equal(expected))
			})
		})

		When("given: matching error", func() {
			It("🧪 should: affirm", func() {
				err := locale.NewInvalidExtGlobFilterMissingSeparatorError("bar")
				Expect(
					locale.IsInvalidExtGlobFilterMissingSeparatorError(err),
				).To(BeTrue(),
					"error does not match InvalidExtGlobFilterMissingSeparator",
				)
			})
		})

		When("given: non matching error", func() {
			It("🧪 should: reject", func() {
				err := errors.New("fake")
				Expect(
					locale.IsInvalidExtGlobFilterMissingSeparatorError(err),
				).To(BeFalse(),
					"not matching error should not match InvalidExtGlobFilterMissingSeparator",
				)
			})
		})
	})

	Context("InvalidIncaseFilterDef error", func() {
		When("variant error created", func() {
			It("should: render translated content", func() {
				const (
					expected = "invalid incase filter definition; pattern is missing separator, pattern: foo"
				)
				text := locale.NewInvalidIncaseFilterDefError(
					"foo",
				).Error()
				Expect(text).To(Equal(expected))
			})
		})

		When("given: matching error", func() {
			It("🧪 should: affirm", func() {
				err := locale.NewInvalidIncaseFilterDefError("bar")
				Expect(
					locale.IsInvalidIncaseFilterDefError(err),
				).To(BeTrue(),
					"error does not match InvalidIncaseFilterDef",
				)
			})
		})

		When("given: non matching error", func() {
			It("🧪 should: reject", func() {
				err := errors.New("fake")
				Expect(
					locale.IsInvalidIncaseFilterDefError(err),
				).To(BeFalse(),
					"not matching error should not match InvalidIncaseFilterDef",
				)
			})
		})
	})
})
