package decks

import (
	"croupier.io/cards"
	"errors"
	"github.com/google/uuid"
	"strings"
	"unicode"
)

// FrenchDeck is the representation of a deck containing French-suited playable cards.
type FrenchDeck struct {
	PlayableDeck
}

var _ Deck = &FrenchDeck{}

// NewFrenchDeck creates and returns a FrenchDeck according to the French-suited card standards.
// A successful NewFrenchDeck returns err == nil.
func NewFrenchDeck(requestedCardCodes []string) (*FrenchDeck, error) {
	playingCards, err := generateFrenchDeckPlayingCards(requestedCardCodes)
	if err != nil {
		return nil, errors.New("french playing cards creation failure on deck generation")
	}
	return &FrenchDeck{
		PlayableDeck: PlayableDeck{
			ID:        uuid.New(),
			Cards:     playingCards,
			Shuffled:  false,
			Remaining: len(playingCards),
		},
	}, nil
}

// generateFrenchDeckPlayingCards generates and return a slice of cards.PlayingCard according to
// the French-suited card standards.
// generateFrenchDeckPlayingCards creates a standard set of French-suited cards if requestedCardCodes is empty.
// If requestedCardCodes contains unrecognizable card codes according to the French-suited card standards,
// an error is returned.
// A successful generateFrenchDeckPlayingCards returns err == nil.
func generateFrenchDeckPlayingCards(requestedCardCodes []string) ([]cards.PlayingCard, error) {
	var playingCards []cards.PlayingCard

	refinedRequestedCardCodes := refineRequestedCardCodes(requestedCardCodes)
	generatedCards := make(map[string]cards.PlayingCard)

	for _, suit := range cards.FrenchCardSuits {
		for _, value := range cards.FrenchCardValues {
			card, err := cards.NewFrenchCard(suit.String(), value)
			if err != nil {
				return nil, errors.New("french playing cards creation failure on deck generation")
			}
			generatedCards[card.Code] = card.PlayingCard
			playingCards = append(playingCards, card.PlayingCard)
		}
	}
	if len(refinedRequestedCardCodes) > 0 {
		var requestedCards []cards.PlayingCard
		for _, cardCode := range refinedRequestedCardCodes {
			_, isPresent := generatedCards[cardCode]
			if !isPresent {
				return nil, errors.New("requested cards code does not exist in the standard deck")
			}
			requestedCards = append(requestedCards, generatedCards[cardCode])
		}
		return requestedCards, nil
	}
	return playingCards, nil
}

// refineRequestedCardCodes returns a processable version of requestedCardCodes
// by removing any duplicates or removing any whitespace contained in the provided card codes.
func refineRequestedCardCodes(requestedCardCodes []string) []string {
	if len(requestedCardCodes) == 0 {
		return []string{}
	}

	refinedRequestedCardCodes := make([]string, 0)
	requestedCardOccurrences := make(map[string]int)
	for _, cardCode := range requestedCardCodes {
		formattedCardCode := strings.Map(func(r rune) rune {
			if unicode.IsSpace(r) {
				return -1
			}
			return r
		}, cardCode)
		if formattedCardCode != "" {
			_, isPresent := requestedCardOccurrences[cardCode]
			if !isPresent {
				refinedRequestedCardCodes = append(refinedRequestedCardCodes, cardCode)
				requestedCardOccurrences[cardCode] = 1
			}
		}
	}
	return refinedRequestedCardCodes
}
