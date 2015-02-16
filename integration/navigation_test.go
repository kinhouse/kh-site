package integration

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/sclevine/agouti/core"
	. "github.com/sclevine/agouti/matchers"
)

var _ = Describe("Navigation", func() {
	var page Page

	BeforeEach(func() {
		var err error

		page, err = agoutiDriver.Page()
		Expect(err).NotTo(HaveOccurred())
	})

	Describe("home page", func() {
		It("should load", func() {
			Expect(page.Navigate(baseUrl)).To(Succeed())
			Expect(page).To(HaveURL(baseUrl + "/"))
		})

		It("should include a link to the RSVP page", func() {
			Expect(page.Navigate(baseUrl)).To(Succeed())
			Expect(page.FindByLink("RSVP").Click()).To(Succeed())
			Expect(page).To(HaveURL(baseUrl + "/rsvp"))
			Expect(page).To(HaveTitle("Alana & Gabe: RSVP"))
		})

		It("should include a link to the event page", func() {
			Expect(page.Navigate(baseUrl)).To(Succeed())
			Expect(page.FindByLink("Event").Click()).To(Succeed())
			Expect(page).To(HaveURL(baseUrl + "/event"))
			Expect(page).To(HaveTitle("Alana & Gabe: Event"))
		})

		It("should include a link to the travel page", func() {
			Expect(page.Navigate(baseUrl)).To(Succeed())
			Expect(page.FindByLink("Travel").Click()).To(Succeed())
			Expect(page).To(HaveURL(baseUrl + "/travel"))
			Expect(page).To(HaveTitle("Alana & Gabe: Travel"))
		})

		It("should include a link to the gifts page", func() {
			Expect(page.Navigate(baseUrl)).To(Succeed())
			Expect(page.FindByLink("Gifts").Click()).To(Succeed())
			Expect(page).To(HaveURL(baseUrl + "/gifts"))
			Expect(page).To(HaveTitle("Alana & Gabe: Gifts"))
		})

	})

	AfterEach(func() {
		page.Destroy()
	})
})
