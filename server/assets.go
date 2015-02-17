package server

import (
	"io/ioutil"
	"path"
	"strings"
)

type AssetProvider struct {
	AssetsDirectory string
}

func (a AssetProvider) GetAssetPath(assetName string) string {
	return path.Join(a.AssetsDirectory, assetName)
}

func (a AssetProvider) ListAllNonHTML() []string {
	files, err := ioutil.ReadDir(a.AssetsDirectory)
	if err != nil {
		panic(err)
	}

	var ret []string
	for _, fi := range files {
		name := fi.Name()
		if !strings.HasSuffix(name, ".html") {
			ret = append(ret, name)
		}
	}
	return ret
}
