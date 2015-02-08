package server_test

import (
	"io/ioutil"
	"os"
	"path"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/kinhouse/kh-site/server"
)

type FakeAssetProvider struct {
	AssetProvider
}

func NewFakeAssetProvider() *FakeAssetProvider {
	dir, err := ioutil.TempDir("", "fake-assets")
	if err != nil {
		panic(err)
	}
	return &FakeAssetProvider{
		AssetProvider: AssetProvider{dir},
	}
}

func (p FakeAssetProvider) Delete() {
	os.RemoveAll(p.AssetsDirectory)
}

func (p FakeAssetProvider) AddFile(name, contents string) {
	ioutil.WriteFile(path.Join(p.AssetsDirectory, name), []byte(contents), os.ModePerm)
}

var _ = Describe("Pages", func() {
	var pageSpecs []PageSpec
	var pageFactory PageFactory

	BeforeEach(func() {
		pageSpecs = []PageSpec{
			PageSpec{
				AssetName: "home",
				Title:     "Home Page",
				Route:     RootRoute,
			},
			PageSpec{
				AssetName: "event",
				Title:     "Event Info",
				Route:     "event/info/for/now",
			},
			PageSpec{
				AssetName: "post_response_page",
				Title:     "Thank you for your submission",
				Route:     NoRoute,
			},
		}

		assetProvider := NewFakeAssetProvider()
		assetProvider.AddFile("home.html", "some html")
		assetProvider.AddFile("event.html", "some html")
		assetProvider.AddFile("post_response_page.html", "some html")

		pageFactory = PageFactory{
			PageTemplateName:       "template.html",
			PageSpecs:              pageSpecs,
			AssetProviderInterface: assetProvider,
		}
	})
	AfterEach(func() {
		pageFactory.AssetProviderInterface.(*FakeAssetProvider).Delete()
	})

	Describe("assembling page data", func() {
		It("returns a PageData for every PageSpec", func() {
			pageDatas := pageFactory.AssemblePageData(pageSpecs)
			Expect(pageDatas).To(HaveLen(3))
			Expect(pageDatas[0].PageSpec).To(Equal(pageSpecs[0]))
			Expect(pageDatas[1].PageSpec).To(Equal(pageSpecs[1]))
			Expect(pageDatas[2].PageSpec).To(Equal(pageSpecs[2]))
		})

		It("includes a page as a NavItem if and only if the page has a valid route", func() {
			navItems := *(pageFactory.AssemblePageData(pageSpecs)[0].NavItems)
			Expect(navItems).To(HaveLen(2))
			Expect(navItems[0]).To(Equal(pageSpecs[0]))
			Expect(navItems[1]).To(Equal(pageSpecs[1]))
		})

		It("uses puts the same NavItems on every page", func() {
			pageDatas := pageFactory.AssemblePageData(pageSpecs)
			Expect(pageDatas[0].NavItems).To(Equal(pageDatas[1].NavItems))
			Expect(pageDatas[0].NavItems).To(Equal(pageDatas[2].NavItems))
		})
	})
})
