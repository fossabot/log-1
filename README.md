Golang Log
=============================

format,colorful,level log and redirect escape log.

[![asciinema](https://asciinema.org/a/NYcaJUSKPOH6xzwLCTG5NBtX2.png)](https://asciinema.org/a/NYcaJUSKPOH6xzwLCTG5NBtX2?autoplay=1)

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

```go
package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/flyaways/log"
)

func init() {
	//Collector for escape ouput
	log.Collector()
}

func main() {
	op := log.New(
		log.LogConfig{
			Level:         log.INFO,
			FirstRollover: true,
			Blocking:      false,
			BufferLength:  10240,
			PrefixName:    "op.log",
			When:          log.Hour,
			BackupCount:   72,
			Format:        log.OpFormat}, //log.AccessFormat
	)

	if op == nil {
		fmt.Fprintf(os.Stderr, "init logger error\n")
		return
	}

	op.Debug("Go is an open source project developed by a team at Google and many contributors from the open source community.")
	op.Info("Go is distributed under a BSD-style license.")
	op.Warn("A low traffic mailing list for important announcements, such as new releases.")
	op.Error("We encourage all Go users to subscribe to golang-announce.")
	//op.Fatal("Go is an open source programming language that makes it easy to build simple, reliable, and efficient software.")

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	<-signals
}

```

* More examples can be found at [github.com/flyaways/log/_examples](https://github.com/flyaways/log/tree/master/_examples).

## Reference

* [github.com/alecthomas/log4go](https://github.com/alecthomas/log4go)
* [github.com/sirupsen/logrus](https://github.com/sirupsen/logrus)

## Lisence

* [Apache License 2.0](https://raw.githubusercontent.com/flyaways/log/master/LICENSE)
