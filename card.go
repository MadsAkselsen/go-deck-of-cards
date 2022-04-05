//go:generate stringer -type=Suit,Rank

package deck

import (
	"fmt"
	"sort"
)

type Suit uint8

const (
	Spade Suit = iota
	Diamond
	Club
	Heart
	Joker // this is a special case
)

var suits = [...]Suit{Spade, Diamond, Club, Heart}

type Rank uint8

const (
	_ Rank = iota
	Ace
	Two
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Jack
	Queen
	King
)

const (
	minRank = Ace
	maxRank = King
)

type Card struct {
	Suit
	Rank
}

func (c Card) String() string {
	if c.Suit == Joker {
		return c.Suit.String()
	}
	return fmt.Sprintf("%s of %ss", c.Rank.String(), c.Suit.String())
}

type NewOpts struct {
	Shuffle bool
}

// Create new deck of cards.
// Cards are ordered by default, but can pass in an
// array of functions to further manipulate the deck.
// Each function takes the deck as argument, and returns
// a new changed deck.
func New(options ...func([]Card) []Card) []Card {
	var cards []Card
	
	for _, suit := range suits {
		for rank := minRank; rank <= maxRank; rank++ {
			cards = append(cards, Card{Suit: suit, Rank: rank})
		}
	}

	for _, option := range options {
		cards = option(cards)
	}

	return cards
}

func DefaultSort(cards []Card) []Card {
	sort.Slice(cards, Less(cards))
	return cards
}

// argument is: func(cards []Card) func(i, j int) bool, which is a less function
// it takes in a slice of cards and returns a Less function.
// It works like the above DefaultSort, but with the difference
// that we pass in a custom less function
func Sort(less func(cards []Card) func(i, j int) bool) func([]Card) []Card {
	return func(cards []Card) []Card {
		sort.Slice(cards, less(cards))
		return cards
	}
}

func Less(cards []Card) func(i, j int) bool {
	return func(i, j int) bool {
		return absRank(cards[i]) < absRank(cards[j])
	}
}

// Creates a unique rank for a card based on its Suit and rank.
func absRank(c Card) int {
	// The 'suit * maxRank' makes sure that Suits are
	// not overlapping.
	// For example,
	// Spade will be 0*13=0
	// Diamond will be 1*13=13
	// Club will be 2*13=26
	// So they cannot overlap because they each have unique range
	//			1		*		13		+	4
	return int(c.Suit) * int(maxRank) + int(c.Rank)
}