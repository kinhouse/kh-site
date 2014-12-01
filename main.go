package main

import (
	"github.com/gin-gonic/gin"
	"os"
)

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
	r.Static("/assets", "assets")
	r.GET("/", func(c *gin.Context) {
		c.File("assets/index.html")
	})
	r.Run(":" + os.Getenv("PORT"))
}
