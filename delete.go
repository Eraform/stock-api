package main

import (
	"github.com/dgraph-io/badger/v3"
)

func (s *Store) Delete(key []byte) error {
	err := s.Badger().Update(
		func(txn *badger.Txn) error {
			err := txn.Delete(key)
			if err != nil {
				return err
			}
			return nil
		})
	return err
}

// DeleteItem retrieves an Item from data store. Throw
func (s *Store) DeleteItem(key string) error {
	encodedKey := encodeKey(key)
	err := s.Delete(encodedKey)
	if err != nil {
		return err
	}
	return nil
}
