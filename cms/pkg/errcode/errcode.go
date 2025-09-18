package errcode

import (
	"fmt"
	"strings"
)

var codes = map[int]string{}

type Error struct {
	// 错误码
	htmlCode int
	code     int
	// 错误消息
	msg string
	// 详细信息
	details []string
}

func newError(htmlCode, code int, msg string) *Error {
	if _, ok := codes[code]; ok {
		panic(fmt.Sprintf("错误码 %d 已经存在，请更换一个", code))
	}
	codes[code] = msg
	return &Error{htmlCode: htmlCode, code: code, msg: msg}
}

// msg + details
func (e *Error) Error() string {
	if len(e.details) == 0 {
		return e.msg
	}
	return e.msg + ": " + strings.Join(e.details, "; ")
}

func (e *Error) HtmlCode() int {
	return e.htmlCode
}

func (e *Error) Code() int {
	return e.code
}

func (e *Error) Msg() string {
	return e.msg
}

// 替换msg
func (e *Error) WithMsg(msg string) *Error {
	newError := *e
	newError.msg = msg

	return &newError
}

func (e *Error) Details() []string {
	return e.details
}

// 追加 details
func (e *Error) WithDetails(details ...string) *Error {
	newError := *e
	newError.details = append(newError.details, details...)

	return &newError
}

// 判断 e 是否是传入的该种 Error
func (e *Error) Is(err *Error) bool {
	return e != nil && e.code == err.code
}
