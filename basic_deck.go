package cards

// BasicDeck is a simple deck of cards
type basicDeck struct {
	Identifier int    `json:"id"`
	Cards      []Card `json:"cards"`
}

// Add adds a card to the deck
func (d *basicDeck) Add(c Card) error {
	d.Cards = append(d.Cards, c)
	return nil
}

// Delete removes a card from the deck
func (d *basicDeck) Delete(c Card) error {
	for i, card := range d.Cards {
		if card.ID() == c.ID() {
			d.Cards = d.Cards[:i+copy(d.Cards[i:], d.Cards[i+1:])]
		}
	}
	return nil
}

// ID returns the ID of the deck
func (d *basicDeck) ID() int {
	return d.Identifier
}

// NoOfCards returns the number of cards in the deck
func (d *basicDeck) NoOfCards() int {
	return len(d.Cards)
}

// NewBasicDeck returns a new BasicDeck
func NewBasicDeck(id int) (Deck, error) {
	d := basicDeck{
		Identifier: id,
	}
	d.Cards = make([]Card, 0)
	return &d, nil
}
