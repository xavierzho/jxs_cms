//go:build windows
// +build windows

package main

import (
	"log"
	"os"
)

// Windows 下暂不支持 syscall.Dup2，可以使用简单写法兼容
func redirectStderr(f *os.File) {
	// 直接替换标准错误输出指针
	err := os.Stderr.Close()
	if err != nil {
		log.Printf("Failed to close old stderr: %v", err)
	}
	os.Stderr = f
}
