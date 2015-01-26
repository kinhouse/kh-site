package main

import (
	"code.google.com/p/goauth2/oauth/jwt"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
	"google.golang.org/cloud"
	"google.golang.org/cloud/datastore"

	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
)

const (
	CloudKeyEnvVar = "GOOGLE_CLOUD_KEY"
	ProjectId      = "kh-site"
)

func loadPage(title string) template.HTML {
	path := "assets/" + title + ".html"
	body, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return template.HTML(body)
}

func loadGoogleCloudKeys() (map[string]string, error) {
	var cloudKeys map[string]string

	jsonString := os.Getenv(CloudKeyEnvVar)
	err := json.Unmarshal([]byte(jsonString), &cloudKeys)
	if err != nil {
		fmt.Printf(err.Error())
		fmt.Printf(jsonString)
		return cloudKeys, errors.New(fmt.Sprintf("error reading cloud key from env var %q", CloudKeyEnvVar))
	}
	return cloudKeys, nil
}

func createGoogleCloudContext() (context.Context, error) {
	cloudKeys, err := loadGoogleCloudKeys()
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
	return cloud.NewContext(ProjectId, oauthHttpClient), nil
}

func datastoreTest(cloudContext context.Context) string {
	type Rsvp struct {
		name   string
		email  string
		guests int `datastore:",noindex"`
	}

	key := datastore.NewIncompleteKey(cloudContext, "rsvp", nil)
	key, err := datastore.Put(cloudContext, key, &Rsvp{
		name:   "Some Person",
		email:  "example@example.com",
		guests: 2,
	})

	if err != nil {
		panic(err)
	}

	return key.Name()
}

func main() {

	cloudContext, err := createGoogleCloudContext()
	if err != nil {
		panic(err)
	}

	fmt.Printf(datastoreTest(cloudContext))

	pageTemplate, err := template.ParseFiles("assets/template.html")
	if err != nil {
		panic(err)
	}

	locationPage := loadPage("event")
	homePage := loadPage("home")

	r := gin.Default()

	r.GET("/datastore", func(c *gin.Context) {
		c.JSON(200, gin.H{"msg": datastoreTest(cloudContext)})
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
