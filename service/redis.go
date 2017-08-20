package service

import (
	"github.com/go-redis/redis"
)

var c *redis.Client

// NewRedisClient returns a client to the Redis Server
func NewRedisClient(uri string) *redis.Client {
	opt, err := redis.ParseURL(uri)
	if err != nil {
		panic(err)
	}

	c = redis.NewClient(opt)
	return c
}
