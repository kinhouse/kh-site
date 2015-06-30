package server

import (
	"github.com/gin-gonic/gin"

	"net/http"
)

type PageFactoryInterface interface {
	GenerateDynamicPage(string, string) []byte
	StaticPages() map[string][]byte
}

type ServerConfig struct {
	AssetNames    []string
	PageFactory   PageFactoryInterface
	AssetProvider AssetProviderInterface
}

func (s ServerConfig) AddPageRoutes(e *gin.Engine) {
	for route, pageContent := range s.PageFactory.StaticPages() {
		d := pageContent // copy variable, for closure
		e.GET("/"+route, func(c *gin.Context) {
			c.Data(http.StatusOK, gin.MIMEHTML, d)
		})
	}
}

func (s ServerConfig) AddStaticAssetRoutes(e *gin.Engine) {
	for _, name := range s.AssetNames {
		assetPath := s.AssetProvider.GetAssetPath(name)
		e.GET("/"+name, func(c *gin.Context) { c.File(assetPath) })
	}
}

func (s ServerConfig) AddRedirects(e *gin.Engine) {
	for _, s := range []string{
		"rsvp", "us", "event", "traditions",
		"travel", "explore", "gifts", "blessings"} {
		route := "/" + s
		e.GET(route, func(c *gin.Context) {
			c.Redirect(http.StatusMovedPermanently, "/")
		})
	}
}

func (s ServerConfig) BuildRouter() *gin.Engine {
	r := gin.Default()

	s.AddStaticAssetRoutes(r)
	s.AddPageRoutes(r)
	s.AddRedirects(r)

	return r
}
