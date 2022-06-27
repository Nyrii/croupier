package decks

import (
	"croupier.io/cards"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateDeckFromUndefinedType(t *testing.T) {
	playingDeck, err := CreateDeck(CreationRequest{PlayingType: cards.PlayingCardType(99999)}, nil)
	assert.Nil(t, playingDeck, "expected no deck")
	assert.NotNil(t, err, "expected error")
}

func TestCreateFrenchDeck(t *testing.T) {
	expectedDeck, _ := NewFrenchDeck([]string{})

	actualDeck, err := CreateDeck(CreationRequest{PlayingType: cards.French}, nil)
	expectedDeck.ID = actualDeck.ID
	assert.Nil(t, err, "expected no error")
	assert.NotNil(t, expectedDeck, "expected generated deck")
	assert.Equal(t, &expectedDeck.PlayableDeck, actualDeck)
	assert.Equal(t, false, expectedDeck.Shuffled)
}

func TestCreateFrenchDeckFailure(t *testing.T) {
	savedFrenchCardSuits := cards.FrenchCardSuits
	savedFrenchCardValues := cards.FrenchCardValues

	cards.FrenchCardSuits = [4]cards.FrenchCardSuit{cards.Spades}
	cards.FrenchCardValues = [13]string{""}
	actualDeck, err := CreateDeck(CreationRequest{PlayingType: cards.French}, nil)
	assert.NotNil(t, err, "expected an error upon playable deck generation")
	assert.Nil(t, actualDeck, "expected no playable deck")

	t.Cleanup(func() {
		cards.FrenchCardSuits = savedFrenchCardSuits
		cards.FrenchCardValues = savedFrenchCardValues
	})
}

func TestCreateShuffledDeck(t *testing.T) {
	sortedDeck, _ := NewFrenchDeck([]string{})

	shuffledDeck, err := CreateDeck(CreationRequest{PlayingType: cards.French, Shuffled: true}, nil)
	sortedDeck.ID = shuffledDeck.ID
	assert.Nil(t, err, "expected no error")
	assert.NotNil(t, sortedDeck, "expected generated deck")
	assert.NotNil(t, shuffledDeck, "expected shuffled deck")
	assert.NotEqual(t, &sortedDeck.PlayableDeck, shuffledDeck, "expected different playable decks")
	assert.Equal(t, true, shuffledDeck.Shuffled)
	assert.Equal(t, len(sortedDeck.PlayableDeck.Cards), len(shuffledDeck.Cards))
}

func TestDrawCard(t *testing.T) {
	requestedCardCodes := []string{"AS", "2S", "3S"}
	aceOfSpades := cards.PlayingCard{Suit: cards.Spades.String(), Value: "ACE", Code: "AS"}
	twoOfSpades := cards.PlayingCard{Suit: cards.Spades.String(), Value: "2", Code: "2S"}
	threeOfSpades := cards.PlayingCard{Suit: cards.Spades.String(), Value: "3", Code: "3S"}

	testRecords := []struct {
		requestedDrawnCardNumber  int
		expectedDrawnPlayingCards []cards.PlayingCard
		expectedRemainingCards    []cards.PlayingCard
	}{
		{0,
			[]cards.PlayingCard{},
			[]cards.PlayingCard{aceOfSpades, twoOfSpades, threeOfSpades}},
		{1,
			[]cards.PlayingCard{aceOfSpades},
			[]cards.PlayingCard{twoOfSpades, threeOfSpades}},

		{3,
			[]cards.PlayingCard{aceOfSpades, twoOfSpades, threeOfSpades},
			[]cards.PlayingCard{}},

		{10,
			[]cards.PlayingCard{aceOfSpades, twoOfSpades, threeOfSpades},
			[]cards.PlayingCard{}},

		{-1,
			[]cards.PlayingCard{},
			[]cards.PlayingCard{aceOfSpades, twoOfSpades, threeOfSpades}},
	}
	for _, testRecord := range testRecords {
		playingDeck, err := CreateDeck(CreationRequest{PlayingType: cards.French}, requestedCardCodes)
		assert.Nil(t, err, "expected no error")
		assert.NotNil(t, playingDeck, "expected a playable deck")

		drawnCards := playingDeck.DrawCard(testRecord.requestedDrawnCardNumber)
		assert.Equal(t, testRecord.expectedDrawnPlayingCards, drawnCards, "expected identical drawn playing cards")
		assert.Equal(t, testRecord.expectedRemainingCards, playingDeck.Cards, "expected identical remaining playing cards")
		assert.Equal(t, len(testRecord.expectedRemainingCards), playingDeck.Remaining, "wrong remaining count")
	}
}
