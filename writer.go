package log

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

type Record struct {
	Level   Level
	Created time.Time
	Source  string
	Message string
}

type Writer interface {
	Write(rec *Record)
	Close()
}

type LogConfig struct {
	When          string
	BackupCount   int
	FirstRollover bool
	Blocking      bool
	BufferLength  int
	PrefixName    string
	Format        string
	Level         Level
}

type LogWriter struct {
	LogConfig

	rec        chan *Record
	isEnd      chan bool
	focus      *os.File
	rolloverAt int64
	interval   int64
	suffix     string
	filter     *regexp.Regexp
}

func (w *LogWriter) Write(rec *Record) {
	if !w.Blocking {
		if len(w.rec) >= w.BufferLength {
			fmt.Println("Log_buffer_overflow", w.BufferLength)
			return
		}
	}

	w.rec <- rec
}

func (w *LogWriter) Close() {
	w.waitForEnd(w.rec)
	close(w.rec)
	w.focus.Sync()
}

func (w *LogWriter) run() {
	defer func() {
		if w.focus != nil {
			w.focus.Close()
		}
	}()

	for {
		select {
		case rec, ok := <-w.rec:
			if !ok {
				return
			}

			if w.endNotify(rec) {
				return
			}

			if time.Now().Unix() >= w.rolloverAt {
				if err := w.rotate(); err != nil {
					fmt.Fprintf(os.Stderr, "%v: %v\n", RotationError, err)
					return
				}
			}

			_, err := fmt.Fprint(w.focus, formatRecord(w.Format, rec))
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v: %v\n", LogWriteError, err)
				return
			}
		}
	}
}

func (w *LogWriter) rollover(now time.Time) int64 {
	if w.FirstRollover {
		w.FirstRollover = false
		return (now.Unix()/w.interval + 1) * w.interval
	}

	if w.When == "MIDNIGHT" {
		return midnight(now)
	}

	return (now.Unix()/w.interval + 1) * w.interval
}

func (w *LogWriter) parameterization() {
	var regRule string

	switch w.When {
	case "M":
		w.interval = 60
		w.suffix = "%Y-%m-%d_%H-%M"
		regRule = `^\d{4}-\d{2}-\d{2}_\d{2}-\d{2}$`
	case "H":
		w.interval = 60 * 60
		w.suffix = "%Y%m%d%H"
		regRule = `^\d{10}$`
	case "D", "MIDNIGHT":
		w.interval = 60 * 60 * 24
		w.suffix = "%Y-%m-%d"
		regRule = `^\d{4}-\d{2}-\d{2}$`
	default:
		w.interval = 60 * 60 * 24
		w.suffix = "%Y-%m-%d"
		regRule = `^\d{4}-\d{2}-\d{2}$`
	}
	w.filter = regexp.MustCompile(regRule)

	fInfo, err := os.Stat(w.PrefixName)

	var t time.Time
	if err == nil {
		t = fInfo.ModTime()
	} else {
		t = time.Now()
	}

	w.rolloverAt = (t.Unix()/w.interval + 1) * w.interval
}

func (w *LogWriter) delList() []string {
	dirName := filepath.Dir(w.PrefixName)
	baseName := filepath.Base(w.PrefixName)

	result := []string{}

	fileInfos, err := ioutil.ReadDir(dirName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v: %v\n", ReadDirError, err)
		return result
	}

	prefix := baseName + "."
	plen := len(prefix)

	for _, fileInfo := range fileInfos {
		prefixname := fileInfo.Name()
		if len(prefixname) >= plen {
			if prefixname[:plen] == prefix {
				suffix := prefixname[plen:]
				if w.filter.MatchString(suffix) {
					result = append(result, filepath.Join(dirName, prefixname))
				}
			}
		}
	}

	sort.Sort(sort.StringSlice(result))

	if len(result) < w.BackupCount {
		result = result[0:0]
	} else {
		result = result[:len(result)-w.BackupCount]
	}
	return result
}

func (w *LogWriter) backup() error {
	_, err := os.Lstat(w.PrefixName)
	if err != nil {
		return nil
	}

	finalname := strings.Join([]string{
		w.PrefixName,
		".",
		w.format(w.suffix,
			time.Unix(w.rolloverAt-w.interval, 0).Local())},
		"")

	if _, err := os.Stat(finalname); err == nil {
		return nil
	}

	if err = os.Rename(w.PrefixName, finalname); err != nil {
		return err
	}

	return nil
}

func (w *LogWriter) rotate() (err error) {
	if w.focus != nil {
		w.focus.Close()
	}

	now := time.Now()

	if now.Unix() >= w.rolloverAt {
		if err = w.backup(); err != nil {
			fmt.Fprintf(os.Stderr, "%v: %v\n", RotationError, err)
			return err
		}
	}

	if w.BackupCount > 0 {
		for _, prefixname := range w.delList() {
			os.Remove(prefixname)
		}
	}

	w.focus, err = os.OpenFile(w.PrefixName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v: %v\n", RotationError, err)
		return err
	}

	newRolloverAt := w.rollover(now)
	for newRolloverAt <= now.Unix() {
		newRolloverAt = newRolloverAt + w.interval
	}

	w.rolloverAt = newRolloverAt

	return nil
}

func (w *LogWriter) format(format string, t time.Time) string {
	var layout []string
	for _, chunk := range strings.Split(format, "%") {
		if len(chunk) == 0 {
			continue
		}

		if layoutCmd, ok := conversion[chunk[0:1]]; ok {
			layout = append(layout, layoutCmd)
			if len(chunk) > 1 {
				layout = append(layout, chunk[1:])
			}
			continue
		}

		layout = append(layout, "%", chunk)
	}
	return t.Format(strings.Join(layout, ""))
}

func (w *LogWriter) endNotify(lr *Record) bool {
	if lr == nil && w.isEnd != nil {
		w.isEnd <- true
		return true
	}
	return false
}

func (w *LogWriter) waitForEnd(rec chan *Record) {
	rec <- nil
	if w.isEnd != nil {
		<-w.isEnd
	}
}
