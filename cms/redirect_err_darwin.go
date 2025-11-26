//go:build darwin
// +build darwin

package main

import (
	"log"
	"os"
	"syscall"
)

// redirectStderr 将标准错误输出重定向到指定文件
func redirectStderr(f *os.File) {
	err := syscall.Dup2(int(f.Fd()), int(os.Stderr.Fd()))
	if err != nil {
		log.Fatalf("Failed to redirect stderr to file: %v", err)
	}
}
