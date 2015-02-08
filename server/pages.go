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
	AssetProviderInterface
	PageTemplateName string
	PageSpecs        []PageSpec
}

const (
	RootRoute = ""
	NoRoute   = "  ~~ NO ROUTE ~~ "
)

type PageSpec struct {
	AssetName, Title, Route string
}

type PageData struct {
	NavItems *[]PageSpec
	Body     string
	PageSpec
}

func (f PageFactory) AssemblePages() map[string]string {
	pageTemplate := f.loadPageTemplate()

	pages := map[string]string{}
	for _, pageData := range f.AssemblePageData(f.PageSpecs) {
		var b bytes.Buffer
		err := pageTemplate.Execute(&b, pageData)
		if err != nil {
			panic(err)
		}
		pages[pageData.Route] = b.String()
	}

	return pages
}

func (f PageFactory) AssemblePageData(pageSpecs []PageSpec) []PageData {
	var pages []PageData
	var navItems []PageSpec

	for _, pageSpec := range pageSpecs {
		pages = append(pages, PageData{
			PageSpec: pageSpec,
			Body:     f.loadPageBody(pageSpec.AssetName),
		})
		if pageSpec.Route != NoRoute {
			navItems = append(navItems, pageSpec)
		}
	}
	for i := range pages {
		pages[i].NavItems = &navItems
	}
	return pages
}

func (f PageFactory) loadPageTemplate() *template.Template {
	pageTemplate, err := template.ParseFiles(f.GetAssetPath(f.PageTemplateName))
	if err != nil {
		panic(err)
	}

	return pageTemplate
}

func (f PageFactory) loadPageBody(pageName string) string {
	body, err := ioutil.ReadFile(f.GetAssetPath(pageName + ".html"))
	if err != nil {
		panic(err)
	}
	return string(body)
}
