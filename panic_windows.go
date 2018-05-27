package log

import (
	"golang.org/x/sys/windows"
)

func dup2(from uintptr) {
	if err := windows.SetStdHandle(windows.STD_OUTPUT_HANDLE, windows.Handle(from)); err != nil {
		panic(err)
	}

	if err := windows.SetStdHandle(windows.STD_ERROR_HANDLE, windows.Handle(from)); err != nil {
		panic(err)
	}
}
