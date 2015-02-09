package server_test

import (
	"github.com/kinhouse/kh-site/server"
	"github.com/kinhouse/kh-site/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("RSVP Handler", func() {
	var rsvp types.Rsvp

	BeforeEach(func() {
		rsvp = types.Rsvp{
			FullName: "A Person",
			Email:    "someone@example.com",
		}
	})

	Context("when the rsvp is a decline", func() {
		It("responds kindly", func() {
			rsvp.Decline = true
			Expect(server.RsvpHandler(rsvp)).To(Equal(
				"We're sorry you won't be making it, but thanks for letting us know!"))
		})
	})

	Context("when the rsvp is an accept", func() {
		Context("when it is only one person", func() {
			It("responds in the singular", func() {
				rsvp.Count = 1
				Expect(server.RsvpHandler(rsvp)).To(Equal(
					"Wonderful!  We're thrilled that you'll be joining us.  Poke around this site for more information about the wedding and our story."))
			})
		})

		Context("when it two people", func() {
			It("responds in to the couple", func() {
				rsvp.Count = 2
				Expect(server.RsvpHandler(rsvp)).To(Equal(
					"Wonderful!  We're thrilled that the two of you will be joining us.  Poke around this site for more information about the wedding and our story."))
			})
		})

		Context("when it is 3 people", func() {
			It("responds to the group", func() {
				rsvp.Count = 3
				Expect(server.RsvpHandler(rsvp)).To(Equal(
					"Wonderful!  We're thrilled that the three of you will be joining us.  Poke around this site for more information about the wedding and our story."))
			})
		})

		Context("when it is 4 people", func() {
			It("responds to the group", func() {
				rsvp.Count = 4
				Expect(server.RsvpHandler(rsvp)).To(Equal(
					"Wonderful!  We're thrilled that the four of you will be joining us.  Poke around this site for more information about the wedding and our story."))
			})
		})

		Context("when it is 5 people", func() {
			It("responds to the group", func() {
				rsvp.Count = 5
				Expect(server.RsvpHandler(rsvp)).To(Equal(
					"Wonderful!  We're thrilled that all five (!) of you will be joining us.  Poke around this site for more information about the wedding and our story."))
			})
		})

		Context("when it is 6 or more people", func() {
			It("responds to generically", func() {
				rsvp.Count = 6
				Expect(server.RsvpHandler(rsvp)).To(Equal(
					"Wonderful!  We're thrilled that you will be joining us.  Poke around this site for more information about the wedding and our story."))
			})
		})

	})

	It("never echos back any user-submitted strings", func() {
		rsvp.FullName = "Evil person tries to put in some <script> tag"
		responseString := server.RsvpHandler(rsvp)
		Expect(responseString).ToNot(ContainSubstring("Evil"))
	})

})
