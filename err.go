package nzrpc

import (
	"fmt"
)

const (
	ErrClassRpc = 1
	ErrClassCli = 2
	// ErrClassSrv = 3
	ErrClassApp = 4
)

const (
	ErrCodeNotLogin    = 1001
	ErrCodeNoUser      = 1002
	ErrCodeAuthFail    = 1003
	ErrCodeNoStream    = 1004
	ErrCodeNotImpl     = 1005
	ErrCodeNoChallenge = 1006
	ErrCodeBadCmd      = 1007
	ErrCodeBadMagic    = 1008

	ErrCodeBadPassword = 1013
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
