package database

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
)

var ErrNotExist = errors.New("resource does not exist")

type DB struct {
	path string
	mu   *sync.RWMutex
}

type DBStructure struct {
	Quotes map[int]Quote `json:"quotes"`
}

type Quote struct {
	Id     int    `json:"id"`
	Quote  string `json:"quote"`
	Author string `json:"author"`
}

func NewDB(path string) (*DB, error) {
	db := &DB{
		path: path,
		mu:   &sync.RWMutex{},
	}
	err := db.ensureDB()
	return db, err 
}

func (db *DB) ensureDB() error {
	_, err := os.ReadFile(db.path)
	if errors.Is(err, os.ErrNotExist) {
		return db.createDB()
}
return err
}

func (db *DB) createDB() error {
	dbStructure := DBStructure{
		Quotes: map[int]Quote{},
	}
	return db.writeDB(dbStructure)
}

func (db *DB) writeDB(dbStructure DBStructure) error {
	db.mu.Lock()
	defer db.mu.Unlock()
	
	dat, err := json.Marshal(dbStructure)
	if err != nil {
		return err
	}
	err = os.WriteFile(db.path, dat, 0600)
	if err != nil {
		return err
	}
	return nil
}

func(db *DB) CreateQuote(quote, author string) (Quote, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return Quote{}, err
	}
	id := len(dbStructure.Quotes) + 1
	newQuote := Quote{
		Id:     id,
		Quote:  quote,
		Author: author,
	}
	dbStructure.Quotes[id] = newQuote

	err = db.writeDB(dbStructure)
	if err != nil {
		return Quote{}, err
	}
	return dbStructure.Quotes[id], nil
}

func (db *DB) loadDB() (DBStructure, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	dbStructure := DBStructure{}
	dat, err := os.ReadFile(db.path)
	if errors.Is(err, os.ErrNotExist) {
		return dbStructure, err
	}
	err = json.Unmarshal(dat, &dbStructure)
	if err != nil {
		return dbStructure, err
	}
	return dbStructure, nil

}

func (db *DB) GetQuotes() ([]Quote, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return nil, err
	}
	quotes := make([]Quote, 0, len(dbStructure.Quotes))
	for _, quote := range dbStructure.Quotes {
		quotes = append(quotes, quote)
	}
	return quotes, nil
}

func (db *DB) GetQuote(id int) (Quote, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return Quote{}, err
	}
	quote, ok := dbStructure.Quotes[id]
	if !ok {
		return Quote{}, ErrNotExist
	}
	return quote, nil
}