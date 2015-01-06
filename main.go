package main

import (
	"code.google.com/p/goauth2/oauth/jwt"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/datastore/v1beta2"
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

func datastoreTest(datastoreService *datastore.Service) string {
	return "Foo"
}

const keypair = `
{
	"private_key_id": "",
	"private_key": "",
	"client_email": "",
	"client_id": "",
	"type": "service_account"
}
`

func NewDatastoreService() (*datastore.Service, error) {
	var cloudKeys map[string]string
	err := json.Unmarshal([]byte(keypair), &cloudKeys)
	if err != nil {
		return nil, err
	}
	iss := cloudKeys["client_email"]
	scope := "https://www.googleapis.com/auth/datastore"
	key := []byte(cloudKeys["private_key"])
	token := jwt.NewToken(iss, scope, key)
	transport, err := jwt.NewTransport(token)
	if err != nil {
		return nil, err
	}
	oauthHttpClient := transport.Client()
	return datastore.New(oauthHttpClient)
}

func main() {

	datastoreService, err := NewDatastoreService()
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

	r.GET("/datastore", func(c *gin.Context) {
		c.JSON(200, gin.H{"msg": datastoreTest(datastoreService)})
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
