package cards_service

// import (
// 	"database/sql"
// 	"errors"

// 	_ "github.com/go-sql-driver/mysql"
// )

// type DB struct {
// 	*sql.DB
// }
// type Tx struct {
// 	*sql.Tx
// }

// // Open returns a DB reference for a data source.
// func Open(dataSourceName string) (*DB, error) {
// 	db, err := sql.Open("mysql", dataSourceName)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &DB{db}, nil
// }

// // Begin starts an returns a new transaction.
// func (db *DB) Begin() (*Tx, error) {
// 	tx, err := db.DB.Begin()
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &Tx{tx}, nil
// }

// // CreateCard creates a new card.
// // Returns the id of the card that was created or an error if the tx fails.
// func (tx *Tx) CreateCard(c *Card) (int64, error) {

// 	if c == nil {
// 		return 0, errors.New("card required")
// 	} else if c.Title == "" {
// 		return 0, errors.New("card.Title required")
// 	}

// 	stmt, err := tx.Prepare("INSERT cards SET title=?")
// 	if err != nil {
// 		return 0, err
// 	}

// 	res, err := stmt.Exec(c.Title)
// 	return res.LastInsertId()
// }

// // DeleteCard deletes a card based on its id.
// // Returns the number of records deleted or an error if the tx fails.
// func (tx *Tx) DeleteCard(c *Card) (int64, error) {

// 	if c == nil {
// 		return 0, errors.New("card required")
// 	} else if c.Id == 0 {
// 		return 0, errors.New("card.Id required")
// 	}

// 	stmt, err := tx.Prepare("DELETE FROM cards WHERE uid=?")
// 	if err != nil {
// 		return 0, err
// 	}

// 	res, err := stmt.Exec(c.Id)
// 	return res.RowsAffected()
// }

// // EmbedCard embeds one card inside another.
// // Returns the number of records ammended or an error if the tx fails.
// func (tx *Tx) EmbedCard(p, c *Card) (int64, error) {

// 	if p == nil || c == nil {
// 		return 0, errors.New("parent and child cards required")
// 	} else if p.Id == 0 || c.Id == 0 {
// 		return 0, errors.New("parent.Id and child.Id are both required")
// 	}

// 	stmt, err := tx.Prepare("INSERT links SET parent=?, child=?")
// 	if err != nil {
// 		return 0, err
// 	}

// 	res, err := stmt.Exec(p.Id, c.Id)
// 	return res.RowsAffected()
// }

// // RemoveCard embeds one card inside another.
// // Returns the number of records ammended or an error if the tx fails.
// func (tx *Tx) RemoveCard(p, c *Card) (int64, error) {

// 	if p == nil || c == nil {
// 		return 0, errors.New("parent and child cards required")
// 	} else if p.Id == 0 || c.Id == 0 {
// 		return 0, errors.New("parent.Id and child.Id are both required")
// 	}

// 	stmt, err := tx.Prepare("DELETE FROM links WHERE parent=? AND child=?")
// 	if err != nil {
// 		return 0, err
// 	}

// 	res, err := stmt.Exec(p.Id, c.Id)
// 	return res.RowsAffected()
// }

// // GetCard returns a card based on its identifier.
// // Returns the card or an error if the tx fails.
// func (tx *Tx) GetCard(c *Card) (*Card, error) {

// 	if c == nil {
// 		return c, errors.New("card required")
// 	} else if c.Id == 0 {
// 		return c, errors.New("card.Id required")
// 	}

// 	stmt, err := tx.Prepare("SELECT uid, title FROM cards WHERE uid=?")
// 	if err != nil {
// 		return c, err
// 	}

// 	row := stmt.QueryRow(c.Id)

// 	var uid int
// 	var title string
// 	err = row.Scan(&uid, &title)
// 	if err != nil {
// 		return c, err
// 	}
// 	c.Id = uint64(uid)
// 	c.Title = title

// 	stmt, err = tx.Prepare("SELECT child FROM links WHERE parent=?")
// 	if err != nil {
// 		return c, err
// 	}

// 	rows, err := stmt.Query(c.Id)
// 	if err != nil {
// 		return c, err
// 	}

// 	for rows.Next() {
// 		rows.Scan(&uid)
// 		c.Cards = append(c.Cards, &Card{Id: uint64(uid)})
// 	}
// 	rows.Close()

// 	return c, nil
// }
