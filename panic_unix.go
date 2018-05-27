// +build !windows

package log

import (
	"os"
	"syscall"
)

func dup2(from uintptr) {
	if err := syscall.Dup2(int(from), int(os.Stderr.Fd())); err != nil {
		panic(err)
	}

	if err := syscall.Dup2(int(from), int(os.Stdout.Fd())); err != nil {
		panic(err)
	}
}
