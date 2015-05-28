package fakes

import (
	"io/ioutil"
	"os"
	"path"

	"github.com/kinhouse/kh-site/server"
	"github.com/kinhouse/kh-site/types"
)

//
// Persist
//
type Persist struct {
	Rsvps []types.Rsvp
}

func (p *Persist) GetAllRSVPs() ([]types.Rsvp, error) {
	return p.Rsvps, nil
}

//
// Asset Provider
//
type AssetProvider struct {
	server.AssetProvider
}

func NewAssetProvider() *AssetProvider {
	dir, err := ioutil.TempDir("", "fake-assets")
	if err != nil {
		panic(err)
	}
	ioutil.WriteFile(
		path.Join(dir, "template.html"),
		[]byte("<html><body>This is a template!: {{.}} </body></html>"),
		os.ModePerm)

	return &AssetProvider{
		AssetProvider: server.AssetProvider{AssetsDirectory: dir},
	}
}

func (p AssetProvider) Delete() {
	os.RemoveAll(p.AssetsDirectory)
}

func (p AssetProvider) AddFile(name, contents string) {
	ioutil.WriteFile(path.Join(p.AssetsDirectory, name), []byte(contents), os.ModePerm)
}

type PageFactory struct {
	FD_StaticPages map[string][]byte
}

func (f *PageFactory) GenerateDynamicPage(title, content string) []byte {
	return []byte("<html><head><title>" + title + "</title></head><body>" + content + "</body></html>")
}
func (f *PageFactory) StaticPages() map[string][]byte {
	return f.FD_StaticPages
}
