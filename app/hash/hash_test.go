package hash

import (
	"math"
	"testing"

	"github.com/matryer/is"
	hashids "github.com/speps/go-hashids"
)

func TestNew(t *testing.T) {
	hd := hashids.NewData()
	hd.Salt = "salt"
	hd.MinLength = 5

	exp := Hash{hashids.NewWithData(hd)}

	is := is.New(t)
	is.Equal(&exp, New("salt", 5))
}

func TestEncode(t *testing.T) {
	hid := New("salt", 0)

	act := hid.DecodeInt64(hid.Encode(math.MaxInt64))

	is := is.New(t)
	is.Equal([]int64{math.MaxInt64}, act)
}
