package integration

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sclevine/agouti"
	. "github.com/sclevine/agouti/matchers"
)

var _ = Describe("Navigation", func() {
	var page *agouti.Page

	BeforeEach(func() {
		var err error

		page, err = agoutiDriver.NewPage()
		Expect(err).NotTo(HaveOccurred())
	})

	Describe("home page", func() {
		It("should load", func() {
			Expect(page.Navigate(baseUrl)).To(Succeed())
			Expect(page).To(HaveURL(baseUrl + "/"))
		})

		It("should not include a link to the RSVP page", func() {
			Expect(page.Navigate(baseUrl)).To(Succeed())
			Expect(page.FindByLink("RSVP").Click()).NotTo(Succeed())
		})

		It("should redirect all requests to /rsvp to the home page", func() {
			Expect(page.Navigate(baseUrl + "/rsvp")).To(Succeed())
			Expect(page).To(HaveURL(baseUrl + "/"))
		})
	})

	AfterEach(func() {
		page.Destroy()
	})
})
