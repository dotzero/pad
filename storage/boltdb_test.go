package storage

import (
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"testing"

	"github.com/matryer/is"
	bolt "go.etcd.io/bbolt"
)

func TestNew(t *testing.T) {
	path := tempfile()
	backend, err := New(path)

	is := is.New(t)
	is.NoErr(err)
	is.Equal([]byte("settings"), backend.bucketSettings)
	is.Equal([]byte("pads"), backend.bucketPads)
}

func suiteSetPad(t *testing.T, b *BoltStorage, key string) {
	is := is.New(t)
	exp := randomString(10)

	err := b.Set(key, exp)
	is.NoErr(err)

	err = b.db.View(func(tx *bolt.Tx) error {
		v := tx.Bucket(b.bucketPads).Get([]byte(key))
		is.Equal([]byte(exp), v)
		return nil
	})
	is.NoErr(err)
}

func TestSet(t *testing.T) {
	backend := newTestBackend()
	defer backend.db.Close()

	suiteSetPad(t, backend, "foo")
}

func TestSetConcurrent(t *testing.T) {
	backend := newTestBackend()
	defer backend.db.Close()

	var wg sync.WaitGroup

	n := rand.Intn(100)
	for i := 0; i < n; i++ {
		wg.Add(1)

		go func(i int, b *BoltStorage) {
			defer wg.Done()
			suiteSetPad(t, b, strconv.Itoa(i))
		}(i, backend)
	}
	wg.Wait()
}

func TestGet(t *testing.T) {
	backend := newTestBackend()
	defer backend.db.Close()

	is := is.New(t)
	exp := randomString(10)

	err := backend.Set("foo", exp)
	is.NoErr(err)

	act, err := backend.Get("foo")
	is.NoErr(err)
	is.Equal(exp, act)
}

func TestGetNotExists(t *testing.T) {
	backend := newTestBackend()
	defer backend.db.Close()

	is := is.New(t)

	act, err := backend.Get("foo")
	is.NoErr(err)
	is.Equal("", act)
}

func TestNextCounter(t *testing.T) {
	backend := newTestBackend()
	defer backend.db.Close()

	is := is.New(t)

	var (
		err error
		cnt uint64
	)

	for n := uint64(1); n <= uint64(10); n++ {
		cnt, err = backend.NextCounter()
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

func newTestBackend() *BoltStorage {
	backend, err := New(tempfile())
	if err != nil {
		panic(err)
	}

	return backend
}

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
