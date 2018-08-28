package errorcode

import "fmt"

var (
	// Global Error
	OK                  = NewError(0, "ok")
	ErrBadRequest       = NewError(4001001, "bad request")
	ErrAuthExpired      = NewError(4001002, "auth expired")
	ErrClientDeprecated = NewError(4001003, "this version of client is deprecated")
	ErrServerError      = NewError(5001001, "server error")
	ErrServerBusy       = NewError(5001002, "server busy")
	ErrNoModuleInfo     = NewError(5001003, "no such module")

	_ = (error)(BaseError{})
)

type BaseError struct {
	Ret int
	Msg string
}

func (e BaseError) Error() string {
	return fmt.Sprintf("(%d)%s", e.Ret, e.Msg)
}

func NewError(ret int, msg string) BaseError {
	return BaseError{
		Ret: ret,
		Msg: msg,
	}
}
