package service

import (
	hashids "github.com/speps/go-hashids"
)

// HashID is a client to the HashID
type HashID struct {
	Client *hashids.HashID
}

// NewHashID returns a HashID client
func NewHashID(salt string, length int) *HashID {
	hd := hashids.NewData()
	hd.Salt = salt
	hd.MinLength = length
	return &HashID{
		Client: hashids.NewWithData(hd),
	}
}

// Encode returns encoded version of number
func (c *HashID) Encode(num int64) string {
	e, _ := c.Client.EncodeInt64([]int64{num})
	return e
}
