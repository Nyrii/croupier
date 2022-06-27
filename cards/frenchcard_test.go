package cards

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFrenchCardStringifiedSuit(t *testing.T) {
	testRecords := []struct {
		suit                    FrenchCardSuit
		expectedStringifiedSuit string
	}{
		{Spades, "SPADES"},
		{Diamonds, "DIAMONDS"},
		{Hearts, "HEARTS"},
		{Clubs, "CLUBS"},
	}
	for _, testRecord := range testRecords {
		assert.Equal(t, testRecord.expectedStringifiedSuit, testRecord.suit.String(), "expected identical enum representation")
	}
}

func TestFrenchCardSuits(t *testing.T) {
	expectedSize := 4
	actualSize := len(FrenchCardSuits)
	expectedSuits := [4]FrenchCardSuit{Spades, Diamonds, Clubs, Hearts}
	assert.Equal(t, expectedSize, actualSize, "expected %d french cards unique suits; has %d", expectedSize, actualSize)
	assert.Equal(t, expectedSuits, FrenchCardSuits)
}

func TestFrenchCardValues(t *testing.T) {
	expectedSize := 13
	actualSize := len(FrenchCardValues)
	expectedValues := [13]string{
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
	assert.Equal(t, expectedSize, actualSize, "expected %d french cards unique values; has %d", expectedSize, actualSize)
	assert.Equal(t, expectedValues, FrenchCardValues)
}

func TestFrenchCard(t *testing.T) {
	card, err := NewFrenchCard(Clubs.String(), "ACE")
	assert.Nil(t, err)
	assert.Equal(t, "AC", card.ComputeCode())
}

func TestFrenchCardWithEmptySuit(t *testing.T) {
	card, err := NewFrenchCard("", "ACE")
	assert.Nil(t, card)
	assert.NotNil(t, err)
}

func TestFrenchCardWithEmptyValue(t *testing.T) {
	card, err := NewFrenchCard(Clubs.String(), "")
	assert.Nil(t, card)
	assert.NotNil(t, err)
}
