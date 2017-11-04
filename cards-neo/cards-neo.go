package cardsneo

import (
	"errors"
	"fmt"
	"io"

	"github.com/johnnadratowski/golang-neo4j-bolt-driver/structures/graph"

	cs "github.com/billglover/cards/cards-service"
	"github.com/billglover/uid"

	bolt "github.com/johnnadratowski/golang-neo4j-bolt-driver"
)

// DB holds a reference to the datbase.
type DB struct {
	bolt.Conn
}

// Tx holds a reference to the transaction.
// With Neo4J this is unused.
type Tx struct {
	bolt.Conn
}

// Open returns a DB reference for a data source.
func Open(dataSourceName string) (*DB, error) {
	driver := bolt.NewDriver()
	conn, err := driver.OpenNeo(dataSourceName)
	if err != nil {
		return nil, err
	}

	return &DB{conn}, nil
}

// Begin starts an returns a new transaction.
func (db *DB) Begin() (*Tx, error) {
	return &Tx{db.Conn}, nil
}

// CreateCard creates a new card.
// Returns the id of the card that was created or an error if the tx fails.
func (tx *Tx) CreateCard(c *cs.Card) (string, error) {

	if c == nil {
		return "", errors.New("card required")
	} else if c.Title == "" {
		return "", errors.New("card.Title required")
	}

	stmt, err := tx.PrepareNeo("CREATE (n:Card {uid: {uid}, title:{title}}) RETURN n.uid, n.title")
	if err != nil {
		return "", err
	}
	defer stmt.Close()

	uid, _ := uid.NextStringID()
	data := map[string]interface{}{"uid": uid, "title": c.Title}
	rows, err := stmt.QueryNeo(data)
	if err != nil {
		return "", err
	}

	// we only expect one row
	row, _, err := rows.NextNeo()
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	_, _, err = rows.NextNeo()
	if err != io.EOF {
		fmt.Println(err)
		return "", err
	}

	return row[0].(string), nil
}

// DeleteCard deletes a card based on its id.
// Returns the number of records deleted or an error if the tx fails.
func (tx *Tx) DeleteCard(c *cs.Card) (int64, error) {

	if c == nil {
		return 0, errors.New("card required")
	} else if c.Id == "" {
		return 0, errors.New("card.Id required")
	}

	stmt, err := tx.PrepareNeo("MATCH (n:Card {uid: {uid}}) DETACH DELETE n")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	data := map[string]interface{}{"uid": c.Id}
	result, err := stmt.ExecNeo(data)

	if err != nil {
		return 0, err
	}

	numResult, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	// Neo4J returns rows for each relationship deleted as well
	if numResult > 1 {
		numResult = 1
	}

	return numResult, nil
}

// EmbedCard embeds one card inside another.
// Returns the number of records ammended or an error if the tx fails.
func (tx *Tx) EmbedCard(p, c *cs.Card) (int64, error) {

	if p == nil || c == nil {
		return 0, errors.New("parent and child cards required")
	} else if p.Id == "" || c.Id == "" {
		return 0, errors.New("parent.Id and child.Id are both required")
	}

	stmt, err := tx.PrepareNeo("MATCH (p:Card {uid: {pid}}), (c:Card {uid: {cid}}) CREATE (p)-[:contains]->(c)")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	data := map[string]interface{}{"pid": p.Id, "cid": c.Id}
	result, err := stmt.ExecNeo(data)
	if err != nil {
		return 0, err
	}

	numResult, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return numResult, nil
}

// RemoveCard embeds one card inside another.
// Returns the number of records amended or an error if the tx fails.
func (tx *Tx) RemoveCard(p, c *cs.Card) (int64, error) {

	if p == nil || c == nil {
		return 0, errors.New("parent and child cards required")
	} else if p.Id == "" || c.Id == "" {
		return 0, errors.New("parent.Id and child.Id are both required")
	}

	stmt, err := tx.PrepareNeo("MATCH (p:Card {uid: {pid}})-[r:contains]->(c:Card {uid: {cid}}) DELETE (r)")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	data := map[string]interface{}{"pid": p.Id, "cid": c.Id}
	result, err := stmt.ExecNeo(data)
	if err != nil {
		return 0, err
	}

	numResult, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("unable to delete relationship %s->%s", p.Id, c.Id)
	}

	return numResult, nil
}

// GetCard returns a card based on its identifier.
// Returns the card or an error if the tx fails.
func (tx *Tx) GetCard(c *cs.Card) (*cs.Card, error) {

	if c == nil {
		return c, errors.New("card required")
	} else if c.Id == "" {
		return c, errors.New("card.Id required")
	}

	stmt, err := tx.PrepareNeo("MATCH (p:Card {uid: {uid}})-[r:contains]->(c:Card) RETURN {parent:p, children: [collect(c)]};")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	data := map[string]interface{}{"uid": c.Id}
	rows, err := stmt.QueryNeo(data)
	if err != nil {
		return nil, err
	}

	// we only expect one row
	row, _, err := rows.NextNeo()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	_, _, err = rows.NextNeo()
	if err != io.EOF {
		fmt.Println(err)
		return nil, err
	}

	r, ok := row[0].(map[string]interface{})
	if ok != true {
		return nil, fmt.Errorf("unable to parse response from database")
	}

	parent, ok := r["parent"].(graph.Node)
	if ok != true {
		return nil, fmt.Errorf("unable to parse parent card in database response")
	}

	c.Id = parent.Properties["uid"].(string)
	c.Title = parent.Properties["title"].(string)

	children, ok := r["children"].([]interface{})
	if ok != true {
		return nil, fmt.Errorf("unable to parse parent card in database response")
	}

	// The data structure returned by the database is a nested array
	// requiring us to unwrap this twice.
	children, ok = children[0].([]interface{})
	if ok != true {
		return nil, fmt.Errorf("unable to parse parent card in database response")
	}

	for _, v := range children {
		child := v.(graph.Node)

		embeddedCard := cs.Card{
			Id:    child.Properties["uid"].(string),
			Title: child.Properties["title"].(string),
		}

		c.Cards = append(c.Cards, &embeddedCard)

	}

	return c, nil
}

// Commit is an empty function as the Neo4J driver doesn't handle
// transactions in a way that we can use.
func (tx *Tx) Commit() error {
	return nil
}

// Rollback is an empty function as the Neo4J driver doesn't handle
// transactions in a way that we can use.
func (tx *Tx) Rollback() error {
	return nil
}
