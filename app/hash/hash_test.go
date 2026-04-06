package hash

import (
	"math"
	"testing"

	"github.com/matryer/is"
	hashids "github.com/speps/go-hashids/v2"
)

func TestNew(t *testing.T) {
	hd := hashids.NewData()
	hd.Salt = "salt"
	hd.MinLength = 5

	hid, err := hashids.NewWithData(hd)
	is := is.New(t)
	is.NoErr(err)

	exp := Hash{hid}

	is.Equal(&exp, New("salt", 5))
}

func TestEncode(t *testing.T) {
	hid := New("salt", 0)

	act := hid.DecodeInt64(hid.Encode(math.MaxInt64))

	is := is.New(t)
	is.Equal([]int64{math.MaxInt64}, act)
}
