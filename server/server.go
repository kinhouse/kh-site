package server

import (
	"github.com/gin-gonic/gin"

	"github.com/kinhouse/kh-site/types"

	"fmt"
	"net/http"
)

type PageFactoryInterface interface {
	GenerateDynamicPage(string, string) []byte
	StaticPages() map[string][]byte
}

type PersistInterface interface {
	GetAllRSVPs() ([]types.Rsvp, error)
	InsertNewRSVP(types.Rsvp) (int64, error)
}

type ServerConfig struct {
	Data          PersistInterface
	AssetNames    []string
	RsvpHandler   func(types.Rsvp) string
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

func (s ServerConfig) AddRsvpPostHandler(e *gin.Engine) {
	e.POST("/rsvp", func(c *gin.Context) {
		var rsvp types.Rsvp
		if !c.Bind(&rsvp) {
			return
		}
		fmt.Printf("Got an RSVP: %+v\n", rsvp)
		responseText := s.RsvpHandler(rsvp)
		responseTitle := "♡"
		if rsvp.Decline {
			responseTitle = "☹"
		}
		responseHtml := s.PageFactory.GenerateDynamicPage(responseTitle, responseText)
		c.Data(http.StatusCreated, gin.MIMEHTML, responseHtml)
	})
}

func (s ServerConfig) BuildRouter() *gin.Engine {
	r := gin.Default()

	s.AddStaticAssetRoutes(r)
	s.AddPageRoutes(r)
	s.AddRsvpPostHandler(r)

	return r
}
