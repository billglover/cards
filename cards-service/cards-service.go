package cards_service

import (
	"database/sql"
	"errors"

	_ "github.com/go-sql-driver/mysql"
)

type DB struct {
	*sql.DB
}
type Tx struct {
	*sql.Tx
}

// Open returns a DB reference for a data source.
func Open(dataSourceName string) (*DB, error) {
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return nil, err
	}
	return &DB{db}, nil
}

// Begin starts an returns a new transaction.
func (db *DB) Begin() (*Tx, error) {
	tx, err := db.DB.Begin()
	if err != nil {
		return nil, err
	}
	return &Tx{tx}, nil
}

// CreateCard creates a new user.
// Returns an error if user is invalid or the tx fails.
func (tx *Tx) CreateCard(c *Card) (int64, error) {
	// Validate the input.
	if c == nil {
		return 0, errors.New("card required")
	} else if c.Title == "" {
		return 0, errors.New("card.Title required")
	}

	// Perform the actual insert and return any errors.
	//return tx.Exec(`INSERT INTO users (...) VALUES`, ...)
	stmt, err := tx.Prepare("INSERT cards SET title=?")
	if err != nil {
		return 0, err
	}

	res, err := stmt.Exec(c.Title)
	return res.LastInsertId()
}
