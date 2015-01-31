package server

import (
	"github.com/gin-gonic/gin"

	"github.com/kinhouse/kh-site/types"

	"fmt"
	"html/template"
	"io/ioutil"
	"os"
)

type PersistInterface interface {
	GetAllRSVPs() ([]types.Rsvp, error)
	InsertNewRSVP(types.Rsvp) (int64, error)
}

type Server struct {
	Data         PersistInterface
	PageTemplate *template.Template
	Pages        map[string]template.HTML
	Assets       []string
}

func BuildServer(persist PersistInterface) Server {
	pageTemplate, err := template.ParseFiles("assets/template.html")
	if err != nil {
		panic(err)
	}

	locationPage := loadPage("event")
	homePage := loadPage("home")

	return Server{
		Data:         persist,
		PageTemplate: pageTemplate,
		Pages: map[string]template.HTML{
			"":      homePage,
			"rsvp":  template.HTML("Coming soon..."),
			"event": locationPage,
		},
		Assets: []string{"main.css", "map.png", "header.png", "favicon.png"},
	}
}

func loadPage(title string) template.HTML {
	path := "assets/" + title + ".html"
	body, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return template.HTML(body)
}

func (s Server) createRoutes() *gin.Engine {
	r := gin.Default()

	r.GET("/api/v0/rsvps", func(c *gin.Context) {
		rsvps, err := s.Data.GetAllRSVPs()
		if err != nil {
			c.Fail(500, err)
		}
		c.JSON(200, rsvps)
	})

	r.POST("/api/v0/rsvps", func(c *gin.Context) {
		var rsvp types.Rsvp
		if c.Bind(&rsvp) {
			fmt.Printf("got rsvp post: %+v\n", rsvp)
			id, err := s.Data.InsertNewRSVP(rsvp)
			if err != nil {
				c.Fail(500, err)
			}
			c.JSON(201, gin.H{"id": id})
		}
	})

	for route, body := range s.Pages {
		pageBody := body
		r.GET("/"+route, func(c *gin.Context) {
			s.PageTemplate.Execute(c.Writer, pageBody)
		})
	}

	for _, filename := range s.Assets {
		asset := filename
		r.GET("/"+asset, func(c *gin.Context) { c.File("assets/" + asset) })
	}

	return r
}

func (s Server) Run() {
	routerGroup := s.createRoutes()
	routerGroup.Run(":" + os.Getenv("PORT"))
}
