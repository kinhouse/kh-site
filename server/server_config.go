package server

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func BuildServer(persist PersistInterface) *gin.Engine {
	assetProvider := &AssetProvider{getAssetsDirectory()}

	pageSpecs := []PageSpec{
		PageSpec{AssetName: "home", Title: "Home", Route: ""},
		PageSpec{AssetName: "event", Title: "Event", Route: "event"},
		PageSpec{AssetName: "rsvp", Title: "RSVP", Route: "rsvp"},
	}

	pageFactory := NewPageFactory(assetProvider, pageSpecs)

	assetNames := []string{"main.css", "map.png", "header.png", "favicon.png"}

	serverConfig := ServerConfig{
		Data:          persist,
		AssetNames:    assetNames,
		PageFactory:   pageFactory,
		RsvpHandler:   RsvpHandler,
		AssetProvider: assetProvider,
	}

	router := serverConfig.BuildRouter()
	return router
}

func getAssetsDirectory() string {
	assetsDirectory := os.Getenv("ASSETS_DIR")
	if assetsDirectory == "" {
		assetsDirectory = "assets"
	}
	path, err := filepath.Abs(assetsDirectory)
	if err != nil {
		panic("could not expand assets directory path: " + err.Error())
	}

	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			panic("Didn't find an assets directory at " + path)
		} else {
			panic("Could not stat assets directory: " + path)
		}
	}

	fmt.Printf("\n Looking for assets in: %s\n", path)
	return path

}
