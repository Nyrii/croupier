package main

import (
	"croupier.io/decks"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"strings"
)

var playingDecks []*decks.PlayableDeck

// createDeck creates and stores a PlayableDeck.
func createDeck(context *gin.Context) {
	var request decks.CreationRequest
	if err := context.BindJSON(&request); err != nil {
		log.Printf("Failed to get the decks creation request: %s", err)
		context.JSON(http.StatusBadRequest, gin.H{"message": "unable to generate the deck"})
		return
	}
	requestedCards := strings.Split(context.Query("cards"), ",")

	playingDeck, err := decks.CreateDeck(request, requestedCards)
	if err != nil {
		log.Printf("Failed to create the decks: %s", err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "unable to generate the deck"})
		return
	}
	playingDecks = append(playingDecks, playingDeck)
	context.JSON(
		http.StatusCreated,
		gin.H{
			"deck_id":   playingDeck.ID,
			"shuffled":  playingDeck.Shuffled,
			"remaining": playingDeck.Remaining,
		})
}

// openDeck finds a PlayableDeck associated with a provided ID, if any.
func openDeck(context *gin.Context) {
	id := context.Param("id")
	playingDeck := findDeck(id)
	if playingDeck == nil {
		context.JSON(http.StatusNotFound, gin.H{"message": "unable to find the deck"})
		return
	}
	context.JSON(http.StatusOK, playingDeck)
}

// drawCard draws cards from a PlayableDeck associated with a provided ID, if applicable.
func drawCard(context *gin.Context) {
	id := context.Param("id")
	requestedDrawCardCount, err := strconv.Atoi(context.Query("count"))
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "unable to find the requested number of cards to draw"})
		return
	}
	playingDeck := findDeck(id)
	if playingDeck == nil {
		context.JSON(http.StatusNotFound, gin.H{"message": "unable to find the deck"})
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"cards": playingDeck.DrawCard(requestedDrawCardCount),
	})
}

// findDeck finds a PlayableDeck associated with id, if any.
// findDeck returns nil if no PlayableDeck is associated with the provided id.
func findDeck(id string) *decks.PlayableDeck {
	for _, playingDeck := range playingDecks {
		if id == playingDeck.ID.String() {
			return playingDeck
		}
	}
	return nil
}
