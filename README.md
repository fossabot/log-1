# Golang Log
=============================

**Note**: The log of level,format(struct),colorful,and can catch escape output.

<!-- TOC -->

    - [Function List](#function-list)
    - [Format Define](#format-define)
    - [Level Define](#level-define)
    - [Log Install](#log-install)
    - [Get started](#get-started)
    - [Examples](#examples)
    - [Reference](#reference)
    - [Lisence](#lisence)
    
<!-- /TOC -->

## Function List

| NO | Function | Details | Remarks|
| :-: | :- | :- |  :- |
| 1 | format| format just like fmt | contain date,time and source(file,lineno,funcname) |
| 2 | color | level color | different level have different keyword of level color|
| 3 | redirect | redirect escape log | panic to panic.log, redirect stdout and stderr to std.log|
| 4 | log level | diffrent case use diffrent level | DEBUG,INFO,WARNING,ERROR|

## Format Define

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
| 3 |WARNING | for warnning output |
| 4 |ERROR   | for error output    |

## Log Install

```bash
go get -u github.com/flyaways/log
```

## Get started

> Create `access` log

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
			Format:        log.AccessFormat},
	)

	if access == nil {
		fmt.Fprintf(os.Stderr, "init logger error\n")
		return
	}

	access.Info("access info")
}
```

> Create `business` log

```go
package main

import (
	"fmt"
	"os"

	"github.com/flyaways/log"
)

func init() {
	log.Collector()
}

func main() {
	operation := log.New(
		log.LogConfig{
			Level:         log.INFO,
			FirstRollover: true,
			Blocking:      false,
			BufferLength:  10240,
			PrefixName:    "operation.log",
			When:          log.Hour,
			BackupCount:   72,
			Format:        log.OprationFormat},
	)

	if operation == nil {
		fmt.Fprintf(os.Stderr, "init logger error\n")
		return
	}

	operation.Debug("Debug")
	operation.Info("Info")
	operation.Warn("Warn")
	operation.Error("Error")
}
```

> Create `messageonly` log

```go
package main

import (
	"fmt"
	"os"

	"github.com/flyaways/log"
)

func main() {
	messageonly := log.New(
		log.LogConfig{
			Level:         log.INFO,
			FirstRollover: true,
			Blocking:      false,
			BufferLength:  10240,
			PrefixName:    "messageonly.log",
			When:          log.Hour,
			BackupCount:   72,
			Format:        log.MessageOnly},
	)

	if messageonly == nil {
		fmt.Fprintf(os.Stderr, "init logger error\n")
		return
	}

	messageonly.Debug("Debug")
	messageonly.Info("Info")
	messageonly.Warn("Warn")
	messageonly.Error("Warn")
}
```

## Examples

* More examples can be found at [github.com/flyaways/log/examples](https://github.com/flyaways/log/tree/master/examples).

## Reference

* [log4go](https://github.com/alecthomas/log4go)
* [logrus](https://github.com/sirupsen/logrus)

## Lisence

* [Apache License 2.0](https://raw.githubusercontent.com/flyaways/log/master/LICENSE)
