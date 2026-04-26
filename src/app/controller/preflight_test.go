package controller_test

import (
	"errors"
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	bedrock "github.com/snivilised/jaywalk/src/app/bedrock"
	"github.com/snivilised/jaywalk/src/app/controller"
	"github.com/snivilised/jaywalk/src/app/shell"
)

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

// stubLocate returns a shell.LocateFunc that succeeds for every name in
// found and fails for every name in notFound. Names not mentioned in
// either set cause a test panic so omissions are caught early.
func stubLocate(found, notFound []string) shell.LocateFunc {
	ok := make(map[string]bool, len(found))
	for _, n := range found {
		ok[n] = true
	}

	fail := make(map[string]bool, len(notFound))
	for _, n := range notFound {
		fail[n] = true
	}

	return func(name string) (string, error) {
		if ok[name] {
			return "/usr/bin/" + name, nil
		}
		if fail[name] {
			return "", fmt.Errorf("shell: %q not found in current environment", name)
		}
		panic(fmt.Sprintf(
			"stubLocate: unexpected executable %q - add it to found or notFound", name,
		))
	}
}

// minimalCfg builds a *bedrock.Config from the supplied actions and
// pipelines maps. All other config fields are left at zero values.
func minimalCfg(
	actions map[string]bedrock.RawAction,
	pipelines map[string]bedrock.RawPipeline,
) *bedrock.Config {
	return &bedrock.Config{
		Raw: bedrock.RawConfig{
			Actions:   actions,
			Pipelines: pipelines,
		},
	}
}

// ---------------------------------------------------------------------------
// Specs
// ---------------------------------------------------------------------------

var _ = Describe("Coordinator.PreFlight", func() {

	// ------------------------------------------------------------------
	// Neutral - no action or pipeline configured
	// ------------------------------------------------------------------

	Context("when neither ActionName nor PipelineName is set", func() {
		It("returns nil without calling locate", func() {
			called := false
			spy := shell.LocateFunc(func(name string) (string, error) {
				called = true
				return "", errors.New("should not be called")
			})

			coord := controller.New(minimalCfg(nil, nil), controller.WithLocate(spy))
			err := coord.PreFlight(&controller.Request{})

			Expect(err).To(BeNil())
			Expect(called).To(BeFalse())
		})
	})

	// ------------------------------------------------------------------
	// Action cases
	// ------------------------------------------------------------------

	Context("when ActionName is set", func() {
		var cfg *bedrock.Config

		BeforeEach(func() {
			cfg = minimalCfg(
				map[string]bedrock.RawAction{
					"good-action":  {Cmd: "ffmpeg -i {{.path}} out.mp4"},
					"bad-action":   {Cmd: "nonexistent-binary arg1"},
					"empty-action": {Cmd: ""},
				},
				nil,
			)
		})

		It("returns nil when the action executable is locatable", func() {
			coord := controller.New(cfg, controller.WithLocate(
				stubLocate([]string{"ffmpeg"}, nil),
			))

			err := coord.PreFlight(&controller.Request{ActionName: "good-action"})

			Expect(err).To(BeNil())
		})

		It("returns an error when the action is not defined in config", func() {
			coord := controller.New(cfg, controller.WithLocate(
				stubLocate(nil, nil),
			))

			err := coord.PreFlight(&controller.Request{ActionName: "undefined-action"})

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("undefined-action"))
			Expect(err.Error()).To(ContainSubstring("not defined in config"))
		})

		It("returns an error when the action cmd is empty", func() {
			coord := controller.New(cfg, controller.WithLocate(
				stubLocate(nil, nil),
			))

			err := coord.PreFlight(&controller.Request{ActionName: "empty-action"})

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("empty-action"))
			Expect(err.Error()).To(ContainSubstring("empty cmd"))
		})

		It("returns an error when the executable is not locatable", func() {
			coord := controller.New(cfg, controller.WithLocate(
				stubLocate(nil, []string{"nonexistent-binary"}),
			))

			err := coord.PreFlight(&controller.Request{ActionName: "bad-action"})

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("bad-action"))
			Expect(err.Error()).To(ContainSubstring("nonexistent-binary"))
			Expect(err.Error()).To(ContainSubstring("not found in current environment"))
		})
	})

	// ------------------------------------------------------------------
	// Pipeline cases
	// ------------------------------------------------------------------

	Context("when PipelineName is set", func() {
		var cfg *bedrock.Config

		BeforeEach(func() {
			cfg = minimalCfg(
				map[string]bedrock.RawAction{
					"encode":   {Cmd: "ffmpeg -i {{.path}} out.mp4"},
					"upload":   {Cmd: "aws s3 cp {{.path}} s3://bucket/"},
					"bad-step": {Cmd: "nonexistent-binary arg1"},
				},
				map[string]bedrock.RawPipeline{
					"good-pipeline":   {Steps: []string{"encode", "upload"}},
					"bad-pipeline":    {Steps: []string{"encode", "bad-step"}},
					"orphan-pipeline": {Steps: []string{"encode", "missing-action"}},
				},
			)
		})

		It("returns nil when all steps resolve and all executables are locatable", func() {
			coord := controller.New(cfg, controller.WithLocate(
				stubLocate([]string{"ffmpeg", "aws"}, nil),
			))

			err := coord.PreFlight(&controller.Request{PipelineName: "good-pipeline"})

			Expect(err).To(BeNil())
		})

		It("returns an error when the pipeline is not defined in config", func() {
			coord := controller.New(cfg, controller.WithLocate(
				stubLocate(nil, nil),
			))

			err := coord.PreFlight(&controller.Request{PipelineName: "undefined-pipeline"})

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("undefined-pipeline"))
			Expect(err.Error()).To(ContainSubstring("not defined in config"))
		})

		It("returns an error when a pipeline step is not defined in actions", func() {
			coord := controller.New(cfg, controller.WithLocate(
				stubLocate([]string{"ffmpeg"}, nil),
			))

			err := coord.PreFlight(&controller.Request{PipelineName: "orphan-pipeline"})

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("orphan-pipeline"))
			Expect(err.Error()).To(ContainSubstring("missing-action"))
			Expect(err.Error()).To(ContainSubstring("not defined in config"))
		})

		It("returns an error when a step executable is not locatable", func() {
			coord := controller.New(cfg, controller.WithLocate(
				stubLocate([]string{"ffmpeg"}, []string{"nonexistent-binary"}),
			))

			err := coord.PreFlight(&controller.Request{PipelineName: "bad-pipeline"})

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("bad-pipeline"))
			Expect(err.Error()).To(ContainSubstring("nonexistent-binary"))
			Expect(err.Error()).To(ContainSubstring("not found in current environment"))
		})

		It("stops at the first failing step and does not check subsequent steps", func() {
			queriedNames := []string{}
			spy := shell.LocateFunc(func(name string) (string, error) {
				queriedNames = append(queriedNames, name)
				if name == "ffmpeg" {
					return "/usr/bin/ffmpeg", nil
				}
				return "", fmt.Errorf("not found")
			})

			coord := controller.New(cfg, controller.WithLocate(spy))
			err := coord.PreFlight(&controller.Request{PipelineName: "bad-pipeline"})

			Expect(err).To(HaveOccurred())
			Expect(queriedNames).NotTo(ContainElement("aws"))
		})
	})
})
