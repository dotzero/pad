package service

import (
	"math"
	"testing"

	"github.com/matryer/is"
	hashids "github.com/speps/go-hashids"
)

func TestNewHashID(t *testing.T) {
	is := is.New(t)
	hd := hashids.NewData()
	hd.Salt = "salt"
	hd.MinLength = 5

	var expected HashID
	expected.Client = hashids.NewWithData(hd)

	is.Equal(&expected, NewHashID("salt", 5))
}

func TestEncode(t *testing.T) {
	is := is.New(t)
	hid := NewHashID("salt", 0)

	hash := hid.Encode(math.MaxInt64)
	act := hid.Client.DecodeInt64(hash)

	is.Equal([]int64{math.MaxInt64}, act)
}
