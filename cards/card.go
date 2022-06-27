// Package cards implements models and actions for retrieving information of a card.
package cards

import (
	"errors"
	"regexp"
)

// PlayingCardType is the representation of a card type.
// PlayingCardType must be created upon the creation of a card extension.
type PlayingCardType int

const (
	French PlayingCardType = iota
)

// String returns a stringified version of a PlayingCardType.
func (cardType PlayingCardType) String() string {
	switch cardType {
	case French:
		return "French"
	}
	return "Undefined"
}

// Card is the interface that wraps the methods applicable from a card.
//
// ComputeCode determines and returns the code of a card.
// Any implementation of the method should take into account the uniqueness of a card code.
type Card interface {
	ComputeCode() string
}

// PlayingCard is the representation of a card entity.
// PlayingCard should be defined in any specific type of card.
type PlayingCard struct {
	Value string `json:"value"`
	Suit  string `json:"suit"`
	Code  string `json:"code"`
}

var _ Card = PlayingCard{}

// NewPlayingCard creates and returns a PlayableCard based on the provided suit and value.
// A successful NewPlayingCard returns err == nil.
func NewPlayingCard(suit string, value string) (*PlayingCard, error) {
	if suit == "" {
		return nil, errors.New("empty suit")
	}
	if value == "" {
		return nil, errors.New("empty value")
	}
	playingCard := PlayingCard{Suit: suit, Value: value}
	playingCard.Code = playingCard.ComputeCode()
	return &playingCard, nil
}

// ComputeCode determines and returns the code of a card.
// Any implementation of the method should take into account the uniqueness of a card code.
func (card PlayingCard) ComputeCode() string {
	isAlphabetical := regexp.MustCompile(`^[a-zA-Z]+$`).MatchString
	cardValue := card.Value
	if isAlphabetical(cardValue) {
		cardValue = card.Value[0:1]
	}
	return cardValue + card.Suit[0:1]
}
