package main

import (
	"bytes"
	"encoding/gob"
	"log"
)

func encodeKey(key string) []byte {
	return []byte(key)
}

func encodeValue(value interface{}) ([]byte, error) {
	var buff bytes.Buffer

	en := gob.NewEncoder(&buff)

	err := en.Encode(value)
	if err != nil {
		return nil, err
	}

	return buff.Bytes(), nil
}

func decodeValue(value []byte) (Stock, error) {
	var stock Stock
	d := gob.NewDecoder(bytes.NewReader(value))
	err := d.Decode(&stock)
	if err != nil {
		log.Println("Decoding error")
		return Stock{}, err
	}
	log.Println("Stock decoded", stock)
	return stock, nil
}
