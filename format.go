package log

import (
	"bytes"
	"fmt"
	"strings"
)

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
