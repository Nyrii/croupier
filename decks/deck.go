// Package decks implements the models used to trigger actions for manipulating a deck.
package decks

import (
	"croupier.io/cards"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"math/rand"
	"time"
)

// Deck is the interface that wraps the methods applicable from a deck.
//
// Shuffle shuffles the cards contained in a deck.
// Shuffle must modify the cards contained in the deck.
//
// DrawCard pulls a specific number of cards from the cards contained in a deck, if any.
// The cards that are drawn must be removed from the deck and must be returned.
type Deck interface {
	Shuffle()
	DrawCard(int) []cards.PlayingCard
}

// PlayableDeck is the representation of a deck entity.
// PlayableDeck should be defined in any specific type of deck.
type PlayableDeck struct {
	ID        uuid.UUID           `json:"deck_id"`
	Cards     []cards.PlayingCard `json:"cards"`
	Shuffled  bool                `json:"shuffled"`
	Remaining int                 `json:"remaining"`
}

// CreationRequest is the representation of a request used to create a PlayableDeck.
type CreationRequest struct {
	PlayingType cards.PlayingCardType `json:"type"`
	Shuffled    bool                  `json:"shuffled"`
}

var _ Deck = &PlayableDeck{}

// Shuffle shuffles the cards contained in a playable deck.
// Shuffle sets Shuffled to true.
func (deck *PlayableDeck) Shuffle() {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(deck.Cards), func(i, j int) { deck.Cards[i], deck.Cards[j] = deck.Cards[j], deck.Cards[i] })
	deck.Shuffled = true
}

// DrawCard pulls a specific number of cards from the cards contained in a deck, if any.
// The cards that are drawn are removed from the deck and are returned.
// DrawCard keeps track of Remaining and sets it to the number of cards which remained in the
// deck after the draw.
func (deck *PlayableDeck) DrawCard(requestedDrawCardCount int) []cards.PlayingCard {
	if requestedDrawCardCount <= 0 || len(deck.Cards) == 0 {
		return make([]cards.PlayingCard, 0)
	}
	if requestedDrawCardCount > len(deck.Cards) {
		requestedDrawCardCount = len(deck.Cards)
	}
	var playingCard cards.PlayingCard
	var playingCards []cards.PlayingCard
	for i := 0; i < requestedDrawCardCount; i++ {
		playingCard, deck.Cards = deck.Cards[0], deck.Cards[1:]
		playingCards = append(playingCards, playingCard)
		deck.Remaining -= 1
	}
	return playingCards
}

// CreateDeck creates a PlayableDeck based on the provided creationRequest and requestedCardCodes.
// If requestedCardCodes is empty, a common PlayableDeck is created, according to the type of deck
// standards.
// CreateDeck can fail to create a PlayableDeck if the requested type is not handled.
func CreateDeck(creationRequest CreationRequest, requestedCardCodes []string) (*PlayableDeck, error) {
	var playingDeck PlayableDeck
	switch creationRequest.PlayingType {
	case cards.French:
		deck, err := NewFrenchDeck(requestedCardCodes)
		if err != nil {
			return nil, err
		}
		playingDeck = deck.PlayableDeck
	default:
		return nil, errors.New(fmt.Sprintf("unsupported operation for cards type '%s'", creationRequest.PlayingType.String()))
	}
	if creationRequest.Shuffled {
		playingDeck.Shuffle()
	}
	return &playingDeck, nil
}
