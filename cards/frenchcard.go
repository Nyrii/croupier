package cards

// FrenchCard is the representation of a French-suited playable card.
type FrenchCard struct {
	PlayingCard
}

var _ Card = &FrenchCard{}

// FrenchCardSuit is the representation of a FrenchCard suit.
type FrenchCardSuit string

const (
	Spades   FrenchCardSuit = "SPADES"
	Diamonds FrenchCardSuit = "DIAMONDS"
	Clubs    FrenchCardSuit = "CLUBS"
	Hearts   FrenchCardSuit = "HEARTS"
)

// String returns a stringified version of a FrenchCardSuit.
func (suit FrenchCardSuit) String() string {
	return string(suit)
}

// FrenchCardSuits is the definition of the FrenchCard suits panel.
// FrenchCardSuits must not be modified to preserve the French-suited playing card standards.
var FrenchCardSuits = [4]FrenchCardSuit{Spades, Diamonds, Clubs, Hearts}

// FrenchCardValues is the definition of the FrenchCard values panel.
// FrenchCardValues must not be modified to preserve the French-suited playing card standards.
var FrenchCardValues = [13]string{
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

// NewFrenchCard creates and returns a FrenchCard based on the provided suit and value.
// A successful NewFrenchCard returns err == nil.
func NewFrenchCard(suit string, value string) (*FrenchCard, error) {
	playingCard, err := NewPlayingCard(suit, value)
	if err != nil {
		return nil, err
	}
	return &FrenchCard{PlayingCard: *playingCard}, nil
}
