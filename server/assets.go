package server

import "path"

type AssetProvider struct {
	AssetsDirectory string
}

func (a AssetProvider) GetAssetPath(assetName string) string {
	return path.Join(a.AssetsDirectory, assetName)
}
