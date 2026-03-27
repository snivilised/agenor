package bedrock_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	bedrock "github.com/snivilised/jaywalk/src/app/bedrock"
	"github.com/snivilised/jaywalk/src/locale"
	"github.com/snivilised/li18ngo"
)

var _ = Describe("Validate", Ordered, func() {
	BeforeAll(func() {
		Expect(li18ngo.Use(
			func(o *li18ngo.UseOptions) {
				o.From.Sources = li18ngo.TranslationFiles{
					locale.SourceID: li18ngo.TranslationSource{Name: "agenor"},
				}
			},
		)).To(Succeed())
	})

	// -----------------------------------------------------------------------
	// LoggingConfig
	// -----------------------------------------------------------------------
	Describe("LoggingConfig.Validate", func() {
		DescribeTable("invalid log levels",
			func(level string) {
				c := bedrock.LoggingConfig{Level: level}
				Expect(c.Validate()).NotTo(Succeed())
			},
			Entry("verbose", "verbose"),
			Entry("ALL", "ALL"),
			Entry("notice", "notice"),
		)

		DescribeTable("valid log levels",
			func(level string) {
				c := bedrock.LoggingConfig{Level: level}
				Expect(c.Validate()).To(Succeed())
			},
			Entry("trace", "trace"),
			Entry("debug", "debug"),
			Entry("info", "info"),
			Entry("warn", "warn"),
			Entry("error", "error"),
			Entry("fatal", "fatal"),
			Entry("panic", "panic"),
			Entry("empty (optional)", ""),
		)

		It("rejects negative max-size-in-mb", func() {
			c := bedrock.LoggingConfig{MaxSizeInMB: -1, Level: "info"}
			Expect(c.Validate()).NotTo(Succeed())
		})

		It("rejects negative max-backups", func() {
			c := bedrock.LoggingConfig{MaxBackups: -1, Level: "info"}
			Expect(c.Validate()).NotTo(Succeed())
		})

		It("rejects negative max-age-in-days", func() {
			c := bedrock.LoggingConfig{MaxAgeInDays: -1, Level: "info"}
			Expect(c.Validate()).NotTo(Succeed())
		})

		It("accepts zero values", func() {
			c := bedrock.LoggingConfig{}
			Expect(c.Validate()).To(Succeed())
		})
	})

	Describe("InteractionConfig.Validate", func() {
		It("rejects bad log level", func() {
			_, err := bedrock.Load(bedrock.LoadOptions{
				ViperInstance: viperFromYAML(badLogLevelYAML),
			})
			Expect(err).NotTo(BeNil())
			Expect(err.Error()).To(ContainSubstring("not a recognised level"))
		})
	})

	// -----------------------------------------------------------------------
	// InteractionConfig
	// -----------------------------------------------------------------------
	Describe("InteractionConfig.Validate", func() {
		It("rejects negative per-item-delay", func() {
			_, err := bedrock.Load(bedrock.LoadOptions{
				ViperInstance: viperFromYAML(negativeDurationYAML),
			})
			Expect(err).NotTo(BeNil())
			Expect(err.Error()).To(ContainSubstring("per-item-delay"))
		})
	})

	// -----------------------------------------------------------------------
	// FlagsConfig
	// -----------------------------------------------------------------------
	Describe("FlagsConfig.Validate", func() {
		It("rejects multi-character short overrides", func() {
			_, err := bedrock.Load(bedrock.LoadOptions{
				ViperInstance: viperFromYAML(badShortYAML),
			})
			Expect(err).NotTo(BeNil())
			Expect(err.Error()).To(ContainSubstring("exactly one character"))
		})
	})

	// -----------------------------------------------------------------------
	// Actions
	// -----------------------------------------------------------------------
	Describe("actions validation", func() {
		It("rejects an action with an empty cmd", func() {
			_, err := bedrock.Load(bedrock.LoadOptions{
				ViperInstance: viperFromYAML(emptyCmdYAML),
			})
			Expect(err).NotTo(BeNil())
			Expect(err.Error()).To(ContainSubstring("bad-action"))
		})
	})

	// -----------------------------------------------------------------------
	// Pipelines
	// -----------------------------------------------------------------------
	Describe("pipeline validation", func() {
		It("rejects a pipeline step referencing an unknown action", func() {
			_, err := bedrock.Load(bedrock.LoadOptions{
				ViperInstance: viperFromYAML(missingActionYAML),
			})
			Expect(err).NotTo(BeNil())
			Expect(err.Error()).To(ContainSubstring("ghost-action"))
		})
	})

	// -----------------------------------------------------------------------
	// Multiple validation failures are returned together
	// -----------------------------------------------------------------------
	Describe("aggregate errors", func() {
		It("reports all failures in one error", func() {
			c := &bedrock.Config{}
			c.Mapped.Logging.Level = "bad-level"
			c.Mapped.Logging.MaxSizeInMB = -5
			err := c.Validate()
			Expect(err).NotTo(BeNil())
			Expect(err.Error()).To(ContainSubstring("bad-level"))
			Expect(err.Error()).To(ContainSubstring("max-size-in-mb"))
		})
	})
})
