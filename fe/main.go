package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"

	"golang.org/x/oauth2"
)

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.GET("/", Homepage)
	r.GET("/callbacks", Callback)
	r.Run(":8080")
}

var Endpoint = oauth2.Endpoint{
	AuthURL:  "http://localhost:4444/oauth2/auth",
	TokenURL: "http://localhost:4444/oauth2/token",
}

var OAuthConf = &oauth2.Config{
	RedirectURL:  "http://localhost:8080/callbacks",
	ClientID:     "loc_client", // TODO from hydra
	ClientSecret: "loc_secret", // TODO from hydra

	Scopes:   []string{"users.write", "users.read", "users.edit", "users.delete", "offline"},
	Endpoint: Endpoint,
}

func Homepage(c *gin.Context) {
	loginURL := OAuthConf.AuthCodeURL("1234567890")
	c.HTML(http.StatusOK, "homepage.html", gin.H{
		"loginURL": loginURL,
	})
}

func Callback(c *gin.Context) {
	ctx := c.Request.Context()

	code := c.Query("code")
	state := c.Query("state")
	fmt.Println("code", code)
	fmt.Println("state", state)

	accessToken, err := OAuthConf.Exchange(ctx, code)
	if err != nil {
		c.JSON(http.StatusOK, err.Error())
	}

	c.HTML(http.StatusOK, "welcome.html", gin.H{
		"accessToken": accessToken,
	})
}
