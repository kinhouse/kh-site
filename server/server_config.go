package server

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func BuildServer(persist PersistInterface, adminPassword string) *gin.Engine {
	assetProvider := &AssetProvider{getAssetsDirectory()}

	pageSpecs := []PageSpec{
		PageSpec{AssetName: "home", Title: "Home", Route: ""},
		PageSpec{AssetName: "us", Title: "Us", Route: "us"},
		PageSpec{AssetName: "event", Title: "Event", Route: "event"},
		PageSpec{AssetName: "travel", Title: "Travel", Route: "travel"},
		PageSpec{AssetName: "explore", Title: "Explore", Route: "explore"},
		PageSpec{AssetName: "gifts", Title: "Gifts", Route: "gifts"},
		PageSpec{AssetName: "rsvp", Title: "RSVP", Route: "rsvp"},
	}

	pageFactory := NewPageFactory(assetProvider, pageSpecs)

	assetNames := assetProvider.ListAllNonHTML()

	serverConfig := ServerConfig{
		Data:                persist,
		AssetNames:          assetNames,
		PageFactory:         pageFactory,
		RsvpHandler:         RsvpHandler,
		RsvpValidator:       ValidateRsvp,
		AssetProvider:       assetProvider,
		RsvpListCredentials: map[string]string{"admin": adminPassword},
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
