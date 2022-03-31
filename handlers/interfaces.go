package handlers

import (
	"io"
)

type hashEncoder interface {
	Encode(num int64) string
}

type padStorage interface {
	Get(name string) (value string, err error)
	Set(name string, value string) error
	NextCounter() (next uint64, err error)
}

type template interface {
	Execute(wr io.Writer, name string, data interface{}) error
}
