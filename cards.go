package cards

// A Card holds an item of learning
type Card interface {
	ID() int
}

// A Deck holds many cards
type Deck interface {
	Card
	Add(Card) error
	Delete(Card) error
	NoOfCards() int
}
