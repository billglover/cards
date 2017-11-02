package cardsneo

import (
	"errors"
	"fmt"
	"io"

	cs "github.com/billglover/cards/cards-service"
	"github.com/billglover/uid"

	bolt "github.com/johnnadratowski/golang-neo4j-bolt-driver"
)

type DB struct {
	bolt.Conn
}
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

	fmt.Println("---")

	fmt.Printf("COLUMNS: %#v\n", rows.Metadata()["fields"].([]interface{}))
	fmt.Printf("FIELDS: %s %s\n", row[0].(string), row[1].(string))
	stmt.Close()

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

	data := map[string]interface{}{"uid": c.Id}
	result, err := stmt.ExecNeo(data)
	stmt.Close()
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

	data := map[string]interface{}{"pid": p.Id, "cid": c.Id}
	result, err := stmt.ExecNeo(data)
	stmt.Close()
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
// Returns the number of records ammended or an error if the tx fails.
func (tx *Tx) RemoveCard(p, c *cs.Card) (int64, error) {

	// if p == nil || c == nil {
	// 	return 0, errors.New("parent and child cards required")
	// } else if p.Id == "" || c.Id == "" {
	// 	return 0, errors.New("parent.Id and child.Id are both required")
	// }

	// stmt, err := tx.Prepare("DELETE FROM links WHERE parent=? AND child=?")
	// if err != nil {
	// 	return 0, err
	// }

	// res, err := stmt.Exec(p.Id, c.Id)
	// return res.RowsAffected()

	return 0, nil
}

// GetCard returns a card based on its identifier.
// Returns the card or an error if the tx fails.
func (tx *Tx) GetCard(c *cs.Card) (*cs.Card, error) {

	// if c == nil {
	// 	return c, errors.New("card required")
	// } else if c.Id == "" {
	// 	return c, errors.New("card.Id required")
	// }

	// stmt, err := tx.Prepare("SELECT uid, title FROM cards WHERE uid=?")
	// if err != nil {
	// 	return c, err
	// }

	// row := stmt.QueryRow(c.Id)

	// var uid int
	// var title string
	// err = row.Scan(&uid, &title)
	// if err != nil {
	// 	return c, err
	// }
	// c.Id = fmt.Sprintf("%d", uid)
	// c.Title = title

	// stmt, err = tx.Prepare("SELECT child FROM links WHERE parent=?")
	// if err != nil {
	// 	return c, err
	// }

	// rows, err := stmt.Query(c.Id)
	// if err != nil {
	// 	return c, err
	// }

	// for rows.Next() {
	// 	rows.Scan(&uid)
	// 	c.Cards = append(c.Cards, &cs.Card{Id: fmt.Sprintf("%d", uid)})
	// }
	// rows.Close()

	return c, nil
}

func (tx *Tx) Commit() error {
	return nil
}

func (tx *Tx) Rollback() error {
	return nil
}
