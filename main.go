package main

import (
	"github.com/gin-gonic/gin"
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

	pageTemplate, err := template.ParseFiles("assets/template.html")
	if err != nil {
		panic(err)
	}

	locationPage := loadPage("event")
	homePage := loadPage("home")

	r := gin.Default()

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
