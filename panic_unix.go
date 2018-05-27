// +build !windows

package log

import (
	"syscall"
)

func dup2(from uintptr) {
	if err := syscall.Dup2(int(from), 1); err != nil {
		panic(err)
	}

	if err := syscall.Dup2(int(from), 2); err != nil {
		panic(err)
	}
}
