package log

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

const (
	Hour          = "H"
	Minute        = "M"
	Day           = "D"
	Midnight      = "MIDNIGHT"
	AccessFormat  = "[%D %T] [Access] %M"
	OpFormat      = "[%D %T] [%L] (%S) %M"
	MessageFormat = "%M"
	Location      = "Asia/Chongqing"
)

const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
	FATAL
)

type Level int

var (
	RotationErr = errors.New("rotation error")
	LogWriteErr = errors.New("logwrite error")
	ReadDirErr  = errors.New("readdir error")

	conversion = map[string]string{
		/*stdLongMonth      */ "B": "January",
		/*stdMonth          */ "b": "Jan",
		// stdNumMonth       */ "m": "1",
		/*stdZeroMonth      */ "m": "01",
		/*stdLongWeekDay    */ "A": "Monday",
		/*stdWeekDay        */ "a": "Mon",
		// stdDay            */ "d": "2",
		// stdUnderDay       */ "d": "_2",
		/*stdZeroDay        */ "d": "02",
		/*stdHour           */ "H": "15",
		// stdHour12         */ "I": "3",
		/*stdZeroHour12     */ "I": "03",
		// stdMinute         */ "M": "4",
		/*stdZeroMinute     */ "M": "04",
		// stdSecond         */ "S": "5",
		/*stdZeroSecond     */ "S": "05",
		/*stdLongYear       */ "Y": "2006",
		/*stdYear           */ "y": "06",
		/*stdPM             */ "p": "PM",
		// stdpm             */ "p": "pm",
		/*stdTZ             */ "Z": "MST",
		// stdISO8601TZ      */ "z": "Z0700",  // prints Z for UTC
		// stdISO8601ColonTZ */ "z": "Z07:00", // prints Z for UTC
		/*stdNumTZ          */ "z": "-0700", // always numeric
		// stdNumShortTZ     */ "b": "-07",    // always numeric
		// stdNumColonTZ     */ "b": "-07:00", // always numeric
	}

	fCache = &formatCache{}

	levelStrings = [...]string{
		"\x1b[36mDEBG\x1b[0m",
		"\x1b[34mINFO\x1b[0m",
		"\x1b[33mWARN\x1b[0m",
		"\x1b[31mEROR\x1b[0m",
		"\x1b[35mFATAL\x1b[0m"}
)

func (l Level) String() string {
	if l < 0 || int(l) > len(levelStrings) {
		return "\x1b[37mUNKN\x1b[0m"
	}
	return levelStrings[int(l)]
}

type formatCache struct {
	LastUpdateSeconds int64
	shortTime         string
	shortDate         string
	longTime          string
	longDate          string
}

func formatRecord(format string, rec *Record) string {
	if rec == nil {
		return "<nil>"
	}

	if len(format) == 0 {
		format = "[%D %T] [%L] (%S) %M"
	}

	out := bytes.NewBuffer(make([]byte, 0, 64))
	secs := rec.Created.UnixNano() / 1e9

	cache := *fCache
	if cache.LastUpdateSeconds != secs {
		month, day, year := rec.Created.Month(), rec.Created.Day(), rec.Created.Year()
		hour, minute, second := rec.Created.Hour(), rec.Created.Minute(), rec.Created.Second()
		zone, _ := rec.Created.Zone()
		updated := &formatCache{
			LastUpdateSeconds: secs,
			shortTime:         fmt.Sprintf("%02d:%02d", hour, minute),
			shortDate:         fmt.Sprintf("%02d/%02d/%02d", day, month, year%100),
			longTime:          fmt.Sprintf("%02d:%02d:%02d %s", hour, minute, second, zone),
			longDate:          fmt.Sprintf("%04d/%02d/%02d", year, month, day),
		}
		cache = *updated
		fCache = updated
	}

	pieces := bytes.Split([]byte(format), []byte{'%'})

	for i, piece := range pieces {
		if i > 0 && len(piece) > 0 {
			switch piece[0] {
			case 'T':
				out.WriteString(cache.longTime)
			case 't':
				out.WriteString(cache.shortTime)
			case 'D':
				out.WriteString(cache.longDate)
			case 'd':
				out.WriteString(cache.shortDate)
			case 'L':
				out.WriteString(rec.Level.String())
			case 'S':
				out.WriteString(rec.Source)
			case 's':
				slice := strings.Split(rec.Source, "/")
				out.WriteString(slice[len(slice)-1])
			case 'M':
				out.WriteString(rec.Message)
			}
			if len(piece) > 1 {
				out.Write(piece[1:])
			}
		} else if len(piece) > 0 {
			out.Write(piece)
		}
	}
	out.WriteByte('\n')

	return out.String()
}

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
					fmt.Fprintf(os.Stderr, "%v: %v\n", RotationErr, err)
					return
				}
			}

			_, err := fmt.Fprint(w.focus, formatRecord(w.Format, rec))
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v: %v\n", LogWriteErr, err)
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
		fmt.Fprintf(os.Stderr, "%v: %v\n", ReadDirErr, err)
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
			fmt.Fprintf(os.Stderr, "%v: %v\n", RotationErr, err)
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
		fmt.Fprintf(os.Stderr, "%v: %v\n", RotationErr, err)
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
