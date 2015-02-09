package server_test

import (
	"io/ioutil"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/kinhouse/kh-site/fakes"

	. "github.com/kinhouse/kh-site/server"
)

var _ = Describe("Pages", func() {
	var pageSpecs []PageSpec
	var pageFactory *PageFactory
	var assetProvider *fakes.AssetProvider

	BeforeEach(func() {
		pageSpecs = []PageSpec{
			PageSpec{
				AssetName: "home",
				Title:     "Home Page",
				Route:     "",
			},
			PageSpec{
				AssetName: "event",
				Title:     "Event Info",
				Route:     "event/info/for/now",
			},
		}

		assetProvider = fakes.NewAssetProvider()
		assetProvider.AddFile("home.html", "some html")
		assetProvider.AddFile("event.html", "some html")

		pageFactory = NewPageFactory(assetProvider, pageSpecs)
	})
	AfterEach(func() {
		assetProvider.Delete()
	})

	Describe("assembling static page data", func() {
		It("returns a PageData for every PageSpec", func() {
			pageDatas := pageFactory.AssemblePageData(pageSpecs)
			Expect(pageDatas).To(HaveLen(2))
			Expect(pageDatas[0].PageSpec).To(Equal(pageSpecs[0]))
			Expect(pageDatas[1].PageSpec).To(Equal(pageSpecs[1]))
		})

		It("uses puts the same NavItems on every page", func() {
			pageDatas := pageFactory.AssemblePageData(pageSpecs)
			Expect(pageDatas[0].NavItems).To(Equal(pageDatas[1].NavItems))
		})
	})

	Describe("generating a dynamic page", func() {
		It("produces a page that contains the input string", func() {
			pageBytes := pageFactory.GenerateDynamicPage("title", "content!")

			Expect(string(pageBytes)).To(ContainSubstring("content!"))
		})

		It("is based on the standard page template", func() {
			content := "we're excited you're attending!"

			pageBytes := pageFactory.GenerateDynamicPage("title", content)

			templateBytes, err := ioutil.ReadFile(assetProvider.GetAssetPath("template.html"))
			if err != nil {
				panic(err)
			}
			N := 10
			Expect(pageBytes[:N]).To(Equal(templateBytes[:N]))
		})

	})
})
