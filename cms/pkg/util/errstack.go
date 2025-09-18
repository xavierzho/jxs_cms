package util

import (
	"bytes"
	"fmt"
	"runtime"
)

func GetStackInfo() string {
	kb := 4

	s := []byte("/src/runtime/panic.go")
	e := []byte("\ngoroutine ")
	line := []byte("\n")
	stack := make([]byte, kb<<10) //4KB
	length := runtime.Stack(stack, true)
	start := bytes.Index(stack, s)
	stack = stack[start:length]
	start = bytes.Index(stack, line) + 1
	stack = stack[start:]
	end := bytes.LastIndex(stack, line)
	if end != -1 {
		stack = stack[:end]
	}
	end = bytes.Index(stack, e)
	if end != -1 {
		stack = stack[:end]
	}
	stack = bytes.TrimRight(stack, "\n")

	return string(stack)
}

func PanicErrStack() error {
	if err := recover(); err != nil {
		errorStr := GetStackInfo()
		return fmt.Errorf(errorStr)
	}

	return nil
}

// 安全调用
func SafeCall(f func()) (err error) {
	defer func() {
		if err1 := recover(); err1 != nil {
			errorStr := GetStackInfo()
			//fmt.Println(errorStr)
			err = fmt.Errorf(errorStr)
		}
	}()

	f()

	return
}

// 异步安全调用
func AsyncSafeCall(f func()) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				errorStr := GetStackInfo()
				fmt.Println(errorStr)
			}
		}()

		f()
	}()

}
