package cards

import "testing"

func TestCreateDeck(t *testing.T) {
	d, err := NewBasicDeck(1)
	if err != nil {
		t.Fatalf("unexpected error returned: %v", err)
	}
	if d == nil {
		t.Fatal("expected a deck, received none")
	}
	if d.ID() != 1 {
		t.Fatalf("unexpected ID: got %d, want %d", d.ID(), 1)
	}
}

func TestAddDeckToDeck(t *testing.T) {
	d1, _ := NewBasicDeck(1)
	d2, _ := NewBasicDeck(1)

	err := d2.Add(d1)
	if err != nil {
		t.Fatalf("unexpected error returned: %v", err)
	}
	if d2.NoOfCards() != 1 {
		t.Fatalf("unexpected number of cards: got %d, want %d", d2.NoOfCards(), 1)
	}
}

func TestDeleteDeckFromDeck(t *testing.T) {
	d1, _ := NewBasicDeck(1)
	d2, _ := NewBasicDeck(1)
	d2.Add(d1)

	err := d2.Delete(d1)
	if err != nil {
		t.Fatalf("unexpected error returned: %v", err)
	}
	if d2.NoOfCards() != 0 {
		t.Fatalf("unexpected number of cards: got %d, want %d", d2.NoOfCards(), 0)
	}
}
