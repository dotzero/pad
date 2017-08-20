package service

import (
	hashids "github.com/speps/go-hashids"
)

var h *hashids.HashID

// NewHash returns a HashID client
func NewHash(salt string) *hashids.HashID {
	hd := hashids.NewData()
	hd.Salt = salt
	hd.MinLength = 3
	h = hashids.NewWithData(hd)
	return h
}
