package core_test

import (
	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok
	"github.com/pkg/errors"

	"github.com/snivilised/traverse/core"
)

var _ = Describe("Variable untranslated Errors", func() {
	// These tests generated using snippet: "Ginko Variable Native Error" (gverrt)
	//
	Context("InvalidNotificationMuteRequested error", func() {
		When("given: matching error", func() {
			It("🧪 should: affirm", func() {
				err := core.NewInvalidNotificationMuteRequestedError("OnBegin")
				Expect(
					core.IsInvalidNotificationMuteRequestedError(err),
				).To(BeTrue(),
					"error does not match InvalidNotificationMuteRequested",
				)
			})
		})

		When("given: non matching error", func() {
			It("🧪 should: reject", func() {
				err := errors.New("fake")
				Expect(
					core.IsInvalidNotificationMuteRequestedError(err),
				).To(BeFalse(),
					"not matching error should not match InvalidNotificationMuteRequested",
				)
			})
		})
	})

	Context("InvalidResumeStateTransition error", func() {
		When("given: matching error", func() {
			It("🧪 should: affirm", func() {
				err := core.NewInvalidResumeStateTransitionError("bar")
				Expect(
					core.IsInvalidResumeStateTransitionError(err),
				).To(BeTrue(),
					"error does not match InvalidResumeStateTransition",
				)
			})
		})

		When("given: non matching error", func() {
			It("🧪 should: reject", func() {
				err := errors.New("fake")
				Expect(
					core.IsInvalidResumeStateTransitionError(err),
				).To(BeFalse(),
					"not matching error should not match InvalidResumeStateTransition",
				)
			})
		})
	})

	Context("NewItemAlreadyExtended error", func() {
		When("given: matching error", func() {
			It("🧪 should: affirm", func() {
				err := core.NewItemAlreadyExtendedError("/some-path")
				Expect(
					core.IsNewItemAlreadyExtendedError(err),
				).To(BeTrue(),
					"error does not match NewItemAlreadyExtended",
				)
			})
		})

		When("given: non matching error", func() {
			It("🧪 should: reject", func() {
				err := errors.New("fake")
				Expect(
					core.IsNewItemAlreadyExtendedError(err),
				).To(BeFalse(),
					"not matching error should not match NewItemAlreadyExtended",
				)
			})
		})
	})

	Context("MissingHibernationDetacherFunction error", func() {
		When("given: matching error", func() {
			It("🧪 should: affirm", func() {
				err := core.NewMissingHibernationDetacherFunctionError("bar")
				Expect(
					core.IsMissingHibernationDetacherFunctionError(err),
				).To(BeTrue(),
					"error does not match MissingHibernationDetacherFunction",
				)
			})
		})

		When("given: non matching error", func() {
			It("🧪 should: reject", func() {
				err := errors.New("fake")
				Expect(
					core.IsMissingHibernationDetacherFunctionError(err),
				).To(BeFalse(),
					"not matching error should not match MissingHibernationDetacherFunction",
				)
			})
		})
	})

	Context("InvalidPeriscopeRootPath error", func() {
		When("given: matching error", func() {
			It("🧪 should: affirm", func() {
				err := core.NewInvalidPeriscopeRootPathError("/some-root", "/come-current")
				Expect(
					core.IsInvalidPeriscopeRootPathError(err),
				).To(BeTrue(),
					"error does not match InvalidPeriscopeRootPath",
				)
			})
		})

		When("given: non matching error", func() {
			It("🧪 should: reject", func() {
				err := errors.New("fake")
				Expect(
					core.IsInvalidPeriscopeRootPathError(err),
				).To(BeFalse(),
					"not matching error should not match InvalidPeriscopeRootPath",
				)
			})
		})
	})

	Context("ResumeControllerNotSet error", func() {
		When("given: matching error", func() {
			It("🧪 should: affirm", func() {
				err := core.NewResumeControllerNotSetError("bar")
				Expect(
					core.IsResumeControllerNotSetError(err),
				).To(BeTrue(),
					"error does not match ResumeControllerNotSet",
				)
			})
		})

		When("given: non matching error", func() {
			It("🧪 should: reject", func() {
				err := errors.New("fake")
				Expect(
					core.IsResumeControllerNotSetError(err),
				).To(BeFalse(),
					"not matching error should not match ResumeControllerNotSet",
				)
			})
		})
	})

	Context("BrokerTopicNotFound error", func() {
		When("given: matching error", func() {
			It("🧪 should: affirm", func() {
				err := core.NewBrokerTopicNotFoundError("/foo/bar")
				Expect(
					core.IsBrokerTopicNotFoundError(err),
				).To(BeTrue(),
					"error does not match BrokerTopicNotFound",
				)
			})
		})

		When("given: non matching error", func() {
			It("🧪 should: reject", func() {
				err := errors.New("fake")
				Expect(
					core.IsBrokerTopicNotFoundError(err),
				).To(BeFalse(),
					"not matching error should not match BrokerTopicNotFound",
				)
			})
		})
	})
})
