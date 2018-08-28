package errorcode

import "fmt"

var (
	// error code format: 0x[1byte,module_no][1byte,reserve][2byte,error_no]
	// db error, module_no=0x01
	ErrInnerDBConfigError        = NewInnerError(0x01000001, "db config error")
	ErrInnerDBQueryUserMsgError  = NewInnerError(0x01000002, "query user msg error")
	ErrInnerDBFormatSQLError     = NewInnerError(0x01000003, "format sql error")
	ErrInnerDBUpdateUserMsgError = NewInnerError(0x01000004, "update user msg error")
	ErrInnerDBInsertUserMsgError = NewInnerError(0x01000005, "insert user msg error")

	// service error, module_no=0x02
	ErrInnerWrongParam        = NewInnerError(0x02000001, "wrong param")
	ErrInnerInvalidToken      = NewInnerError(0x02000002, "invalid token")
	ErrInnerHTTPReqError      = NewInnerError(0x02000003, "http request error")
	ErrInnerHTTPRespError     = NewInnerError(0x02000004, "http response error")
	ErrInnerQueryUserError    = NewInnerError(0x02000005, "query user error")
	ErrInnerChannelNotSupport = NewInnerError(0x02000006, "channel not support")
	ErrInnerJSONEncodeError   = NewInnerError(0x02000007, "json encode error")
	ErrInnerJSONDecodeError   = NewInnerError(0x02000008, "json decode error")

	_ error = InnerError{}
)

type InnerError struct {
	BaseError
}

func (e InnerError) Error() string {
	return fmt.Sprintf("[0x%08x] %s", e.Ret, e.Msg)
}

func NewInnerError(ret int, msg string) InnerError {
	return InnerError{
		BaseError{
			Ret: ret,
			Msg: msg,
		},
	}
}
