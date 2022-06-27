package main

import (
	"github.com/gin-gonic/gin"
)

// NewRouter adds all the routes and route handlers necessary for the API and returns a router.
func NewRouter() *gin.Engine {
	router := gin.Default()

	AddDeckApi(router)

	return router
}

// AddDeckApi attaches the routes and route handlers associated with decks.
func AddDeckApi(router *gin.Engine) {
	deckApi := router.Group("/decks")
	{
		deckApi.POST("", createDeck)
		deckApi.GET("/:id", openDeck)
		deckApi.POST("/:id/cards/draw", drawCard)
	}
}
