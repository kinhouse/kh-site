package integration

import (
	"github.com/kinhouse/kh-site/fakes"
	"github.com/kinhouse/kh-site/server"

	"fmt"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/sclevine/agouti/core"
	. "github.com/sclevine/agouti/matchers"
)

var _ = Describe("Navigation", func() {
	var page Page

	var baseUrl string
	const port = 5000

	BeforeEach(func() {
		var err error

		err = os.Chdir("../")
		if err != nil {
			panic("chdir failed")
		}

		persist := fakes.Persist{}
		server := server.BuildServer(persist)
		go server.Run(port)

		baseUrl = fmt.Sprintf("http://localhost:%d", port)
		WaitToBoot(baseUrl)

		page, err = agoutiDriver.Page()
		Expect(err).NotTo(HaveOccurred())
	})

	It("should load", func() {
		Expect(page.Navigate(baseUrl)).To(Succeed())
		Expect(page).To(HaveURL(baseUrl + "/"))
	})

	AfterEach(func() {
		page.Destroy()
	})
})
