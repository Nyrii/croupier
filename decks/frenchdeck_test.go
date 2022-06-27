package decks

import (
	"croupier.io/cards"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewFrenchDeck(t *testing.T) {
	var expectedFrenchCardSuits = [4]cards.FrenchCardSuit{cards.Spades, cards.Diamonds, cards.Clubs, cards.Hearts}
	var expectedFrenchCardValues = [13]string{
		"ACE",
		"2",
		"3",
		"4",
		"5",
		"6",
		"7",
		"8",
		"9",
		"10",
		"JACK",
		"QUEEN",
		"KING"}

	actualDeck, err := NewFrenchDeck([]string{})
	playingCards := actualDeck.Cards
	assert.Nil(t, err, "expected no error when generating the deck")
	assert.NotNil(t, actualDeck, "expected a valid deck")

	var i = 0
	for _, suit := range expectedFrenchCardSuits {
		for _, value := range expectedFrenchCardValues {
			card, err := cards.NewFrenchCard(suit.String(), value)
			assert.Nil(t, err, "expected no error upon French cards creation")
			assert.Equal(t, card.PlayingCard, playingCards[i], "expected identical cards")
			i += 1
		}
	}
	assert.Equal(t, len(cards.FrenchCardSuits)*len(cards.FrenchCardValues), actualDeck.Remaining)
	assert.Equal(t, false, actualDeck.Shuffled)
}

func TestGenerateFrenchDeckPlayingCards(t *testing.T) {
	testRecords := []struct {
		requestedCardCodes       []string
		expectedRefinedCardCodes []string
		expectedCardsSize        int
	}{
		{[]string{}, []string{}, len(cards.FrenchCardSuits) * len(cards.FrenchCardValues)},
		{[]string{"AS", "5S", "10S", "AH", "KH", "2D"}, []string{"AS", "5S", "10S", "AH", "KH", "2D"}, 6},
		{[]string{"AS", "    ", "5S", ""}, []string{"AS", "5S"}, 2},
		{[]string{"X", "AS"}, []string{}, -1},
	}

	for _, testRecord := range testRecords {
		playingCards, err := generateFrenchDeckPlayingCards(testRecord.requestedCardCodes)
		if testRecord.expectedCardsSize == -1 {
			assert.Nil(t, playingCards, "expected no playing cards")
			assert.NotNil(t, err, "expected an error when generating the deck")
		} else {
			assert.NotNil(t, playingCards, "expected playing cards")
			assert.Nil(t, err, "expected no error")
			assert.Equal(t, testRecord.expectedCardsSize, len(playingCards))
			if len(testRecord.expectedRefinedCardCodes) > 0 {
				for i, card := range playingCards {
					assert.Equal(t, testRecord.expectedRefinedCardCodes[i], card.ComputeCode(), "requested cards code and expected cards code do not match")
				}
			}
		}
	}
}

func TestNewFrenchDeckWithEmptyProperties(t *testing.T) {
	testRecords := []struct {
		cardSuits  [4]cards.FrenchCardSuit
		cardValues [13]string
	}{
		{[4]cards.FrenchCardSuit{cards.Spades}, [13]string{""}},
		{[4]cards.FrenchCardSuit{""}, [13]string{"ACE"}},
	}
	for _, testRecord := range testRecords {
		savedFrenchCardSuits := cards.FrenchCardSuits
		savedFrenchCardValues := cards.FrenchCardValues

		cards.FrenchCardSuits = testRecord.cardSuits
		cards.FrenchCardValues = testRecord.cardValues
		deck, err := NewFrenchDeck([]string{})
		assert.NotNil(t, err, "expected an error when generating the deck")
		assert.Nil(t, deck, "expected no deck")

		t.Cleanup(func() {
			cards.FrenchCardSuits = savedFrenchCardSuits
			cards.FrenchCardValues = savedFrenchCardValues
		})
	}
}
