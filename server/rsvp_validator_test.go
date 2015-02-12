package server_test

import (
	. "github.com/kinhouse/kh-site/server"
	"github.com/kinhouse/kh-site/types"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("RsvpValidator", func() {

	Context("when the rsvp is missing a name", func() {
		It("returns an error", func() {
			Expect(ValidateRsvp(types.Rsvp{
				FullName: "   ",
			})).To(MatchError("Please provide your name!"))
		})
	})
	Context("when the rsvp is missing an email", func() {
		It("returns an error", func() {
			Expect(ValidateRsvp(types.Rsvp{
				FullName: "Someone",
				Email:    "nope@   ",
			})).To(MatchError("Please provide a valid email address."))
		})
	})
	Context("when the rsvp has decline=false but count=0", func() {
		It("returns an error", func() {
			Expect(ValidateRsvp(types.Rsvp{
				FullName: "Someone",
				Email:    "someone@example.com",
			})).To(MatchError("Please indicate if you're a Yes or a No."))
		})
	})
	Context("when the rsvp has count > 9", func() {
		It("returns an error", func() {
			Expect(ValidateRsvp(types.Rsvp{
				FullName: "Someone",
				Email:    "someone@example.com",
				Count:    10,
			})).To(MatchError("That's quite a group!  Please get in touch with us."))
		})
	})
	Context("when the rsvp has count < 0", func() {
		It("returns an error", func() {
			Expect(ValidateRsvp(types.Rsvp{
				FullName: "Someone",
				Email:    "someone@example.com",
				Count:    -1,
			})).To(MatchError("Your guest count is...negative?"))
		})
	})
	Context("when decline is true but the count != 0", func() {
		It("returns an error", func() {
			Expect(ValidateRsvp(types.Rsvp{
				FullName: "Someone",
				Email:    "someone@example.com",
				Decline:  true,
				Count:    1,
			})).To(MatchError("If you're declining, please set your group size to 0."))
		})

	})
})
