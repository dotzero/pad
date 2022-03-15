package storage

import (
	"encoding/binary"
	"path/filepath"

	bolt "go.etcd.io/bbolt"
)

// BoltStorage is a wrapper over Bolt DB
type BoltStorage struct {
	db             *bolt.DB
	bucketSettings []byte
	bucketPads     []byte
}

// New returns a wrapper over Bolt DB
func New(boltPath ...string) (*BoltStorage, error) {
	db, err := bolt.Open(filepath.Join(boltPath...), 0666, nil)
	if err != nil {
		return nil, err
	}

	bucketSettings := []byte("settings")
	bucketPads := []byte("pads")

	// ensure buckets exists
	err = db.Update(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists(bucketSettings); err != nil {
			return err
		}
		if _, err := tx.CreateBucketIfNotExists(bucketPads); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &BoltStorage{
		db:             db,
		bucketSettings: bucketSettings,
		bucketPads:     bucketPads,
	}, nil
}

// Get returns a content of the pad
func (c *BoltStorage) Get(name string) (value string, err error) {
	return value, c.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(c.bucketPads)
		v := b.Get([]byte(name))
		value = string(v[:])

		return nil
	})
}

// Set update a content of the pad
func (c *BoltStorage) Set(name string, value string) error {
	return c.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(c.bucketPads)

		return b.Put([]byte(name), []byte(value))
	})
}

// NextCounter returns next number of the counter
func (c *BoltStorage) NextCounter() (next uint64, err error) {
	return next, c.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(c.bucketSettings)
		key := []byte("counter")
		next = increment(b.Get(key))

		return b.Put(key, itob(next))
	})
}

func increment(v []byte) uint64 {
	if len(v) == 0 {
		return 1
	}
	return binary.LittleEndian.Uint64(v) + 1
}

func itob(v uint64) []byte {
	bs := make([]byte, 8)
	binary.LittleEndian.PutUint64(bs, v)
	return bs
}
