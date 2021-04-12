package main

import (
	"errors"
	"github.com/dgraph-io/badger/v3"
)

// ErrDataAlreadyExist is returned when no data is found for the given key
var ErrDataAlreadyExist = errors.New("data has already exist for this key")

// Upsert inserts the record into the badger if it doesn't exist.  If it does already exist, then it updates
// the existing record
func (s *Store) Upsert(key string, data interface{}) error {
	encodedKey := encodeKey(key)
	encodedValue, err := encodeValue(data)
	if err != nil {
		return err
	}
	err = s.Badger().Update(func(txn *badger.Txn) error {
		txn.Set(encodedKey, encodedValue)
		return nil
	})
	return err
}

// Add inserts the record into the badger if it doesn't exist.  If it does already exist, then it fail
func (s *Store) Add(key string, data interface{}) error {
	value, _ := s.getValue(key)
	if value == nil {
		encodedKey := encodeKey(key)
		encodedValue, err := encodeValue(data)
		if err != nil {
			return err
		}
		err = s.Badger().Update(func(txn *badger.Txn) error {
			txn.Set(encodedKey, encodedValue)
			return nil
		})
		return err
	} else {
		return ErrDataAlreadyExist
	}
}
