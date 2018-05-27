package log

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type Logger struct {
	Level
	Writer
}

func (log Logger) high(lvl Level, arg0 interface{}, args ...interface{}) {
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
	log.high(DEBUG, arg0, args)
}

func (log Logger) Info(arg0 interface{}, args ...interface{}) {
	log.high(INFO, arg0, args)
}

func (log Logger) Warn(arg0 interface{}, args ...interface{}) {
	log.high(WARN, arg0, args)
}

func (log Logger) Error(arg0 interface{}, args ...interface{}) {
	log.high(ERROR, arg0, args)
}

func (log Logger) Fatal(arg0 interface{}, args ...interface{}) {
	log.high(FATAL, arg0, args)
}
