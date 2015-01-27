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

	r := gin.Default()

	r.GET("/api/v0/rsvps", func(c *gin.Context) {
		rsvps, err := data.GetAllRSVPs()
		if err != nil {
			c.Fail(500, err)
		}
		c.JSON(200, rsvps)
	})

	r.POST("/api/v0/rsvps", func(c *gin.Context) {
		var rsvp persist.Rsvp
		if c.Bind(&rsvp) {
			fmt.Printf("got rsvp post: %+v\n", rsvp)
			id, err := data.InsertNewRSVP(rsvp)
			if err != nil {
				c.Fail(500, err)
			}
			c.JSON(201, gin.H{"id": id})
		}
	})

	r.GET("/event", func(c *gin.Context) {
		pageTemplate.Execute(c.Writer, locationPage)
	})

	r.GET("/", func(c *gin.Context) {
		pageTemplate.Execute(c.Writer, homePage)
	})

	r.GET("/rsvp", func(c *gin.Context) {
		pageTemplate.Execute(c.Writer, template.HTML("Coming soon..."))
	})

	staticAssets := []string{"main.css", "map.png", "header.png", "favicon.png"}
	for _, filename := range staticAssets {
		asset := filename
		r.GET("/"+asset, func(c *gin.Context) {
			c.File("assets/" + asset)
		})
	}

	r.Run(":" + os.Getenv("PORT"))
}
