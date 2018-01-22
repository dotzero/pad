package service

import (
	"log"

	"github.com/go-redis/redis"
)

// Redis is a client to the Redis Server
type Redis struct {
	Prefix string
	Client *redis.Client
}

// NewRedisClient returns a client to the Redis Server
func NewRedisClient(uri string, prefix string) *Redis {
	opt, err := redis.ParseURL(uri)
	if err != nil {
		log.Fatal(err)
	}

	return &Redis{
		Prefix: prefix,
		Client: redis.NewClient(opt),
	}
}

// GetNextCounter returns next number of counter from the Redis Server
func (c *Redis) GetNextCounter() (int64, error) {
	val, err := c.Client.Incr(c.prefixed("#counter#")).Result()
	if err != nil {
		return 0, err
	}

	return val, nil
}

// GetPad returns a content of pad from the Redis Server
func (c *Redis) GetPad(name string) string {
	value, _ := c.Client.Get(c.prefixed(name)).Result()
	return value
}

// SetPad update a content of pad in the Redis Server
func (c *Redis) SetPad(name string, value string) error {
	err := c.Client.Set(c.prefixed(name), value, 0).Err()
	return err
}

func (c *Redis) prefixed(key string) string {
	return c.Prefix + ":" + key
}
