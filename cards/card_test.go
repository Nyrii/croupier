package cards

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPlayingCardStringifiedTypes(t *testing.T) {
	testRecords := []struct {
		playingCardType         PlayingCardType
		expectedStringifiedType string
	}{
		{PlayingCardType(9999), "Undefined"},
		{French, "French"},
	}
	for _, testRecord := range testRecords {
		assert.Equal(t, testRecord.expectedStringifiedType, testRecord.playingCardType.String(), "expected identical enum representation")
	}
}

func TestNewPlayingCard(t *testing.T) {
	testRecords := []struct {
		suit         string
		value        string
		expectedCard *PlayingCard
	}{
		{Clubs.String(), "ACE", &PlayingCard{Value: "ACE", Suit: "CLUBS", Code: "AC"}},
		{Diamonds.String(), "2", &PlayingCard{Value: "2", Suit: "DIAMONDS", Code: "2D"}},
		{Hearts.String(), "10", &PlayingCard{Value: "10", Suit: "HEARTS", Code: "10H"}},
		{Spades.String(), "JACK", &PlayingCard{Value: "JACK", Suit: "SPADES", Code: "JS"}},
		{Clubs.String(), "QUEEN", &PlayingCard{Value: "QUEEN", Suit: "CLUBS", Code: "QC"}},
		{Clubs.String(), "KING", &PlayingCard{Value: "KING", Suit: "CLUBS", Code: "KC"}},
		{"", "KING", nil},
		{Clubs.String(), "", nil},
	}
	for _, testRecord := range testRecords {
		result, err := NewPlayingCard(testRecord.suit, testRecord.value)
		if testRecord.expectedCard == nil {
			assert.Nil(t, result, "expected no result")
			assert.NotNil(t, err, "expected error")
		} else {
			assert.NotNil(t, result, "expected result")
			assert.Equal(t, testRecord.expectedCard, result, "expected result and expected cards to be identical")
			assert.Nil(t, err, "expected no error")
		}
	}
}
