package service

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/boltdb/bolt"
)

func TestNewBoltBackend(t *testing.T) {
	path := tempfile()
	backend, err := NewBoltBackend(path)

	ok(t, err)
	equals(t, []byte("settings"), backend.bucketSettings)
	equals(t, []byte("pads"), backend.bucketPads)
}

func TestSetPad(t *testing.T) {
	backend := newTestBackend()
	defer backend.db.Close()

	err := backend.SetPad("foo", "bar")
	ok(t, err)

	if err := backend.db.View(func(tx *bolt.Tx) error {
		v := tx.Bucket(backend.bucketPads).Get([]byte("foo"))
		equals(t, []byte("bar"), v)
		return nil
	}); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
}

func TestGetPad_Exists(t *testing.T) {
	backend := newTestBackend()
	defer backend.db.Close()

	err := backend.SetPad("foo", "bar")
	ok(t, err)

	v, err := backend.GetPad("foo")
	ok(t, err)
	equals(t, "bar", v)
}

func TestGetPad_NotExists(t *testing.T) {
	backend := newTestBackend()
	defer backend.db.Close()

	v, err := backend.GetPad("foo")
	ok(t, err)
	equals(t, "", v)
}

func TestGetNextCounter(t *testing.T) {
	backend := newTestBackend()
	defer backend.db.Close()

	var (
		err error
		cnt uint64
	)

	for n := uint64(1); n <= uint64(10); n++ {
		cnt, err = backend.GetNextCounter()
		ok(t, err)
		equals(t, n, cnt)
	}
}

func TestIncrement(t *testing.T) {
	equals(t, uint64(1), increment([]byte{}))
	equals(t, uint64(1), increment(itob(0)))
	equals(t, uint64(100), increment(itob(99)))
	equals(t, uint64(1000000000), increment(itob(999999999)))
}

func newTestBackend() *BoltBackend {
	backend, err := NewBoltBackend(tempfile())
	if err != nil {
		panic(err)
	}

	return backend
}

// tempfile returns a temporary file path.
func tempfile() string {
	f, err := ioutil.TempFile("", "bolt-")
	if err != nil {
		panic(err)
	}
	if err := f.Close(); err != nil {
		panic(err)
	}
	if err := os.Remove(f.Name()); err != nil {
		panic(err)
	}
	return f.Name()
}
