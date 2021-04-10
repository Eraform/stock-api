package main

import (
	"github.com/dgraph-io/badger/v3"
	"log"
)

// Badger is a simple helper for db access
type Badger struct {
	db *badger.DB
}

// KVP simple named key value pair storage
type KVP struct {
	Key   []byte
	Value []byte
}

// NewBadgerDB returns a new badgerdb
func NewBadgerDB() (Badger, error) {
	b := Badger{}
	opts := badger.DefaultOptions("./tmp")
	db, err := badger.Open(opts)
	if err != nil {
		log.Fatal(err)
		return b, err
	}
	b.db = db
	return b, nil
}

// Update update a set of key/values in the db
func (b Badger) Update(key, value []byte) error {
	err := b.db.Update(func(txn *badger.Txn) error {
		txn.Set(key, value)
		return nil
	})
	if err != nil {
		return err
	} else {
		return nil
	}
}

func (b Badger) Get(key []byte) ([]byte, error) {
	var valueCopy []byte
	err := b.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			return err
		}
		err = item.Value(func(val []byte) error {
			valueCopy = append([]byte{}, val...)
			return nil
		})
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	} else {
		return valueCopy, nil
	}
}

func (b Badger) Delete(key []byte) error {
	err := b.db.Update(
		func(txn *badger.Txn) error {
			err := txn.Delete(key)
			if err != nil {
				return err
			}
			return nil
		})
	return err
}
