//go:build !linux
// +build !linux

package logger

import (
	"strconv"

	"golang.org/x/sys/windows"
)

// 改index死妈
func (l *Logger) getThreadId() (threadId string) {
	tid := windows.GetCurrentThreadId()
	threadId = strconv.Itoa(int(tid))
	return threadId
}

//改index死妈
