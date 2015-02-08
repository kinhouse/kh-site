package integration

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/sclevine/agouti/core"
	. "github.com/sclevine/agouti/matchers"
)

var _ = Describe("RSVPing to an invite", func() {
	var page Page

	BeforeEach(func() {
		var err error

		page, err = agoutiDriver.Page()
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		page.Destroy()
	})

	It("is reachable from the home page", func() {
		Expect(page.Navigate(baseUrl)).To(Succeed())
		Expect(page.FindByLink("RSVP").Click()).To(Succeed())
		Expect(page).To(HaveURL(baseUrl + "/rsvp"))
		Expect(page).To(HaveTitle("Alana & Gabe: RSVP"))
	})

	It("may be Declined", func() {
		By("Navigating to the RSVP page", func() {
			Expect(page.Navigate(baseUrl + "/rsvp")).To(Succeed())
		})
		By("entering her name and email", func() {
			Expect(page.FindByLabel("Name:").Fill("Someone not attending")).To(Succeed())
			Expect(page.FindByLabel("Email:").Fill("someone@example.com")).To(Succeed())
		})
		By("and selecting the 'No' option", func() {
			Expect(page.FindByLabel("No, I will not be attending").Click()).To(Succeed())
		})
		By("and pressing 'Submit'", func() {
			Expect(page.Find("#rsvp_form").Submit()).To(Succeed())
		})
		By("showing a friendly response", func() {
			Expect(page).To(HaveURL(baseUrl + "/rsvp"))
			// Expect(page).To(HaveTitle("RSVP: Declined"))
			// Expect(page.Find(".acknowledgement")).To(HaveText("Too bad!"))
		})
	})

	It("may be accepted", func() {
		By("Navigating to the RSVP page", func() {
			Expect(page.Navigate(baseUrl + "/rsvp")).To(Succeed())
		})

		By("entering her name and email", func() {
			Expect(page.FindByLabel("Name:").Fill("Someone attending")).To(Succeed())
			Expect(page.FindByLabel("Email:").Fill("someone@example.com")).To(Succeed())
		})
		By("and selecting the 'Yes' option", func() {
			Expect(page.FindByLabel("Yes, looking forward to it!").Click()).To(Succeed())
		})
		By("selecting how big the party is", func() {

		})
		By("and pressing 'Submit'", func() {
			Expect(page.Find("#rsvp_form").Submit()).To(Succeed())
		})
		By("showing a friendly response", func() {
			Expect(page).To(HaveURL(baseUrl + "/rsvp"))
			// Expect(page).To(HaveTitle("RSVP: Accepted"))
			// Expect(page.Find(".acknowledgement")).To(HaveText("Sweet!"))
		})
	})

	Describe("Error cases", func() {
		XIt("does not allow submission when the email field is blank or malformed", func() {

		})
		XIt("does not allow submission when the attendance options are both unselected", func() {

		})
	})

})
