package nzrpc

import (
	"fmt"
)

type Error struct {
	Class int
	Code  int
}

func NewError(class, code int) error {
	return &Error{Class: class, Code: code}
}

func (e Error) Error() string {
	return fmt.Sprintf("nzrpc-error: class=%d, code=%d", e.Class, e.Code)
}
