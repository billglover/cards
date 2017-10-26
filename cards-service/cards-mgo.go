package cards_service

import (
	"gopkg.in/mgo.v2/bson"

	"gopkg.in/mgo.v2"

	_ "gopkg.in/mgo.v2"
)

type DB struct {
	*mgo.Session
}
type Tx struct {
	*mgo.Session
}

type CardDoc struct {
	UID   bson.ObjectId   `bson:"_id"`
	Title string          `bson:"title"`
	Cards []bson.ObjectId `bson:"cards"`
}

// Open returns a DB reference for a data source.
func Open(dataSourceName string) (*DB, error) {
	db, err := mgo.Dial(dataSourceName)
	if err != nil {
		return nil, err
	}
	return &DB{db}, nil
}

// Begin returns the DB instance as MongoDB doesn't support transactions.
func (db *DB) Begin() (*Tx, error) {
	return &Tx{db.Session}, nil
}

// CreateCard creates a new card.
// Returns the id of the card that was created or an error if the tx fails.
func (tx *Tx) CreateCard(c *Card) (string, error) {

	col := tx.DB("cards").C("cards")

	cDoc := CardDoc{
		UID:   bson.NewObjectId(),
		Title: c.Title,
	}

	err := col.Insert(cDoc)
	if err != nil {
		return "", nil
	}

	return cDoc.UID.Hex(), err
}

// DeleteCard deletes a card based on its id.
// Returns the number of records deleted or an error if the tx fails.
func (tx *Tx) DeleteCard(c *Card) (int, error) {
	return 0, nil
}

// EmbedCard embeds one card inside another.
// Returns the number of records amended or an error if the tx fails.
func (tx *Tx) EmbedCard(p, c *Card) (int, error) {
	return 0, nil
}

// RemoveCard embeds one card inside another.
// Returns the number of records amended or an error if the tx fails.
func (tx *Tx) RemoveCard(p, c *Card) (int, error) {
	return 0, nil
}

// GetCard returns a card based on its identifier.
// Returns the card or an error if the tx fails.
func (tx *Tx) GetCard(c *Card) (*Card, error) {
	return c, nil
}

func (tx *Tx) Commit() error {
	return nil
}

func (tx *Tx) Rollback() error {
	return nil
}
