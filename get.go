package main

import (
	"bytes"
	"encoding/gob"
	"github.com/dgraph-io/badger/v3"
	"log"
)

// KVP simple named key value pair storage
type KVP struct {
	Key   []byte
	Value []byte
}

// GetValue retrieves a value []byte from data store and return a copy of it
func (s *Store) getValue(key string) ([]byte, error) {
	encodedKey := encodeKey(key)

	var valueCopy []byte
	err := s.Badger().View(func(txn *badger.Txn) error {
		item, err := txn.Get(encodedKey)
		if err != nil {
			return ErrNotFound
		}
		err = item.Value(func(val []byte) error {
			valueCopy = append([]byte{}, val...)
			return nil
		})
		return err
	})
	return valueCopy, err
}

// GetStock retrieves an Stock from data store
func (s *Store) GetStock(key string) (Stock, error) {
	valueBytes, _ := s.getValue(key)

	if valueBytes == nil {
		return Stock{}, ErrNotFound
	}

	var stock Stock
	d := gob.NewDecoder(bytes.NewReader(valueBytes))
	if err := d.Decode(&stock); err != nil {
		//log.Println("Decoding error")
		return Stock{}, err
	}
	//log.Println("Item decoded", stock)
	return stock, nil
}

// GetValues retrieves a value from data store
func (s *Store) getValues() ([]KVP, error) {
	var results []KVP

	err := s.Badger().View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			k := item.Key()
			err := item.Value(func(v []byte) error {
				//fmt.Printf("key=%s, value=%s\n", k, v)
				res := KVP{k, v}
				results = append(results, res)
				return nil
			})

			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return []KVP{}, err
	}
	return results, nil
}

func (s *Store) GetStocks() ([]Stock, error) {
	valuesBytes, err := s.getValues()
	if err != nil {
		return nil, err
	}
	var results []Stock

	for _, row := range valuesBytes {
		item, err := decodeValue(row.Value)
		if err != nil {
			return nil, err
		}
		results = append(results, item)
	}
	return results, nil
}

func (s *Store) GetAllKeys() ([]string, error) {
	var results []string

	err := s.Badger().View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			row := it.Item()
			k := row.Key()
			results = append(results, string(k))
			log.Printf("key=%s\n", k)
		}
		return nil
	})
	if err != nil {
		return []string{}, err
	}
	return results, nil
}
