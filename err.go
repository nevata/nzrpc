package nzrpc

import (
	"fmt"
)

const (
	Success = 0

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
	Call  string
	Class int
	Code  int
	Msg   string
}

func NewError(call string, class, code int, msg string) error {
	return &Error{Call: call, Class: class, Code: code, Msg: msg}
}

func (e Error) Error() string {
	return fmt.Sprintf("nzrpc-error: class=%d, code=%d", e.Class, e.Code)
}

func checkError(result []interface{}) error {
	call := result[0].(string)
	flag := int(result[1].(float64))
	switch flag {
	case Success:
		return nil
	case ErrClassRpc, ErrClassCli, ErrClassApp:
		class := flag
		code := int(result[2].(float64))
		msg := ""
		if len(result) > 3 {
			msg = result[3].(string)
		}
		return NewError(
			call,
			class,
			code,
			msg,
		)
	default:
		code := int(result[1].(float64))
		msg := ""
		if len(result) > 2 {
			msg = result[2].(string)
		}
		return NewError(
			call,
			0,
			code,
			msg,
		)
	}
}
