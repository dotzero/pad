package service

import (
	"fmt"
	"testing"

	"github.com/alicebob/miniredis"
	"github.com/go-redis/redis"
)

func TestNewRedisClient(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer s.Close()

	exp := redis.NewClient(&redis.Options{
		Addr:     s.Addr(),
		Password: "",
		DB:       0,
	})

	act := NewRedisClient(fmt.Sprintf("redis://%s/0", s.Addr()), "prefix")

	equals(t, "prefix", act.Prefix)
	equals(t, exp.String(), act.Client.String())
}

func TestGetNextCounter(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer s.Close()

	s.Set("prefix:#counter#", "99")

	r := NewRedisClient(fmt.Sprintf("redis://%s/0", s.Addr()), "prefix")

	act, _ := r.GetNextCounter()
	equals(t, int64(100), act)
}

func TestGetPad(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer s.Close()

	s.Set("prefix:foo", "bar")

	r := NewRedisClient(fmt.Sprintf("redis://%s/0", s.Addr()), "prefix")

	act, err := r.GetPad("foo")
	ok(t, err)
	equals(t, "bar", act)
}

func TestSetPad(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer s.Close()

	r := NewRedisClient(fmt.Sprintf("redis://%s/0", s.Addr()), "prefix")
	err = r.SetPad("foo", "bar")
	ok(t, err)

	act, _ := s.Get("prefix:foo")
	equals(t, "bar", act)
}

func TestPrefixed(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer s.Close()

	r := NewRedisClient(fmt.Sprintf("redis://%s/0", s.Addr()), "prefix")
	equals(t, "prefix:foo", r.prefixed("foo"))
}
