package main

import (
	"bytes"
	"croupier.io/cards"
	"croupier.io/decks"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type CreateResponse struct {
	DeckID    uuid.UUID `json:"deck_id"`
	Shuffled  bool      `json:"shuffled"`
	Remaining int       `json:"remaining"`
}

type DrawCardResponse struct {
	Cards []cards.PlayingCard `json:"cards"`
}

func TestCreateDeck(t *testing.T) {
	router := NewRouter()

	statusCode, creationResponse := requestCreateDeck(t, router, "", nil)
	assert.Equal(t, http.StatusCreated, statusCode)
	assert.NotEqual(t, creationResponse.DeckID, uuid.Nil, "expected a non-empty decks ID")
	assert.NotZero(t, creationResponse.Remaining, "expected a non-empty decks upon creation")
}

func TestCreateDeckWithInvalidBodyRequest(t *testing.T) {
	router := NewRouter()

	statusCode, _ := requestCreateDeck(t, router, "", "invalid body")
	assert.Equal(t, http.StatusBadRequest, statusCode)
}

func TestCreateDeckWithUnhandledType(t *testing.T) {
	router := NewRouter()
	request := decks.CreationRequest{
		PlayingType: cards.PlayingCardType(99999),
	}

	statusCode, _ := requestCreateDeck(t, router, "", &request)
	assert.Equal(t, http.StatusInternalServerError, statusCode)
}

func TestCreateShuffledDeck(t *testing.T) {
	router := NewRouter()
	request := decks.CreationRequest{
		Shuffled: true,
	}

	statusCode, creationResponse := requestCreateDeck(t, router, "", &request)
	assert.NotEqual(t, creationResponse.DeckID, uuid.Nil, "expected a non-empty decks ID")
	assert.NotZero(t, creationResponse.Remaining, "expected a non-empty decks upon creation")
	assert.Equal(t, http.StatusCreated, statusCode)
}

func TestCreateCustomDeck(t *testing.T) {
	router := NewRouter()
	requestedCardCodes := []string{"AS", "KD", "AC", "2C", "KH"}

	statusCode, sortedDeck := requestCreateDeck(t, router, "?cards="+strings.Join(requestedCardCodes, ","), nil)
	assert.Equal(t, http.StatusCreated, statusCode)
	assert.NotEqual(t, sortedDeck.DeckID, uuid.Nil, "expected a non-empty decks ID")
	assert.Equal(t, 5, sortedDeck.Remaining, "expected a non-empty decks upon creation")
}

func TestOpenDeck(t *testing.T) {
	router := NewRouter()

	_, creationResponse := requestCreateDeck(t, router, "", nil)

	statusCode, actualPlayingDeck := requestOpenDeck(t, router, creationResponse.DeckID.String())
	assert.NotNil(t, actualPlayingDeck, "expected a decks")
	assert.Equal(t, creationResponse.DeckID, actualPlayingDeck.ID)
	assert.Equal(t, http.StatusOK, statusCode)
}

func TestOpenUnknownDeck(t *testing.T) {
	router := NewRouter()

	statusCode, _ := requestOpenDeck(t, router, "unknown_id")
	assert.Equal(t, http.StatusNotFound, statusCode)
}

func TestDrawCard(t *testing.T) {
	router := NewRouter()

	requestedCardCodes := []string{"AS", "2S", "3S"}
	aceOfSpades := cards.PlayingCard{Suit: cards.Spades.String(), Value: "ACE", Code: "AS"}
	twoOfSpades := cards.PlayingCard{Suit: cards.Spades.String(), Value: "2", Code: "2S"}
	threeOfSpades := cards.PlayingCard{Suit: cards.Spades.String(), Value: "3", Code: "3S"}

	_, creationResponse := requestCreateDeck(t, router, "?cards="+strings.Join(requestedCardCodes, ","), nil)
	_, playingDeck := requestOpenDeck(t, router, creationResponse.DeckID.String())
	assert.Equal(
		t,
		[]cards.PlayingCard{aceOfSpades, twoOfSpades, threeOfSpades},
		playingDeck.Cards)

	statusCode, drawCardResponse := requestDrawCard(t, router, playingDeck.ID.String(), "?count=1")
	assert.Equal(t, []cards.PlayingCard{aceOfSpades}, drawCardResponse.Cards)
	assert.Equal(t, http.StatusOK, statusCode)

	_, playingDeck = requestOpenDeck(t, router, creationResponse.DeckID.String())
	assert.Equal(
		t,
		[]cards.PlayingCard{twoOfSpades, threeOfSpades},
		playingDeck.Cards)
}

func TestDrawCardFromUnknownDeck(t *testing.T) {
	router := NewRouter()

	statusCode, _ := requestDrawCard(t, router, "2", "?count=1")
	assert.Equal(t, http.StatusNotFound, statusCode)
}

func TestDrawCardFromDeckWithInvalidCount(t *testing.T) {
	router := NewRouter()

	requestedDrawCardCount := []string{"", "a12"}
	for _, requestedDrawCardCount := range requestedDrawCardCount {
		_, creationResponse := requestCreateDeck(t, router, "", nil)
		statusCode, _ := requestDrawCard(t, router, creationResponse.DeckID.String(), "?count="+requestedDrawCardCount)
		assert.Equal(t, http.StatusBadRequest, statusCode)
	}
}

func requestCreateDeck(t *testing.T, router *gin.Engine, queryParameters string, deckCreationRequest interface{}) (int, CreateResponse) {
	byteBody, err := json.Marshal(deckCreationRequest)
	if err != nil {
		t.Fail()
	}
	responseWriter := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/decks"+queryParameters, bytes.NewBuffer(byteBody))
	router.ServeHTTP(responseWriter, request)

	var response CreateResponse
	if err := json.Unmarshal(responseWriter.Body.Bytes(), &response); err != nil {
		t.Fail()
	}
	return responseWriter.Code, response
}

func requestOpenDeck(t *testing.T, router *gin.Engine, id string) (int, decks.PlayableDeck) {
	responseWriter := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/decks/"+id, nil)
	router.ServeHTTP(responseWriter, request)

	var actualPlayingDeck decks.PlayableDeck
	if err := json.Unmarshal(responseWriter.Body.Bytes(), &actualPlayingDeck); err != nil {
		t.Fail()
	}
	return responseWriter.Code, actualPlayingDeck
}

func requestDrawCard(t *testing.T, router *gin.Engine, id string, queryParameters string) (int, DrawCardResponse) {
	responseWriter := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", fmt.Sprintf("/decks/%s/cards/draw"+queryParameters, id), nil)
	router.ServeHTTP(responseWriter, request)

	var drawCardResponse DrawCardResponse
	if err := json.Unmarshal(responseWriter.Body.Bytes(), &drawCardResponse); err != nil {
		t.Fail()
	}
	return responseWriter.Code, drawCardResponse
}
