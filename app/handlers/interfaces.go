package handlers

//go:generate moq -skip-ensure -out mock.go . encoder storage

import (
	"io"
)

type encoder interface {
	Encode(num int64) string
}

type storage interface {
	Get(name string) (value string, err error)
	Set(name string, value string) error
	NextCounter() (next uint64, err error)
}

type tpl interface {
	Execute(wr io.Writer, data interface{}) error
}
