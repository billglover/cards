package cards_service

import (
	"log"

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

	col := tx.DB("cards").C("cards")
	id := bson.ObjectIdHex(c.Id)

	err := col.RemoveId(id)
	if err != nil {
		return 0, err
	}

	return 1, nil
}

// EmbedCard embeds one card inside another.
// Returns the number of records amended or an error if the tx fails.
func (tx *Tx) EmbedCard(p, c *Card) (int, error) {
	col := tx.DB("cards").C("cards")
	cID := bson.ObjectIdHex(c.Id)
	pID := bson.ObjectIdHex(p.Id)

	change := bson.M{"$addToSet": bson.M{"cards": cID}}

	err := col.UpdateId(pID, change)
	if err != nil {
		return 0, err
	}

	return 1, nil
}

// RemoveCard embeds one card inside another.
// Returns the number of records amended or an error if the tx fails.
func (tx *Tx) RemoveCard(p, c *Card) (int, error) {
	col := tx.DB("cards").C("cards")
	cID := bson.ObjectIdHex(c.Id)
	pID := bson.ObjectIdHex(p.Id)

	change := bson.M{"$pull": bson.M{"cards": cID}}

	err := col.UpdateId(pID, change)
	if err != nil {
		return 0, err
	}

	return 1, nil
}

// GetCard returns a card based on its identifier.
// Returns the card or an error if the tx fails.
func (tx *Tx) GetCard(c *Card) (*Card, error) {

	col := tx.DB("cards").C("cards")
	id := bson.ObjectIdHex(c.Id)

	cDoc := &CardDoc{}

	err := col.FindId(id).One(cDoc)
	if err != nil {
		return c, err
	}

	c.Title = cDoc.Title
	log.Println(c)
	for _, v := range cDoc.Cards {
		c.Cards = append(c.Cards, &Card{Id: v.Hex()})
	}
	return c, nil
}

func (tx *Tx) Commit() error {
	return nil
}

func (tx *Tx) Rollback() error {
	return nil
}
