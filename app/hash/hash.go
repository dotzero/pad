package hash

import (
	hashids "github.com/speps/go-hashids"
)

// Hash is a wrapper over HashID
type Hash struct {
	*hashids.HashID
}

// New returns a wrapper over HashID
func New(salt string, length int) *Hash {
	hd := hashids.NewData()
	hd.Salt = salt
	hd.MinLength = length

	return &Hash{
		hashids.NewWithData(hd),
	}
}

// Encode returns encoded version of number
func (c *Hash) Encode(n int64) string {
	e, _ := c.EncodeInt64([]int64{n})

	return e
}
