package mocks

import (
	"testing"

	"github.com/matryer/is"
)

func TestNewStorage(t *testing.T) {
	is := is.New(t)
	storage := NewStorage()

	is.Equal(uint64(0), storage.Counter)
	is.Equal(0, len(storage.Content))
}

func TestSetPad(t *testing.T) {
	is := is.New(t)
	storage := NewStorage()

	err := storage.SetPad("foo", "bar")
	is.NoErr(err)

	is.Equal(1, len(storage.Content))
	is.Equal("bar", storage.Content["foo"])
}

func TestGetPad_Exists(t *testing.T) {
	is := is.New(t)
	storage := NewStorage()

	err := storage.SetPad("foo", "bar")
	is.NoErr(err)

	act, err := storage.GetPad("foo")
	is.NoErr(err)
	is.Equal("bar", act)
}

func TestGetPad_NotExists(t *testing.T) {
	is := is.New(t)
	storage := NewStorage()

	act, err := storage.GetPad("foo")
	is.NoErr(err)
	is.Equal("", act)
}

func TestGetNextCounter(t *testing.T) {
	is := is.New(t)
	storage := NewStorage()

	var (
		err error
		cnt uint64
	)

	for n := uint64(1); n <= uint64(10); n++ {
		cnt, err = storage.GetNextCounter()
		is.NoErr(err)
		is.Equal(n, cnt)
	}
}
