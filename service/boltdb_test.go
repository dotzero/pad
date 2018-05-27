package service

import (
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"testing"

	"github.com/coreos/bbolt"
	"github.com/matryer/is"
)

func SuiteSetPad(t *testing.T, b *BoltBackend, key string) {
	is := is.New(t)
	exp := randomString(10)
	err := b.SetPad(key, exp)
	is.NoErr(err)

	if err := b.db.View(func(tx *bolt.Tx) error {
		v := tx.Bucket(b.bucketPads).Get([]byte(key))
		is.Equal([]byte(exp), v)
		return nil
	}); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
}

func TestNewBoltBackend(t *testing.T) {
	is := is.New(t)
	path := tempfile()
	backend, err := NewBoltBackend(path)

	is.NoErr(err)
	is.Equal([]byte("settings"), backend.bucketSettings)
	is.Equal([]byte("pads"), backend.bucketPads)
}

func TestSetPad(t *testing.T) {
	backend := newTestBackend()
	defer backend.db.Close()

	SuiteSetPad(t, backend, "foo")
}

func TestSetPad_Concurrent(t *testing.T) {
	backend := newTestBackend()
	defer backend.db.Close()

	var wg sync.WaitGroup

	n := rand.Intn(100)
	for i := 0; i < n; i++ {
		wg.Add(1)

		go func(i int, b *BoltBackend) {
			defer wg.Done()
			SuiteSetPad(t, b, strconv.Itoa(i))
		}(i, backend)
	}
	wg.Wait()
}

func TestGetPad_Exists(t *testing.T) {
	is := is.New(t)
	backend := newTestBackend()
	defer backend.db.Close()

	exp := randomString(10)
	err := backend.SetPad("foo", exp)
	is.NoErr(err)

	act, err := backend.GetPad("foo")
	is.NoErr(err)
	is.Equal(exp, act)
}

func TestGetPad_NotExists(t *testing.T) {
	is := is.New(t)
	backend := newTestBackend()
	defer backend.db.Close()

	act, err := backend.GetPad("foo")
	is.NoErr(err)
	is.Equal("", act)
}

func TestGetNextCounter(t *testing.T) {
	is := is.New(t)
	backend := newTestBackend()
	defer backend.db.Close()

	var (
		err error
		cnt uint64
	)

	for n := uint64(1); n <= uint64(10); n++ {
		cnt, err = backend.GetNextCounter()
		is.NoErr(err)
		is.Equal(n, cnt)
	}
}

func TestIncrement(t *testing.T) {
	is := is.New(t)
	is.Equal(uint64(1), increment([]byte{}))
	is.Equal(uint64(1), increment(itob(0)))
	is.Equal(uint64(100), increment(itob(99)))
	is.Equal(uint64(1000000000), increment(itob(999999999)))
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

func randomString(n int) string {
	letterBytes := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
