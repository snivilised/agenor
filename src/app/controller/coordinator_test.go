package controller_test

import (
	"log/slog"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/snivilised/jaywalk/src/app/bedrock"
	jac "github.com/snivilised/jaywalk/src/app/controller"
)

var _ = Describe("Coordinator options", func() {
	var cfg *bedrock.Config

	BeforeEach(func() {
		cfg = &bedrock.Config{}
	})

	Describe("WithAdminPath", func() {
		It("sets the admin path on the Coordinator", func() {
			coord := jac.New(cfg, jac.WithAdminPath("/custom/admin/resume"))
			Expect(coord.AdminPath()).To(Equal("/custom/admin/resume"))
		})

		It("defaults to empty when not provided", func() {
			coord := jac.New(cfg)
			Expect(coord.AdminPath()).To(BeEmpty())
		})
	})

	Describe("WithLogger", func() {
		It("sets the logger on the Coordinator", func() {
			logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
			coord := jac.New(cfg, jac.WithLogger(logger))
			Expect(coord.Logger()).NotTo(BeNil())
			Expect(coord.Logger()).To(Equal(logger))
		})

		It("defaults to nil when not provided", func() {
			coord := jac.New(cfg)
			Expect(coord.Logger()).To(BeNil())
		})
	})

	Describe("Combined options", func() {
		It("sets both adminPath and logger when both options are provided", func() {
			logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
			coord := jac.New(cfg,
				jac.WithAdminPath("/combined/admin"),
				jac.WithLogger(logger),
			)
			Expect(coord.AdminPath()).To(Equal("/combined/admin"))
			Expect(coord.Logger()).NotTo(BeNil())
		})
	})
})
