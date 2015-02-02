package server

import (
	"bytes"
	"io/ioutil"
	"text/template"
)

type PageFactory struct {
	*AssetProvider
	PageTemplateName string
	PageSpecs        []PageSpec
}

type PageSpec struct {
	AssetName, Title, Route string
}

type pageData struct {
	NavItems *[]pageData
	ID       int
	Body     string
	PageSpec
}

func (f PageFactory) AssemblePages() map[string]string {
	pageTemplate := f.loadPageTemplate()

	pages := map[string]string{}
	for _, pageData := range f.assemblePageData(f.PageSpecs) {
		var b bytes.Buffer
		err := pageTemplate.Execute(&b, pageData)
		if err != nil {
			panic(err)
		}
		pages[pageData.Route] = b.String()
	}

	return pages
}

func (f PageFactory) assemblePageData(pageSpecs []PageSpec) []pageData {
	var pages []pageData

	for _, pageSpec := range pageSpecs {
		pages = append(pages, pageData{
			PageSpec: pageSpec,
			Body:     f.loadPageBody(pageSpec.AssetName),
		})
	}
	for i := range pages {
		pages[i].ID = i
		pages[i].NavItems = &pages
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
