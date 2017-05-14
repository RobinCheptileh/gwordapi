package main

import (
	"github.com/gin-gonic/gin"
)

var router *gin.Engine

func main() {

	gin.SetMode(gin.ReleaseMode)
	// Set the router as the default one provided by Gin
	router = gin.Default()

	router.Static("/css", "templates/css")
	router.Static("/fonts", "templates/fonts")
	router.Static("/js", "templates/js")
	router.Static("/img", "templates/img")

	// Process the templates at the start so that they don't have to be loaded
	// from the disk again. This makes serving HTML pages very fast.
	router.LoadHTMLGlob("templates/index.html")

	// Define the route for the index page and display the index.html template
	// To start with, we'll use an inline route handler. Later on, we'll create
	// standalone functions that will be used as route handlers.
	router.GET("/", index)

	router.GET("/ws", func(c *gin.Context) {
		wshandler(c.Writer, c.Request)
	})

	// Start serving the application
	router.Run(":5000")

}
