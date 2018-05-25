package service

import (
	"encoding/binary"
	"path/filepath"

	"github.com/boltdb/bolt"
)

// BoltBackend is a client to the Bolt DB
type BoltBackend struct {
	db             *bolt.DB
	bucketSettings []byte
	bucketPads     []byte
}

// NewBoltBackend returns a client to the Bolt DB
func NewBoltBackend(boltPath string) (*BoltBackend, error) {
	db, err := bolt.Open(filepath.Join(boltPath, "/pad.db"), 0666, nil)
	if err != nil {
		return nil, err
	}

	bucketSettings := []byte("settings")
	bucketPads := []byte("pads")

	// Ensure buckets exists
	if err = db.Update(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists(bucketSettings); err != nil {
			return err
		}
		if _, err := tx.CreateBucketIfNotExists(bucketPads); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return &BoltBackend{
		db:             db,
		bucketSettings: bucketSettings,
		bucketPads:     bucketPads,
	}, nil
}

// SetPad update a content of the pad in BoltBackend
func (c *BoltBackend) SetPad(name string, value string) error {
	return c.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(c.bucketPads)
		return b.Put([]byte(name), []byte(value))
	})
}

// GetPad returns a content of pad from BoltBackend
func (c *BoltBackend) GetPad(name string) (value string, err error) {
	return value, c.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(c.bucketPads)
		v := b.Get([]byte(name))
		value = string(v[:])
		return nil
	})
}

// GetNextCounter returns next number of counter from BoltBackend
func (c *BoltBackend) GetNextCounter() (next uint64, err error) {
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
