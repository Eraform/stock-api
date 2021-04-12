package main

import "errors"

// ErrNotFound is returned when no data is found for the given key
var ErrNotFound = errors.New("no data found for this key")
