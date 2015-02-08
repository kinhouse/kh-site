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
	})

	AfterEach(func() {
		page.Destroy()
	})
})
