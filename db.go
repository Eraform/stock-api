package main

import (
	"github.com/dgraph-io/badger/v3"
	"log"
)

// Store is a badger
type Store struct {
	db *badger.DB
}

// NewBadger returns a new badger with path
func NewBadger(path string) (Store, error) {
	b := Store{}
	opts := badger.DefaultOptions(path)
	db, err := badger.Open(opts)
	if err != nil {
		log.Fatal(err)
		return b, err
	}
	b.db = db
	return b, nil
}

// Badger returns the underlying Badger DB
func (s *Store) Badger() *badger.DB {
	return s.db
}
