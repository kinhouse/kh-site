package integration

import (
	"github.com/kinhouse/kh-site/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sclevine/agouti"
	. "github.com/sclevine/agouti/matchers"
)

var _ = Describe("RSVPing to an invite", func() {
	var page *agouti.Page

	BeforeEach(func() {
		var err error
		page, err = agoutiDriver.NewPage()
		Expect(err).NotTo(HaveOccurred())
		localPersist.Rsvps = []types.Rsvp{}
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
			Expect(page.FindByLabel("Name:").Fill("someone not attending")).To(Succeed())
			Expect(page.FindByLabel("Email:").Fill("not-attending@example.com")).To(Succeed())
		})
		By("and selecting the 'No' option", func() {
			Expect(page.FindByLabel("No, I will not be attending").Click()).To(Succeed())
		})
		By("and pressing 'Submit'", func() {
			Expect(page.Find("#rsvp_form").Submit()).To(Succeed())
		})
		By("showing a friendly response", func() {
			Expect(page).To(HaveURL(baseUrl + "/rsvp"))
			Expect(page.Find(".page")).To(HaveText("We're sorry you won't be making it, but thanks for letting us know!"))
		})
		By("verifying the RSVP was persisted", func() {
			Expect(localPersist.Rsvps).To(Equal([]types.Rsvp{
				types.Rsvp{
					FullName: "someone not attending",
					Email:    "not-attending@example.com",
					Decline:  true,
					Count:    0,
				},
			}))
		})
	})

	It("may be accepted", func() {
		By("Navigating to the RSVP page", func() {
			Expect(page.Navigate(baseUrl + "/rsvp")).To(Succeed())
		})

		By("entering her name and email", func() {
			Expect(page.FindByLabel("Name:").Fill("Someone attending")).To(Succeed())
			Expect(page.FindByLabel("Email:").Fill("attending@example.com")).To(Succeed())
		})
		By("and selecting the 'Yes' option", func() {
			Expect(page.FindByLabel("Yes, looking forward to it!").Click()).To(Succeed())
		})
		By("selecting how big the party is", func() {
			Expect(page.Find("#count").Select("Group of 4")).To(Succeed())
		})
		By("and pressing 'Submit'", func() {
			Expect(page.Find("#rsvp_form").Submit()).To(Succeed())
		})
		By("showing a friendly response", func() {
			Expect(page).To(HaveURL(baseUrl + "/rsvp"))
			Expect(page.Find(".page")).To(HaveText("Wonderful! We're thrilled that the four of you will be joining us. Poke around this site for more information about the wedding and our story."))
		})
		By("verifying the RSVP was persisted", func() {
			Expect(localPersist.Rsvps).To(Equal([]types.Rsvp{
				types.Rsvp{
					FullName: "Someone attending",
					Email:    "attending@example.com",
					Decline:  false,
					Count:    4,
				},
			}))
		})
	})

	Describe("dynamic UI", func() {
		It("Only displays the count label and dropdown after the user has selected 'Yes'", func() {
			Expect(page.Find("#count").Select("Group of 3")).NotTo(Succeed())
		})
	})

	Describe("Error cases", func() {
		XIt("does not allow submission when the email field is blank or malformed", func() {
			Expect(page.Navigate(baseUrl + "/rsvp")).To(Succeed())
			Expect(page.FindByLabel("Name:").Fill("Someone without an email")).To(Succeed())
			Expect(page.Find("#rsvp_form").Submit()).NotTo(Succeed())
		})
		XIt("does not allow submission when the attendance options are both unselected", func() {

		})
	})
})
