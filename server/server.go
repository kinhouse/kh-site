package server

import (
	"github.com/gin-gonic/gin"

	"github.com/kinhouse/kh-site/types"

	"errors"
	"fmt"
	"os"
)

type PersistInterface interface {
	GetAllRSVPs() ([]types.Rsvp, error)
	InsertNewRSVP(types.Rsvp) (int64, error)
}

type ServerConfig struct {
	Data       PersistInterface
	AssetNames []string
	*PageFactory
}

func BuildServerConfig(persist PersistInterface) ServerConfig {
	assetsDirectory := os.Getenv("ASSETS_DIR")
	if assetsDirectory == "" {
		assetsDirectory = "assets"
	}
	fmt.Printf("\n Assets dir: %s\n", assetsDirectory)

	return ServerConfig{
		Data:       persist,
		AssetNames: []string{"main.css", "map.png", "header.png", "favicon.png"},
		PageFactory: &PageFactory{
			AssetProvider:    &AssetProvider{assetsDirectory},
			PageTemplateName: "template.html",
			PageSpecs: []PageSpec{
				PageSpec{"home", "Home", ""},
				PageSpec{"event", "Event", "event"},
				PageSpec{"rsvp", "RSVP", "rsvp"},
			},
		},
	}
}

func (s ServerConfig) AddPageRoutes(e *gin.Engine) {
	pages := s.PageFactory.AssemblePages()

	for route, pageContent := range pages {
		data := []byte(pageContent)
		e.GET("/"+route, func(c *gin.Context) {
			c.Data(200, gin.MIMEHTML, data)
		})
	}
}

func (s ServerConfig) AddAPIRoutes(e *gin.Engine) {
	e.GET("/api/v0/rsvps", func(c *gin.Context) {
		rsvps, err := s.Data.GetAllRSVPs()
		if err != nil {
			c.Fail(500, err)
			return
		}
		c.JSON(200, rsvps)
	})

	e.POST("/api/v0/rsvps", func(c *gin.Context) {
		var rsvp types.Rsvp
		if !c.Bind(&rsvp) {
			return
		}

		fmt.Printf("got rsvp post: %+v\n", rsvp)

		if rsvp.FullName == "" || rsvp.Email == "" {
			c.Fail(400, errors.New("Missing required fields"))
			return
		}

		id, err := s.Data.InsertNewRSVP(rsvp)
		if err != nil {
			c.Fail(500, err)
			return
		}

		c.JSON(201, gin.H{"id": id})
	})
}

func (s ServerConfig) AddStaticAssetRoutes(e *gin.Engine) {
	for _, name := range s.AssetNames {
		assetPath := s.GetAssetPath(name)
		e.GET("/"+name, func(c *gin.Context) { c.File(assetPath) })
	}
}

func (s ServerConfig) BuildRouter() *gin.Engine {
	r := gin.Default()

	s.AddAPIRoutes(r)
	s.AddStaticAssetRoutes(r)
	s.AddPageRoutes(r)

	return r
}

func (s ServerConfig) BuildAndRun(port int) {
	router := s.BuildRouter()
	router.Run(fmt.Sprintf(":%d", port))
}
