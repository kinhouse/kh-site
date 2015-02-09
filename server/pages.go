package server

import (
	"bytes"
	"io/ioutil"
	"text/template"
)

type AssetProviderInterface interface {
	GetAssetPath(string) string
}

type PageFactory struct {
	assetProvider AssetProviderInterface
	pageTemplate  *template.Template
	staticPages   []PageSpec
}

func NewPageFactory(
	assetProvider AssetProviderInterface,
	staticPages []PageSpec) *PageFactory {

	pageTemplate, err := template.ParseFiles(assetProvider.GetAssetPath("template.html"))
	if err != nil {
		panic(err)
	}

	return &PageFactory{
		assetProvider,
		pageTemplate,
		staticPages,
	}
}

type PageSpec struct {
	AssetName, Title, Route string
}

type PageData struct {
	NavItems *[]PageSpec
	Body     string
	PageSpec
}

func (f PageFactory) GenerateDynamicPage(title, content string) []byte {
	pageData := PageData{
		NavItems: &f.staticPages,
		Body:     content,
		PageSpec: PageSpec{
			AssetName: "",
			Title:     title,
			Route:     "~~NOROUTE~~",
		},
	}
	return f.assemblePage(pageData)
}

func (f PageFactory) assemblePage(pageData PageData) []byte {
	var b bytes.Buffer
	err := f.pageTemplate.Execute(&b, pageData)
	if err != nil {
		panic(err)
	}
	ret, err := ioutil.ReadAll(&b)
	if err != nil {
		panic(err)
	}
	return ret
}

func (f PageFactory) StaticPages() map[string][]byte {
	pages := map[string][]byte{}
	for _, pageData := range f.AssemblePageData(f.staticPages) {
		pages[pageData.Route] = f.assemblePage(pageData)
	}

	return pages
}

func (f PageFactory) AssemblePageData(pageSpecs []PageSpec) []PageData {
	var pages []PageData

	for _, pageSpec := range pageSpecs {
		pages = append(pages, PageData{
			PageSpec: pageSpec,
			Body:     f.loadPageBody(pageSpec.AssetName),
			NavItems: &f.staticPages,
		})
	}
	return pages
}

func (f PageFactory) loadPageBody(pageName string) string {
	body, err := ioutil.ReadFile(f.assetProvider.GetAssetPath(pageName + ".html"))
	if err != nil {
		panic(err)
	}
	return string(body)
}
