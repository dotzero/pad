package service

import (
	"math"
	"testing"

	hashids "github.com/speps/go-hashids"
)

func TestNewHashID(t *testing.T) {
	hd := hashids.NewData()
	hd.Salt = "salt"
	hd.MinLength = 5

	var expected HashID
	expected.Client = hashids.NewWithData(hd)

	equals(t, &expected, NewHashID("salt", 5))
}

func TestEncode(t *testing.T) {
	hid := NewHashID("salt", 0)

	hash := hid.Encode(math.MaxInt64)
	act := hid.Client.DecodeInt64(hash)

	equals(t, []int64{math.MaxInt64}, act)
}
