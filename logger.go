package log

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

func New(logconfig LogConfig) *Logger {
	w := &LogWriter{
		rec:       make(chan *Record, logconfig.BufferLength),
		isEnd:     make(chan bool),
		LogConfig: logconfig,
	}

	w.parameterization()

	if err := w.rotate(); err != nil {
		return nil
	}

	go w.run()

	return &Logger{
		Level:  logconfig.Level,
		Writer: w,
	}
}

type Logger struct {
	Level
	Writer
}

func (log Logger) low(lvl Level, arg0 interface{}, args ...interface{}) {
	if lvl < log.Level {
		return
	}

	format := ""
	switch first := arg0.(type) {
	case string:
		format = first
	case func() string:
		format = first()
	default:
		format = fmt.Sprint(arg0) + strings.Repeat(" %v", len(args))
	}

	pc, file, lineno, ok := runtime.Caller(2)
	src := ""
	if ok {
		src = fmt.Sprintf("%s:%s:%d", filepath.Base(file), runtime.FuncForPC(pc).Name(), lineno)
	}

	msg := format
	if len(args) > 0 {
		msg = fmt.Sprintf(format, args...)
	}

	fmt.Println(args, len(args))

	log.Write(&Record{
		Level:   lvl,
		Created: time.Now(),
		Source:  src,
		Message: msg,
	})

	if lvl == FATAL {
		panic(msg)
	}
}

func (log Logger) Debug(arg0 interface{}, args ...interface{}) {
	log.low(DEBUG, arg0, args...)
}

func (log Logger) Info(arg0 interface{}, args ...interface{}) {
	log.low(INFO, arg0, args...)
}

func (log Logger) Warn(arg0 interface{}, args ...interface{}) {
	log.low(WARN, arg0, args...)
}

func (log Logger) Error(arg0 interface{}, args ...interface{}) {
	log.low(ERROR, arg0, args...)
}

func (log Logger) Fatal(arg0 interface{}, args ...interface{}) {
	log.low(FATAL, arg0, args...)
}

func midnight(now time.Time) int64 {
	loc, _ := time.LoadLocation(Location)
	last, _ := time.ParseInLocation("2006-01-02", now.Format("2006-01-02"), loc)
	return last.AddDate(0, 0, 1).Unix()
}

// Collector for redirect escape log.
func Collector() {
	escapePanic()
	escapeStd()
}

func escapeStd() {
	fd, _ := os.OpenFile("std.log", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	os.Stdout = fd
	os.Stderr = fd
	log.SetOutput(fd)
}

func escapePanic() {
	fd, err := os.OpenFile("panic.log", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}

	dup2(fd.Fd())
}
