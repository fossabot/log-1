package log

import (
	"log"
	"os"
	"syscall"
	"time"
)

func midnight(now time.Time) int64 {
	loc, _ := time.LoadLocation(Location)
	last, _ := time.ParseInLocation("2006-01-02", now.Format("2006-01-02"), loc)
	return last.AddDate(0, 0, 1).Unix()
}

// Collector for redirect escape log.
func Collector() {
	Panic()
	Std()
}

func Std() {
	fd, _ := os.OpenFile("std.log", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	os.Stdout = fd
	os.Stderr = fd
	log.SetOutput(fd)
}

func Panic() {
	fd, _ := os.OpenFile("panic.log", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	syscall.Dup2(int(fd.Fd()), 1)
	syscall.Dup2(int(fd.Fd()), 2)
}
