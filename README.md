Golang Log
=============================

format,colorful,level log and redirect escape log.

```bash
 Recommended: [%D %T] [%L] (%S) %M
```

| NO | Format | Details | Output Example|
| :-: | :-: | :- |  :- |
| 1 | %T | Time | 15:04:05 MST |
| 2 | %t | Time | 15:04|
| 3 | %D | Date | 2006/01/02 |
| 4 | %d | Date | 01/02/06 |
| 5 | %L | Level |DEBUG,INFO,WARNING,ERROR|
| 6 | %S | Source | filename,lineno,funcname |
| 7 | %M | Message | output text|

## Level Define

| NO | Level | Details |
| :-: | :-: | :- |
| 1 |DEBUG   | for debug output    |
| 2 |INFO    | for general output  |
| 3 |WARN    | for warnning output |
| 4 |ERROR   | for error output    |
| 5 |FATAL   | for fatal output    |

## Log Install

```bash
go get -u github.com/flyaways/log
```

## Get started

> `Access Format`

```go
package main

import (
	"fmt"
	"os"

	"github.com/flyaways/log"
)

func main() {
	access := log.New(
		log.LogConfig{
			Level:         log.INFO,
			FirstRollover: true,
			Blocking:      false,
			BufferLength:  10240,
			PrefixName:    "access.log",
			When:          log.Hour,
			BackupCount:   72,
			Format:        log.AccessFormat},//log.OprationFormat
	)

	if access == nil {
		fmt.Fprintf(os.Stderr, "init logger error\n")
		return
	}

	access.Info("access info")
}
```

## Examples

* More examples can be found at [github.com/flyaways/log/_examples](https://github.com/flyaways/log/tree/master/_examples).

## Reference

* [github.com/alecthomas/log4go](https://github.com/alecthomas/log4go)
* [github.com/sirupsen/logrus](https://github.com/sirupsen/logrus)

## Lisence

* [Apache License 2.0](https://raw.githubusercontent.com/flyaways/log/master/LICENSE)
