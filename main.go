package main

import (
	"github.com/gin-gonic/gin"

	"github.com/kinhouse/kh-site/persist"

	"fmt"
	"html/template"
	"io/ioutil"
	"os"
)

func loadPage(title string) template.HTML {
	path := "assets/" + title + ".html"
	body, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return template.HTML(body)
}

type PersistInterface interface {
	GetAllRSVPs() ([]persist.Rsvp, error)
	InsertNewRSVP(persist.Rsvp) (int64, error)
}

type Config struct {
	data         PersistInterface
	pageTemplate *template.Template
	pages        map[string]template.HTML
	assets       []string
}

func BuildRouterGroup(config Config) *gin.Engine {
	r := gin.Default()

	r.GET("/api/v0/rsvps", func(c *gin.Context) {
		rsvps, err := config.data.GetAllRSVPs()
		if err != nil {
			c.Fail(500, err)
		}
		c.JSON(200, rsvps)
	})

	r.POST("/api/v0/rsvps", func(c *gin.Context) {
		var rsvp persist.Rsvp
		if c.Bind(&rsvp) {
			fmt.Printf("got rsvp post: %+v\n", rsvp)
			id, err := config.data.InsertNewRSVP(rsvp)
			if err != nil {
				c.Fail(500, err)
			}
			c.JSON(201, gin.H{"id": id})
		}
	})

	for route, body := range config.pages {
		pageBody := body
		r.GET("/"+route, func(c *gin.Context) {
			config.pageTemplate.Execute(c.Writer, pageBody)
		})
	}

	for _, filename := range config.assets {
		asset := filename
		r.GET("/"+asset, func(c *gin.Context) { c.File("assets/" + asset) })
	}

	return r
}

func main() {

	data, err := persist.NewPersist()
	if err != nil {
		panic(err)
	}

	pageTemplate, err := template.ParseFiles("assets/template.html")
	if err != nil {
		panic(err)
	}

	locationPage := loadPage("event")
	homePage := loadPage("home")

	config := Config{
		data:         data,
		pageTemplate: pageTemplate,
		pages: map[string]template.HTML{
			"":      homePage,
			"rsvp":  template.HTML("Coming soon..."),
			"event": locationPage,
		},
		assets: []string{"main.css", "map.png", "header.png", "favicon.png"},
	}

	routerGroup := BuildRouterGroup(config)
	routerGroup.Run(":" + os.Getenv("PORT"))
}
